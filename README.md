# ArgoCD RBAC Controller

> Are you tired of creating custom tooling around the [pesky CSV based ConfigMap that manages RBAC in ArgoCD](https://argoproj.github.io/argo-cd/operator-manual/rbac/). Do you also yearn for a declarative way of creating roles and permissions in ArgoCD like the Kubernetes RBAC? If yes, you've found the right repo.

`argocd-rbac-controller` is a Kubernetes operator that lets you declaratively define the ArgoCD `groups`, `roles` and `permissions` using Kubernetes custom resources. 

## Installation

Use the [helm chart](https://github.com/codemug/argocd-rbac-controller/tree/main/helm/argocd-rbac-controller) in this repo to deploy the controller on your kubernetes cluster:

```shell
helm install argocd-rbac-controller helm/argocd-rbac-controller
```

**Note**: If you intend to build the image yourself and push it to your own registry, you can update the `REGISTRY` value in the `Makefile` and do a:

```shell
make docker-build
make docker-push
```

And then, when deploying the helm chart, you can set the image name as follows:

```shell
helm install argocd-rbac-controller helm/argocd-rbac-controller --set image.registry=name-of-your-registry
```

## Usage

The operator installs two namespace-scoped `CustomResourceDefinitions` on your cluster:

#### [GroupMapping](https://github.com/codemug/argocd-rbac-controller/blob/main/helm/argocd-rbac-controller/crds/argocd.codemug.io_groupmappings.yaml)

This translates to the `g` statements in the `argocd-rbac-cm` `ConfigMap`. For example, consider the following entry:

```csv
g, bar, role:foo
```

This would be created through:

```yaml
apiVersion: argocd.codemug.io/v1beta1
kind: GroupMapping
metadata:
  name: groupmapping-sample
spec:
  mappings:
    - roleName: foo
      groupName: bar
```

#### [RoleMapping](https://github.com/codemug/argocd-rbac-controller/blob/main/helm/argocd-rbac-controller/crds/argocd.codemug.io_rolemappings.yaml)

This translates to the `p` statements in the `argocd-rbac-cm` `ConfigMap`. For example, consider the following entry:

```csv
p, role:foo, applications, get, *, allow
```

This would be created through:

```yaml
apiVersion: argocd.codemug.io/v1beta1
kind: RoleMapping
metadata:
  name: rolemapping-sample
spec:
  roles:
    - name: foo
      permissions:
        - resource: applications
          actions:
          - get
          instance: "*"
```

## Configuration

The name and namespace of the `argocd-rbac-cm` can be changed/configured at the time of the helm chart deployment:

```shell
helm install argocd-rbac-controller helm/argocd-rbac-controller --set controller.rbacConfigMapName rbac-cm --set controller.rbacConfigMapNamespace cd-system
```

Similarly, the value for `policy.default` in this `ConfigMap` can also be configured:

```shell
helm install argocd-rbac-controller helm/argocd-rbac-controller --set controller.defaultPolicy role:admin
```
