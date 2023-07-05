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
## Running with Docker

Build the Docker image using the provided Dockerfile:

   ```shell
   docker build -t k8s-pod-cpu-stressor .
  ```
Run the Docker container, specifying the desired CPU usage and stress duration as arguments:
```shell
docker run --rm k8s-pod-cpu-stressor -cpu=0.2 -duration=10s
```
Replace 0.2 and 10s with the desired CPU usage (fraction) and duration, respectively.

## Contributing
Contributions are welcome! If you find a bug or have a suggestion, please open an issue or submit a pull request. For major changes, please discuss them first in the issue tracker.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.





