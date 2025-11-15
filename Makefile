APP_NAME=resume-app
IMAGE_REPO?=your-account-id.dkr.ecr.us-east-1.amazonaws.com/$(APP_NAME)
VERSION?=latest

build:
	go build -o bin/server ./cmd/server

docker-build:
	docker build -t $(IMAGE_REPO):$(VERSION) .

docker-run:
	docker run --rm -p 8080:8080 $(IMAGE_REPO):$(VERSION)

dev:
	docker compose up --build

dev-local:
	chmod +x scripts/dev.sh && ./scripts/dev.sh

k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/service.yaml

k8s-delete:
	kubectl delete namespace resume-app --ignore-not-found

test:
	go test ./...
