build:
	@go build

run:
	@go run main.go start

deps:
	@go get

default: @build
