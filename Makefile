.PHONY: docker-build
docker-build:
	docker build -t k8s-pod-cpu-stressor:1.0.0 .

.PHONY: build
build:
	go build -o cpu-stress .


