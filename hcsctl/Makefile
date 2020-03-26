.DEFAULT_GOAL := build
.PHONY: check build clean

check:
	golangci-lint run
build:
	pkger -include /sample-inventory
	go build -o build/hcsctl
	cd e2e && ginkgo build && mv e2e.test ../build/hcsctl.test && cd ..
clean:
	go clean ./...
	rm -rf build/
