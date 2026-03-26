FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Fix: Build from correct path
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gdt ./cmd/gitlab-deployment-toolkit

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gdt .

RUN chmod +x ./gdt

# Fix: Entrypoint should point to the binary, not the source path
ENTRYPOINT ["./gdt"]