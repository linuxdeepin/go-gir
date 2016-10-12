/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

const g_object_ref_unref = `extern GObject *g_object_ref_sink(GObject*);
extern void g_object_unref(GObject*);`

const g_error_free = `extern void g_error_free(GError*);`

const g_free = `extern void g_free(void*);`

const g_list_funcs = `
GList* g_list_append(GList*, void*);
void g_list_free(GList*);
`

var go_utils_template = must_template(`
const alot = 999999

type _GSList struct {
	data unsafe.Pointer
	next *_GSList
}

type _GList struct {
	data unsafe.Pointer
	next *_GList
	prev *_GList
}

type _GError struct {
	domain uint32
	code int32
	message *C.char
}
func (e _GError) ToGError() GError {
	return GError{e.domain, e.code, C.GoString(e.message)}
}

type GError struct {
	Domain uint32
	Code int32
	Message string
}
func (e GError) Error() string {
	return e.Message
}

func _GoStringToGString(x string) *C.char {
	if x == "\x00" {
		return nil
	}
	return C.CString(x)
}

func _GoBoolToCBool(x bool) C.int {
	if x { return 1 }
	return 0
}

func _CInterfaceToGoInterface(iface [2]unsafe.Pointer) interface{} {
	return *(*interface{})(unsafe.Pointer(&iface))
}

func _GoInterfaceToCInterface(iface interface{}) *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Pointer(&iface))
}

[<if not .nocallbacks>]
//export _[<.namespace>]_go_callback_cleanup
func _[<.namespace>]_go_callback_cleanup(gofunc unsafe.Pointer) {
	[<.gobjectns>]Holder.Release(gofunc)
}
[<end>]
`)

var object_template = must_template(`
type [<.name>]Like interface {
	[<.parentlike>]
	InheritedFrom[<.cprefix>][<.name>]() [<.cgotype>]
}

type [<.name>] struct {
	[<.parent>]
	[<.interfaces>]
}

func To[<.name>](objlike [<.gobjectns>]ObjectLike) *[<.name>] {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*[<.name>])(nil).GetStaticType()
	obj := [<.gobjectns>]ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*[<.name>])(obj)
	}
	panic("cannot cast to [<.name>]")
}

func (this0 *[<.name>]) InheritedFrom[<.cprefix>][<.name>]() [<.cgotype>] {
	if this0 == nil {
		return nil
	}
	return ([<.cgotype>])(this0.C)
}

func (this0 *[<.name>]) GetStaticType() [<.gobjectns>]Type {
	return [<.gobjectns>]Type(C.[<.typeinit>]())
}

func [<.name>]GetType() [<.gobjectns>]Type {
	return (*[<.name>])(nil).GetStaticType()
}
`)

// XXX: uses gc specific hack, expect problems on gccgo and/or ask developers
// about the address of an empty embedded struct
var interface_template = must_template(`
type [<.name>]Like interface {
	Implements[<.cprefix>][<.name>]() [<.cgotype>]
}

type [<.name>] struct {
	[<.gobjectns>]Object
	[<.name>]Impl
}

func (*[<.name>]) GetStaticType() [<.gobjectns>]Type {
	return [<.gobjectns>]Type(C.[<.typeinit>]())
}


type [<.name>]Impl struct {}

func To[<.name>](objlike [<.gobjectns>]ObjectLike) *[<.name>] {
	c := objlike.InheritedFromGObject()
	obj := [<.gobjectns>]ObjectGrabIfType(unsafe.Pointer(c), [<.gobjectns>]Type(C.[<.typeinit>]()))
	if obj != nil {
		return (*[<.name>])(obj)
	}
	panic("cannot cast to [<.name>]")
}

func (this0 *[<.name>]Impl) Implements[<.cprefix>][<.name>]() [<.cgotype>] {
	base := unsafe.Pointer(uintptr(unsafe.Pointer(this0)) - unsafe.Sizeof(uintptr(0)))
	return ([<.cgotype>])((*[<.gobjectns>]Object)(base).C)
}

`)

const c_header = `#pragma once
#include <stdlib.h>
#include <stdint.h>

typedef size_t GType;
typedef void *GVaClosureMarshal;
static unsigned int _array_length(void* _array)
{
	void** array = (void**)_array;
	unsigned int i=0;
	while (array && array[i] != 0) i++;
	return i;
}

`

var c_template = must_template(`
#include "[<.package>].gen.h"
#include "_cgo_export.h"

static void _c_callback_cleanup(void *userdata)
{
	_[<.namespace>]_go_callback_cleanup(userdata);
}

`)
