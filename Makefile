BINARY := af

build-mac:
	env GOOS=darwin GOARCH=amd64 go build -v -o $(BINARY)-darwin-amd64

build-linux:
	env GOOS=linux GOARCH=amd64 go build -v -o $(BINARY)-linux-amd64

install-mac: build-mac
	mv -f ./$(BINARY)-darwin-amd64 /usr/local/bin/$(BINARY)
	chmod 755 /usr/local/bin/$(BINARY)

install-linux: build-linux
	mv -f ./$(BINARY)-linux-amd64 /usr/local/bin/$(BINARY)
	chmod 755 /usr/local/bin/$(BINARY)

uninstall:
	rm -v /usr/local/bin/$(BINARY)

clean:
	find . -name '$(BINARY)[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: all deps clean
