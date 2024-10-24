[![CI Status](https://github.com/narmidm/k8s-pod-cpu-stressor/actions/workflows/trivy-image-scan.yml/badge.svg)](https://github.com/narmidm/k8s-pod-cpu-stressor/actions/workflows/trivy-image-scan.yml)
[![CD Status](https://github.com/narmidm/k8s-pod-cpu-stressor/actions/workflows/docker-publish-image.yml/badge.svg)](https://github.com/narmidm/k8s-pod-cpu-stressor/actions/workflows/docker-publish-image.yml)
[![Docker Image Version](https://img.shields.io/docker/v/narmidm/k8s-pod-cpu-stressor?sort=semver)](https://hub.docker.com/repository/docker/narmidm/k8s-pod-cpu-stressor)
[![Docker Pulls](https://img.shields.io/docker/pulls/narmidm/k8s-pod-cpu-stressor)](https://hub.docker.com/repository/docker/narmidm/k8s-pod-cpu-stressor)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/narmidm/k8s-pod-cpu-stressor)](https://raw.githubusercontent.com/narmidm/k8s-pod-cpu-stressor/refs/heads/master/go.mod)
[![GitHub License](https://img.shields.io/github/license/narmidm/k8s-pod-cpu-stressor)](https://raw.githubusercontent.com/narmidm/k8s-pod-cpu-stressor/refs/heads/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/narmidm/k8s-pod-cpu-stressor)](https://goreportcard.com/report/github.com/narmidm/k8s-pod-cpu-stressor)
![Contributors](https://img.shields.io/github/contributors/narmidm/k8s-pod-cpu-stressor)
[![GitHub Issues](https://img.shields.io/github/issues/narmidm/k8s-pod-cpu-stressor)](https://github.com/narmidm/k8s-pod-cpu-stressor/issues)
[![GitHub Stars](https://img.shields.io/github/stars/narmidm/k8s-pod-cpu-stressor)](https://github.com/narmidm/k8s-pod-cpu-stressor/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/narmidm/k8s-pod-cpu-stressor)](https://github.com/narmidm/k8s-pod-cpu-stressor/forks)
[![Last Commit](https://img.shields.io/github/last-commit/narmidm/k8s-pod-cpu-stressor)](https://github.com/narmidm/k8s-pod-cpu-stressor/commits/master/)

### Connect with me
[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/that_imran)](https://x.com/intent/user?screen_name=that_imran)
[![Linktree](https://img.shields.io/badge/linktree-1de9b6?style=plastic&logo=linktree&logoColor=white)](https://linktr.ee/narmidm)
<a href="https://www.linkedin.com/comm/mynetwork/discovery-see-all?usecase=PEOPLE_FOLLOWS&followMember=narmidm" target="blank"><img src="https://img.shields.io/badge/LinkedIn-Connect-blue" alt="narmidm" /></a>


# k8s-pod-cpu-stressor

The `k8s-pod-cpu-stressor` is a tool designed to simulate CPU stress on Kubernetes pods. It allows you to specify the desired CPU usage and stress duration, helping you test the behavior of your Kubernetes cluster under different CPU load scenarios.

## Features

- Simulates CPU stress on Kubernetes pods.
- Configurable CPU usage (in millicores) and stress duration.
- Option to run CPU stress indefinitely.
- Helps evaluate Kubernetes cluster performance and resource allocation.

## Getting Started

### Prerequisites

To use the `k8s-pod-cpu-stressor`, you need to have the following installed:

- Go (version 1.19 or higher)
- Docker

### Building the Binary

1. Clone this repository to your local machine.
2. Navigate to the repository directory.
3. Build the binary using the following command:

   ```shell
   go build -o cpu-stress .
   ```

## Running with Docker

Build the Docker image using the provided Dockerfile:

   ```shell
   docker build -t k8s-pod-cpu-stressor .
   ```

Run the Docker container, specifying the desired CPU usage, stress duration, and optionally whether to run CPU stress indefinitely:

```shell
docker run --rm k8s-pod-cpu-stressor -cpu=0.2 -duration=10s -forever
```

Replace `0.2` and `10s` with the desired CPU usage (fraction) and duration, respectively. Add `-forever` flag to run CPU stress indefinitely.

## CPU Usage and Duration

The `k8s-pod-cpu-stressor` allows you to specify the desired CPU usage and stress duration using the following parameters:

- **CPU Usage**: The CPU usage is defined as a fraction of CPU resources. It is specified using the `-cpu` argument. For example, `-cpu=0.2` represents a CPU usage of 20% or 200 milliCPU (mCPU).

- **Stress Duration**: The stress duration defines how long the CPU stress operation should run. It is specified using the `-duration` argument, which accepts a duration value with a unit. Supported units include seconds (s), minutes (m), hours (h), and days (d). For example, `-duration=10s` represents a stress duration of 10 seconds, `-duration=5m` represents 5 minutes, `-duration=2h` represents 2 hours, and `-duration=1d` represents 1 day.

- **Run Indefinitely**: To run CPU stress indefinitely, include the `-forever` flag.

Adjust these parameters according to your requirements to simulate different CPU load scenarios.

### Kubernetes Resource Requests and Limits

It is recommended to specify Kubernetes resource requests and limits to control the amount of CPU resources consumed by the pod, and to prevent overloading your cluster. For example:

- **Requests**: This defines the minimum amount of CPU that the pod is guaranteed to have.
- **Limits**: This defines the maximum amount of CPU that the pod can use.

Adding requests and limits helps Kubernetes manage resources efficiently and ensures that your cluster remains stable during stress testing.

Example:

```yaml
resources:
  requests:
    cpu: "100m"
  limits:
    cpu: "200m"
```

## Check the Public Docker Image

The [`k8s-pod-cpu-stressor`](https://hub.docker.com/r/narmidm/k8s-pod-cpu-stressor "Docker Hub - narmidm/k8s-pod-cpu-stressor") Docker image is publicly available on Docker Hub. You can check and pull the image using the following command:

```shell
docker pull narmidm/k8s-pod-cpu-stressor:latest
```

## Sample Deployment Manifest

Use the following deployment manifest as a starting point to deploy the k8s-pod-cpu-stressor image in your Kubernetes cluster:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cpu-stressor-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cpu-stressor
  template:
    metadata:
      labels:
        app: cpu-stressor
    spec:
      containers:
        - name: cpu-stressor
          image: narmidm/k8s-pod-cpu-stressor:latest
          args:
            - "-cpu=0.2"
            - "-duration=10s"
            - "-forever"
          resources:
            limits:
              cpu: "200m"
            requests:
              cpu: "100m"
```

## Sample Job Manifest

If you want to run the CPU stressor for a fixed duration as a one-time job, you can use the following Kubernetes Job manifest:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: cpu-stressor-job
spec:
  template:
    metadata:
      labels:
        app: cpu-stressor
    spec:
      containers:
        - name: cpu-stressor
          image: narmidm/k8s-pod-cpu-stressor:latest
          args:
            - "-cpu=0.5"
            - "-duration=5m"
          resources:
            limits:
              cpu: "500m"
            requests:
              cpu: "250m"
      restartPolicy: Never
  backoffLimit: 3
```

This manifest runs the `k8s-pod-cpu-stressor` as a Kubernetes Job, which will execute the stress test once for 5 minutes and then stop. The `backoffLimit` specifies the number of retries if the job fails.

## Contributing

Contributions are welcome! If you find a bug or have a suggestion, please open an issue or submit a pull request. For major changes, please discuss them first in the issue tracker.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
