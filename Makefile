.PHONY: build build-image build-image-arm clean coverage coverage-report lint test

build: clean
	go build -o bin/ ./...


build-image: build
	ko publish --base-import-paths --local ./cmd/jwks-server/


build-image-arm: build
	GOOS=linux GOARCH=arm64 ko publish --base-import-paths --local ./cmd/jwks-server/


clean:
	rm -rf bin/*
	go mod tidy


coverage: test
	go tool cover -func=coverage.out


coverage-report: test
	go tool cover -html=coverage.out


lint:
	${GOPATH}/bin/staticcheck ./...


test: build
	go test ./... -coverpkg=./internal... -race -covermode=atomic -coverprofile=coverage.out
