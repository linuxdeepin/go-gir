PREFIX = /usr

ifndef USE_GCCGO
	GOBUILD = go build
	GOTEST = go build
	GORUN = go run
else
	LDFLAGS = $(shell pkg-config --libs gobject-introspection-1.0 gio-2.0 gudev-1.0 gdk-3.0)
	GOBUILD = go build -compiler gccgo -gccgoflags "${LDFLAGS}"
	GOTEST = go test -compiler gccgo -gccgoflags "${LDFLAGS}"
	GORUN = go run -compiler gccgo -gccgoflags "${LDFLAGS}"
endif



all: build

CURRENT_DIR = $(shell pwd)
export GOPATH = $(CURDIR):$(CURDIR)/vendor:$(CURDIR)/out

GENERATOR = out/gir-generator

build: glib-2.0 gobject-2.0 gio-2.0 gudev-1.0

generator:
	mkdir -p out/src/gir
	cd src/gir-generator && ${GOBUILD}  -o $(CURRENT_DIR)/${GENERATOR}

copyfile:
	cp -r  lib.in/gobject-2.0/*   out/src/gir/gobject-2.0
	cp -r  lib.in/gio-2.0/*       out/src/gir/gio-2.0
	cp -r  lib.in/glib-2.0/*      out/src/gir/glib-2.0
	cp -r  lib.in/gudev-1.0/*     out/src/gir/gudev-1.0

glib-2.0: lib.in/glib-2.0/glib.go.in lib.in/glib-2.0/config.json generator
	${GENERATOR} -o  out/src/gir/$@ $<

gobject-2.0: lib.in/gobject-2.0/gobject.go.in lib.in/gobject-2.0/config.json generator
	${GENERATOR} -o out/src/gir/$@ $<

gio-2.0:  lib.in/gio-2.0/gio.go.in lib.in/gio-2.0/config.json generator
	${GENERATOR} -o out/src/gir/$@ $<

gudev-1.0: lib.in/gudev-1.0/gudev.go.in lib.in/gudev-1.0/config.json generator
	${GENERATOR} -o out/src/gir/$@ $<

test: copyfile glib-2.0 gobject-2.0 gio-2.0 gudev-1.0
	cd out/src/gir/gobject-2.0 && ${GOTEST}
	cd out/src/gir/gio-2.0 && ${GOTEST}
	cd out/src/gir/glib-2.0 && ${GOTEST}
	cd out/src/gir/gudev-1.0 && ${GOTEST}
	@echo "Memory Testing"
	GOPATH=`pwd`/out ${GORUN} test/memory.go

install: copyfile
	install -d  $(DESTDIR)$(PREFIX)/share/gocode/src/gir $(DESTDIR)$(PREFIX)/bin
	cp -r  out/src/gir/*   $(DESTDIR)$(PREFIX)/share/gocode/src/gir
	cp     out/gir-generator $(DESTDIR)$(PREFIX)/bin/

clean:
	rm -fr out
