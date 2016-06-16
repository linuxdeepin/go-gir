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
	"bufio"
	"bytes"
	"fmt"
	"gi"
	"os"
	"strings"
)

type binding_generator struct {
	file_go *os.File
	file_c  *os.File
	file_h  *os.File
	out_go  *bufio.Writer
	out_c   *bufio.Writer
	out_h   *bufio.Writer

	// temporary buffer for 'go_bindings'
	go_bindings bytes.Buffer

	// map for C header generator to avoid duplicates
	declared_c_funcs map[string]bool
}

func new_binding_generator(out_base string) *binding_generator {
	var err error

	this := new(binding_generator)
	this.declared_c_funcs = make(map[string]bool)
	this.file_go, err = os.Create(out_base + ".go")
	panic_if_error(err)
	this.file_c, err = os.Create(out_base + ".gen.c")
	panic_if_error(err)
	this.file_h, err = os.Create(out_base + ".gen.h")
	panic_if_error(err)
	this.out_go = bufio.NewWriter(this.file_go)
	this.out_c = bufio.NewWriter(this.file_c)
	this.out_h = bufio.NewWriter(this.file_h)

	return this
}

func (this *binding_generator) release() {
	this.out_go.Flush()
	this.out_c.Flush()
	this.out_h.Flush()
	this.file_go.Close()
	this.file_c.Close()
	this.file_h.Close()
}

func (this *binding_generator) generate(go_template string) {
	// this will fill the 'go_bindings' buffer
	repo := gi.DefaultRepository()
	for i, n := 0, repo.NumInfo(config.namespace); i < n; i++ {
		this.process_base_info(repo.Info(config.namespace, i))
	}

	t := must_template(go_template)
	t.Execute(this.out_go, map[string]interface{}{
		"g_list_funcs":       g_list_funcs,
		"g_object_ref_unref": g_object_ref_unref,
		"go_utils":           go_utils(true),
		"go_utils_no_cb":     go_utils(false),
		"go_bindings":        this.go_bindings.String(),
		"g_error_free":       g_error_free,
		"g_free":             g_free,
	})

	// write source/header preambles
	p := printer_to(this.out_h)
	p(c_header)

	// TODO: using config.pkg here is probably incorrect, we should use the
	// filename
	c_template.Execute(this.out_c, map[string]interface{}{
		"namespace": config.namespace,
		"package":   config.pkg,
	})
	p = printer_to(this.out_c)
	p("\n\n")

	// this will write the rest of .c/.h files
	this.c_forward_declarations()
}

func (this *binding_generator) c_func_forward_declaration(fi *gi.FunctionInfo) {
	symbol := fi.Symbol()
	if _, declared := this.declared_c_funcs[symbol]; declared {
		return
	}

	this.declared_c_funcs[symbol] = true
	container := fi.Container()
	ph := printer_to(this.out_h)
	pc := printer_to(this.out_c)

	flags := fi.Flags()
	narg := fi.NumArg()
	if flags&gi.FUNCTION_THROWS != 0 {
		narg++
	}

	ph("extern %s %s(", c_type(fi.ReturnType(), type_none), symbol)
	if flags&gi.FUNCTION_IS_METHOD != 0 {
		ph("%s", c_type_for_interface(container, type_pointer))
		if narg > 0 {
			ph(", ")
		}
	}

	var closure_scope gi.ScopeType = gi.SCOPE_TYPE_INVALID
	var closure_userdata int = -1
	var closure_destroy int = -1
	var closure_callback int = -1

	for i, n := 0, narg; i < n; i++ {
		if i != 0 {
			ph(", ")
		}

		if i == n-1 && flags&gi.FUNCTION_THROWS != 0 {
			ph("GError**")
			continue
		}

		arg := fi.Arg(i)
		t := arg.Type()

		if uarg := arg.Closure(); uarg != -1 {
			if t.Tag() == gi.TYPE_TAG_INTERFACE {
				if t.Interface().Type() == gi.INFO_TYPE_CALLBACK {
					if darg := arg.Destroy(); darg != -1 {
						closure_destroy = darg
					}
					closure_userdata = uarg
					closure_callback = i
					closure_scope = arg.Scope()
				}
			}
		}

		ph("%s", c_type(t, type_none))

		switch arg.Direction() {
		case gi.DIRECTION_INOUT, gi.DIRECTION_OUT:
			ph("*")
		}
	}
	ph(");\n")

	if closure_scope == gi.SCOPE_TYPE_INVALID {
		return
	}

	// Also if this function is actually blacklisted, we don't want to
	// generate wrappers
	if config.is_object_blacklisted(gi.ToBaseInfo(fi)) {
		return
	}

	// Or maybe this is a method and method's owner is blacklisted
	if container != nil && config.is_object_blacklisted(container) {
		return
	}

	// in case if the function takes callback, generate appropriate wrappers
	// for Go (gc compiler mainly)
	var tmp bytes.Buffer
	p := printer_to(&tmp)
	p("%s _%s(", c_type(fi.ReturnType(), type_none), symbol)
	if flags&gi.FUNCTION_IS_METHOD != 0 {
		p("%s this", c_type_for_interface(container, type_pointer))
		if narg > 0 {
			p(", ")
		}
	}

	for i, n := 0, narg; i < n; i++ {
		if i == closure_userdata || i == closure_destroy {
			// skip userdata and destroy func in the wrapper
			continue
		}

		if i != 0 {
			p(", ")
		}

		if i == closure_callback {
			// replace callback argument with Go function pointer
			// and optional unique id (if scope is not CALL)
			p("void* gofunc")
			continue
		}

		if i == n-1 && flags&gi.FUNCTION_THROWS != 0 {
			p("GError** arg%d", i)
			continue
		}

		arg := fi.Arg(i)
		p("%s", c_type(arg.Type(), type_none))

		switch arg.Direction() {
		case gi.DIRECTION_INOUT, gi.DIRECTION_OUT:
			p("*")
		}
		p(" arg%d", i)
	}
	p(")")
	ph("extern %s;\n", tmp.String())
	pc("%s {\n", tmp.String())

	// body
	defcall := func(argprint func(i int, t *gi.TypeInfo)) {
		pc("\t\t")
		if ret := fi.ReturnType(); !(ret.Tag() == gi.TYPE_TAG_VOID && !ret.IsPointer()) {
			pc("return ")
		}
		pc("%s(", symbol)
		if flags&gi.FUNCTION_IS_METHOD != 0 {
			pc("this")
			if narg > 0 {
				pc(", ")
			}
		}
		for i, n := 0, narg; i < n; i++ {
			arg := fi.Arg(i)
			t := arg.Type()

			if i != 0 {
				pc(", ")
			}

			argprint(i, t)
		}
		pc(");\n")
	}
	pc("\tif (gofunc) {\n")
	defcall(func(i int, t *gi.TypeInfo) {
		switch i {
		case closure_userdata:
			pc("gofunc")
		case closure_destroy:
			pc("_c_callback_cleanup")
		case closure_callback:
			pc("_%s_c_wrapper", c_type(t, type_none))
			if closure_scope == gi.SCOPE_TYPE_ASYNC {
				pc("_once")
			}
		default:
			pc("arg%d", i)
		}
	})
	pc("\t} else {\n")
	defcall(func(i int, t *gi.TypeInfo) {
		switch i {
		case closure_userdata, closure_destroy, closure_callback:
			pc("0")
		default:
			pc("arg%d", i)
		}
	})
	pc("\t}\n")
	pc("}\n")
}

// generating forward C declarations properly:
// 10 - various typedefs
// 20 - functions and methods
// 30 - struct definitions

func (this *binding_generator) c_forward_declaration10(bi *gi.BaseInfo) {
	p := printer_to(this.out_h)
	switch bi.Type() {
	case gi.INFO_TYPE_OBJECT:
		cctype := c_type_for_interface(bi, type_none)
		p("typedef struct _%s %s;\n", cctype, cctype)
	case gi.INFO_TYPE_FLAGS, gi.INFO_TYPE_ENUM:
		ei := gi.ToEnumInfo(bi)
		p("typedef %s %s;\n", c_type_for_tag(ei.StorageType(), type_none),
			c_type_for_interface(bi, type_none))
	case gi.INFO_TYPE_STRUCT, gi.INFO_TYPE_INTERFACE, gi.INFO_TYPE_UNION:
		fullnm := strings.ToLower(bi.Namespace()) + "." + bi.Name()
		cctype := c_type_for_interface(bi, type_none)
		if config.is_disguised(fullnm) {
			p("typedef void *%s;\n", cctype)
		} else {
			p("typedef struct _%s %s;\n", cctype, cctype)
		}
	case gi.INFO_TYPE_CALLBACK:
		pc := printer_to(this.out_c)
		ctype := c_type_for_interface(bi, type_none)

		// type doesn't matter here, it's just a pointer after all
		p("typedef void* %s;\n", ctype)

		// and wrapper declarations for .c file only (cgo has problems
		// with that)
		pc("extern void _%s_c_wrapper();\n", ctype)
		pc("extern void _%s_c_wrapper_once();\n", ctype)
	}
}

func (this *binding_generator) c_forward_declaration20(bi *gi.BaseInfo) {
	p := printer_to(this.out_h)
	switch bi.Type() {
	case gi.INFO_TYPE_FUNCTION:
		fi := gi.ToFunctionInfo(bi)
		this.c_func_forward_declaration(fi)
	case gi.INFO_TYPE_OBJECT:
		oi := gi.ToObjectInfo(bi)
		for i, n := 0, oi.NumMethod(); i < n; i++ {
			meth := oi.Method(i)
			this.c_func_forward_declaration(meth)
		}
		p("extern GType %s();\n", oi.TypeInit())
	case gi.INFO_TYPE_INTERFACE:
		ii := gi.ToInterfaceInfo(bi)
		for i, n := 0, ii.NumMethod(); i < n; i++ {
			meth := ii.Method(i)
			this.c_func_forward_declaration(meth)
		}
		p("extern GType %s();\n", ii.TypeInit())
	case gi.INFO_TYPE_STRUCT:
		si := gi.ToStructInfo(bi)
		for i, n := 0, si.NumMethod(); i < n; i++ {
			meth := si.Method(i)
			this.c_func_forward_declaration(meth)
		}
	case gi.INFO_TYPE_UNION:
		ui := gi.ToUnionInfo(bi)
		for i, n := 0, ui.NumMethod(); i < n; i++ {
			meth := ui.Method(i)
			this.c_func_forward_declaration(meth)
		}
	}
}

func (this *binding_generator) c_forward_declaration30(bi *gi.BaseInfo) {
	p := printer_to(this.out_h)
	fullnm := strings.ToLower(bi.Namespace()) + "." + bi.Name()

	switch bi.Type() {
	case gi.INFO_TYPE_STRUCT:
		si := gi.ToStructInfo(bi)
		size := si.Size()
		if config.is_disguised(fullnm) {
			return
		}
		cctype := c_type_for_interface(bi, type_none)
		if size == 0 {
			p("struct _%s {};\n", cctype)
		} else {
			p("struct _%s { uint8_t _data[%d]; };\n", cctype, size)
		}
	case gi.INFO_TYPE_UNION:
		ui := gi.ToUnionInfo(bi)
		size := ui.Size()
		cctype := c_type_for_interface(bi, type_none)
		if size == 0 {
			p("struct _%s {};\n", cctype)
		} else {
			p("struct _%s { uint8_t _data[%d]; };\n", cctype, size)
		}
	}
}

func get_dependcies_declarations(namespace string) (ret []string) {
	repo := gi.DefaultRepository()
	deps := repo.Dependencies(namespace)
	for _, dep := range deps {
		depv := strings.Split(dep, "-")
		deps = append(deps, get_dependcies_declarations(depv[0])...)
	}
	tmp := list_to_map(deps)
	for k, _ := range tmp {
		ret = append(ret, k)
	}
	return
}

func (this *binding_generator) c_forward_declarations() {
	repo := gi.DefaultRepository()
	deps := get_dependcies_declarations(config.namespace)
	for _, dep := range deps {
		depv := strings.Split(dep, "-")
		_, err := repo.Require(depv[0], depv[1], 0)
		if err != nil {
			panic(err)
		}

		for i, n := 0, repo.NumInfo(depv[0]); i < n; i++ {
			this.c_forward_declaration10(repo.Info(depv[0], i))
			this.c_forward_declaration30(repo.Info(depv[0], i))
		}
	}

	for i, n := 0, repo.NumInfo(config.namespace); i < n; i++ {
		this.c_forward_declaration10(repo.Info(config.namespace, i))
	}
	for i, n := 0, repo.NumInfo(config.namespace); i < n; i++ {
		this.c_forward_declaration20(repo.Info(config.namespace, i))
	}
	for i, n := 0, repo.NumInfo(config.namespace); i < n; i++ {
		this.c_forward_declaration30(repo.Info(config.namespace, i))
	}
}

func go_utils(cb bool) string {
	var out bytes.Buffer

	go_utils_template.Execute(&out, map[string]interface{}{
		"gobjectns":   config.gns,
		"namespace":   config.namespace,
		"nocallbacks": !cb,
	})

	return out.String()
}

//------------------------------------------------------------------------
//------------------------------------------------------------------------

func (this *binding_generator) process_object_info(oi *gi.ObjectInfo) {
	p := printer_to(&this.go_bindings)

	parent := "C unsafe.Pointer"
	parentlike := ""
	if p := oi.Parent(); p != nil {
		parent = ""
		if ns := p.Namespace(); ns != config.namespace {
			parent += strings.ToLower(ns) + "."
		}
		parent += p.Name()
		parentlike = parent + "Like"
	}

	// interface that this class and its subclasses implement
	cprefix := gi.DefaultRepository().CPrefix(config.namespace)
	name := oi.Name()
	cgotype := cgo_type_for_interface(gi.ToBaseInfo(oi), type_pointer)

	var interfaces bytes.Buffer
	pi := printer_to(&interfaces)
	for i, n := 0, oi.NumInterface(); i < n; i++ {
		ii := oi.Interface(i)
		name := ii.Name()
		ns := ii.Namespace()
		if i != 0 {
			pi("\n\t")
		}
		if ns != config.namespace {
			pi("%s.", strings.ToLower(ns))
		}
		pi("%sImpl", name)
	}

	p("%s\n", execute_template(object_template, map[string]string{
		"name":       name,
		"cprefix":    cprefix,
		"cgotype":    cgotype,
		"parent":     parent,
		"parentlike": parentlike,
		"typeinit":   oi.TypeInit(),
		"gobjectns":  config.gns,
		"interfaces": interfaces.String(),
	}))

	for i, n := 0, oi.NumMethod(); i < n; i++ {
		meth := oi.Method(i)
		if config.is_method_blacklisted(name, meth.Name()) {
			p("// blacklisted: %s.%s (method)\n", name, meth.Name())
			continue
		}
		this.process_function_info(meth)
	}
}

func (this *binding_generator) process_struct_info(si *gi.StructInfo) {
	p := printer_to(&this.go_bindings)

	name := si.Name()
	size := si.Size()

	if si.IsGTypeStruct() {
		return
	}
	if strings.HasSuffix(name, "Private") {
		return
	}

	fullnm := strings.ToLower(si.Namespace()) + "." + name
	if config.is_disguised(fullnm) {
		size = -1
	}

	if !config.is_blacklisted("structdefs", name) {
		switch size {
		case -1:
			p("type %s struct { Pointer unsafe.Pointer }\n", name)
		case 0:
			p("type %s struct {}\n", name)
		default:
			p("type %s struct {\n", name)
			offset := 0
			for i, n := 0, si.NumField(); i < n; i++ {
				field := si.Field(i)
				fo := field.Offset()
				ft := field.Type()
				nm := field.Name()
				if fo != offset {
					pad := fo - offset
					p("\t_ [%d]byte\n", pad)
					offset += pad
				}
				if type_needs_wrapper(ft) {
					p("\t%s0 %s\n", nm, cgo_type(ft, type_exact))
				} else {
					p("\t%s %s\n", lower_case_to_camel_case(nm),
						go_type(ft, type_exact))
				}
				offset += type_size(ft, type_exact)
			}
			if size != offset {
				p("\t_ [%d]byte\n", size-offset)
			}
			p("}\n")
			//printf("type %s struct { data [%d]byte }\n", name, size)
		}

		// for each field that needs a wrapper, generate it
		for i, n := 0, si.NumField(); i < n; i++ {
			field := si.Field(i)
			ft := field.Type()
			nm := field.Name()

			if !type_needs_wrapper(ft) {
				continue
			}

			gotype := go_type(ft, type_return)
			p("func (this0 *%s) %s() %s {\n",
				name, lower_case_to_camel_case(nm), gotype)
			p("\tvar %s1 %s\n", nm, gotype)
			conv := cgo_to_go(ft, "this0."+nm+"0", nm+"1",
				conv_own_none, 0)
			p("%s", print_lines_with_indent(conv))
			p("\treturn %s1\n", nm)
			p("}\n")
		}
	}

	for i, n := 0, si.NumMethod(); i < n; i++ {
		meth := si.Method(i)
		if config.is_method_blacklisted(name, meth.Name()) {
			continue
		}

		this.process_function_info(meth)
	}
}

func (this *binding_generator) process_union_info(ui *gi.UnionInfo) {
	p := printer_to(&this.go_bindings)

	name := ui.Name()
	p("type %s struct {\n", name)
	p("\t_data [%d]byte\n", ui.Size())
	p("}\n")

	for i, n := 0, ui.NumMethod(); i < n; i++ {
		meth := ui.Method(i)
		if config.is_method_blacklisted(name, meth.Name()) {
			continue
		}

		this.process_function_info(meth)
	}
}

func (this *binding_generator) process_enum_info(ei *gi.EnumInfo) {
	p := printer_to(&this.go_bindings)

	p("type %s %s\n", ei.Name(), cgo_type_for_tag(ei.StorageType(), type_none))
	p("const (\n")
	for i, n := 0, ei.NumValue(); i < n; i++ {
		val := ei.Value(i)
		p("\t%s%s %s = %d\n", ei.Name(),
			lower_case_to_camel_case(val.Name()), ei.Name(), val.Value())
	}
	p(")\n")
}

func (this *binding_generator) process_constant_info(ci *gi.ConstantInfo) {
	p := printer_to(&this.go_bindings)

	name := ci.Name()
	if config.namespace == "Gdk" && strings.HasPrefix(name, "KEY_") {
		// KEY_ constants deserve a special treatment
		p("const Key_%s = %#v\n", name[4:], ci.Value())
		return
	}
	p("const %s = %#v\n", lower_case_to_camel_case(strings.ToLower(name)), ci.Value())
}

func (this *binding_generator) process_callback_info(ci *gi.CallableInfo) {
	p := printer_to(&this.go_bindings)

	// list of args
	var args []*gi.ArgInfo
	for i, n := 0, ci.NumArg(); i < n; i++ {
		arg := ci.Arg(i)
		args = append(args, arg)
	}

	userdata := -1
	for i, arg := range args {
		if arg.Closure() != -1 {
			userdata = i
			break
		}

		// treat any void* as userdata O_o
		t := arg.Type()
		if t.Tag() == gi.TYPE_TAG_VOID && t.IsPointer() {
			userdata = i
			break
		}
	}
	if userdata == -1 {
		p("// blacklisted (no userdata): ")
	}

	name := ci.Name()
	p("type %s func(", name)
	for i, ri, n := 0, 0, len(args); i < n; i++ {
		if i == userdata {
			continue
		}

		// I use here TypeReturn because for closures it's inverted
		// C code calls closure and it has to be concrete and Go code
		// returns stuff to C (like calling a C function)
		arg := args[i]
		if ri != 0 {
			p(", ")
		}
		p("%s %s", arg.Name(), go_type(arg.Type(), type_return))
		ri++
	}
	rt := ci.ReturnType()
	if !rt.IsPointer() && rt.Tag() == gi.TYPE_TAG_VOID {
		p(")\n")
	} else {
		p(") %s\n", go_type(rt, type_none))
	}
	if userdata == -1 {
		return
	}

	// now we need to generate two wrappers, it's a bit tricky due to cgo
	// ugliness.
	// 1. Function signature should consist of basic types and unsafe.Pointer
	//    for any pointer type, cgo cannot export other kinds of functions.
	// 2. Convert these types to C.* ones.
	// 3. Convert C.* types to Go values.
	// 4. Call Go callback (function pointer is in the userdata).
	// 5. Convert all returned Go values to C.* types and then to basic cgo
	//    friendly types.

	// signature
	ctype := c_type_for_interface(gi.ToBaseInfo(ci), type_none)
	p("//export _%s_c_wrapper\n", ctype)
	p("func _%s_c_wrapper(", ctype)
	for i, arg := range args {
		if i != 0 {
			p(", ")
		}
		p("%s0 %s", arg.Name(), simple_cgo_type(arg.Type(), type_none))
	}
	p(") ")
	if ret := ci.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		p("%s ", simple_cgo_type(ret, type_none))
	}

	p("{\n")
	// --- body stage 1 (C to Go conversions)

	// var declarations
	for i, arg := range args {
		gotype := go_type(arg.Type(), type_return)
		if i == userdata {
			gotype = name
		}
		p("\tvar %s1 %s\n", arg.Name(), gotype)
	}

	// conversions
	for i, arg := range args {
		t := arg.Type()
		aname := arg.Name()
		if i == userdata {
			p("\t%s1 = *(*%s)(%s0)\n", aname, name, aname)
			continue
		}
		ownership := ownership_to_conv_flags(arg.OwnershipTransfer())
		conv := simple_cgo_to_go(t, aname+"0", aname+"1", ownership)
		p("%s", print_lines_with_indent(conv))
	}

	// --- body stage 2 (the callback call)
	p("\t")
	if ret := ci.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		p("ret1 := ")
	}
	p("%s1(", args[userdata].Name())
	for i, ri, n := 0, 0, len(args); i < n; i++ {
		if i == userdata {
			continue
		}
		if ri != 0 {
			p(", ")
		}
		p("%s1", args[i].Name())
		ri++
	}
	p(")\n")

	// --- body stage 3 (return value)
	if ret := ci.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		p("\tvar ret2 %s\n", cgo_type(ret, type_none))
		ownership := ownership_to_conv_flags(ci.CallerOwns())
		conv := go_to_cgo(ret, "ret1", "ret2", ownership)
		p("%s", print_lines_with_indent(conv))
		p("\treturn (%s)(ret2)\n", simple_cgo_type(ret, type_none))
	}
	p("}\n")

	// and finally add "_once" wrapper
	p("//export _%s_c_wrapper_once\n", ctype)
	p("func _%s_c_wrapper_once(", ctype)
	for i, arg := range args {
		if i != 0 {
			p(", ")
		}
		p("%s0 %s", arg.Name(), simple_cgo_type(arg.Type(), type_none))
	}
	p(") ")
	if ret := ci.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		p("%s ", simple_cgo_type(ret, type_none))
	}
	p("{\n\t")
	if ret := ci.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		p("ret := ")
	}
	p("_%s_c_wrapper(", ctype)
	for i, arg := range args {
		if i != 0 {
			p(", ")
		}
		p("%s0", arg.Name())
	}
	p(")\n")
	p("\t%sHolder.Release(%s0)\n", config.gns,
		args[userdata].Name())
	if ret := ci.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		p("\treturn ret\n")
	}
	p("}\n")
}

func (this *binding_generator) process_function_info(fi *gi.FunctionInfo) {
	p := printer_to(&this.go_bindings)

	var fullnm string
	flags := fi.Flags()
	name := fi.Name()

	container := fi.Container()

	// --- header
	fb := new_function_builder(fi)
	p("func ")
	if flags&gi.FUNCTION_IS_METHOD != 0 {
		// add receiver if it's a method
		p("(this0 %s) ", go_type_for_interface(container, type_pointer|type_receiver))
		fullnm = container.Name() + "."
	}
	switch {
	case flags&gi.FUNCTION_IS_CONSTRUCTOR != 0:
		// special names for constructors
		name = fmt.Sprintf("New%s%s", container.Name(), ctor_suffix(name))
	case flags&gi.FUNCTION_IS_METHOD == 0 && container != nil:
		name = fmt.Sprintf("%s%s", container.Name(), lower_case_to_camel_case(name))
	default:
		name = fmt.Sprintf("%s", lower_case_to_camel_case(name))
	}
	fullnm += name
	name = config.rename(fullnm, name)
	p("%s(", name)
	for i, arg := range fb.args {
		if i != 0 {
			p(", ")
		}
		p("%s0 %s", arg.arg_info.Name(), go_type(arg.type_info, type_none))
	}
	p(")")
	switch len(fb.rets) {
	case 0:
		// do nothing if there are not return values
	case 1:
		if flags&gi.FUNCTION_IS_CONSTRUCTOR != 0 {
			// override return types for constructors, We can't
			// return generic widget here as C does. Go's type
			// system is stronger.
			p(" %s", go_type_for_interface(container, type_pointer|type_return))
			break
		}
		if fb.rets[0].index == -2 {
			p(" error")
		} else {
			p(" %s", go_type(fb.rets[0].type_info, type_return))
		}
	default:
		p(" (")
		for i, ret := range fb.rets {
			if ret.index == -2 {
				// special error type in Go to represent GError
				p("error")
				continue
			}
			p(go_type(ret.type_info, type_return))
			if i != len(fb.rets)-1 {
				p(", ")
			}
		}
		p(")")
	}
	p(" {\n")

	// --- body stage 1 (Go to C conversions)

	// var declarations
	if flags&gi.FUNCTION_IS_METHOD != 0 {
		p("\tvar this1 %s\n", cgo_type_for_interface(container, type_pointer))
	}
	for _, arg := range fb.args {
		p("\tvar %s1 %s\n", arg.arg_info.Name(), cgo_type(arg.type_info, type_none))
		if al := arg.type_info.ArrayLength(); al != -1 {
			arg := fb.orig_args[al]
			p("\tvar %s1 %s\n", arg.Name(), cgo_type(arg.Type(), type_none))
		}
	}

	for _, ret := range fb.rets {
		if ret.index == -1 {
			continue
		}

		if ret.index == -2 {
			p("\tvar err1 *C.GError\n")
			continue
		}

		if ret.arg_info.Direction() == gi.DIRECTION_INOUT {
			continue
		}

		p("\tvar %s1 %s\n", ret.arg_info.Name(), cgo_type(ret.type_info, type_none))
		if al := ret.type_info.ArrayLength(); al != -1 {
			arg := fb.orig_args[al]
			p("\tvar %s1 %s\n", arg.Name(), cgo_type(arg.Type(), type_none))
		}
	}

	// conversions
	if flags&gi.FUNCTION_IS_METHOD != 0 {
		conv := go_to_cgo_for_interface(container, "this0", "this1", conv_pointer)
		p("%s", print_lines_with_indent(conv))
	}
	for _, arg := range fb.args {
		nm := arg.arg_info.Name()
		conv := go_to_cgo(arg.type_info, nm+"0", nm+"1", conv_none)
		p("%s", print_lines_with_indent(conv))

		// register callback in the global map
		if arg.type_info.Tag() == gi.TYPE_TAG_INTERFACE {
			bi := arg.type_info.Interface()
			if bi.Type() == gi.INFO_TYPE_CALLBACK {
				if arg.arg_info.Scope() != gi.SCOPE_TYPE_CALL {
					p("\t%sHolder.Grab(%s1)\n", config.gns, nm)
				}
			}
		}

		// array length
		if len := arg.type_info.ArrayLength(); len != -1 {
			lenarg := fb.orig_args[len]
			conv = go_to_cgo(lenarg.Type(), "len("+nm+"0)", lenarg.Name()+"1", conv_none)
			p("\t%s\n", conv)
		}
	}

	// --- body stage 2 (the function call)
	p("\t")
	if fb.has_return_value() {
		p("ret1 := ")
	}

	userdata, destroy, scope := fb.has_closure_argument()
	p("C.")
	if scope != gi.SCOPE_TYPE_INVALID {
		p("_")
	}
	p("%s(", fi.Symbol())
	if flags&gi.FUNCTION_IS_METHOD != 0 {
		p("this1")
		if len(fb.orig_args) > 0 {
			p(", ")
		}
	}
	for i, ri, n := 0, 0, len(fb.orig_args); i < n; i++ {
		if i == userdata || i == destroy {
			continue
		}

		oarg := fb.orig_args[i]

		var arg string
		dir := oarg.Direction()
		if dir == gi.DIRECTION_INOUT || dir == gi.DIRECTION_OUT {
			arg = fmt.Sprintf("&%s1", oarg.Name())
		} else {
			arg = fmt.Sprintf("%s1", oarg.Name())
		}

		if ri != 0 {
			p(", ")
		}
		p("%s", arg)
		ri++
	}
	if flags&gi.FUNCTION_THROWS != 0 {
		p(", &err1")
	}
	p(")\n")

	// --- body stage 3 (C to Go conversions)

	// var declarations
	for _, ret := range fb.rets {
		switch ret.index {
		case -1:
			if flags&gi.FUNCTION_IS_CONSTRUCTOR != 0 {
				p("\tvar ret2 %s\n", go_type_for_interface(container, type_pointer|type_return))
			} else {
				p("\tvar ret2 %s\n", go_type(ret.type_info, type_return))
			}
		case -2:
			p("\tvar err2 error\n")
		default:
			p("\tvar %s2 %s\n", ret.arg_info.Name(), go_type(ret.type_info, type_return))
		}
	}

	// conversions
	for _, ret := range fb.rets {
		if ret.index == -2 {
			p("\tif err1 != nil {\n")
			p("\t\terr2 = ((*_GError)(unsafe.Pointer(err1))).ToGError()\n")
			p("\t\tC.g_error_free(err1)\n")
			p("\t}\n")
			continue
		}

		var nm string
		if ret.index == -1 {
			nm = "ret"
		} else {
			nm = ret.arg_info.Name()
		}

		// array length
		if len := ret.type_info.ArrayLength(); len != -1 {
			// TODO: move this to cgo_to_go
			lenarg := fb.orig_args[len]
			p("\t%s2 = make(%s, %s)\n",
				nm, go_type(ret.type_info, type_return),
				lenarg.Name()+"1")
		} else if ret.type_info.IsZeroTerminated() {
			//p("\t%s2 = make(%s, uint(C._array_length(unsafe.Pointer(%s1))))\n", nm, go_type(ret.type_info, type_return), nm)
		}

		var ownership gi.Transfer
		if ret.index == -1 {
			ownership = fi.CallerOwns()
		} else {
			ownership = ret.arg_info.OwnershipTransfer()
			if ret.arg_info.Direction() == gi.DIRECTION_INOUT {
				// not sure if it's true in all cases, but so far
				ownership = gi.TRANSFER_NOTHING
			}
		}

		p("\n//DEBUG: %s1(%s):flags = %q\n", nm, ret.type_info.Tag(), ownership_to_conv_flags(ownership))

		var conv string
		if flags&gi.FUNCTION_IS_CONSTRUCTOR != 0 && ret.index == -1 {
			conv = cgo_to_go_for_interface(container, "ret1", "ret2",
				conv_pointer|ownership_to_conv_flags(ownership))
		} else {
			conv = cgo_to_go(ret.type_info, nm+"1", nm+"2",
				ownership_to_conv_flags(ownership), 0)
		}
		p("%s", print_lines_with_indent(conv))
	}

	// --- body stage 4 (return)
	if len(fb.rets) == 0 {
		p("}\n")
		return
	}

	p("\treturn ")
	for i, ret := range fb.rets {
		var nm string
		switch ret.index {
		case -1:
			nm = "ret2"
		case -2:
			nm = "err2"
		default:
			nm = ret.arg_info.Name() + "2"
		}
		if i != 0 {
			p(", ")
		}
		p(nm)
	}

	p("\n}\n")
}

func (this *binding_generator) process_interface_info(ii *gi.InterfaceInfo) {
	p := printer_to(&this.go_bindings)

	name := ii.Name()
	cprefix := gi.DefaultRepository().CPrefix(ii.Namespace())
	cgotype := cgo_type_for_interface(gi.ToBaseInfo(ii), type_pointer)

	p("%s\n", execute_template(interface_template, map[string]string{
		"name":      name,
		"cprefix":   cprefix,
		"cgotype":   cgotype,
		"typeinit":  ii.TypeInit(),
		"gobjectns": config.gns,
	}))

	for i, n := 0, ii.NumMethod(); i < n; i++ {
		meth := ii.Method(i)
		if config.is_method_blacklisted(name, meth.Name()) {
			p("// blacklisted: %s.%s (method)\n", name, meth.Name())
			continue
		}
		this.process_function_info(meth)
	}
}

func (this *binding_generator) process_base_info(bi *gi.BaseInfo) {
	p := printer_to(&this.go_bindings)

	if config.is_object_blacklisted(bi) {
		p("// blacklisted: %s (%s)\n", bi.Name(), bi.Type())
		return
	}

	switch bi.Type() {
	case gi.INFO_TYPE_UNION:
		this.process_union_info(gi.ToUnionInfo(bi))
	case gi.INFO_TYPE_STRUCT:
		this.process_struct_info(gi.ToStructInfo(bi))
	case gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS:
		this.process_enum_info(gi.ToEnumInfo(bi))
	case gi.INFO_TYPE_CONSTANT:
		this.process_constant_info(gi.ToConstantInfo(bi))
	case gi.INFO_TYPE_CALLBACK:
		this.process_callback_info(gi.ToCallableInfo(bi))
	case gi.INFO_TYPE_FUNCTION:
		this.process_function_info(gi.ToFunctionInfo(bi))
	case gi.INFO_TYPE_INTERFACE:
		this.process_interface_info(gi.ToInterfaceInfo(bi))
	case gi.INFO_TYPE_OBJECT:
		this.process_object_info(gi.ToObjectInfo(bi))
	default:
		p("// TODO: %s (%s)\n", bi.Name(), bi.Type())
	}
}
