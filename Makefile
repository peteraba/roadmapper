default: build

test:
	go test .
	golangci-lint run

generate:
	go generate

integration:
	go test -tags=integration .

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

release: generate integration
	$(eval GIT_REV=$(shell git rev-parse HEAD | cut -c1-8))
	$(eval GIT_TAG=$(shell git describe --exact-match --tags $(git log -n1 --pretty='%h')))
	go build -o ./build/roadmapper -ldflags "-X main.version=${GIT_REV}" -ldflags "-X main.tag=${GIT_TAG}" .
	GOOS=linux GOARCH=386 go build -o ./docker/roadmapper -ldflags "-X main.version=${GIT_REV}" -ldflags "-X main.tag=${GIT_TAG}" .
	docker build -t peteraba/roadmapper:latest -t "peteraba/roadmapper:${GIT_TAG}" docker
	docker push peteraba/roadmapper:latest
	docker push "peteraba/roadmapper:${GIT_TAG}"

.PHONY: default test generate integration build docker install update release
