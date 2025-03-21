repo_version=$(shell cat VERSION)
build_time=$(shell LC_TIME=C date -u +"%a %b %d %H:%M:%S UTC %Y")
revision=$(shell git rev-parse --short HEAD)
branch=$(shell git rev-parse --abbrev-ref HEAD)

ifeq ($(branch), HEAD)
    branch1=$(shell git show -s --pretty=%d HEAD | awk '{ print $$3 }')
    branch=$(shell echo $(branch1) | sed 's/[),]//g')
endif

modified_version=$(repo_version)
ifeq ($(branch), main)
    modified_version=$(repo_version)-$(revision)
endif

LDFLAGS = "\
    -X 'github.com/AndreyShep2012/go-company-handler/internal/version.Version=$(modified_version)' \
    -X 'github.com/AndreyShep2012/go-company-handler/internal/version.Revision=$(revision)' \
    -X 'github.com/AndreyShep2012/go-company-handler/internal/version.Branch=$(branch)' \
    -X 'github.com/AndreyShep2012/go-company-handler/internal/version.BuildTime=$(build_time)' \
    -s -w"

dev:
	go run cmd/main.go

build:
	@echo "Building app"
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o ./bin/main ./cmd/main.go

run: build
	./bin/main

run-docker:
	docker-compose up --build -d

stop-docker:
	docker-compose down

fmt:
	go fmt ./...

lint:
	golangci-lint run --timeout 5m -c .golangci.yml

test-unit:
	go test -v -race -count=1 -cover \
		-coverpkg github.com/AndreyShep2012/go-company-handler/internal... \
		-coverprofile="./coverage.out" ./internal/...

test-itg:
	docker-compose -f ./tests/integration/docker-compose.yaml up -d
	go test -v -race -count=1 -cover \
		-coverpkg github.com/AndreyShep2012/go-company-handler/internal... \
		-coverprofile="./coverage.out" ./tests/integration/...
	docker-compose -f ./tests/integration/docker-compose.yaml down

test-all:
	docker-compose -f ./tests/integration/docker-compose.yaml up -d
	go test -v -race -count=1 -cover \
		-coverpkg github.com/AndreyShep2012/go-company-handler/internal... \
		-coverprofile="./coverage.out" ./internal/... ./tests/integration/...
	docker-compose -f ./tests/integration/docker-compose.yaml down
	go tool cover -html=coverage.out