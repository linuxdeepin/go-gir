// Tree View/List Store
//
// The GtkListStore is used to store data in list form, to be used
// later on by a GtkTreeView to display it. This demo builds a
// simple GtkListStore and displays it. See the Stock Browser
// demo for a more advanced example.
package list_store

import (
	"gobject/gdk-3.0"
	"gobject/gobject-2.0"
	"gobject/gtk-3.0"
	"time"
)

var window *gtk.Window

type bug struct {
	fixed       bool
	number      int
	severity    string
	description string
}

const (
	column_fixed = iota
	column_number
	column_severity
	column_description
	column_pulse
	column_icon
	column_active
	column_sensitive
)

var data = []bug{
	{false, 60482, "Normal", "scrollable notebooks and hidden tabs"},
	{false, 60620, "Critical", "gdk_window_clear_area (gdkwindow-win32.c) is not thread-safe"},
	{false, 50214, "Major", "Xft support does not clean up correctly"},
	{true, 52877, "Major", "GtkFileSelection needs a refresh method. "},
	{false, 56070, "Normal", "Can't click button after setting in sensitive"},
	{true, 56355, "Normal", "GtkLabel - Not all changes propagate correctly"},
	{false, 50055, "Normal", "Rework width/height computations for TreeView"},
	{false, 58278, "Normal", "gtk_dialog_set_response_sensitive () doesn't work"},
	{false, 55767, "Normal", "Getters for all setters"},
	{false, 56925, "Normal", "Gtkcalender size"},
	{false, 56221, "Normal", "Selectable label needs right-click copy menu"},
	{true, 50939, "Normal", "Add shift clicking to GtkTextView"},
	{false, 6112, "Enhancement", "netscape-like collapsable toolbars"},
	{false, 1, "Normal", "First bug :=)"},
}

func pulse(ticker *time.Ticker, cancel chan int, list_store *gtk.ListStore) {
	for {
		select {
		case <-ticker.C:
			gdk.ThreadsEnter()

			var pulse int
			iter, _ := list_store.GetIterFirst()
			list_store.Get(&iter, column_pulse, &pulse)

			if pulse == 99999 {
				pulse = 0
			} else {
				pulse++
			}

			list_store.Set(&iter,
				column_pulse, pulse,
				column_active, true)

			gdk.ThreadsLeave()
		case <-cancel:
			return
		}
	}
}

func create_model() *gtk.ListStore {
	store := gtk.NewListStore(gobject.Boolean,
		gobject.Int,
		gobject.String,
		gobject.String,
		gobject.Int,
		gobject.String,
		gobject.Boolean,
		gobject.Boolean)

	for i := range data {
		var icon_name string
		var sensitive = true

		if i == 1 || i == 3 {
			icon_name = "battery-caution-charging-symbolic"
		}

		if i == 3 {
			sensitive = false
		}

		store.Append(data[i].fixed,
			data[i].number,
			data[i].severity,
			data[i].description,
			0,
			icon_name,
			false,
			sensitive)
	}

	return store
}

func add_columns(treeview *gtk.TreeView, model *gtk.ListStore) {
	var r gtk.CellRendererLike
	var c *gtk.TreeViewColumn

	// column for fixed toggles
	toggle := gtk.NewCellRendererToggle()
	toggle.Connect("toggled", func(cell *gtk.CellRendererToggle, path_str string) {
		path := gtk.NewTreePathFromString(path_str)
		iter, _ := model.GetIter(path)
		var checked bool
		model.Get(&iter, column_fixed, &checked)
		model.Set(&iter, column_fixed, !checked)
	})
	c = gtk.NewTreeViewColumnWithAttributes("Fixed?", toggle, "active", column_fixed)

	// set this column to a fixed sizing (of 50 pixels)
	c.SetSizing(gtk.TreeViewColumnSizingFixed)
	c.SetFixedWidth(50)
	treeview.AppendColumn(c)

	// column for bug numbers
	r = gtk.NewCellRendererText()
	c = gtk.NewTreeViewColumnWithAttributes("Bug number", r, "text", column_number)
	c.SetSortColumnID(column_number)
	treeview.AppendColumn(c)

	// column for severities
	r = gtk.NewCellRendererText()
	c = gtk.NewTreeViewColumnWithAttributes("Severity", r, "text", column_severity)
	c.SetSortColumnID(column_severity)
	treeview.AppendColumn(c)

	// column for description
	r = gtk.NewCellRendererText()
	c = gtk.NewTreeViewColumnWithAttributes("Description", r, "text", column_description)
	c.SetSortColumnID(column_description)
	treeview.AppendColumn(c)

	// column for spinner
	r = gtk.NewCellRendererSpinner()
	c = gtk.NewTreeViewColumnWithAttributes("Spinning", r,
		"pulse", column_pulse,
		"active", column_active)
	c.SetSortColumnID(column_pulse)
	treeview.AppendColumn(c)

	// column for symbolic icon
	pixbuf := gtk.NewCellRendererPixbuf()
	pixbuf.SetProperty("follow-state", true)
	c = gtk.NewTreeViewColumnWithAttributes("Symbolic icon", pixbuf,
		"icon-name", column_icon,
		"sensitive", column_sensitive)
	c.SetSortColumnID(column_icon)
	treeview.AppendColumn(c)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		ticker := time.NewTicker(80 * time.Millisecond)
		cancel := make(chan int)

		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetTitle("gtk.ListStore demo")
		window.Connect("destroy", func() { window = nil })
		window.SetBorderWidth(8)

		vbox := gtk.NewBox(gtk.OrientationVertical, 8)
		window.Add(vbox)

		label := gtk.NewLabel("This is the bug list (note: not based on real data, it would be nice to have a nice ODBC interface to bugzilla or so, though).")
		vbox.PackStart(label, false, false, 0)

		sw := gtk.NewScrolledWindow(nil, nil)
		sw.SetShadowType(gtk.ShadowTypeEtchedIn)
		sw.SetPolicy(gtk.PolicyTypeNever, gtk.PolicyTypeAutomatic)
		vbox.PackStart(sw, true, true, 0)

		// create tree model
		model := create_model()

		// create tree view
		treeview := gtk.NewTreeViewWithModel(model)
		treeview.SetRulesHint(true)
		treeview.SetSearchColumn(column_description)
		sw.Add(treeview)

		// add column to the tree view
		add_columns(treeview, model)

		// finish & show
		window.SetDefaultSize(280, 250)
		window.Connect("delete-event", func() {
			ticker.Stop()
			cancel <- 0
		})

		go pulse(ticker, cancel, model)
	}
	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
