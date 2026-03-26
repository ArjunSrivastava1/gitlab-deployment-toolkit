package terraform

import (
	"fmt"
	"strings"
)

func GenerateDeploymentSnippet(name string) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`resource "kubernetes_deployment" "%s" {
  metadata {
    name = "%s"
    labels = {
      app = "%s"
    }
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        app = "%s"
      }
    }

    template {
      metadata {
        labels = {
          app = "%s"
        }
      }

      spec {
        container {
          image = "nginx:latest"
          name  = "%s-container"

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }

          port {
            container_port = 80
          }
        }
      }
    }
  }
}
`, name, name, name, name, name, name))

	sb.WriteString(fmt.Sprintf(`
resource "kubernetes_service" "%s" {
  metadata {
    name = "%s-service"
  }

  spec {
    selector = {
      app = "%s"
    }

    port {
      port        = 80
      target_port = 80
    }

    type = "ClusterIP"
  }
}
`, name, name, name))

	return sb.String(), nil
}
