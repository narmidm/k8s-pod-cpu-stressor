.PHONY: docker-build
docker-build:
	docker build -t k8s-pod-cpu-stressor:latest .

.PHONY: build
build:
	go build -o cpu-stress .


