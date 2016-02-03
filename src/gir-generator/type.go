/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
// Converting GIR type to Go/C representation

package main

import (
	"bytes"
	"fmt"
	"gi"
	"strings"
	"unsafe"
)

func force_pointer(x string) string {
	if x == "unsafe.Pointer" {
		return x
	}
	if !strings.HasPrefix(x, "*") {
		return "*" + x
	}
	return x
}

type type_flags int

const (
	type_none    type_flags = 0
	type_pointer type_flags = 1 << iota
	type_return
	type_list_member
	type_receiver
	type_exact
)

//------------------------------------------------------------------
// Cgo Type (C type in Go)
//------------------------------------------------------------------

func cgo_type(ti *gi.TypeInfo, flags type_flags) string {
	var out bytes.Buffer

	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			out.WriteString("unsafe.Pointer")
			break
		}
		panic("Non-pointer void type is not supported in cgo")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME:
		out.WriteString("*C.char")
	case gi.TYPE_TAG_ARRAY:
		switch ti.ArrayType() {
		case gi.ARRAY_TYPE_C:
			out.WriteString("*")
			out.WriteString(cgo_type(ti.ParamType(0), flags))
		case gi.ARRAY_TYPE_ARRAY:
			out.WriteString("*C.GArray")
		case gi.ARRAY_TYPE_PTR_ARRAY:
			out.WriteString("*C.GPtrArray")
		case gi.ARRAY_TYPE_BYTE_ARRAY:
			out.WriteString("*C.GByteArray")
		}
	case gi.TYPE_TAG_GLIST:
		out.WriteString("*C.GList")
	case gi.TYPE_TAG_GSLIST:
		out.WriteString("*C.GSList")
	case gi.TYPE_TAG_GHASH:
		out.WriteString("*C.GHashTable")
	case gi.TYPE_TAG_ERROR:
		out.WriteString("*C.GError")
	case gi.TYPE_TAG_INTERFACE:
		if ti.IsPointer() {
			flags |= type_pointer
		}
		out.WriteString(cgo_type_for_interface(ti.Interface(), flags))
	default:
		if ti.IsPointer() {
			out.WriteString("*")
		}
		out.WriteString(cgo_type_for_tag(tag, flags))
	}

	return out.String()
}

func cgo_type_for_interface(bi *gi.BaseInfo, flags type_flags) string {
	var out bytes.Buffer

	switch bi.Type() {
	case gi.INFO_TYPE_CALLBACK:
		out.WriteString("unsafe.Pointer")
	default:
		ns := bi.Namespace()
		nm := bi.Name()
		fullnm := strings.ToLower(ns) + "." + nm

		if flags&type_pointer != 0 && !config.is_disguised(fullnm) {
			out.WriteString("*")
		}

		out.WriteString("C.")
		out.WriteString(gi.DefaultRepository().CPrefix(ns))
		out.WriteString(bi.Name())
	}
	return out.String()
}

func cgo_type_for_tag(tag gi.TypeTag, flags type_flags) string {
	var out bytes.Buffer
	p := printer_to(&out)

	if flags&type_pointer != 0 {
		p("*")
	}

	switch tag {
	case gi.TYPE_TAG_BOOLEAN:
		p("C.int")
	case gi.TYPE_TAG_INT8:
		p("C.int8_t")
	case gi.TYPE_TAG_UINT8:
		p("C.uint8_t")
	case gi.TYPE_TAG_INT16:
		p("C.int16_t")
	case gi.TYPE_TAG_UINT16:
		p("C.uint16_t")
	case gi.TYPE_TAG_INT32:
		p("C.int32_t")
	case gi.TYPE_TAG_UINT32:
		p("C.uint32_t")
	case gi.TYPE_TAG_INT64:
		p("C.int64_t")
	case gi.TYPE_TAG_UINT64:
		p("C.uint64_t")
	case gi.TYPE_TAG_FLOAT:
		p("C.float")
	case gi.TYPE_TAG_DOUBLE:
		p("C.double")
	case gi.TYPE_TAG_GTYPE:
		p("C.GType")
	case gi.TYPE_TAG_UNICHAR:
		p("C.uint32_t")
	default:
		panic("unreachable")
	}

	return out.String()
}

//------------------------------------------------------------------
// C Type
//------------------------------------------------------------------

func c_type_for_tag(tag gi.TypeTag, flags type_flags) string {
	return cgo_type_for_tag(tag, flags)[2:]
}

func c_type(ti *gi.TypeInfo, flags type_flags) string {
	var out bytes.Buffer

	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			out.WriteString("void*")
			break
		}
		out.WriteString("void")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME:
		out.WriteString("char*")
	case gi.TYPE_TAG_ARRAY:
		switch ti.ArrayType() {
		case gi.ARRAY_TYPE_C:
			out.WriteString(c_type(ti.ParamType(0), flags))
			out.WriteString("*")
		case gi.ARRAY_TYPE_ARRAY:
			out.WriteString("GArray*")
		case gi.ARRAY_TYPE_PTR_ARRAY:
			out.WriteString("GPtrArray*")
		case gi.ARRAY_TYPE_BYTE_ARRAY:
			out.WriteString("GByteArray*")
		}
	case gi.TYPE_TAG_GLIST:
		out.WriteString("GList*")
	case gi.TYPE_TAG_GSLIST:
		out.WriteString("GSList*")
	case gi.TYPE_TAG_GHASH:
		out.WriteString("GHashTable*")
	case gi.TYPE_TAG_ERROR:
		out.WriteString("GError*")
	case gi.TYPE_TAG_INTERFACE:
		if ti.IsPointer() {
			flags |= type_pointer
		}
		out.WriteString(c_type_for_interface(ti.Interface(), flags))
	default:
		out.WriteString(c_type_for_tag(tag, flags))
		if ti.IsPointer() {
			out.WriteString("*")
		}
	}

	return out.String()
}

func c_type_for_interface(bi *gi.BaseInfo, flags type_flags) string {
	var out bytes.Buffer

	ns := bi.Namespace()
	nm := bi.Name()
	fullnm := strings.ToLower(ns) + "." + nm
	out.WriteString(gi.DefaultRepository().CPrefix(ns))
	out.WriteString(bi.Name())

	if flags&type_pointer != 0 && !config.is_disguised(fullnm) {
		out.WriteString("*")
	}

	return out.String()
}

//------------------------------------------------------------------
// Go Type
//------------------------------------------------------------------

func go_type_for_interface(bi *gi.BaseInfo, flags type_flags) string {
	var out bytes.Buffer
	printf := printer_to(&out)
	ns := bi.Namespace()
	fullnm := strings.ToLower(ns) + "." + bi.Name()

	if flags&type_list_member != 0 {
		switch bi.Type() {
		case gi.INFO_TYPE_OBJECT, gi.INFO_TYPE_INTERFACE:
			return go_type_for_interface(bi, type_pointer|type_return)
		default:
			return go_type_for_interface(bi, type_return)
		}
	}

	switch t := bi.Type(); t {
	case gi.INFO_TYPE_OBJECT, gi.INFO_TYPE_INTERFACE:
		if flags&type_exact != 0 {
			// exact type for object/interface is always an unsafe.Pointer
			printf("unsafe.Pointer")
			break
		}

		if flags&(type_return|type_receiver) != 0 && flags&type_pointer != 0 {
			// receivers and return values are actual types,
			// and a pointer most likely
			printf("*")
		}
		if ns != config.namespace {
			// prepend foreign types with appropriate namespace
			printf("%s.", strings.ToLower(ns))
		}
		printf(bi.Name())
		if flags&(type_return|type_receiver) == 0 {
			// ordinary function arguments are substituted by their *Like
			// counterparts
			printf("Like")
		}
		if flags&type_receiver != 0 && t == gi.INFO_TYPE_INTERFACE {
			// special case for interfaces, we use *Impl structures
			// as receivers
			printf("Impl")
		}
	case gi.INFO_TYPE_CALLBACK:
		if flags&type_exact != 0 {
			printf("unsafe.Pointer")
			break
		}
		goto handle_default
	case gi.INFO_TYPE_STRUCT:
		if ns == "cairo" {
			printf(cairo_go_type_for_interface(bi, flags))
			break
		}
		goto handle_default
	default:
		goto handle_default
	}
	return out.String()
handle_default:
	if flags&type_pointer != 0 && !config.is_disguised(fullnm) {
		printf("*")
	}
	if ns != config.namespace {
		printf("%s.", strings.ToLower(ns))
	}
	printf(bi.Name())
	return out.String()
}

func go_type(ti *gi.TypeInfo, flags type_flags) string {
	var out bytes.Buffer

	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			out.WriteString("unsafe.Pointer")
			break
		}
		panic("Non-pointer void type is not supported")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME:
		if flags&type_exact != 0 {
			out.WriteString("unsafe.Pointer")
		} else {
			out.WriteString("string")
		}
	case gi.TYPE_TAG_ARRAY:
		size := ti.ArrayFixedSize()
		if size != -1 {
			fmt.Fprintf(&out, "[%d]", size)
		} else {
			if flags&type_exact != 0 {
				out.WriteString("unsafe.Pointer")
			} else {
				out.WriteString("[]")
			}
		}
		out.WriteString(go_type(ti.ParamType(0), flags))
	case gi.TYPE_TAG_GLIST:
		if flags&type_exact != 0 {
			out.WriteString("unsafe.Pointer")
		} else {
			out.WriteString("[]")
			out.WriteString(go_type(ti.ParamType(0), flags|type_list_member))
		}
	case gi.TYPE_TAG_GSLIST:
		if flags&type_exact != 0 {
			out.WriteString("unsafe.Pointer")
		} else {
			out.WriteString("[]")
			out.WriteString(go_type(ti.ParamType(0), flags|type_list_member))
		}
	case gi.TYPE_TAG_GHASH:
		if flags&type_exact != 0 {
			out.WriteString("unsafe.Pointer")
		} else {
			out.WriteString("map[")
			out.WriteString(go_type(ti.ParamType(0), flags))
			out.WriteString("]")
			out.WriteString(go_type(ti.ParamType(1), flags))
		}
	case gi.TYPE_TAG_ERROR:
		// not used?
		out.WriteString("error")
	case gi.TYPE_TAG_INTERFACE:
		if ti.IsPointer() {
			flags |= type_pointer
		}
		out.WriteString(go_type_for_interface(ti.Interface(), flags))
	default:
		if ti.IsPointer() {
			flags |= type_pointer
		}
		out.WriteString(go_type_for_tag(tag, flags))
	}

	return out.String()
}

func go_type_for_tag(tag gi.TypeTag, flags type_flags) string {
	var out bytes.Buffer
	p := printer_to(&out)

	if flags&type_pointer != 0 {
		p("*")
	}

	if flags&type_exact != 0 {
		switch tag {
		case gi.TYPE_TAG_BOOLEAN:
			p("int32") // sadly
		case gi.TYPE_TAG_INT8:
			p("int8")
		case gi.TYPE_TAG_UINT8:
			p("uint8")
		case gi.TYPE_TAG_INT16:
			p("int16")
		case gi.TYPE_TAG_UINT16:
			p("uint16")
		case gi.TYPE_TAG_INT32:
			p("int32")
		case gi.TYPE_TAG_UINT32:
			p("uint32")
		case gi.TYPE_TAG_INT64:
			p("int64")
		case gi.TYPE_TAG_UINT64:
			p("uint64")
		case gi.TYPE_TAG_FLOAT:
			p("float32")
		case gi.TYPE_TAG_DOUBLE:
			p("float64")
		case gi.TYPE_TAG_GTYPE:
			if config.namespace != "GObject" {
				p("gobject.Type")
			} else {
				p("Type")
			}
		case gi.TYPE_TAG_UNICHAR:
			p("rune")
		default:
			panic("unreachable")
		}
	} else {
		switch tag {
		case gi.TYPE_TAG_BOOLEAN:
			p("bool")
		case gi.TYPE_TAG_INT8:
			p("int8")
		case gi.TYPE_TAG_UINT8:
			p("uint8")
		case gi.TYPE_TAG_INT16:
			p("int16")
		case gi.TYPE_TAG_UINT16:
			p("uint16")
		case gi.TYPE_TAG_INT32:
			p("int32")
		case gi.TYPE_TAG_UINT32:
			p("uint32")
		case gi.TYPE_TAG_INT64:
			p("int64")
		case gi.TYPE_TAG_UINT64:
			p("uint64")
		case gi.TYPE_TAG_FLOAT:
			p("float32")
		case gi.TYPE_TAG_DOUBLE:
			p("float64")
		case gi.TYPE_TAG_GTYPE:
			if config.namespace != "GObject" {
				p("gobject.Type")
			} else {
				p("Type")
			}
		case gi.TYPE_TAG_UNICHAR:
			p("rune")
		default:
			panic("unreachable")
		}
	}

	return out.String()
}

//------------------------------------------------------------------
// Simple Cgo Type (for exported functions)
//------------------------------------------------------------------

func simple_cgo_type(ti *gi.TypeInfo, flags type_flags) string {
	tag := ti.Tag()
	switch tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			return "unsafe.Pointer"
		}
		panic("Non-pointer void type is not supported")
	case gi.TYPE_TAG_INTERFACE:
		bi := ti.Interface()
		switch bi.Type() {
		case gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS:
			ei := gi.ToEnumInfo(bi)
			return go_type_for_tag(ei.StorageType(), flags|type_exact)
		case gi.INFO_TYPE_STRUCT:
			ns := bi.Namespace()
			nm := bi.Name()
			fullnm := strings.ToLower(ns) + "." + nm
			if config.is_disguised(fullnm) {
				return "unsafe.Pointer"
			}
		}
	}
	if !strings.HasPrefix(cgo_type(ti, flags), "*") {
		return go_type_for_tag(tag, flags|type_exact)
	}
	return "unsafe.Pointer"
}

//------------------------------------------------------------------
// Type sizes
//------------------------------------------------------------------

func type_size_for_interface(bi *gi.BaseInfo, flags type_flags) int {
	ptrsize := int(unsafe.Sizeof(unsafe.Pointer(nil)))
	if flags&type_pointer != 0 {
		return ptrsize
	}

	switch t := bi.Type(); t {
	case gi.INFO_TYPE_OBJECT, gi.INFO_TYPE_INTERFACE:
		return ptrsize
	case gi.INFO_TYPE_STRUCT:
		si := gi.ToStructInfo(bi)
		return si.Size()
	case gi.INFO_TYPE_UNION:
		ui := gi.ToUnionInfo(bi)
		return ui.Size()
	case gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS:
		ei := gi.ToEnumInfo(bi)
		return type_size_for_tag(ei.StorageType(), flags)
	case gi.INFO_TYPE_CALLBACK:
		return ptrsize
	}
	panic("unreachable: " + bi.Type().String())
}

// returns the size of a type, works only for TypeExact
func type_size(ti *gi.TypeInfo, flags type_flags) int {
	ptrsize := int(unsafe.Sizeof(unsafe.Pointer(nil)))
	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			return ptrsize
		}
		panic("Non-pointer void type is not supported")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME, gi.TYPE_TAG_GLIST,
		gi.TYPE_TAG_GSLIST, gi.TYPE_TAG_GHASH:
		return ptrsize
	case gi.TYPE_TAG_ARRAY:
		size := ti.ArrayFixedSize()
		if size != -1 {
			return size * type_size(ti.ParamType(0), flags)
		}
		return ptrsize
	case gi.TYPE_TAG_INTERFACE:
		if ti.IsPointer() {
			flags |= type_pointer
		}
		return type_size_for_interface(ti.Interface(), flags)
	default:
		if ti.IsPointer() {
			flags |= type_pointer
		}
		return type_size_for_tag(tag, flags)
	}
	panic("unreachable: " + ti.Tag().String())
}

func type_size_for_tag(tag gi.TypeTag, flags type_flags) int {
	ptrsize := int(unsafe.Sizeof(unsafe.Pointer(nil)))
	if flags&type_pointer != 0 {
		return ptrsize
	}

	switch tag {
	case gi.TYPE_TAG_BOOLEAN:
		return 4
	case gi.TYPE_TAG_INT8:
		return 1
	case gi.TYPE_TAG_UINT8:
		return 1
	case gi.TYPE_TAG_INT16:
		return 2
	case gi.TYPE_TAG_UINT16:
		return 2
	case gi.TYPE_TAG_INT32:
		return 4
	case gi.TYPE_TAG_UINT32:
		return 4
	case gi.TYPE_TAG_INT64:
		return 8
	case gi.TYPE_TAG_UINT64:
		return 8
	case gi.TYPE_TAG_FLOAT:
		return 4
	case gi.TYPE_TAG_DOUBLE:
		return 8
	case gi.TYPE_TAG_GTYPE:
		return ptrsize
	case gi.TYPE_TAG_UNICHAR:
		return 4
	}
	panic("unreachable: " + tag.String())
}

//------------------------------------------------------------------
// Type needs wrapper?
//------------------------------------------------------------------

func type_needs_wrapper(ti *gi.TypeInfo) bool {
	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			return false
		}
		panic("Non-pointer void type is not supported")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME, gi.TYPE_TAG_GLIST,
		gi.TYPE_TAG_GSLIST, gi.TYPE_TAG_GHASH:
		return true
	case gi.TYPE_TAG_ARRAY:
		size := ti.ArrayFixedSize()
		if size != -1 {
			return type_needs_wrapper(ti.ParamType(0))
		}
		return true
	case gi.TYPE_TAG_ERROR:
		panic("not implemented")
	case gi.TYPE_TAG_INTERFACE:
		switch ti.Interface().Type() {
		case gi.INFO_TYPE_CALLBACK, gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS,
			gi.INFO_TYPE_STRUCT, gi.INFO_TYPE_UNION:
			return false
		}
		return true
	}
	return false
}
