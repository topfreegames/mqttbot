PACKAGES = $(shell glide novendor)
GODIRS = $(shell go list ./... | grep -v /vendor/ | sed s@github.com/topfreegames/mqttbot@.@g | egrep -v "^[.]$$")

setup:
	@go get -u github.com/Masterminds/glide/...
	@glide install

setup-ci:
	@sudo add-apt-repository -y ppa:masterminds/glide && sudo apt-get update
	@sudo apt-get install -y glide
	@go get github.com/mattn/goveralls
	@cd mosquitto_test && ./run.sh && cd ..
	@glide install

build:
	@go build $(PACKAGES)
	@go build

coverage:
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
run:
	@go run main.go start

deps:
	@glide install

default: @build
