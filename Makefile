PREFIX = /usr
GOBUILD = go build
GOTEST = go build
GORUN = go run
OUT_GIR_DIR = out/src/github.com/linuxdeepin/go-gir/
export GO111MODULE=off

all: build

CURRENT_DIR = $(shell pwd)
export GOPATH = $(CURDIR):$(CURDIR)/vendor:$(CURDIR)/out

GENERATOR = out/gir-generator

build: glib-2.0 gobject-2.0 gio-2.0 gudev-1.0 copyfile

generator:
	mkdir -p $(OUT_GIR_DIR)
	cd src/gir-generator && ${GOBUILD}  -o $(CURRENT_DIR)/${GENERATOR}

copyfile:
	cp -r  lib.in/gobject-2.0   $(OUT_GIR_DIR)
	cp -r  lib.in/gio-2.0       $(OUT_GIR_DIR)
	cp -r  lib.in/glib-2.0      $(OUT_GIR_DIR)
	cp -r  lib.in/gudev-1.0     $(OUT_GIR_DIR)

glib-2.0: lib.in/glib-2.0/glib.go.in lib.in/glib-2.0/config.json generator
	${GENERATOR} -o  $(OUT_GIR_DIR)$@ $<

gobject-2.0: lib.in/gobject-2.0/gobject.go.in lib.in/gobject-2.0/config.json generator
	${GENERATOR} -o $(OUT_GIR_DIR)$@ $<

gio-2.0:  lib.in/gio-2.0/gio.go.in lib.in/gio-2.0/config.json generator
	${GENERATOR} -o $(OUT_GIR_DIR)$@ $<

gudev-1.0: lib.in/gudev-1.0/gudev.go.in lib.in/gudev-1.0/config.json generator
	${GENERATOR} -o $(OUT_GIR_DIR)$@ $<

test:
	cd $(OUT_GIR_DIR)gobject-2.0 && ${GOTEST}
	cd $(OUT_GIR_DIR)gio-2.0 && ${GOTEST}
	cd $(OUT_GIR_DIR)glib-2.0 && ${GOTEST}
	cd $(OUT_GIR_DIR)gudev-1.0 && ${GOTEST}
	@echo "Memory Testing"
	#${GORUN} test/memory.go  阻塞打包

install:
	install -d  $(DESTDIR)$(PREFIX)/share/gocode/src/github.com/linuxdeepin/go-gir $(DESTDIR)$(PREFIX)/bin
	cp -r  $(OUT_GIR_DIR)*   $(DESTDIR)$(PREFIX)/share/gocode/src/github.com/linuxdeepin/go-gir
	cp     out/gir-generator $(DESTDIR)$(PREFIX)/bin/

clean:
	rm -fr out

.PHONY: test
