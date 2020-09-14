## Alfred Kubernetes
Alfred workflow to operate Kubernetes resources.


## Features
- List Kubernetes resources and copy them to clipboard (e.g. pod. deployment, ingress etc..).
- Switch Context/Namespace

## Install
- Download the workflow form [latest release](https://github.com/konoui/alfred-k8s/releases).
- Build Workflow on your computer.
```
$ make package
$ ls
alfred-k8s.alfredworkflow (snip)
```

## Usage
Kyeword is `kube`.

<img src="./usage.png" width="50%">

### List resources in current namespace
Please type `kube <resource-name>`.  
e.g.) `kube pod`

### List resources in all namespaces
Please add `-A` option.  
e.g.) `kube pod -A`

### List specific resources
`kube obj <resource-name>` is for other resources not supported.  
For example, you can list replicaset resources by `kube obj rs` as the workflow does not support `kube rs`.


### Switch Context/Namespace
Please type `kube context` or `kube ns`.  
The following is default key mapping.

|  Key Combination  |  Action  |
| ---- | ---- |
| `enter` | Copy to clipboard |
|  `control` + `enter`  |  Switch context/namespace  |
|  `shift` + `enter`  |  Delete context  |


## Configuration
### kubectl command path
The workflow will try to use `/usr/local/bin/kubectl` by default.  
When your kube config has [client-go credential plugins](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins) command as relative path, the workflow will search `/usr/local/bin/` path.

If you change above default values, please create configuration file as `.alfred-k8s` in home directory (~/).  
e.g.) `aws` command for EKS exists in `${HOME}/.pyenv/shims/`.

```yaml
kubectl:
    kubectl_absolute_path: "/usr/local/bin/kubectl"
    plugin_paths:
    - "/usr/local/bin/"
    - "${HOME}/.pyenv/shims/"
```

### Workflow Key Mapping
The workflow key mapping is changed with config file as bellow.

```yaml
kubectl:
    kubectl_absolute_path: "/usr/local/bin/kubectl"

key_maps:
    context_key_map:
        enter: "copy"
        ctrl: "use"
        shift: "delete"
        cmd: ""
        alt: ""
    namespace_key_map:
        enter: "copy"
        ctrl: "use"
    pod_key_map:
        enter: "copy"
        ctrl: "delete"
        shift: "stern_copy"
    deployment_key_map:
        enter: "copy"
        shift: "stern_copy"
    service_key_map:
        enter: "copy"
        shift: "stern_copy"
        alt: "port_forward_copy"
```


## License
MIT License.
