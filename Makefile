build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/godoc-operator ./cmd/godoc-operator

clean:
	rm -rf ./bin

docker-build: build
	docker build \
		-t rpflynn22/godoc-operator:latest \
		--build-arg BIN=./bin/godoc-operator \
		-f docker/godoc-operator/Dockerfile .
	docker build \
		-t rpflynn22/godoc-server:latest \
		-f docker/godoc-server/Dockerfile .

deploy: docker-build deploy-ns create-crd create-secret
	kubectl -n godoc apply -f k8s/godoc-operator.yaml

undeploy: delete-crd delete-secret
	kubectl -n godoc delete -f k8s/godoc-operator.yaml

deploy-ns:
	kubectl apply -f k8s/namespace.yaml

create-crd:
	kubectl apply -f k8s/crd.yaml

delete-crd:
	kubectl delete -f k8s/crd.yaml

generate: install-controller-gen
	controller-gen object paths=./internal/api/...

install-controller-gen:
	which controller-gen || go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

create-secret:
	kubectl -n godoc create secret generic github --from-literal=pat=$(PERSONAL_GITHUB_TOKEN)

delete-secret:
	kubectl -n godoc delete secret github
