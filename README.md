# 🚀 GitLab Deployment Toolkit (GDT)

[![CI Pipeline](https://github.com/ArjunSrivastava1/gitlab-deployment-toolkit/actions/workflows/ci.yml/badge.svg)](https://github.com/ArjunSrivastava1/gitlab-deployment-toolkit/actions)
[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.29-blue)](https://kubernetes.io/)
> A production-ready Go CLI tool for validating Kubernetes deployments and generating Infrastructure as Code. Built for platform engineers who need reliable deployment automation.

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **☸️ Kubernetes Validation** | Real-time deployment health checks using Kubernetes client-go |
| **📋 Manifest Validation** | Validate YAML manifests against best practices |
| **🏗️ Terraform Generation** | Generate production-ready Terraform HCL for Kubernetes resources |
| **🐳 Container Ready** | Multi-stage Docker build for minimal production images |
| **⚡ CI/CD Pipeline** | Full GitHub Actions pipeline with 7 automated jobs |
| **🔒 Security Scanned** | Automated vulnerability scanning with Trivy and GoSec |

---

## 🎯 Quick Start

### Prerequisites
- Go 1.21+
- kubectl (optional, for cluster validation)
- Docker (optional, for containerized usage)

### Installation

```bash
# Clone the repository
git clone https://github.com/ArjunSrivastava1/gitlab-deployment-toolkit.git
cd gitlab-deployment-toolkit

# Build from source
go build -o gdt ./cmd/gdt

# Run it!
./gdt -version
```

### Using Docker

```bash
# Build the Docker image
docker build -t gdt .

# Run in container
docker run --rm gdt -version
```

---

## 📖 Usage Examples

### 1. Check Kubernetes Deployment Health

```bash
# List all deployments in a namespace
./gdt -list-deployments default

# Check a specific deployment
./gdt -check-k8s default/nginx-example

# Output:
# 🔍 Checking deployment nginx-example in namespace default...
# ✅ Deployment is healthy and all replicas are ready
# 📊 Replicas: 3/3 ready
```

### 2. Validate Kubernetes Manifests

```bash
# Validate a deployment YAML file
./gdt -validate examples/example-config.yaml

# Output:
# 📋 Validating manifest: examples/example-config.yaml
# ✅ Validation passed!
# Errors found: 0
```

### 3. Generate Terraform Code

```bash
# Generate Terraform configuration for a deployment
./gdt -generate-tf my-awesome-app

# Outputs ready-to-use Terraform HCL:
# resource "kubernetes_deployment" "my-awesome-app" {
#   metadata { ... }
#   spec { ... }
# }
```

---

## 🏗️ Project Architecture

```
gitlab-deployment-toolkit/
├── cmd/gdt/              # CLI entrypoint (cobra-style)
├── pkg/
│   ├── k8s/              # Kubernetes client integration
│   │   ├── check.go      # Deployment health validation
│   │   └── check_test.go # Unit tests with Kind
│   ├── validate/         # YAML manifest validation
│   │   └── validate.go   # Schema and best practices
│   └── terraform/        # Infrastructure as Code generation
│       └── generate.go   # Terraform HCL generator
├── examples/             # Sample configurations
├── .github/workflows/    # CI/CD pipeline
├── Dockerfile            # Multi-stage container build
├── kind-config.yaml      # Kind cluster configuration
└── README.md            # You are here!
```

---

## 🤖 CI/CD Pipeline

This project uses **GitHub Actions** for complete automation:

| Job | Description | Status |
|-----|-------------|--------|
| **lint** | Code formatting, go vet, staticcheck | ✅ |
| **unit-tests** | Unit tests with coverage reporting | ✅ |
| **build** | Multi-platform binaries (Linux/macOS/Windows) | ✅ |
| **docker** | Container image build and test | ✅ |
| **kubernetes-tests** | Integration tests with Kind cluster | ✅ |
| **security-scan** | Trivy + GoSec vulnerability scanning | ✅ |
| **release** | Automatic GitHub release on tags | 🚀 |

[View the pipeline →](https://github.com/ArjunSrivastava1/gitlab-deployment-toolkit/actions)

---

## 🧪 Testing

### Unit Tests
```bash
go test -v ./...
```

### Integration Tests (requires Kind)
```bash
# Create a Kind cluster
kind create cluster --config kind-config.yaml

# Run integration tests
go test -v -tags=integration ./pkg/k8s/
```

---

## 🐳 Docker Support

```dockerfile
# Multi-stage build for minimal images
FROM golang:1.21 AS builder
# ... build stage

FROM alpine:latest
# ... production stage (only ~15MB)
```

**Benefits:**
- ✅ Small image size (~15MB)
- ✅ No build tools in production
- ✅ Security optimized
- ✅ Quick pull and deploy

---

## 🛣️ Roadmap

- [x] Basic CLI structure
- [x] Kubernetes deployment validation
- [x] Manifest YAML validation
- [x] Terraform code generation
- [x] GitHub Actions CI/CD
- [x] Kind integration tests
- [x] Docker containerization
- [ ] Prometheus metrics export
- [ ] Multi-cluster support
- [ ] Slack/Webhook notifications
- [ ] Helm chart generation

---

## 🤝 Contributing

This is a portfolio project, but contributions are welcome! Feel free to:
1. Fork the repository
2. Create a feature branch
3. Submit a Pull Request

---

## 📄 License

MIT License - feel free to use this for your own projects!

---

## 🙏 Acknowledgments

Built with:
- [Kubernetes client-go](https://github.com/kubernetes/client-go)
- [GitHub Actions](https://github.com/features/actions)
- [Kind](https://kind.sigs.k8s.io/)
- [Docker](https://www.docker.com/)

---

## 📞 Connect

**Arjun Srivastava**
- GitHub: [@ArjunSrivastava1](https://github.com/ArjunSrivastava1)
- Project Link: [gitlab-deployment-toolkit](https://github.com/ArjunSrivastava1/gitlab-deployment-toolkit)
---
