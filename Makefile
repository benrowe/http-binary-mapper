.PHONY: help

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## setup the project for exection
init:
	go get ./...

## generate the binary
build: init
	go build -o proxy

## generate the binary
build-nix:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy .

## execute the program from the source files
run: build
	./proxy -cfg=config.example.yaml

## Run the unit tests
test: init
	go test

## Run the coverage report and output it as html
coverage: init
	go test -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html

## Build the docker image and tag it as "benrowe/http-binary-mapper"
docker-build:
	docker build -t benrowe/http-binary-mapper:latest .

## Run the service through docker
docker-run: docker-build
	docker run --rm -it -v "${PWD}:/src/data" -p 8000:8000 benrowe/http-binary-mapper:latest
