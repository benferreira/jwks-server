.PHONY: build
build: clean
	go build -o bin/ ./...

.PHONY: build-image
build-image: build
	ko publish --base-import-paths --local ./cmd/jwks-server/

.PHONY: build-image-arm
build-image-arm: build
	GOOS=linux GOARCH=arm64 ko publish --base-import-paths --local ./cmd/jwks-server/

.PHONY: clean
clean:
	rm -rf bin/*
	go mod tidy

.PHONY: coverage
coverage: test
	go tool cover -func=coverage.out

.PHONY: coverage-report
coverage-report: test
	go tool cover -html=coverage.out

.PHONY: test
test: build
	go test ./... -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.out
