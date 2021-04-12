export APP_CMD_NAME = userservice
export REGISTRY = vadimmakerov/music-streaming
export APP_PROTO_FILES = \
	api/userservice/userservice.proto
export DOCKER_IMAGE_NAME = $(REGISTRY)-$(APP_CMD_NAME):master

all: build check test

.PHONY: build
build: generate modules
	bin/go-build.sh "cmd" "bin/$(APP_CMD_NAME)" $(APP_CMD_NAME)

.PHONY: generate
generate:
	bin/generate-grpc.sh $(foreach path,$(APP_PROTO_FILES),"$(path)")

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run

.PHONY: publish
publish:
	docker build . --tag=$(DOCKER_IMAGE_NAME)