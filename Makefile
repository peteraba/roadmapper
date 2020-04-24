VERSION			:= snapshot
NAME			:= roadmapper

GIT_REV			:= $(shell git rev-parse HEAD | cut -c1-8)
#GIT_TAG			:= $(shell git describe --exact-match --tags $(git log -n1 --pretty='%h'))
PACKAGES		:= $(shell find -name "*.go" 2>&1 | grep -v "Permission denied" | grep -v -e bindata | xargs -n1 dirname | sort -u)
MAIN_DIR		:= ./cmd/$(NAME)
BUILD_OUTPUT	:= ./build/$(NAME)
DOCKER_OUTPUT	:= ./docker/$(NAME)
DOCKER_DIR		:= ./docker
LINK_FLAGS		:= -X main.AppVersion=$(GIT_REV) -X main.GitTag=$(GIT_TAG)
DOCKER_IMAGE	:= peteraba/$(NAME)

default: build

debug:
	echo $(PACKAGES)

test:
	golangci-lint --version
	golangci-lint run $(PACKAGES)
	go test -race -bench=. $(PACKAGES)

generate:
	go generate $(PACKAGES)
	find pkg -name "mocks" -type d -exec rm -rf {} +

e2e:
	go test -race -tags=e2e $(PACKAGES)

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
	go build -o ./build/roadmapper -ldflags="$(LINK_FLAGS)" $(MAIN_DIR)
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(DOCKER_IMAGE) -ldflags="$(LINK_FLAGS)" $(MAIN_DIR)
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

.PHONY: default debug test generate e2e build docker install update release deploy
