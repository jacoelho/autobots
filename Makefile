.PHONY: crosscompile

crosscompile:
	GOARCH="amd64" GOOS="linux" go build
