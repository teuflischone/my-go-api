.PHONY: tidy run test fmt
tidy:
	go mod tidy
fmt:
	go fmt ./...
test:
	go test ./... -race
run:
	go run ./cmd/api
