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

# Produces artifacts in the dist directory
# DOES NOT push release artifacts
test-release:
	goreleaser --snapshot --skip-publish --rm-dist

# Requires a GITHUB_TOKEN to be set in the environment
release:
	goreleaser --rm-dist

# Convenient in dev to rebuild the plugin, re-init TF, and run a plan
bounce: build
	rm *tfstate* && terraform init && terraform plan


