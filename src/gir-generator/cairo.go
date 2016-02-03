/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import (
	"bytes"
	"gi"
)

func cairo_go_type_for_interface(bi *gi.BaseInfo, flags type_flags) string {
	var out bytes.Buffer
	p := printer_to(&out)
	name := bi.Name()

	switch name {
	case "Surface", "Pattern":
		if flags&type_return == 0 {
			p("cairo.%sLike", name)
			break
		}
		fallthrough
	default:
		if flags&type_pointer != 0 {
			p("*")
		}
		p("cairo.%s", name)
	}

	return out.String()
}

func cairo_go_to_cgo_for_interface(bi *gi.BaseInfo, arg0, arg1 string, flags conv_flags) string {
	var out bytes.Buffer
	p := printer_to(&out)
	name := bi.Name()

	switch name {
	case "Surface", "Pattern":
		p("if %s != nil {\n", arg0)
		p("\t%s = (*C.cairo%s)(%s.InheritedFromCairo%s().C)\n", arg1, name, arg0, name)
		p("}")
	case "RectangleInt", "Rectangle", "TextCluster", "Matrix":
		ctype := cgo_type_for_interface(bi, type_none)
		if flags&conv_pointer != 0 {
			p("%s = (*%s)(unsafe.Pointer(%s))",
				arg1, ctype, arg0)
		} else {
			p("%s = *(*%s)(unsafe.Pointer(&%s))",
				arg1, ctype, arg0)
		}
	default:
		p("if %s != nil {\n", arg0)
		p("\t%s = (*C.cairo%s)(%s.C)\n", arg1, name, arg0)
		p("}")
	}

	return out.String()
}

func cairo_cgo_to_go_for_interface(bi *gi.BaseInfo, arg1, arg2 string, flags conv_flags) string {
	var out bytes.Buffer
	p := printer_to(&out)
	name := bi.Name()

	switch name {
	case "Path", "FontOptions":
		if flags&conv_pointer == 0 {
			panic("unexpected non-pointer type")
		}

		p("%s = (*cairo.%s)(cairo.%sWrap(unsafe.Pointer(%s)))", arg2, name, name, arg1)
	case "Surface", "Region", "Pattern", "Context", "FontFace", "ScaledFont":
		if flags&conv_pointer == 0 {
			panic("unexpected non-pointer type")
		}

		grab := "true"
		if flags&conv_own_everything != 0 {
			grab = "false"
		}
		p("%s = (*cairo.%s)(cairo.%sWrap(unsafe.Pointer(%s), %s))", arg2, name, name, arg1, grab)
	default:
		gotype := go_type_for_interface(bi, type_return)
		if flags&conv_pointer != 0 {
			p("%s = (*%s)(unsafe.Pointer(%s))",
				arg2, gotype, arg1)
		} else {
			p("%s = *(*%s)(unsafe.Pointer(&%s))",
				arg2, gotype, arg1)
		}
	}

	return out.String()
}
