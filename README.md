# tpl2git-operator

## generate project

must be in an empty directory (will not even let you have a .git directory)

initialize using [operator-framework/operator-sdk](https://github.com/operator-framework/operator-sdk)

```
operator-sdk init --domain "stephanhcb.github.io" --repo "github.com/StephanHCB/tpl2git-operator" \
--license none --owner "StephanHCB" --plugins "go.kubebuilder.io/v3" --project-name "tpl2git-operator"
```

create a resource and a controller scaffold

```
operator-sdk create api --group tpl2git --version v1alpha1 --kind Renderer --resource --controller
```




