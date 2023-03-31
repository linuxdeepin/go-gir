PREFIX = /usr

all: build

build: test

generate:
	go generate ./...

test:
	go test github.com/linuxdeepin/go-gir/...

install:
	cp -r * $(DESTDIR)$(PREFIX)/share/gocode/src/github.com/linuxdeepin/go-gir

.PHONY: test
