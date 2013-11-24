package main

import (
	"bytes"
	"fmt"
	"github.com/snyh/gir-generator/gi"
	"strings"
)

func cgo_array_to_go_array(elem *gi.TypeInfo, name string) string {
	return fmt.Sprintf("(*(*[999999]%s)(unsafe.Pointer(%s)))",
		cgo_type(elem, type_none), name)
}

type conv_flags int

const (
	conv_none    conv_flags = 0
	conv_pointer conv_flags = 1 << iota
	conv_list_member
	conv_own_none
	conv_own_container
	conv_own_everything
)

func ownership_to_conv_flags(t gi.Transfer) conv_flags {
	switch t {
	case gi.TRANSFER_NOTHING:
		return conv_own_none
	case gi.TRANSFER_CONTAINER:
		return conv_own_container
	case gi.TRANSFER_EVERYTHING:
		return conv_own_everything
	}
	return 0
}

//------------------------------------------------------------------
// Go to Cgo Converter
//------------------------------------------------------------------

func go_to_cgo_for_interface(bi *gi.BaseInfo, arg0, arg1 string, flags conv_flags) string {
	var out bytes.Buffer
	printf := printer_to(&out)

	switch bi.Type() {
	case gi.INFO_TYPE_OBJECT:
		prefix := gi.DefaultRepository().CPrefix(bi.Namespace())
		printf("if %s != nil {\n", arg0)
		printf("\t%s = %s.InheritedFrom%s%s()\n",
			arg1, arg0, prefix, bi.Name())
		printf("}")
	case gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS:
		ctype := cgo_type_for_interface(bi, type_none)
		printf("%s = %s(%s)", arg1, ctype, arg0)
	case gi.INFO_TYPE_INTERFACE:
		prefix := gi.DefaultRepository().CPrefix(bi.Namespace())
		printf("if %s != nil {\n", arg0)
		printf("\t%s = %s.Implements%s%s()\n",
			arg1, arg0, prefix, bi.Name())
		printf("}")
	case gi.INFO_TYPE_STRUCT:
		ns := bi.Namespace()
		if ns == "cairo" {
			printf(cairo_go_to_cgo_for_interface(bi, arg0, arg1, flags))
			break
		}

		fullnm := strings.ToLower(ns) + "." + bi.Name()
		if config.is_disguised(fullnm) {
			flags &^= conv_pointer
		}
		ctype := cgo_type_for_interface(bi, type_none)
		if flags&conv_pointer != 0 {
			printf("%s = (*%s)(unsafe.Pointer(%s))",
				arg1, ctype, arg0)
		} else {
			printf("%s = *(*%s)(unsafe.Pointer(&%s))",
				arg1, ctype, arg0)
		}
	case gi.INFO_TYPE_CALLBACK:
		printf("if %s != nil {\n", arg0)
		printf("\t%s = unsafe.Pointer(&%s)", arg1, arg0)
		printf("}")
	}

	return out.String()
}

func go_to_cgo(ti *gi.TypeInfo, arg0, arg1 string, flags conv_flags) string {
	var out bytes.Buffer
	printf := printer_to(&out)

	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			printf("%s = unsafe.Pointer(%s)", arg1, arg0)
			break
		}
		printf("<ERROR: void>")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME:
		printf("%s = _GoStringToGString(%s)", arg1, arg0)
		if flags&conv_own_everything == 0 {
			printf("\ndefer C.free(unsafe.Pointer(%s))", arg1)
		}
	case gi.TYPE_TAG_ARRAY:
		switch ti.ArrayType() {
		case gi.ARRAY_TYPE_C:
			var nelem string
			if ti.IsZeroTerminated() {
				nelem = fmt.Sprintf("(len(%s) + 1)", arg0)
			} else {
				nelem = fmt.Sprintf("len(%s)", arg0)
			}

			// alloc memory
			printf("%s = (%s)(C.malloc(C.size_t(int(unsafe.Sizeof(*%s)) * %s)))\n",
				arg1, cgo_type(ti, type_none), arg1, nelem)
			printf("defer C.free(unsafe.Pointer(%s))\n", arg1)

			// convert elements
			printf("for i, e := range %s {\n", arg0)
			array := cgo_array_to_go_array(ti.ParamType(0), arg1)
			conv := go_to_cgo(ti.ParamType(0), "e", array+"[i]", flags)
			printf(print_lines_with_indent(conv))
			printf("}")

			// write a trailing zero if necessary (TODO: buggy)
			if ti.IsZeroTerminated() {
				printf("\n%s[len(%s)] = nil", array, arg0)
			}

		}
	case gi.TYPE_TAG_GLIST:
	case gi.TYPE_TAG_GSLIST:
	case gi.TYPE_TAG_GHASH:
	case gi.TYPE_TAG_ERROR:
	case gi.TYPE_TAG_INTERFACE:
		if ti.IsPointer() {
			flags |= conv_pointer
		}
		printf(go_to_cgo_for_interface(ti.Interface(), arg0, arg1, flags))
	default:
		if ti.IsPointer() {
			flags |= conv_pointer
		}
		printf(go_to_cgo_for_tag(tag, arg0, arg1, flags))
	}

	return out.String()
}

func go_to_cgo_for_tag(tag gi.TypeTag, arg0, arg1 string, flags conv_flags) string {
	switch tag {
	case gi.TYPE_TAG_BOOLEAN:
		return fmt.Sprintf("%s = _GoBoolToCBool(%s)", arg1, arg0)
	default:
		if flags&conv_pointer == 0 {
			return fmt.Sprintf("%s = %s(%s)", arg1,
				cgo_type_for_tag(tag, type_none), arg0)
		} else {
			return fmt.Sprintf("%s = (%s)(unsafe.Pointer(%s))", arg1,
				cgo_type_for_tag(tag, type_pointer), arg0)
		}
	}

	panic("unreachable")
	return ""
}

//------------------------------------------------------------------
// Cgo to Go Converter
//------------------------------------------------------------------

func cgo_to_go_for_interface(bi *gi.BaseInfo, arg1, arg2 string, flags conv_flags) string {
	var out bytes.Buffer
	printf := printer_to(&out)

	switch bi.Type() {
	case gi.INFO_TYPE_OBJECT, gi.INFO_TYPE_INTERFACE:
		gotype := go_type_for_interface(bi, type_return)
		if flags&conv_own_everything != 0 {
			printf("%s = (*%s)(%sObjectWrap(unsafe.Pointer(%s), false))",
				arg2, gotype, config.gns, arg1)
		} else {
			printf("%s = (*%s)(%sObjectWrap(unsafe.Pointer(%s), true))",
				arg2, gotype, config.gns, arg1)
		}
	case gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS:
		gotype := go_type_for_interface(bi, type_return)
		printf("%s = %s(%s)", arg2, gotype, arg1)
	case gi.INFO_TYPE_STRUCT, gi.INFO_TYPE_UNION:
		ns := bi.Namespace()
		if ns == "cairo" {
			printf(cairo_cgo_to_go_for_interface(bi, arg1, arg2, flags))
			break
		}

		fullnm := strings.ToLower(ns) + "." + bi.Name()
		gotype := go_type_for_interface(bi, type_return)

		if flags&conv_list_member != 0 {
			printf("%s = *(*%s)(unsafe.Pointer(%s))",
				arg2, gotype, arg1)
			break
		}

		if config.is_disguised(fullnm) {
			printf("%s = %s{unsafe.Pointer(%s)}",
				arg2, gotype, arg1)
			break
		}

		if flags&conv_pointer != 0 {
			printf("%s = (*%s)(unsafe.Pointer(%s))",
				arg2, gotype, arg1)
		} else {
			printf("%s = *(*%s)(unsafe.Pointer(&%s))",
				arg2, gotype, arg1)
		}
	}
	return out.String()
}

func cgo_to_go(ti *gi.TypeInfo, arg1, arg2 string, flags conv_flags) string {
	var out bytes.Buffer
	printf := printer_to(&out)

	switch tag := ti.Tag(); tag {
	case gi.TYPE_TAG_VOID:
		if ti.IsPointer() {
			printf("%s = %s", arg2, arg1)
			break
		}
		printf("<ERROR: void>")
	case gi.TYPE_TAG_UTF8, gi.TYPE_TAG_FILENAME:
		printf("%s = C.GoString(%s)", arg2, arg1)
		if flags&conv_own_everything != 0 {
			printf("\nC.g_free(unsafe.Pointer(%s))", arg1)
		}
	case gi.TYPE_TAG_ARRAY:
		switch ti.ArrayType() {
		case gi.ARRAY_TYPE_C:
			// array was allocated already at this point
			printf("for i := range %s {\n", arg2)
			array := cgo_array_to_go_array(ti.ParamType(0), arg1)
			conv := cgo_to_go(ti.ParamType(0),
				array+"[i]", arg2+"[i]", flags)
			printf(print_lines_with_indent(conv))
			printf("}")
		}
	case gi.TYPE_TAG_GLIST:
		ptype := ti.ParamType(0)
		printf("for iter := (*_GList)(unsafe.Pointer(%s)); iter != nil; iter = iter.next {\n", arg1)
		elt := fmt.Sprintf("(%s)(iter.data)",
			force_pointer(cgo_type(ptype, type_return|type_list_member)))
		printf("\tvar elt %s\n", go_type(ptype, type_return|type_list_member))
		conv := cgo_to_go(ptype, elt, "elt", flags|conv_list_member)
		printf(print_lines_with_indent(conv))
		printf("\t%s = append(%s, elt)\n", arg2, arg2)
		printf("}")
	case gi.TYPE_TAG_GSLIST:
		ptype := ti.ParamType(0)
		printf("for iter := (*_GSList)(unsafe.Pointer(%s)); iter != nil; iter = iter.next {\n", arg1)
		elt := fmt.Sprintf("(%s)(iter.data)",
			force_pointer(cgo_type(ptype, type_return|type_list_member)))
		printf("\tvar elt %s\n", go_type(ptype, type_return|type_list_member))
		conv := cgo_to_go(ptype, elt, "elt", flags|conv_list_member)
		printf(print_lines_with_indent(conv))
		printf("\t%s = append(%s, elt)\n", arg2, arg2)
		printf("}")
	case gi.TYPE_TAG_GHASH:
	case gi.TYPE_TAG_ERROR:
		printf("if %s != nil {\n", arg1)
		printf("\t%s = errors.New(C.GoString(((*_GError)(unsafe.Pointer(%s))).message))\n", arg2, arg1)
		printf("\tC.g_error_free(%s)\n", arg1)
		printf("}\n")
	case gi.TYPE_TAG_INTERFACE:
		if ti.IsPointer() {
			flags |= conv_pointer
		}
		printf(cgo_to_go_for_interface(ti.Interface(), arg1, arg2, flags))
	default:
		if ti.IsPointer() {
			flags |= conv_pointer
		}
		printf(cgo_to_go_for_tag(tag, arg1, arg2, flags))
	}

	return out.String()
}

func cgo_to_go_for_tag(tag gi.TypeTag, arg1, arg2 string, flags conv_flags) string {
	switch tag {
	case gi.TYPE_TAG_BOOLEAN:
		return fmt.Sprintf("%s = %s != 0", arg2, arg1)
	case gi.TYPE_TAG_INT8:
		return fmt.Sprintf("%s = int(%s)", arg2, arg1)
	case gi.TYPE_TAG_UINT8:
		return fmt.Sprintf("%s = int(%s)", arg2, arg1)
	case gi.TYPE_TAG_INT16:
		return fmt.Sprintf("%s = int(%s)", arg2, arg1)
	case gi.TYPE_TAG_UINT16:
		return fmt.Sprintf("%s = int(%s)", arg2, arg1)
	case gi.TYPE_TAG_INT32:
		return fmt.Sprintf("%s = int(%s)", arg2, arg1)
	case gi.TYPE_TAG_UINT32:
		return fmt.Sprintf("%s = int(%s)", arg2, arg1)
	case gi.TYPE_TAG_INT64:
		return fmt.Sprintf("%s = int64(%s)", arg2, arg1)
	case gi.TYPE_TAG_UINT64:
		return fmt.Sprintf("%s = uint64(%s)", arg2, arg1)
	case gi.TYPE_TAG_FLOAT:
		return fmt.Sprintf("%s = float64(%s)", arg2, arg1)
	case gi.TYPE_TAG_DOUBLE:
		return fmt.Sprintf("%s = float64(%s)", arg2, arg1)
	case gi.TYPE_TAG_GTYPE:
		if config.namespace != "GObject" {
			return fmt.Sprintf("%s = gobject.Type(%s)", arg2, arg1)
		}
		return fmt.Sprintf("%s = Type(%s)", arg2, arg1)
	case gi.TYPE_TAG_UNICHAR:
		return fmt.Sprintf("%s = rune(%s)", arg2, arg1)
	}

	panic("unreachable")
	return ""
}

//------------------------------------------------------------------
// Simple Cgo to Go Converter
//------------------------------------------------------------------

func simple_cgo_to_go(ti *gi.TypeInfo, arg0, arg1 string, flags conv_flags) string {
	cgotype := cgo_type(ti, type_none)
	arg0 = fmt.Sprintf("(%s)(%s)", cgotype, arg0)
	return cgo_to_go(ti, arg0, arg1, flags)
}
