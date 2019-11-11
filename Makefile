run:
	GORACE="halt_on_error=1" go run -race cmd/recipes/main.go
build:
	go build -ldflags "-s -w" cmd/recipes/main.go
test:
	go test ./... -v -race -cover -count=1