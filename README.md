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
          image: narmidm/k8s-pod-cpu-stressor:1.0.0
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

## Contributing

Contributions are welcome! If you find a bug or have a suggestion, please open an issue or submit a pull request. For major changes, please discuss them first in the issue tracker.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
