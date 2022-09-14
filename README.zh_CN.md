# go-gir-generator

## 描述

go-gir-generator项目是用Go语言实现的基于Object-introspection的绑定生成器。

有许多基于GObject/Gtk库实现的Go绑定生成器，但几乎所有这些都是手工编写的。编写的过程无聊且容易出错，因为相应的C库可能会经常更改。 go-gir-generator项目使得开发者无需修改一行代码即可适配最新版本的 GObject，减少了用 GObject 编写 Go 代码的痛苦。

*注: 目前只支持GObject-2.0、Glib-2.0、Gio-2.0，对Gdk/Gtk的支持还没有完成。*

非常感谢创建 [GObject Introspection](https://wiki.gnome.org/action/show/Projects/GObjectIntrospection) 的天才们和原作者 [nsf](https://github.com/nsf/gogobject）。

## 依赖

### 编译依赖

- golang 1.3 above
- pkg-config
- gobject-introspection-1.0

### 运行依赖

- gobject-introspection-1.0

## 安装

go-gir-generator需要预安装以下包

```
$ sudo apt-get install libgirepository1.0-dev libgudev-1.0-dev
$ make install
```

## 使用方法

可执行程序gir-generator是一个静态绑定的代码生成器。
它从 lib.in 目录中读取 GIRepository 和模板文件。

例如，如果我们要生成gobject-2.0时，

```
cd lib.in/gobject-2.0
gir-generator -o . gobject.go.in
```

*注：项目生成的所有绑定在`$GOPATH/src/gir/$PackageName`这个目录下。*

## 获得帮助

如果您遇到任何其他问题，您可能还会发现这些渠道很有用：

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org)
* [WiKi](http://wiki.deepin.org/)

## 贡献指南

我们鼓励您报告问题并做出更改

* [Contribution guide for users](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Users)
* [Contribution guide for developers](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Developers)

## 开源许可
go-gir-generator项目在LGPL-3.0-or-later开源协议下发布。

