default: build

test:
	go test -race -bench=. .
	golangci-lint run

generate:
	go generate

e2e:
	go test -race -bench=. -tags=e2e .

build: test
	mkdir -p ./airtmp
	go build -o ./build/roadmapper .

docker: test
	GOOS=linux GOARCH=386 go build -o ./docker/roadmapper .
	docker build -t peteraba/roadmapper docker
	rm -f docker/roadmapper

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
	go build -o ./build/roadmapper -ldflags "-X main.version=${GIT_REV}" -ldflags "-X main.tag=${GIT_TAG}" .
	GOOS=linux GOARCH=386 go build -o ./docker/roadmapper -ldflags "-X main.version=${GIT_REV}" -ldflags "-X main.tag=${GIT_TAG}" .
	docker build -t peteraba/roadmapper:latest -t "peteraba/roadmapper:${GIT_TAG}" docker
	docker push peteraba/roadmapper:latest
	docker push "peteraba/roadmapper:${GIT_TAG}"

deploy:
	git pull
	docker pull peteraba/roadmapper
	docker-compose stop roadmapper
	docker-compose rm -f roadmapper
	docker-compose up -d roadmapper
	docker-compose exec roadmapper /roadmapper mu

.PHONY: default test generate e2e build docker install update release deploy
