.PHONY: build run

build:
	@go build -o bin/go-login-api

run: build
	@./bin/go-login-api
