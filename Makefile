IMAGE_NAME ?= k8s-pod-cpu-stressor
IMAGE_TAG ?= latest
PLATFORMS ?= linux/amd64,linux/arm64,linux/s390x,linux/ppc64le

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o cpu-stress .

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: docker-buildx
docker-buildx:
	docker buildx build --platform $(PLATFORMS) -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: test
test:
	go test -v ./...
