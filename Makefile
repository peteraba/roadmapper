build:
	go test .
	golangci-lint run
	mkdir -p ./airtmp
	go build -o ./build/roadmapper .

install:
	# Install [golangci-lint](https://github.com/golangci/golangci-lint)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.23.8

.PHONY: build install dev
