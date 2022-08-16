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
	kubectl rollout restart deploy go-app
	kubectl port-forward svc/go-app 8080:80

builddeploy: build deploy

all: prep lint test build deploy

