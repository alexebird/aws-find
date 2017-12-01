BINARY := af
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SOURCE_FILES := main.go

.PHONY: all
all: pkg/darwin_amd64/$(BINARY) pkg/linux_amd64/$(BINARY)

pkg/darwin_amd64/$(BINARY): $(SOURCE_FILES)
	GOOS=darwin GOARCH=amd64 \
	go build -v -o "$@"

pkg/linux_amd64/$(BINARY): $(SOURCE_FILES)
	GOOS=linux GOARCH=amd64 \
	go build -v -o "$@"

install:
	cp pkg/linux_amd64/$(BINARY) /usr/local/bin/$(BINARY)
	chmod 755 /usr/local/bin/$(BINARY)

.PHONY: deps
deps:
	go get -v -d

.PHONY: clean
clean:
	go clean -i -x
