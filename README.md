# Godoc Operator

A Kubernetes controller/operator that deploys godoc servers for private
Github repos.

## Use it

I've so far only tested on minikube. In order for minikube to get locally-built
docker images, run `eval $(minikube docker-env)`. The deployment files are
already setup to not pull remote images.

Run
```sh
$ kubectl create ns godoc
$ make deploy
```
to deploy everything. It
- builds the operator docker image
- builds the godoc server docker image
- deploys the operator along with rbac permissions it needs
- deploys the CRD

To see it work, run

```sh
$ kubectl -n godoc apply -f k8s/sample-repo.yaml
```

Note that this won't work for you as the sample references _my_ private repo.
You'll need to change the repo to a URL you have access to (public or private)
and create a Secret containing the PAT that can be used to access the repo.

## Current State

It only creates missing things right now. It does not deal with spec drift or
deletion.
