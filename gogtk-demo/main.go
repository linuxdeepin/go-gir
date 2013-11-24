package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"gobject/gdk-3.0"
	"gobject/gdkpixbuf-2.0"
	"gobject/gobject-2.0"
	"gobject/gtk-3.0"
	"gobject/pango-1.0"
	"io/ioutil"
	"os"
	"strings"
	"./gogtk-demo/common"
)

const (
	title_column = iota
	filename_column
	func_column
	style_column
)

var infobuf *gtk.TextBuffer
var sourcebuf *gtk.TextBuffer

func create_tree_view() *gtk.Widget {
	model := gtk.NewTreeStore(
		gobject.String,      // title
		gobject.String,      // filename
		gobject.GoInterface, // app
		gobject.Int,         // style
	)

	tree_view := gtk.NewTreeViewWithModel(model)
	selection := tree_view.GetSelection()
	selection.SetMode(gtk.SelectionModeBrowse)
	tree_view.SetSizeRequest(200, -1)

	for _, demo := range demos {
		iter := model.Append(nil, demo.Title, demo.Filename, demo.Func, pango.StyleNormal)
		for _, cdemo := range demo.Children {
			model.Append(&iter, cdemo.Title, cdemo.Filename, cdemo.Func, pango.StyleNormal)
		}
	}

	r := gtk.NewCellRendererText()
	c := gtk.NewTreeViewColumnWithAttributes("Widget (double click for demo)", r,
		"text", title_column,
		"style", style_column)
	tree_view.AppendColumn(c)

	iter, _ := model.GetIterFirst()
	selection.SelectIter(&iter)

	selection.Connect("changed", func(selection *gtk.TreeSelection) {
		_, iter, ok := selection.GetSelected()
		if !ok {
			return
		}

		var filename string
		model.Get(&iter, filename_column, &filename)
		if filename != "" {
			load_file(filename)
		}
	})
	tree_view.Connect("row-activated", func(tree_view *gtk.TreeView, path *gtk.TreePath) {
		iter, _ := model.GetIter(path)
		var app interface{}
		var style pango.Style
		model.Get(&iter, func_column, &app, style_column, &style)
		if style == pango.StyleItalic {
			style = pango.StyleNormal
		} else {
			style = pango.StyleItalic
		}

		if app.(DemoFunc) == nil {
			return
		}

		model.Set(&iter, style_column, style)
		w := app.(DemoFunc)(gtk.ToWindow(tree_view.GetToplevel()))
		if w != nil {
			w.Connect("destroy", func() {
				var style pango.Style
				model.Get(&iter, style_column, &style)
				if style == pango.StyleItalic {
					model.Set(&iter, style_column, pango.StyleNormal)
				}
			})
		}
	})
	tree_view.CollapseAll()
	tree_view.SetHeadersVisible(false)

	scrolled_window := gtk.NewScrolledWindow(nil, nil)
	scrolled_window.SetPolicy(gtk.PolicyTypeNever, gtk.PolicyTypeAutomatic)
	scrolled_window.Add(tree_view)

	label := gtk.NewLabel("Widget (double click for demo)")

	nb := gtk.NewNotebook()
	nb.AppendPage(scrolled_window, label)
	tree_view.GrabFocus()

	return gtk.ToWidget(nb)
}

func create_text(is_source bool) (*gtk.ScrolledWindow, *gtk.TextBuffer) {
	sw := gtk.NewScrolledWindow(nil, nil)
	sw.SetPolicy(gtk.PolicyTypeAutomatic, gtk.PolicyTypeAutomatic)
	sw.SetShadowType(gtk.ShadowTypeIn)

	tv := gtk.NewTextView()
	sw.Add(tv)

	buf := gtk.NewTextBuffer(nil)
	tv.SetBuffer(buf)
	tv.SetEditable(false)
	tv.SetCursorVisible(false)
	if is_source {
		font_desc := pango.FontDescriptionFromString(
			"BitStream Vera Sans Mono, Monaco, Consolas, Courier New, monospace 9")
		tv.OverrideFont(font_desc)
		tv.SetWrapMode(gtk.WrapModeNone)
	} else {
		tv.SetWrapMode(gtk.WrapModeWord)
		tv.SetPixelsAboveLines(2)
		tv.SetPixelsBelowLines(2)
	}

	return sw, buf
}

func setup_default_icon() {
	filename := common.FindFile("gtk-logo-rgb.gif")
	pixbuf, err := gdkpixbuf.NewPixbufFromFile(filename)
	if err != nil {
		dialog := gtk.NewMessageDialog(nil, 0, gtk.MessageTypeError, gtk.ButtonsTypeClose,
			"Failed to read icon file: %s", err)
		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.ShowAll()
	} else {
		pixbuf = pixbuf.AddAlpha(true, 0xFF, 0xFF, 0xFF)
		gtk.WindowSetDefaultIcon(pixbuf)
	}
}

var current_file string

var go_highlighter_idents = map[string]string{
	"true":  "predefined",
	"false": "predefined",
	"iota":  "predefined",
	"nil":   "predefined",
}

type go_highlighter struct {
	fset *token.FileSet
	buf  *gtk.TextBuffer
	file *ast.File
	data []byte
}

func (this *go_highlighter) highlight(tag string, beg, end token.Pos) {
	begp := this.fset.Position(beg)
	endp := this.fset.Position(end)

	begc, endc := begp.Column-1, endp.Column-1
	begl, endl := begp.Line-1, endp.Line-1
	if begl < 0 {
		return
	}

	begi := this.buf.GetIterAtLineOffset(begl, begc)
	endi := this.buf.GetIterAtLineOffset(endl, endc)
	this.buf.ApplyTagByName(tag, &begi, &endi)
}

func (this *go_highlighter) highlight_file() {
	var s scanner.Scanner
	fset := token.NewFileSet()
	s.Init(fset.AddFile(current_file, fset.Base(), len(this.data)), this.data, nil, 0)
	for {
		pos, tok, str := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok.IsKeyword() {
			this.highlight("keyword", pos, pos+token.Pos(len(str)))
		}
	}

	ast.Inspect(this.file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.BasicLit:
			switch n.Kind {
			case token.STRING, token.CHAR:
				this.highlight("string", n.Pos(), n.End())
			case token.INT, token.FLOAT, token.IMAG:
				this.highlight("number", n.Pos(), n.End())
			}
		case *ast.Ident:
			if tag, ok := go_highlighter_idents[n.Name]; ok {
				this.highlight(tag, n.Pos(), n.End())
				break
			}

			if n.Obj != nil && n.Obj.Pos() == n.Pos() {
				if n.Obj.Kind == ast.Fun {
					this.highlight("function", n.Pos(), n.End())
				} else {
					this.highlight("declaration", n.Pos(), n.End())
				}
			}
		case *ast.CallExpr:
			switch f := n.Fun.(type) {
			case *ast.Ident:
				this.highlight("funcall", f.Pos(), f.End())
			case *ast.SelectorExpr:
				this.highlight("funcall", f.Sel.Pos(), f.Sel.End())
			}
		}

		return true
	})

	for _, cg := range this.file.Comments {
		this.highlight("comment", cg.Pos(), cg.End())
	}
}

func fontify(data []byte) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, current_file, data, parser.ParseComments)
	if err != nil {
		println("failed to parse file: ", err.Error())
		return
	}

	x := go_highlighter{fset, sourcebuf, file, data}
	x.highlight_file()
}

func set_info(info string) {
	var para bytes.Buffer

	iter := infobuf.GetStartIter()
	lines := strings.Split(info, "\n")
	for i, line := range lines {
		switch i {
		case 0:
			infobuf.InsertWithTagsByName(&iter, line, -1, "title")
			infobuf.Insert(&iter, "\n", -1)
		case 1:
			continue
		default:
			if line == "" {
				// flush paragraph on empty lines
				para.WriteString("\n")
				infobuf.Insert(&iter, para.String(), -1)
				para.Reset()
				continue
			}

			// by default append to paragraph buffer
			if para.Len() != 0 {
				para.WriteString(" ")
			}
			para.WriteString(line)

			// flush on last like as well:
			if line != "" && i == len(lines)-1 {
				para.WriteString("\n")
				infobuf.Insert(&iter, para.String(), -1)
				para.Reset()
			}
		}
	}
}

func load_file(filename string) {
	if current_file == filename {
		return
	}

	current_file = filename

	// clear info and source buffers
	beg, end := infobuf.GetBounds()
	infobuf.Delete(&beg, &end)

	beg, end = sourcebuf.GetBounds()
	sourcebuf.Delete(&beg, &end)

	// find file
	filename_full := common.FindFile(filename)
	if filename_full == "" {
		println("failed to find file: ", filename)
		return
	}

	// load file
	data, err := ioutil.ReadFile(filename_full)
	if err != nil {
		println("failed to read file: ", err.Error())
		return
	}

	// figure out package info and starting offset
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, current_file, data,
		parser.ParseComments|parser.PackageClauseOnly)
	if err != nil {
		println("failed to parse package clause: ", err.Error())
		return
	}

	var offset int
	if file.Doc != nil {
		set_info(strings.TrimSpace(file.Doc.Text()))
		pos := fset.Position(file.Doc.End())
		offset = pos.Offset + 1
	}

	beg = sourcebuf.GetStartIter()
	sourcebuf.Insert(&beg, string(data[offset:]), -1)

	fontify(data[offset:])
}

func main() {
	gdk.ThreadsInit()
	gtk.Init(os.Args)
	setup_default_icon()
	window := gtk.NewWindow(gtk.WindowTypeToplevel)
	window.SetTitle("GoGTK Code Demos")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	window.SetDefaultSize(600, 400)

	hbox := gtk.NewHBox(false, 3)
	window.Add(hbox)

	// treeview
	tree_view := create_tree_view()
	hbox.PackStart(tree_view, false, false, 0)

	// notebook
	notebook := gtk.NewNotebook()
	hbox.PackStart(notebook, true, true, 0)

	// info
	var sw *gtk.ScrolledWindow
	sw, infobuf = create_text(false)
	notebook.AppendPage(sw, gtk.NewLabelWithMnemonic("_Info"))
	infobuf.CreateTag("title",
		"font", "Sans 18")

	// source
	sw, sourcebuf = create_text(true)
	notebook.AppendPage(sw, gtk.NewLabelWithMnemonic("_Source"))

	sourcebuf.CreateTag("comment",
		"foreground", "#0066FF")
	sourcebuf.CreateTag("declaration",
		"foreground", "#318495")
	sourcebuf.CreateTag("funcall",
		"foreground", "#3C4C72",
		"weight", pango.WeightBold)
	sourcebuf.CreateTag("string",
		"foreground", "#036A07")
	sourcebuf.CreateTag("keyword",
		"weight", pango.WeightBold,
		"foreground", "#0707FF")
	sourcebuf.CreateTag("function",
		"weight", pango.WeightBold,
		"foreground", "#0000A2")
	sourcebuf.CreateTag("number",
		"weight", pango.WeightBold,
		"foreground", "#C5060B")
	sourcebuf.CreateTag("predefined",
		"weight", pango.WeightBold,
		"foreground", "#585CF6")

	window.ShowAll()

	load_file(demos[0].Filename)

	gdk.ThreadsEnter()
	gtk.Main()
	gdk.ThreadsLeave()
}
