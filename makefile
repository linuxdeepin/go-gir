#!/usr/bin/make -f


PREFIX = /usr

all: build

build:
	mkdir -p bin
	GOPATH=`pwd`:`pwd`/vendor go build -o bin/gir-generator gir-generator

install:
	mkdir -p ${DESTDIR}${PREFIX}/bin/
	cp bin/gir-generator ${DESTDIR}${PREFIX}/bin/
clean:
	rm -f bin

