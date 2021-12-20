# Godoc Operator

This repository contains a Kubernetes operator that deploys
[Godoc servers](https://pkg.go.dev/golang.org/x/tools/cmd/godoc) to a k8s
cluster and sets up Services and Ingresses to access them. pkg.go.dev already
does this (with a better UI), but it cannot generate documentation for private
repositories on Github. By using this operator, you can write a few lines of
yaml and deploy a godoc server that _can_ access your private repos and serve
documentation for them.

## Use it

I've so far only tested on minikube. In order for minikube to get locally-built
docker images, run `eval $(minikube docker-env)`. The deployment files are
already setup to not pull remote images.

Run
```sh
$ make deploy
```
to deploy everything. It
- builds the operator docker image
- builds the godoc server docker image
- creates a `godoc` namespace where everything is placed
- deploys the CRD
- deploys the operator along with rbac permissions it needs
- creates a Secret containing your Github PAT; uses the environment variable
`$PERSONAL_GITHUB_TOKEN`

To see it work, run

```sh
$ kubectl -n godoc apply -f k8s/sample-repo.yaml
```

Note that this won't work for you as the sample references _my_ private repo.
You'll need to change the repo to a URL you have access to (public or private).

