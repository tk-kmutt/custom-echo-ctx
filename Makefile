code-check: mod imports fmt vet
vet:
	go vet ./...
fmt:
	gofmt -d -s .
imports:
	goimports -w .
mod:
	go mod tidy
	go mod verify
	go mod download

install-oapi-codegen:
	go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen

oapi-codegen:
	mkdir -p internal/http/gen
	oapi-codegen -generate "types" -package gen openapi.yaml > ./internal/http/gen/model.go
	oapi-codegen -generate "server,spec" -package gen openapi.yaml > ./internal/http/gen/server.go

install-wire:
	go get github.com/google/wire/cmd/wire
wire-cec:
	go run github.com/google/wire/cmd/wire ./cmd/cec/...
