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
	oapi-codegen -generate "types" -package gen twitterCloneApi.yaml > ./internal/http/gen/model.go
	oapi-codegen -generate "server,spec" -package gen twitterCloneApi.yaml > ./internal/http/gen/server.go