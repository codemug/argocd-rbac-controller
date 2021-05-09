# ArgoCD RBAC Controller

> Are you tired of creating custom tooling around the [pesky CSV based ConfigMap that manages RBAC in ArgoCD](https://argoproj.github.io/argo-cd/operator-manual/rbac/). Do you also yearn for a declarative way of creating roles and permissions in ArgoCD like the Kubernetes RBAC? If yes, you've found the right repo.

`argocd-rbac-controller` is a Kubernetes operator that lets you declaratively define the ArgoCD `groups`, `roles` and `permissions` using Kubernetes custom resources. 

