TAG?=latest

build:
	docker build -t functions/openfaas-operator:$(TAG) . -f Dockerfile

build-armhf:
	docker build -t functions/openfaas-operator:$(TAG)-armhf . -f Dockerfile.armhf

push:
	docker push functions/openfaas-operator:$(TAG)

test:
	go test ./...
