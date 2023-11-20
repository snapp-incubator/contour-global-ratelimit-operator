
# Contour Global Rate Limit Operator.

## Description

The Contour Global Rate Limit Operator is a Kubernetes operator designed to read and parse Contour HTTPProxy objects, extract rate limit settings from the global rate limit section of HTTPProxy, and serve these configurations via the xDS protocol to a rate limit server.

## Features

- Parses Contour HTTPProxy objects to extract rate limit settings.
- Serves rate limit configurations via the xDS protocol to a rate limit server.
- Configures rate limits based on various descriptors and conditions.

## Example

### Contour HTTPProxy Object

```yaml
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echo
  namespace: test
spec:
  ingressClassName: private
  virtualhost:
    fqdn: my-test.example.com
  routes:
  - conditions:
    - prefix: /
    services:
    - name: ingress-conformance-echo
      port: 80
    rateLimitPolicy:
      global:
        descriptors:
          - entries:
             - genericKey:
                key: test.echo.limit1 # Namespace.Name.optional_Name
                value: "3/m"
             - requestHeaderValueMatch:
                 headers:
                 - name: foo
                   exact: bar
                 value: bar  
          - entries:
             - genericKey:
                key: test.echo.limit2 # Namespace.Name.optional_Name
                value: "4/m"
             - requestHeader:
                  headerName: id
                  descriptorKey: id
          - entries:
              - genericKey:
                  key: test.echo.limit3 # Namespace.Name.optional_Name
                  value: "30/m"
```

### Generated Rate Limit Server Configuration

```yaml
name: contour
domain: contour
descriptors:
    - key: test.echo.limit1
      value: 3/m
      descriptors:
        - key: header_match
          value: bar
          ratelimit:
            unit: minute
            requestsperunit: 3
            unlimited: false
    - key: test.echo.limit2
      value: 4/m
      descriptors:
        - key: id
          value: ""
          ratelimit:
            unit: minute
            requestsperunit: 4
            unlimited: false
    - key: test.echo.limit3
      value: 30/m
      ratelimit:
        unit: minute
        requestsperunit: 30
```

## Usage

To use the Contour Global Rate Limit Operator, you'll need a Kubernetes cluster. You can set up a local cluster for testing using [KIND](https://sigs.k8s.io/kind) or use a remote cluster.

### Installation

1. Install instances of custom resources:

   ```sh
   kubectl apply -f config/samples/
   ```

2. Build and push your operator image to a container registry:

   ```sh
   make docker-build docker-push IMG=<some-registry>/contour-global-ratelimit-operator:tag
   ```

3. Deploy the operator to the cluster with the image specified by `IMG`:

   ```sh
   make deploy IMG=<some-registry>/contour-global-ratelimit-operator:tag
   ```

### Uninstallation

To delete the custom resource definitions (CRDs) from the cluster:

```sh
make uninstall
```

To undeploy the operator from the cluster:

```sh
make undeploy
```

### How It Works

This project follows the Kubernetes Operator pattern, utilizing controllers to manage resources and synchronize them to the desired state on the cluster. For detailed information, refer to the [Kubernetes Operator documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

### Testing

To test the operator:

1. Install the CRDs into the cluster:

   ```sh
   make install
   ```

2. Run the operator (it will run in the foreground):

   ```sh
   make run
   ```

You can also install and run the operator in one step using `make install run`.

### Modifying the API Definitions

If you need to edit the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

For more information about available `make` targets, run `make --help`.

## License

Copyright © 2023.

Licensed under the Apache License, Version 2.0. For more details, see the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).
