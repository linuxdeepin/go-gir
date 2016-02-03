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
	"encoding/json"
	"flag"
	"fmt"
	"gi"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var config config_type

type config_type struct {
	out_go           *bufio.Writer
	out_c            *bufio.Writer
	out_h            *bufio.Writer
	namespace        string
	version          string
	pkg              string
	gns              string
	renames          map[string]string
	blacklist        map[string]map[string]bool
	whitelist        map[string]map[string]bool
	method_blacklist map[string]map[string]bool
	method_whitelist map[string]map[string]bool
	disguised_types  map[string]bool
	word_subst       map[string]string
}

func (this *config_type) load(path string) {
	var tmp struct {
		Namespace       string              `json:"namespace"`
		Version         string              `json:"version"`
		Blacklist       map[string][]string `json:"blacklist"`
		Whitelist       map[string][]string `json:"whitelist"`
		MethodBlacklist map[string][]string `json:"method-blacklist"`
		MethodWhitelist map[string][]string `json:"method-whitelist"`
		Renames         map[string]string   `json:"renames"`
	}

	err := parse_json_with_comments(path, &tmp)
	if err != nil {
		panic(err)
	}

	this.namespace = tmp.Namespace
	this.version = tmp.Version
	this.blacklist = map_list_to_map_map(tmp.Blacklist)
	this.whitelist = map_list_to_map_map(tmp.Whitelist)
	this.method_blacklist = map_list_to_map_map(tmp.MethodBlacklist)
	this.method_whitelist = map_list_to_map_map(tmp.MethodWhitelist)
	this.renames = tmp.Renames

	this.pkg = strings.ToLower(tmp.Namespace)
	if this.namespace != "GObject" {
		this.gns = "gobject."
	}
}

func (this *config_type) load_sys(path string) {
	var tmp struct {
		DisguisedTypes []string          `json:"disguised-types"`
		WordSubst      map[string]string `json:"word-subst"`
	}

	err := parse_json_with_comments(path, &tmp)
	if err != nil {
		panic(err)
	}

	this.disguised_types = list_to_map(tmp.DisguisedTypes)
	this.word_subst = tmp.WordSubst
}

func (this *config_type) rename(path, oldname string) string {
	if newname, ok := this.renames[path]; ok {
		return newname
	}
	return oldname
}

func (this *config_type) is_disguised(name string) bool {
	_, ok := this.disguised_types[name]
	return ok
}

func (this *config_type) is_object_blacklisted(bi *gi.BaseInfo) bool {
	switch bi.Type() {
	case gi.INFO_TYPE_UNION:
		return config.is_blacklisted("unions", bi.Name())
	case gi.INFO_TYPE_STRUCT:
		return config.is_blacklisted("structs", bi.Name())
	case gi.INFO_TYPE_ENUM, gi.INFO_TYPE_FLAGS:
		return config.is_blacklisted("enums", bi.Name())
	case gi.INFO_TYPE_CONSTANT:
		return config.is_blacklisted("constants", bi.Name())
	case gi.INFO_TYPE_CALLBACK:
		return config.is_blacklisted("callbacks", bi.Name())
	case gi.INFO_TYPE_FUNCTION:
		c := bi.Container()
		if c != nil {
			return config.is_method_blacklisted(c.Name(), bi.Name())
		}
		return config.is_blacklisted("functions", bi.Name())
	case gi.INFO_TYPE_INTERFACE:
		return config.is_blacklisted("interfaces", bi.Name())
	case gi.INFO_TYPE_OBJECT:
		return config.is_blacklisted("objects", bi.Name())
	default:
		println("TODO: %s (%s)\n", bi.Name(), bi.Type())
		return true
	}
	panic("unreachable")
}

func (this *config_type) is_blacklisted(section, entry string) bool {
	// check if the entry is in the blacklist
	if section_map, ok := this.blacklist[section]; ok {
		if _, ok := section_map[entry]; ok {
			return true
		}
	}

	// check if the entry is missing from the whitelist
	if section_map, ok := this.whitelist[section]; ok {
		if _, ok := section_map[entry]; !ok {
			return true
		}
	}

	return false
}

func (this *config_type) is_method_blacklisted(class, method string) bool {
	// don't want to see these
	if method == "ref" || method == "unref" {
		return true
	}

	if class_map, ok := this.method_blacklist[class]; ok {
		if _, ok := class_map[method]; ok {
			return true
		}
	}

	if class_map, ok := this.method_whitelist[class]; ok {
		if _, ok := class_map[method]; !ok {
			return true
		}
	}

	return false
}

func parse_json_with_comments(filename string, data interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(new_comment_skipper(f))
	err = d.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	sysconfig_path := flag.String("config", "", "specify global config file")
	output_dir := flag.String("o", "", "override output directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <dir>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	// figure in/out paths
	in_dir, in_base := filepath.Split(flag.Arg(0))
	in_path := flag.Arg(0)

	out_dir := in_dir
	if *output_dir != "" {
		out_dir = *output_dir
		os.MkdirAll(out_dir, 0755)
	}
	out_base := filepath.Join(out_dir, in_base[:len(in_base)-6])

	// parse system config file
	if *sysconfig_path != "" {
		config.load_sys(*sysconfig_path)
	}

	// parse local config
	config.load(filepath.Join(in_dir, "config.json"))

	// load namespace
	_, err := gi.DefaultRepository().Require(config.namespace, config.version, 0)
	panic_if_error(err)

	// load go template
	go_template, err := ioutil.ReadFile(in_path)
	panic_if_error(err)

	// generate bindings
	bg := new_binding_generator(out_base)
	defer bg.release()
	bg.generate(string(go_template))
}
