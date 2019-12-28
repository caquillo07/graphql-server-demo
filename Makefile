# Makefile
GORUN_CMD=go run -mod=vendor
GQLGEN_DIR=vendor/github.com/99designs/gqlgen/gqlgen

generate: tools
	${GQLGEN_DIR} -v

run-dev:
	${GORUN_CMD} main.go gql --dev-log --config example-config.yaml

tools:
	cd vendor/github.com/99designs/gqlgen && go build

.PHONY: tools migrate-dev generate run-dev
