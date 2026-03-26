package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gitlab-deployment-toolkit/pkg/k8s"
	"gitlab-deployment-toolkit/pkg/terraform"
	"gitlab-deployment-toolkit/pkg/validate"
)

var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print the version of the toolkit")
	k8sFlag := flag.String("check-k8s", "", "Check readiness of a Kubernetes deployment (format: namespace/deployment-name)")
	listFlag := flag.String("list-deployments", "", "List all deployments in a namespace")
	validateFlag := flag.String("validate", "", "Validate a deployment manifest file")
	terraformFlag := flag.String("generate-tf", "", "Generate Terraform snippet for a deployment (provide name)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "GitLab Deployment Toolkit (GDT) - Kubernetes Deployment Validator\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -check-k8s default/nginx-deployment\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -list-deployments default\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -validate examples/example-config.yaml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -generate-tf my-app\n", os.Args[0])
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("gitlab-deployment-toolkit version %s\n", Version)
		os.Exit(0)
	}

	if *listFlag != "" {
		names, err := k8s.ListDeployments(*listFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deployments in namespace '%s':\n", *listFlag)
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
		os.Exit(0)
	}

	if *k8sFlag != "" {
		parts := strings.SplitN(*k8sFlag, "/", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Invalid format. Use namespace/deployment-name\n")
			os.Exit(1)
		}
		namespace, deploymentName := parts[0], parts[1]

		fmt.Printf("🔍 Checking deployment %s in namespace %s...\n\n", deploymentName, namespace)

		status, err := k8s.CheckDeploymentReadiness(namespace, deploymentName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		if status.Ready {
			fmt.Printf("✅ %s\n", status.Message)
		} else {
			fmt.Printf("⚠️  %s\n", status.Message)
		}
		fmt.Printf("📊 Replicas: %d/%d ready\n", status.ReadyReplicas, status.DesiredReplicas)

		os.Exit(0)
	}

	if *validateFlag != "" {
		result, err := validate.ValidateManifest(*validateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Validation error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("📋 Validating manifest: %s\n\n", *validateFlag)
		if result.Valid {
			fmt.Printf("✅ Validation passed!\n")
		} else {
			fmt.Printf("❌ Validation failed!\n")
		}
		fmt.Printf("Errors found: %d\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
		os.Exit(0)
	}

	if *terraformFlag != "" {
		tfCode, err := terraform.GenerateDeploymentSnippet(*terraformFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Generation error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("// Terraform configuration for deployment: %s\n", *terraformFlag)
		fmt.Println(tfCode)
		os.Exit(0)
	}

	flag.Usage()
	os.Exit(1)
}
