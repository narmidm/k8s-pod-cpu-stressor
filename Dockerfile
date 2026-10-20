# Build stage using official multi-arch Golang base
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /app

# Copy module manifests first for efficient Docker layer caching
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

# Statically compile pure Go binary for target architecture (s390x, amd64, arm64, ppc64le)
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o /app/cpu-stressor .

# Minimal runtime image based on UBI minimal (multi-arch supported by Red Hat)
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

WORKDIR /

# Copy binary to root AND to standard PATH
COPY --from=builder /app/cpu-stressor /cpu-stressor
COPY --from=builder /app/cpu-stressor /usr/local/bin/cpu-stressor
COPY --from=builder /app/cpu-stressor /usr/local/bin/cpu-stress

# Ensure binary is executable by any arbitrary UID/GID allocated by OpenShift SCC
RUN chmod 755 /cpu-stressor /usr/local/bin/cpu-stressor /usr/local/bin/cpu-stress

# Default non-root user for vanilla Kubernetes / Docker
USER 1001

ENTRYPOINT ["/cpu-stressor"]