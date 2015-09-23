.PHONY: crosscompile lint

all: crosscompile

lint:
	gofmt -w $$(pwd)

crosscompile:
	GOARCH="amd64" GOOS="linux" go build
