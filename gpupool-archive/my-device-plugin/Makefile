IMAGE_VERSION = latest
IMAGE = irenetht/my-device-plugin

.PHONY: build deploy

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/falcon cmd/server/app.go

buildImage:
	docker build --no-cache -t ${IMAGE} .

pushImage:
	docker push ${IMAGE}

deploy:
	helm install falcon deploy/helm/falcon
	
remove:
	helm uninstall falcon
