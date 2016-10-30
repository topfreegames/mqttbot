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

create-es-index-template:
	@bash create_es_index_template.sh

run-tests: run-containers
	@/bin/bash -c "until docker exec testcontainers_elasticsearch_1 curl localhost:9200; do echo 'Waiting for Elasticsearch...' && sleep 1; done"
	@make create-es-index-template
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

cross: cross-linux cross-darwin

cross-linux:
	@mkdir -p ./bin
	@echo "Building for linux-i386..."
	@env GOOS=linux GOARCH=386 go build -o ./bin/mqttbot-linux-i386 ./main.go
	@echo "Building for linux-x86_64..."
	@env GOOS=linux GOARCH=amd64 go build -o ./bin/mqttbot-linux-x86_64 ./main.go
	@$(MAKE) cross-exec

cross-darwin:
	@mkdir -p ./bin
	@echo "Building for darwin-i386..."
	@env GOOS=darwin GOARCH=386 go build -o ./bin/mqttbot-darwin-i386 ./main.go
	@echo "Building for darwin-x86_64..."
	@env GOOS=darwin GOARCH=amd64 go build -o ./bin/mqttbot-darwin-x86_64 ./main.go
	@$(MAKE) cross-exec

cross-exec:
	@chmod +x bin/*

default: @build
