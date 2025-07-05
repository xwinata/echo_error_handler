.PHONY: test tidy

test:
	go test ./... -coverprofile=coverage.out -coverpkg=./...
	go tool cover -html=coverage.out

tidy:
	go mod tidy