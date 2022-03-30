.PHONY: build test testacc vet lint release

plugin_name=terraform-provider-appoptics

default: build

build:
	go build -o $(plugin_name)

buildall:
	env CGO_ENABLED= GOOS="linux" GOARCH="amd64" go build -trimpath -buildmode=pie -ldflags "-s -w" -o  build/$(plugin_name)-linux-amd64
	env CGO_ENABLED= GOOS="linux" GOARCH="arm64" go build -trimpath -buildmode=pie -ldflags "-s -w" -o  build/$(plugin_name)-linux-arm64
	env CGO_ENABLED= GOOS="darwin" GOARCH="amd64" go build -trimpath -buildmode=pie -ldflags "-s -w" -o  build/$(plugin_name)-darwin-amd64
	env CGO_ENABLED= GOOS="darwin" GOARCH="arm64" go build -trimpath -buildmode=pie -ldflags "-s -w" -o  build/$(plugin_name)-darwin-arm64
	env CGO_ENABLED= GOOS="windows" GOARCH="amd64" go build -trimpath -buildmode=pie -ldflags "-s -w" -o  build/$(plugin_name)-windows-amd64
	cd build && sha256sum -b * > checksums_sha256.txt

test:
	go test ./...

testacc:
	TF_ACC=1 go test -v -timeout 120m

vet:
	go vet ./...

lint:
	"$$(go env GOPATH)/bin/golangci-lint" run

# Convenient in dev to rebuild the plugin, re-init TF, and run a plan
bounce: build
	rm -f *.tfstate* && terraform init && terraform plan


