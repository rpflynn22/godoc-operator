DOCKER_USE_MK ?= eval $$(minikube docker-env 2>&1) &&

build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/godoc-operator ./cmd/godoc-operator

clean:
	rm -rf ./bin

docker-build: build
	$(DOCKER_USE_MK) docker build \
		-t rpflynn22/godoc-operator:0.0.4 \
		--build-arg BIN=./bin/godoc-operator \
		-f docker/godoc-operator/Dockerfile . && \
	$(DOCKER_USE_MK) docker build \
		-t rpflynn22/godoc-server:0.0.4 \
		-f docker/godoc-server/Dockerfile .

delete-crd:
	kubectl delete -f helm/godoc-operator/crds/repo.yaml

generate: install-controller-gen
	controller-gen object paths=./internal/api/...

install-controller-gen:
	which controller-gen || go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

create-secret:
	kubectl -n godoc create secret generic github --from-literal=pat=$(PERSONAL_GITHUB_TOKEN)

delete-secret:
	kubectl -n godoc delete secret github

helm-template:
	helm template -n godoc -f helm/minikube.yaml godoc ./helm/godoc-operator

helm-install: docker-build create-secret
	helm install -n godoc -f helm/minikube.yaml godoc ./helm/godoc-operator

helm-uninstall: delete-secret
	helm uninstall -n godoc godoc
