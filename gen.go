package gir

//go:generate go run ./gir-generator -o glib-2.0 lib.in/glib-2.0/glib.go.in
//go:generate go run ./gir-generator -o gobject-2.0 lib.in/gobject-2.0/gobject.go.in
//go:generate go run ./gir-generator -o gio-2.0 lib.in/gio-2.0/gio.go.in
//go:generate go run ./gir-generator -o gudev-1.0 lib.in/gudev-1.0/gudev.go.in
