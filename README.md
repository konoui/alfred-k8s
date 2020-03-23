## alfred k8s (Kubernetes)
Alfred Workflow to operate k8s resources.

## Install
Download the workflow form [latest release](https://github.com/konoui/alfred-k8s/releases).

## Configuration
The workflow will try to use `/usr/local/bin/kubectl` by default.  
If your kube config has [client-go credential plugins](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins) command as relative path, the workflow will search `/usr/local/bin/` path.

If you change above default values, please create configuration file as `.alfred-k8s` in home directory (~/).  
e.g.) `aws` command for EKS exists in `${HOME}/.pyenv/shims/`.
```yaml
kubectl:
    kubectl_absolute_path: "/usr/local/bin/kubectl"
    plugin_paths:
    - "/usr/local/bin/"
    - "${HOME}/.pyenv/shims/"
```

## Feature
#### List k8s resources and copy to clipboard
- node, pod, deployment, service, ingress, namespace, context, etc...

#### Switch Context/Namespace
- `kube context` or `kube ns`

|  Key Combination  |  Action  |
| ---- | ---- |
| `enter` | Copy to clipboard |
|  `control` + `enter`  |  Switch context/namespace  |
