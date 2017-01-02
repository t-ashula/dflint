

COMMIT = $$(git describe --always)
DEBUG_FLAG = $(if $(DEBUG),-debug)

updatedeps: devdeps
	@echo "====> Install & Update depedencies..."
	glide up

devdeps:
	@echo "====> Install depedencies for development..."
	go get github.com/Masterminds/glide

deps: devdeps
	@echo "====> Install depedencies..."
	glide -q install

build: deps
	@echo "====> Build dflint in . "
	go build

install: build
	@echo "====> Install dflint in $(GOPATH)/bin ..."
	@go install

all: build
	@echo "====> All "

test: build devdeps
	@echo "====> Run test"
	go test $$(glide novendor)

test-cover: build devdeps
	go test -cover $$(glide novendor)
