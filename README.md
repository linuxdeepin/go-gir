# Go-gir-generator

**Description**

Go-gir-generator imeplement static golang bindings for GObject.

There has many go bindings for GObject/Gtk libraries, but almost all of them
are written by hand. It's bored and error-prone when the binding C libaray changed.

Go-gir-geneator's object is like python-gobject's that binding the newest library
without need change binding codes.

Currently it only official support Gobject-2.0, Glib-2.0, Gio-2.0.
Because generate the gdkpixbuf binding hasn't completed, so Gdk/Gtk were
also in blocking.


Thanks the genius guys who created the
[Gobject Introspection](https://wiki.gnome.org/action/show/Projects/GObjectIntrospection)
and the orignal author [nsf](https://github.com/nsf/gogobject).

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
$ apt-get install libgirepository1.0-dev
$ make install
```

## Usage

The binary program gir-generator is the static binding
code generator.
It read the GIRepository and tempalte files from lib.in direactory.

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

Startdde is licensed under [GPLv3](LICENSE).

