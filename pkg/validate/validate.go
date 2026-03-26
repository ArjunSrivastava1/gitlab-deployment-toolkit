package validate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ValidationResult struct {
	Valid  bool
	Errors []string
}

type DeploymentManifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas *int32 `yaml:"replicas"`
		Selector struct {
			MatchLabels map[string]string `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels map[string]string `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []struct {
					Name  string `yaml:"name"`
					Image string `yaml:"image"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

func ValidateManifest(filePath string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:  true,
		Errors: []string{},
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var manifest DeploymentManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if manifest.Kind != "Deployment" {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid kind: %s, expected Deployment", manifest.Kind))
	}

	if manifest.Metadata.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "metadata.name is required")
	}

	if manifest.Spec.Replicas == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "spec.replicas is required")
	} else if *manifest.Spec.Replicas < 1 {
		result.Valid = false
		result.Errors = append(result.Errors, "spec.replicas must be at least 1")
	}

	if len(manifest.Spec.Template.Spec.Containers) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "At least one container is required")
	}

	for i, container := range manifest.Spec.Template.Spec.Containers {
		if container.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Container %d: name is required", i))
		}
		if container.Image == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Container %d: image is required", i))
		}
	}

	return result, nil
}
