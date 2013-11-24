glib-2.0: glib-2.0/glib.go.in glib-2.0/config.json
	gir-generator -o ${GOPATH}/src/gobject/$@ $<
	go install gobject/$@

gobject-2.0: gobject-2.0/gobject.go.in gobject-2.0/config.json
	gir-generator -o ${GOPATH}/src/gobject/$@ $<
	cp gobject-2.0/gobject.c ${GOPATH}/src/gobject/$@
	cp gobject-2.0/gobject.h ${GOPATH}/src/gobject/$@
	go install gobject/$@

gio-2.0: gio-2.0/gio.go.in gio-2.0/config.json
	gir-generator -o ${GOPATH}/src/gobject/$@ $<
	go install gobject/$@

pango-1.0: pango-1.0/pango.go.in pango-1.0/config.json
	gir-generator -o ${GOPATH}/src/gobject/$@ $<
	go install gobject/$@

gdkpixbuf-2.0: gdkpixbuf-2.0/gdkpixbuf.go.in gdkpixbuf-2.0/config.json
	gir-generator -o ${GOPATH}/src/gobject/$@ $<
	go install gobject/$@

cairo-1.0: cairo-1.0/cairo.go
	cp -rf cairo-1.0/ ${GOPATH}/src/gobject/
	go install gobject/$@

gdk-3.0: gdk-3.0/gdk.go.in gdk-3.0/config.json
	gir-generator -o ${GOPATH}/src/gobject/$@ $<
	go install gobject/$@

all:  glib-2.0 gobject-2.0 gio-2.0 gdkpixbuf-2.0
