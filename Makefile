VERSION			:= snapshot
NAME			:= roadmapper

PACKAGES		:= $(shell find -name "*.go" 2>&1 | grep -v "Permission denied" | grep -v -e bindata | grep -v -e server | xargs -n1 dirname | uniq | sort -u)
MAIN_DIR		:= ./cmd/$(NAME)
BUILD_OUTPUT	:= ./build/$(NAME)
DOCKER_OUTPUT	:= ./docker/$(NAME)
DOCKER_DIR		:= ./docker
DOCKER_IMAGE	:= peteraba/$(NAME)

default: build

debug:
	@ echo $(PACKAGES)

generate:
	go generate $(PACKAGES)
	find pkg -name "mocks" -type d -exec rm -rf {} +

goldenfiles:
	go test -mod=readonly -tags=e2e,integration ./cmd/roadmapper -update
	go test -mod=readonly ./pkg/roadmap -update

test:
	golangci-lint --version
	golangci-lint run $(PACKAGES)
	go test -race -bench=. $(PACKAGES)

integration:
	go test -race -tags=integration $(PACKAGES)

e2e:
	go test -race -v -tags=e2e,integration $(PACKAGES)

codecov:
ifndef CODECOV_TOKEN
	$(error CODECOV_TOKEN is not set)
endif
	# Download codecov
	curl -o b.sh https://codecov.io/bash
	chmod +x b.sh
	# Code coverage for Go unit tests
	go test -race -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)
	./b.sh -c -F go_unittests
	# Code coverage for All tests
	go test -race -count=1 -coverprofile=coverage.txt -covermode=atomic -tags=e2e,integration ./...
	./b.sh -c -F alltests
	rm -f b.sh

build: test
	mkdir -p ./airtmp
	go build -o $(BUILD_OUTPUT) $(MAIN_DIR)

docker: test
	GOOS=linux GOARCH=386 go build -o $(DOCKER_OUTPUT) $(MAIN_DIR)
	docker build -t $(DOCKER_IMAGE) $(DOCKER_DIR)
	rm -f $(DOCKER_OUTPUT)

install:
	# Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.23.8
	# Install goreleaser
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh

update:
	go get -u ./...

release: e2e
	$(eval GIT_REV=$(shell git rev-parse HEAD | cut -c1-8))
	$(eval GIT_TAG=$(shell git describe --exact-match --tags $(git log -n1 --pretty='%h')))
	go build -o ./build/roadmapper -ldflags="-X main.AppVersion=${GIT_REV} -X main.GitTag=${GIT_TAG}" $(MAIN_DIR)
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(DOCKER_OUTPUT) -ldflags="-X main.AppVersion=${GIT_REV} -X main.GitTag=${GIT_TAG}" $(MAIN_DIR)
	docker build -t "${DOCKER_IMAGE}:latest" -t "${DOCKER_IMAGE}:${GIT_TAG}" $(DOCKER_DIR)
	docker push "${DOCKER_IMAGE}:latest"
	docker push "${DOCKER_IMAGE}:${GIT_TAG}"

deploy:
	git pull
	docker pull $(DOCKER_IMAGE)
	docker-compose stop $(NAME)
	docker-compose rm -f $(NAME)
	docker-compose up -d $(NAME)
	docker-compose exec $(NAME) "/${NAME}" mu

clean:
	rm -rvf coverfileprofile.txt

.PHONY: default debug generate goldenfiles test e2e codecov build docker install update release deploy clean
