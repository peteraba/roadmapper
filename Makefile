test:
	go test .
	golangci-lint run
	mkdir -p ./airtmp
	go build -o ./build/roadmapper .

build: test
	mkdir -p ./airtmp
	go build -o ./build/roadmapper .

docker: test
	GOOS=linux GOARCH=386 go build -o ./docker/roadmapper .
	docker build -t peteraba/roadmapper docker
	rm -f docker/roadmapper

install:
	# Install [golangci-lint](https://github.com/golangci/golangci-lint)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.23.8

update:
	go get -u ./...

.PHONY: build docker install update
