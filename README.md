# k8s-pod-cpu-stressor

The `k8s-pod-cpu-stressor` is a tool designed to simulate CPU stress on Kubernetes pods. It allows you to specify the desired CPU usage and stress duration, helping you test the behavior of your Kubernetes cluster under different CPU load scenarios.

## Features

- Simulates CPU stress on Kubernetes pods.
- Configurable CPU usage (in millicores) and stress duration.
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
