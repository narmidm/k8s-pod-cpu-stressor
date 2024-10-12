# Security Policy

## Supported Versions

We maintain security updates and support for the following versions of `k8s-pod-cpu-stressor`:

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

Please ensure you are running a supported version to benefit from security patches.

## Reporting a Vulnerability

If you discover a vulnerability in this project, please follow these steps to report it securely:

1. **Do not open a public issue** on GitHub, as this may expose the vulnerability to others before it can be addressed.
2. Contact us by sending an email to [imranaec@outlook.com](mailto:imranaec@outlook.com) with the details of the vulnerability, including steps to reproduce it, affected versions, and potential impact.
3. Please allow us **at least 90 days** to investigate and apply a fix before disclosing the issue publicly.

We will work to acknowledge your report within **7 days** and provide an estimated timeline for a fix.

## Security Best Practices for Users

To help ensure the security of your Kubernetes environment, consider the following when using `k8s-pod-cpu-stressor`:

- **Namespace Isolation**: Run the tool in a dedicated namespace to limit any potential impact.
- **Permissions**: Grant minimal permissions needed for the pod to run. Avoid giving it elevated privileges unless explicitly necessary.
- **Network Policies**: Apply appropriate network policies to restrict access to and from the pods running this tool.

## Responsible Disclosure Policy

We believe in and support responsible disclosure. If you report a vulnerability and work with us constructively, we are committed to acknowledging your contributions in the release notes or other appropriate acknowledgments (with your permission).

## Contact

If you have general security concerns or questions, please contact the maintainers at [imranaec@outlook.com](mailto:imranaec@outlook.com).
