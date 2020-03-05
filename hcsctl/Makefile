.DEFAULT_GOAL := build
.PHONY: build clean

build:
	go build -o build/hcsctl
	cd e2e && ginkgo build && mv e2e.test ../build/hcsctl.test && cd ..
clean:
	go clean ./...
	rm -rf build/
