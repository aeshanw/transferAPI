build:
	go build cmd/.

run:
	go run ./cmd/main.go

update-vendor:
	go mod tidy
	go mod vendor

test:
	go test -cover ./...