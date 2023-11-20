# contour-global-ratelimit-operator

## Usage

[Helm](https://helm.sh) must be installed to use the charts. Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

```shell
helm repo add rls-operator https://snapp-incubator.github.io/contour-global-ratelimit-operator
```

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages. You can then run `helm search repo
rls-operator` to see the charts.

To install the rls-operator chart:

```shell
helm install my-release rls-operator/rls-operator
```

To uninstall the chart:

```shell
helm delete my-release
```
