run:
	GORACE="halt_on_error=1" go run -race cmd/recipes/main.go
build:
	go build -ldflags "-s -w" cmd/recipes/main.go
test:
	go test ./... -v -race -cover -count=1
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.6 golangci-lint run
lint-insecure:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.6 git config --global http.sslVerify false && golangci-lint run