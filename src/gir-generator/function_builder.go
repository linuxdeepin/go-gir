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
	"gi"
)

type function_builder struct {
	function  *gi.FunctionInfo
	args      []function_builder_arg
	rets      []function_builder_arg
	orig_args []*gi.ArgInfo
}

type function_builder_arg struct {
	index     int
	arg_info  *gi.ArgInfo
	type_info *gi.TypeInfo
}

func int_slice_contains(haystack []int, needle int) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

func new_function_builder(fi *gi.FunctionInfo) *function_builder {
	fb := new(function_builder)
	fb.function = fi

	// prepare an array of ArgInfos
	for i, n := 0, fi.NumArg(); i < n; i++ {
		arg := fi.Arg(i)
		fb.orig_args = append(fb.orig_args, arg)
	}

	// build skip list
	var skiplist []int
	for _, arg := range fb.orig_args {
		ti := arg.Type()

		len := ti.ArrayLength()
		if len != -1 {
			skiplist = append(skiplist, len)
		}

		clo := arg.Closure()
		if clo != -1 {
			skiplist = append(skiplist, clo)
		}

		des := arg.Destroy()
		if des != -1 {
			skiplist = append(skiplist, des)
		}
	}

	// then walk over arguments
	for i, ai := range fb.orig_args {
		if int_slice_contains(skiplist, i) {
			continue
		}

		ti := ai.Type()

		switch ai.Direction() {
		case gi.DIRECTION_IN:
			fb.args = append(fb.args, function_builder_arg{i, ai, ti})
		case gi.DIRECTION_INOUT:
			fb.args = append(fb.args, function_builder_arg{i, ai, ti})
			fb.rets = append(fb.rets, function_builder_arg{i, ai, ti})
		case gi.DIRECTION_OUT:
			fb.rets = append(fb.rets, function_builder_arg{i, ai, ti})
		}
	}

	// add return value if it exists to 'rets'
	if ret := fi.ReturnType(); ret != nil && ret.Tag() != gi.TYPE_TAG_VOID {
		fb.rets = append(fb.rets, function_builder_arg{-1, nil, ret})
	}

	// add GError special argument (if any)
	if fi.Flags()&gi.FUNCTION_THROWS != 0 {
		fb.rets = append(fb.rets, function_builder_arg{-2, nil, nil})
	}

	return fb
}

func (fb *function_builder) has_return_value() bool {
	return (len(fb.rets) > 0 && fb.rets[len(fb.rets)-1].index == -1) ||
		(len(fb.rets) > 1 && fb.rets[len(fb.rets)-2].index == -1)
}

func (fb *function_builder) has_closure_argument() (int, int, gi.ScopeType) {
	for _, arg := range fb.args {
		userdata := arg.arg_info.Closure()
		if userdata == -1 {
			continue
		}

		if arg.type_info.Tag() != gi.TYPE_TAG_INTERFACE {
			continue
		}

		if arg.type_info.Interface().Type() != gi.INFO_TYPE_CALLBACK {
			continue
		}

		destroy := arg.arg_info.Destroy()
		scope := arg.arg_info.Scope()
		return userdata, destroy, scope
	}
	return -1, -1, gi.SCOPE_TYPE_INVALID
}
