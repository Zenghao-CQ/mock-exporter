USER=zenghaocq
TAG=0.0.4
IMG=${USER}/failure-exporter:${TAG}

all: build docker-build docker-push

build:
	go build mock-exporter.go

clean:
	rm -r mock-exporter

docker-build:  ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}