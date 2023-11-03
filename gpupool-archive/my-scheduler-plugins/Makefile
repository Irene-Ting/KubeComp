.PHONY: build deploy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=bin/my-scheduler ./cmd/scheduler

buildImage:
	docker build --no-cache . -t irenetht/my-scheduler

pushImage: 
	docker push irenetht/my-scheduler

deploy:
	helm install scheduler-plugins charts/ --create-namespace --namespace scheduler-plugins

remove:
	helm uninstall scheduler-plugins -n scheduler-plugins

clean:
	rm -rf bin/
