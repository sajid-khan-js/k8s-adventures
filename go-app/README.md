# Basic Go API

## Description

Creating using `gin` <https://github.com/gin-gonic/gin>

Let's you list Namespaces and Pods, and create Namespaces in a Kubernetes cluster.

:memo: Check out of the `Makefile`

## Running locally

:memo: Will use your current `kubeconfig` (`~/.kube/config`) context

```sh
# Run server
go run main.go

# Interact with API with curl or Postman
curl -v -L localhost:8080/namespaces/
```

OR

`gin --appPort 8080 --port 5000 --immediate` (coincidently, this hot-reloading Go tool is also called gin: <https://github.com/codegangsta/gin>)

then do something like: `curl -v -L localhost:5000/namespaces/`

## API Spec

To update API specs docs run: `swag init -g main.go --output docs`

Access Swagger UI here: `http://localhost:5000/swagger/index.html`
