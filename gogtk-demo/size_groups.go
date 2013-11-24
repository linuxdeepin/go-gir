// Size Groups
//
// GtkSizeGroup provides a mechanism for grouping a number of
// widgets together so they all request the same amount of space.
// This is typically useful when you want a column of widgets to
// have the same size, but you can't use a GtkTable widget.
//
// Note that size groups only affect the amount of space requested,
// not the size that the widgets finally receive. If you want the
// widgets in a GtkSizeGroup to actually be the same size, you need
// to pack them in such a way that they get the size they request
// and not more. For example, if you are packing your widgets
// into a table, you would not include the GTK_FILL flag.
package size_groups

import "gobject/gtk-3.0"

var dialog *gtk.Dialog

// Convenience function to create a combo box holding a number of strings
func create_combo_box(strings []string) *gtk.ComboBoxText {
	combo_box := gtk.NewComboBoxText()
	for _, s := range strings {
		combo_box.AppendText(s)
	}
	combo_box.SetActive(0)
	return combo_box
}

func add_row(table *gtk.Grid, row int, size_group *gtk.SizeGroup, label_text string, options []string) {
	label := gtk.NewLabelWithMnemonic(label_text)
	label.SetHAlign(gtk.AlignStart)
	label.SetVAlign(gtk.AlignEnd)
	label.SetHExpand(true)
	table.Attach(label, 0, row, 1, 1)

	combo_box := create_combo_box(options)
	label.SetMnemonicWidget(combo_box)
	size_group.AddWidget(combo_box)
	table.Attach(combo_box, 1, row, 1, 1)
}

var color_options = []string{
	"Red", "Green", "Blue",
}

var dash_options = []string{
	"Solid", "Dashed", "Dotted",
}

var end_options = []string{
	"Square", "Round", "Arrow",
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if dialog == nil {
		dialog = gtk.NewDialogWithButtons("GtkSizeGroup", mainwin,
			0, gtk.StockClose, gtk.ResponseTypeNone)
		dialog.SetResizable(false)

		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.Connect("destroy", func() { dialog = nil })

		content_area := gtk.ToBox(dialog.GetContentArea())

		vbox := gtk.NewBox(gtk.OrientationVertical, 5)
		content_area.PackStart(vbox, true, true, 0)
		vbox.SetBorderWidth(5)

		size_group := gtk.NewSizeGroup(gtk.SizeGroupModeHorizontal)

		// Create one frame holding color options
		frame := gtk.NewFrame("Color Options")
		vbox.PackStart(frame, true, true, 0)

		table := gtk.NewGrid()
		table.SetBorderWidth(5)
		table.SetRowSpacing(5)
		table.SetColumnSpacing(10)
		frame.Add(table)

		add_row(table, 0, size_group, "_Foreground", color_options)
		add_row(table, 1, size_group, "_Background", color_options)

		// And another frame holding line style options
		frame = gtk.NewFrame("Line Options")
		vbox.PackStart(frame, false, false, 0)

		table = gtk.NewGrid()
		table.SetBorderWidth(5)
		table.SetRowSpacing(5)
		table.SetColumnSpacing(10)
		frame.Add(table)

		add_row(table, 0, size_group, "_Dashing", dash_options)
		add_row(table, 1, size_group, "_Line ends", end_options)

		// And a check button to turn grouping on and off
		check_button := gtk.NewCheckButtonWithMnemonic("_Enable grouping")
		vbox.PackStart(check_button, false, false, 0)

		check_button.SetActive(true)
		check_button.Connect("toggled", func() {
			var new_mode gtk.SizeGroupMode

			// GTK_SIZE_GROUP_NONE is not generally useful, but is useful
			// here to show the effect of GTK_SIZE_GROUP_HORIZONTAL by
			// contrast.
			if check_button.GetActive() {
				new_mode = gtk.SizeGroupModeHorizontal
			} else {
				new_mode = gtk.SizeGroupModeNone
			}

			size_group.SetMode(new_mode)
		})
	}

	if !dialog.GetVisible() {
		dialog.ShowAll()
	} else {
		dialog.Destroy()
	}
	return gtk.ToWindow(dialog)
}
