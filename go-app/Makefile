prep:
	swag init -g main.go --output docs/
	go mod tidy

run: prep
	go run main.go

gin: prep
	gin --appPort 8080 --port 5000 --immediate

lint:
	golangci-lint run --enable-all

test:
	go test ./modules/k8sclient
	go test

localbuild: prep
	go build
	
build: prep
	docker build --tag go-app:v1.0 .

deploy:
	kubectl apply -f manifests/
	kubectl rollout restart deploy go-app -n go-app

builddeploy: build deploy

connect:
	kubectl port-forward -n go-app svc/go-app 8080:80 

all: prep lint test build deploy

load:
	ab -n 100 -c 10 http://localhost:8080/namespaces

logs:
	kubectl logs -f -n go-app -l app=go-app --tail -1

watch:
	watch kubectl get pods -n go-app

.PHONY: prep run gin lint test localbuild build deploy builddeploy connect all load logs watch
