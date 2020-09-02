.DEFAULT_GOAL := build
.PHONY: check build clean

check:
	golangci-lint run
build:
	pkger -include /sample-inventory
	go build -o build/hcsctl
	cd e2e/cdi && ginkgo build && mv cdi.test ../../build/cdi.test && cd ..
	cd e2e/ceph && ginkgo build && mv ceph.test ../../build/ceph.test && cd ..
clean:
	go clean ./...
	rm -rf build/
