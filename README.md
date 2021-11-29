# go-gir-generator

## Description

Project go-gir-generator implements GObject-introspection based bindings generator for Go.

There are many Go bindings of GObject/Gtk libraries out there, but almost all of them are written by hand. It's boring and error-prone since the corresponding C library may change very often. go-gir-generator can adapt to the newest version of GObject without changing one single line of code, which makes it less painful to write Go code with GObject.

*NOTE: Currently it only supports GObject-2.0, Glib-2.0, Gio-2.0, support of Gdk/Gtk is not completed.*


Many thanks to the genius guys who created the [GObject Introspection](https://wiki.gnome.org/action/show/Projects/GObjectIntrospection) and the original author [nsf](https://github.com/nsf/gogobject).

## Dependencies

### Build dependencies

- golang 1.3 above
- pkg-config
- gobject-introspection-1.0

### Runtime dependencies

- gobject-introspection-1.0

## Installation

Install prerequisites

```
$ sudo apt-get install libgirepository1.0-dev libgudev-1.0-dev
$ make install
```

## Usage

The binary program gir-generator is the static binding
code generator.
It read the GIRepository and template files from lib.in directory.

For example, we need generate gobject-2.0 binding,

```
cd lib.in/gobject-2.0
gir-generator -o . gobject.go.in
```

Note: deepin generate all bindings under $GOPATH/src/gir/$PackageName

## Getting help

Any usage issues can ask for help via

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org)
* [WiKi](http://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes

* [Contribution guide for users](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Users)
* [Contribution guide for developers](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Developers)

## License

go-gir-generator is licensed under [GPLv3](LICENSE).

