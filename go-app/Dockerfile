# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-buster AS build

WORKDIR /

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . . 

RUN go get -u github.com/gin-gonic/gin && go build

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /go-app .

EXPOSE 8080

# i.e. USER 65535:65535
USER nonroot:nonroot

ENTRYPOINT ["./go-app"]
