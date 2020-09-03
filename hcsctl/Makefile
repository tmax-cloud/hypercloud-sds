.DEFAULT_GOAL := build
.PHONY: check build clean

check:
	golangci-lint run
build:
	pkger -include /sample-inventory
	go build -o build/hcsctl
	cd e2e/cdi && ginkgo build && mv cdi.test ../../build/cdi.test && cd ..
	cd e2e/rook && ginkgo build && mv rook.test ../../build/rook.test && cd ..
	cp -r e2e/rook/manifests build/test-manifests
clean:
	go clean ./...
	rm -rf build/
