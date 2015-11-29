all:  go-dir glib-2.0 gobject-2.0 gio-2.0

go-dir: 
	mkdir -p go/src/pkg.deepin.io

glib-2.0: glib-2.0/glib.go.in glib-2.0/config.json
	gir-generator -config $@/config.json -o go/src/pkg.deepin.io/$@ $@/glib.go.in
	bash -c 'cp -rf $@/* go/src/pkg.deepin.io/$@'
	bash -c 'GOPATH=$$GOPATH:`pwd`/go go test pkg.deepin.io/$@'

gobject-2.0: gobject-2.0/gobject.go.in gobject-2.0/config.json
	gir-generator -config $@/config.json -o go/src/pkg.deepin.io/$@ $@/gobject.go.in
	bash -c 'cp -rf $@/* go/src/pkg.deepin.io/$@'
	bash -c 'GOPATH=$$GOPATH:`pwd`/go go test pkg.deepin.io/$@'

gio-2.0: gio-2.0/gio.go.in gio-2.0/config.json
	gir-generator -config $@/config.json -o go/src/pkg.deepin.io/$@ $@/gio.go.in
	bash -c 'cp -rf $@/* go/src/pkg.deepin.io/$@'
	bash -c 'GOPATH=$$GOPATH:`pwd`/go go test pkg.deepin.io/$@'

