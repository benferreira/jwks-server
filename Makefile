
build: clean
	go build -o bin/ ./...

.PHONY: clean
clean:
	rm -rf bin/*
	go mod tidy

coverage: test
	go tool cover -func=coverage.out

coverage-report: test
	go tool cover -html=coverage.out

test: build
	go test ./... -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out
