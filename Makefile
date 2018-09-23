TAG?=latest

build:
	docker build -t openfaas/openfaas-operator:$(TAG) . -f Dockerfile

build-armhf:
	docker build -t openfaas/openfaas-operator:$(TAG)-armhf . -f Dockerfile.armhf

push:
	docker push openfaas/openfaas-operator:$(TAG)

test:
	go test ./...

verify-codegen:
	./hack/verify-codegen.sh

ci-armhf-build:
	docker build -t openfaas/openfaas-operator:$(TAG)-armhf . -f Dockerfile.armhf

ci-armhf-push:
	docker push openfaas/openfaas-operator:$(TAG)-armhf
