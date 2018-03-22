TAG?=latest

build:
	docker build -t functions/faas-o6s:$(TAG) . -f Dockerfile

build-armhf:
	docker build -t functions/faas-o6s:$(TAG)-armhf . -f Dockerfile.armhf

push:
	docker push functions/faas-o6s:$(TAG)
