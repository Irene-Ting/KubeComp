IMAGE_VERSION = latest
IMAGE = irenetht/reconfig-daemon

.PHONY: build

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/falcon main.go

buildImage:
	docker build -t ${IMAGE} . --no-cache

pushImage:
	docker push ${IMAGE}

deploy:
	helm install reconfig-daemon chart/
	
remove:
	helm uninstall reconfig-daemon
