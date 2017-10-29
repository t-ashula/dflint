NAME   := dflint
VERSION  := v0.0.4
REVISION := $(shell git rev-parse --short HEAD)

SRCS     := *.go rule/*.go
PKGS     := $(shell go list ./...)
LDFLAGS  := "-X main.Version=$(VERSION)"

all: $(NAME)

$(NAME): $(SRCS)
	go build -ldflags $(LDFLAGS) -o $(NAME)

dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get github.com/golang/dep/cmd/dep
endif

deps: dep
	dep ensure

test:
	go test -cover $(PKGS)

test-cover:
	echo 'mode: atomic' > cover-all.out
	$(foreach pkg, $(PKGS), \
		go test -coverprofile=cover.out -covermode=atomic -v $(pkg); \
		tail -n +2 cover.out >> cover-all.out; \
	)
	go tool cover -func=cover-all.out

clean:
	rm -rf $(NAME) cover-all.out cover.out bin vendor

force: clean all

install:
	go install

.PHONY: force clean test-cover test all deps dep
