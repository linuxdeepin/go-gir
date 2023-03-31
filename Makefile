PREFIX = /usr
GOPATH_DIR = gopath
GOPKG_PREFIX = github.com/linuxdeepin/go-gir
export GO111MODULE=off

# FIXME: not working! export GOPATH= $(shell go env GOPATH)
export GOPATH = /usr/share/gocode

all: build

build: test

generate:
	go generate ./...

prepare:
	@mkdir -p ${GOPATH_DIR}/src/$(dir ${GOPKG_PREFIX});
	@ln -snf ../../../.. ${GOPATH_DIR}/src/${GOPKG_PREFIX};

test: clean prepare
	env GOPATH="${CURDIR}/${GOPATH_DIR}:${GOPATH}" go test ./... -coverpkg=${GOPATH_DIR}

install: clean
	cp -r * $(DESTDIR)$(PREFIX)/share/gocode/src/github.com/linuxdeepin/go-gir

clean:
	rm -rf ${GOPATH_DIR}

.PHONY: test
