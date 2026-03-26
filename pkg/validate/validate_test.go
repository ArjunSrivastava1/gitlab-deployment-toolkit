package validate

import (
	"os"
	"testing"
)

func TestValidateManifest_Valid(t *testing.T) {
	validYAML := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: test
        image: nginx:latest
`

	tmpfile, err := os.CreateTemp("", "valid-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(validYAML)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateManifest(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid manifest, got errors: %v", result.Errors)
	}
}

func TestValidateManifest_MissingReplicas(t *testing.T) {
	invalidYAML := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
spec:
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: test
        image: nginx:latest
`

	tmpfile, err := os.CreateTemp("", "invalid-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(invalidYAML)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateManifest(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Valid {
		t.Errorf("Expected invalid manifest, but validation passed")
	}

	found := false
	for _, errMsg := range result.Errors {
		if errMsg == "spec.replicas is required" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected error about missing replicas, got: %v", result.Errors)
	}
}
