fmt:
	go fmt ./...

test:
	go test ./...

dev:
	PORT=3000 go run ./src/main.go
