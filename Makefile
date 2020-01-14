.PHONY: build test testacc vet lint release

plugin_name=terraform-provider-appoptics

default: build

build:
	go build -o $(plugin_name)

test:
	go test ./...

testacc:
	TF_ACC=1 go test -v -timeout 120m

vet:
	go vet ./...

lint:
	"$$(go env GOPATH)/bin/golangci-lint" run

