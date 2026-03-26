package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// DeploymentStatus represents the health status of a deployment
type DeploymentStatus struct {
	Name            string
	Namespace       string
	Ready           bool
	DesiredReplicas int32
	ReadyReplicas   int32
	Conditions      []appsv1.DeploymentCondition
	Message         string
}

// CheckDeploymentReadiness validates a Kubernetes deployment
func CheckDeploymentReadiness(namespace, name string) (*DeploymentStatus, error) {
	// Get Kubernetes client
	clientset, err := getKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Get the deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s/%s: %w", namespace, name, err)
	}

	// Build status
	status := &DeploymentStatus{
		Name:            deployment.Name,
		Namespace:       deployment.Namespace,
		DesiredReplicas: *deployment.Spec.Replicas,
		ReadyReplicas:   deployment.Status.ReadyReplicas,
		Conditions:      deployment.Status.Conditions,
	}

	// Check if deployment is ready
	if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
		status.Ready = true
		status.Message = "Deployment is healthy and all replicas are ready"
	} else {
		status.Ready = false
		status.Message = fmt.Sprintf("Deployment not ready: %d/%d replicas ready",
			deployment.Status.ReadyReplicas,
			*deployment.Spec.Replicas,
		)
	}

	// Check for failing conditions
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == corev1.ConditionFalse {
			status.Message += fmt.Sprintf(" | Progressing condition: %s", condition.Message)
		}
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == corev1.ConditionFalse {
			status.Message += fmt.Sprintf(" | Not available: %s", condition.Message)
		}
	}

	return status, nil
}

// ListDeployments returns all deployments in a namespace
func ListDeployments(namespace string) ([]string, error) {
	clientset, err := getKubernetesClient()
	if err != nil {
		return nil, err
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(
		context.TODO(),
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, deploy := range deployments.Items {
		names = append(names, deploy.Name)
	}
	return names, nil
}

// getKubernetesClient creates a Kubernetes client using kubeconfig
func getKubernetesClient() (*kubernetes.Clientset, error) {
	// Try in-cluster config first (for when running inside a pod)
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig file
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		// Allow override via environment variable
		if envConfig := os.Getenv("KUBECONFIG"); envConfig != "" {
			kubeconfig = envConfig
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return clientset, nil
}
