# tpl2git-operator

## generate project

must be in an empty directory (will not even let you have a .git directory)

initialize using [operator-framework/operator-sdk](https://github.com/operator-framework/operator-sdk)

```
operator-sdk init --domain "stephanhcb.github.io" --repo "github.com/StephanHCB/tpl2git-operator" \
--license none --owner "StephanHCB" --plugins "go.kubebuilder.io/v3" --project-name "tpl2git-operator"
```

create a custom resource and a controller scaffold

```
operator-sdk create api --group tpl2git --version v1alpha1 --kind Renderer --resource --controller
```

## implementation

edit `api/v1alpha1/renderer_types.go` to define the fields of your custom resource. 

run `make` to adapt the code to reflect your changes.

run `make manifests` to create the custom resource definition manifests.

implement the reconcile loop, then build the binary

## execution for debugging

- build the binary
- transfer it to a system with an admin kubeconfig
- transfer config/crd directory
- `kubectl apply -k config/crd` (uses built-in kustomize)
- run the binary with access to an admin kubeconfig

For deployment into the cluster instead
- build the docker image
- push it to a registry the cluster can reach
- `kubectl apply -k config/default` (uses built-in kustomize) 

## example resource

See `example.yaml`.

`kubectl apply -f example.yaml`
