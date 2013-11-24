// Paned Widgets
//
// The GtkHPaned and GtkVPaned Widgets divide their content
// area into two panes with a divider in between that the
// user can adjust. A separate child is placed into each
// pane.
//
// There are a number of options that can be set for each pane.
// This test contains both a horizontal (HPaned) and a vertical
// (VPaned) widget, and allows you to adjust the options for
// each side of each widget.
package paned_widgets

import "gobject/gobject-2.0"
import "gobject/gtk-3.0"

var window *gtk.Window

func toggle(child *gtk.Widget, toggle_shrink, toggle_resize bool) {
	paned := gtk.ToPaned(child.GetParent())

	// we can do that because there's only one Go representation
	// of every gobject
	is_child1 := child == paned.GetChild1()

	var resize, shrink bool
	paned.ChildGet(child, "resize", &resize, "shrink", &shrink)
	if toggle_resize {
		resize = !resize
	}
	if toggle_shrink {
		shrink = !shrink
	}

	paned.Remove(child)
	if is_child1 {
		paned.Pack1(child, resize, shrink)
	} else {
		paned.Pack2(child, resize, shrink)
	}
}

func create_pane_options(paned *gtk.Paned, frame_label, label1, label2 string) *gtk.Widget {
	child1 := paned.GetChild1()
	child2 := paned.GetChild2()

	frame := gtk.NewFrame(frame_label)
	frame.SetBorderWidth(4)

	table := gtk.NewGrid()
	frame.Add(table)

	label := gtk.NewLabel(label1)
	table.Attach(label, 0, 0, 1, 1)

	check_button := gtk.NewCheckButtonWithMnemonic("_Resize")
	table.Attach(check_button, 0, 1, 1, 1)
	check_button.Connect("toggled", func() { toggle(child1, false, true) })

	check_button = gtk.NewCheckButtonWithMnemonic("_Shrink")
	table.Attach(check_button, 0, 2, 1, 1)
	check_button.SetActive(true)
	check_button.Connect("toggled", func() { toggle(child1, true, false) })

	label = gtk.NewLabel(label2)
	table.Attach(label, 1, 0, 1, 1)

	check_button = gtk.NewCheckButtonWithMnemonic("_Resize")
	table.Attach(check_button, 1, 1, 1, 1)
	check_button.SetActive(true)
	check_button.Connect("toggled", func() { toggle(child2, false, true) })

	check_button = gtk.NewCheckButtonWithMnemonic("_Shrink")
	table.Attach(check_button, 1, 2, 1, 1)
	check_button.SetActive(true)
	check_button.Connect("toggled", func() { toggle(child2, true, false) })

	return gtk.ToWidget(frame)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.Connect("destroy", func() { window = nil })
		window.SetTitle("Panes")
		window.SetBorderWidth(0)

		vbox := gtk.NewBox(gtk.OrientationVertical, 0)
		window.Add(vbox)

		vpaned := gtk.NewPaned(gtk.OrientationVertical)
		vbox.PackStart(vpaned, true, true, 0)
		vpaned.SetBorderWidth(5)

		hpaned := gtk.NewPaned(gtk.OrientationHorizontal)
		vpaned.Add1(hpaned)

		frame := gtk.NewFrame(gobject.NilString)
		frame.SetShadowType(gtk.ShadowTypeIn)
		frame.SetSizeRequest(60, 60)
		hpaned.Add1(frame)

		button := gtk.NewButtonWithMnemonic("_Hi there")
		frame.Add(button)

		frame = gtk.NewFrame(gobject.NilString)
		frame.SetShadowType(gtk.ShadowTypeIn)
		frame.SetSizeRequest(80, 60)
		hpaned.Add2(frame)

		frame = gtk.NewFrame(gobject.NilString)
		frame.SetShadowType(gtk.ShadowTypeIn)
		frame.SetSizeRequest(60, 80)
		vpaned.Add2(frame)

		vbox.PackStart(create_pane_options(hpaned, "Horizontal", "Left", "Right"),
			false, false, 0)
		vbox.PackStart(create_pane_options(vpaned, "Vertical", "Top", "Bottom"),
			false, false, 0)
		vbox.ShowAll()
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
