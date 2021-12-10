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
	"fmt"
	"io"
	"strings"
	"text/template"
)

func panic_if_error(err error) {
	if err != nil {
		panic(err)
	}
}

func printer_to(w io.Writer) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}
}

func lower_case_to_camel_case(name string) string {
	var out bytes.Buffer
	for _, word := range strings.Split(name, "_") {
		word = strings.ToLower(word)
		if subst, ok := config.word_subst[word]; ok {
			out.WriteString(subst)
			continue
		}

		if word == "" {
			out.WriteString("_")
			continue
		}
		out.WriteString(strings.ToUpper(word[0:1]))
		out.WriteString(word[1:])
	}
	return out.String()
}

func must_template(tpl string) *template.Template {
	tpl = strings.TrimSpace(tpl)
	return template.Must(
		template.New("").
			Delims("[<", ">]").
			Parse(tpl),
	)
}

func execute_template(tpl *template.Template, args interface{}) string {
	var out bytes.Buffer
	tpl.Execute(&out, args)
	return out.String()
}

func ctor_suffix(name string) string {
	if len(name) > 4 {
		return lower_case_to_camel_case(name[4:])
	}
	return ""
}

func print_lines_with_indent(str string) string {
	var out bytes.Buffer
	if str == "" {
		return ""
	}
	for _, line := range strings.Split(str, "\n") {
		fmt.Fprintf(&out, "\t%s\n", line)
	}

	return out.String()
}

func map_list_to_map_map(maplist map[string][]string) map[string]map[string]bool {
	out := make(map[string]map[string]bool)
	for section, list := range maplist {
		out[section] = list_to_map(list)
	}
	return out
}

func list_to_map(list []string) map[string]bool {
	m := make(map[string]bool)
	for _, entry := range list {
		m[entry] = true
	}
	return m
}
