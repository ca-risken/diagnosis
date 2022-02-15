# RISKEN Diagnosis

![Build Status](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiMW9Vd0RVRGNncWsrM09jTm4wdW5NVHRJYXl2TWJUMzFPVEh6UkxXaFJsa2hacGV6ZEY0T1l2bXA1akw0MmkwVi8yaXFjeTV0YXM0czFpVHdnWU9zYVQwPSIsIml2UGFyYW1ldGVyU3BlYyI6ImJ3U2xVNE85SFlGem9zOWwiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN Diagnosis` is a security monitoring system for fuzzing and network scanning.

Please check [RISKEN Documentation](https://docs.security-hub.jp/).

## Installation

### Requirements

This module requires the following modules:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Protocol Buffer](https://grpc.io/docs/protoc-installation/)

### Install packages

This module is developed in the `Go language`, please run the following command after installing the `Go`.

```bash
$ make install
```

### Building

Build the containers on your machine with the following command

```bash
$ make build
```

### Running Apps

Deploy the pre-built containers to the Kubernetes environment on your local machine.

- Follow the [documentation](https://docs.security-hub.jp/admin/infra_local/#risken) to download the Kubernetes manifest sample.
- Fix the Kubernetes object specs of the manifest file as follows and deploy it.

`k8s-sample/overlays/local/diagnosis.yaml`

| service         | spec                                | before (public images)                                   | after (pre-build images on your machine) |
| --------------- | ----------------------------------- | -------------------------------------------------------- | ---------------------------------------- |
| diagnosis       | spec.template.spec.containers.image | `public.ecr.aws/risken/diagnosis/diagnosis:latest`       | `diagnosis/diagnosis:latest`             |
| portscan        | spec.template.spec.containers.image | `public.ecr.aws/risken/diagnosis/portscan:latest`        | `diagnosis/portscan:latest`              |
| applicationscan | spec.template.spec.containers.image | `public.ecr.aws/risken/diagnosis/applicationscan:latest` | `diagnosis/applicationscan:latest`       |

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
