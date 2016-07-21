PACKAGES = $(shell glide novendor)
GODIRS = $(shell go list ./... | grep -v /vendor/ | sed s@github.com/topfreegames/mqttbot@.@g | egrep -v "^[.]$$")

setup:
	@go get -u github.com/Masterminds/glide/...
	@glide install

setup-ci:
	@sudo add-apt-repository -y ppa:masterminds/glide && sudo apt-get update
	@sudo apt-get install -y glide
	@go get github.com/mattn/goveralls
	@glide install

build:
	@go build $(PACKAGES)
	@go build

run-containers:
	@cd test_containers && docker-compose up -d && cd ..

kill-containers:
	@cd test_containers && docker-compose stop && cd ..

run-tests: run-containers
	@sleep 5 #wait for es to start
	@make coverage
	@make kill-containers

test: run-tests

coverage:
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg) || exit 1 &&\
		tail -n +2 coverage.out >> coverage-all.out;)

run:
	@go run main.go start

deps:
	@glide install

default: @build
