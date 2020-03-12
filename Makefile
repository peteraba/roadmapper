test:
	go test .
	golangci-lint run

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

release:
	$(eval GIT_REV=$(shell git rev-parse HEAD | cut -c1-8))
	$(eval GIT_TAG=$(shell git describe --exact-match --tags $(git log -n1 --pretty='%h')))
	go build -o ./build/roadmapper -ldflags "-X main.version=${GIT_REV}" -ldflags "-X main.tag=${GIT_TAG}" .

.PHONY: build docker install update release
