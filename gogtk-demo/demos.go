package main

import "gobject/gtk-3.0"
import (
	"./gogtk-demo/assistant"
	"./gogtk-demo/builder"
	"./gogtk-demo/button_boxes"
	"./gogtk-demo/color_selector"
	"./gogtk-demo/entry_buffer"
	"./gogtk-demo/expander"
	"./gogtk-demo/info_bar"
	"./gogtk-demo/links"
	"./gogtk-demo/list_store"
	"./gogtk-demo/main_window"
	"./gogtk-demo/paned_widgets"
	"./gogtk-demo/pickers"
	"./gogtk-demo/size_groups"
	"./gogtk-demo/spinner"
	"./gogtk-demo/stock_browser"
	"./gogtk-demo/menus"
	"./gogtk-demo/entry_completion"
	"./gogtk-demo/drawing_area"
	"./gogtk-demo/combo_boxes"
	"./gogtk-demo/dialog"
	"./gogtk-demo/search_entry"
	"./gogtk-demo/iconview_edit"
	"./gogtk-demo/iconview"
	"./gogtk-demo/pixbufs"
)

type DemoFunc func(mainwin *gtk.Window) *gtk.Window

type DemoDesc struct {
	Title    string
	Filename string
	Func     DemoFunc
	Children []*DemoDesc
}

var demos = []*DemoDesc{
	{Title: "Application main window", Filename: "main_window.go", Func: main_window.Do},
	{Title: "Assistant", Filename: "assistant.go", Func: assistant.Do},
	{Title: "Builder", Filename: "builder.go", Func: builder.Do},
	{Title: "Button Boxes", Filename: "button_boxes.go", Func: button_boxes.Do},
	{Title: "Color Selector", Filename: "color_selector.go", Func: color_selector.Do},
	{Title: "Combo boxes", Filename: "combo_boxes.go", Func: combo_boxes.Do},
	{Title: "Dialog and Message Boxes", Filename: "dialog.go", Func: dialog.Do},
	{Title: "Drawing Area", Filename: "drawing_area.go", Func: drawing_area.Do},
	{Title: "Entry", Children: []*DemoDesc{
		{Title: "Entry Buffer", Filename: "entry_buffer.go", Func: entry_buffer.Do},
		{Title: "Entry Completion", Filename: "entry_completion.go", Func: entry_completion.Do},
		{Title: "Search Entry", Filename: "search_entry.go", Func: search_entry.Do},
	}},
	{Title: "Expander", Filename: "expander.go", Func: expander.Do},
	{Title: "Icon View", Children: []*DemoDesc{
		{Title: "Icon View Basics", Filename: "iconview.go", Func: iconview.Do},
		{Title: "Editing and Drag-and-Drop", Filename: "iconview_edit.go", Func: iconview_edit.Do},
	}},
	{Title: "Info bar", Filename: "info_bar.go", Func: info_bar.Do},
	{Title: "Links", Filename: "links.go", Func: links.Do},
	{Title: "Menus", Filename: "menus.go", Func: menus.Do},
	{Title: "Paned Widgets", Filename: "paned_widgets.go", Func: paned_widgets.Do},
	{Title: "Pickers", Filename: "pickers.go", Func: pickers.Do},
	{Title: "Pixbufs", Filename: "pixbufs.go", Func: pixbufs.Do},
	{Title: "Size Groups", Filename: "size_groups.go", Func: size_groups.Do},
	{Title: "Spinner", Filename: "spinner.go", Func: spinner.Do},
	{Title: "Stock Item and Icon Browser", Filename: "stock_browser.go", Func: stock_browser.Do},
	{Title: "Tree View", Children: []*DemoDesc{
		{Title: "List Store", Filename: "list_store.go", Func: list_store.Do},
	}},
}
