package terraform

import (
	"strings"
	"testing"
)

func TestGenerateDeploymentSnippet(t *testing.T) {
	name := "test-app"
	output, err := GenerateDeploymentSnippet(name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the output contains expected strings
	expectedStrings := []string{
		`resource "kubernetes_deployment" "test-app"`,
		`resource "kubernetes_service" "test-app"`,
		`replicas = 3`,
		`image = "nginx:latest"`,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q", expected)
		}
	}
}

func TestGenerateDeploymentSnippet_EmptyName(t *testing.T) {
	output, err := GenerateDeploymentSnippet("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if output == "" {
		t.Errorf("Expected non-empty output for empty name")
	}
}
