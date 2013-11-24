// Combo boxes
//
// The ComboBox widget allows to select one option out of a list.
// The ComboBoxEntry additionally allows the user to enter a value
// that is not in the list of options.
//
// How the options are displayed is controlled by cell renderers.
package combo_boxes

import "gobject/gobject-2.0"
import "gobject/gtk-3.0"
import "gobject/gdk-3.0"
import "gobject/gdkpixbuf-2.0"
import "strings"
import "regexp"

var window *gtk.Window

const (
	pixbuf_column = iota
	text_column
)

func create_stock_icon_store() *gtk.TreeModel {
	stock_ids := []string{
		gtk.StockDialogWarning,
		gtk.StockStop,
		gtk.StockNew,
		gtk.StockClear,
		gobject.NilString,
		gtk.StockOpen,
	}


	cellview := gtk.NewCellView()
	store := gtk.NewListStore(gdkpixbuf.PixbufGetType(), gobject.String)

	for _, id := range stock_ids {
		if id != gobject.NilString {
			pixbuf := cellview.RenderIconPixbuf(id, int(gtk.IconSizeButton))
			item, _ := gtk.StockLookup(id)
			label := strings.Replace(item.Label, "_", "", -1)
			store.Append(pixbuf, label)
		} else {
			store.Append((*gdkpixbuf.Pixbuf)(nil), "separator")
		}
	}

	return gtk.ToTreeModel(store)
}

// A GtkCellLayoutDataFunc that demonstrates how one can control
// sensitivity of rows. This particular function does nothing
// useful and just makes the second row insensitive.
func set_sensitive(cell_layout *gtk.CellLayout, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	path := model.GetPath(iter)
	_, indices := path.GetIndices()
	sensitive := indices[0] != 1
	cell.SetProperty("sensitive", sensitive)
}

// A GtkTreeViewRowSeparatorFunc that demonstrates how rows can be
// rendered as separators. This particular function does nothing
// useful and just turns the fourth row into a separator.
func is_separator(model *gtk.TreeModel, iter *gtk.TreeIter) bool {
	path := model.GetPath(iter)
	_, indices := path.GetIndices()
	result := indices[0] == 4
	return result
}

func create_capital_store() *gtk.TreeModel {
	capitals := []struct{
		name string
		capitals []string
	}{
		{"A - B", []string{
			"Albany",
			"Annapolis",
			"Atlanta",
			"Augusta",
			"Austin",
			"Baton Rouge",
			"Bismarck",
			"Boise",
			"Boston",
		}},
		{"C - D", []string{
			"Carson City",
			"Charleston",
			"Cheyenne",
			"Columbia",
			"Columbus",
			"Concord",
			"Denver",
			"Des Moines",
			"Dover",
		}},
		{"E - J", []string{
			"Frankfort",
			"Harrisburg",
			"Hartford",
			"Helena",
			"Honolulu",
			"Indianapolis",
			"Jackson",
			"Jefferson City",
			"Juneau",
		}},
		{"K - O", []string{
			"Lansing",
			"Lincoln",
			"Little Rock",
			"Madison",
			"Montgomery",
			"Montpelier",
			"Nashville",
			"Oklahoma City",
			"Olympia",
		}},
		{"P - S", []string{
			"Phoenix",
			"Pierre",
			"Providence",
			"Raleigh",
			"Richmond",
			"Sacramento",
			"Salem",
			"Salt Lake City",
			"Santa Fe",
			"Springfield",
			"St. Paul",
		}},
		{"T - Z", []string{
			"Tallahassee",
			"Topeka",
			"Trenton",
		}},
	}

	store := gtk.NewTreeStore(gobject.String)
	for _, group := range capitals {
		iter := store.Append(nil, group.name)
		for _, capital := range group.capitals {
			store.Append(&iter, capital)
		}
	}
	return gtk.ToTreeModel(store)
}

func is_capital_sensitive(cell_layout *gtk.CellLayout, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	sensitive := !model.IterHasChild(iter)
	cell.SetProperty("sensitive", sensitive)
}

func fill_combo_entry(combo *gtk.ComboBoxText) {
	combo.AppendText("One")
	combo.AppendText("Two")
	combo.AppendText("2½")
	combo.AppendText("Three")
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Combo boxes")
		window.Connect("destroy", func() { window = nil })

		window.SetBorderWidth(10)

		vbox := gtk.NewBox(gtk.OrientationVertical, 2)
		window.Add(vbox)

		// A combobox demonstrating cell renderers, separators and
		// insensitive rows
		frame := gtk.NewFrame("Some stock icons")
		vbox.PackStart(frame, false, false, 0)

		box := gtk.NewBox(gtk.OrientationVertical, 0)
		box.SetBorderWidth(5)
		frame.Add(box)

		model := create_stock_icon_store()
		combo := gtk.NewComboBoxWithModel(model)
		box.Add(combo)

		rp := gtk.NewCellRendererPixbuf()
		combo.PackStart(rp, false)
		combo.SetAttributes(rp, "pixbuf", pixbuf_column)
		combo.SetCellDataFunc(rp, set_sensitive)

		rt := gtk.NewCellRendererText()
		combo.PackStart(rt, true)
		combo.SetAttributes(rt, "text", text_column)
		combo.SetCellDataFunc(rt, set_sensitive)

		combo.SetRowSeparatorFunc(is_separator)
		combo.SetActive(0)

		// A combobox demonstrating trees.
		frame = gtk.NewFrame("Where are we ?")
		vbox.PackStart(frame, false, false, 0)

		box = gtk.NewBox(gtk.OrientationVertical, 0)
		box.SetBorderWidth(5)
		frame.Add(box)

		model = create_capital_store()
		combo = gtk.NewComboBoxWithModel(model)
		box.Add(combo)

		rt = gtk.NewCellRendererText()
		combo.PackStart(rt, true)
		combo.SetAttributes(rt, "text", 0)
		combo.SetCellDataFunc(rt, is_capital_sensitive)

		path := gtk.NewTreePathFromIndices(0, 8)
		iter, _ := model.GetIter(path)
		path.Free()
		combo.SetActiveIter(&iter)

		// A GtkComboBoxEntry with validation.
		frame = gtk.NewFrame("Editable")
		vbox.PackStart(frame, false, false, 0)

		box = gtk.NewBox(gtk.OrientationVertical, 0)
		box.SetBorderWidth(5)
		frame.Add(box)

		combotext := gtk.NewComboBoxTextWithEntry()
		fill_combo_entry(combotext)
		box.Add(combotext)

		// A simple validating entry
		entry := gtk.ToEntry(combotext.GetChild())
		entry.Connect("changed", func(entry *gtk.Entry) {
			error_color := gdk.RGBA{0.8, 0, 0, 1}
			text := entry.GetText()
			re := regexp.MustCompile(`^([0-9]*|One|Two|2½|Three)$`)
			if re.MatchString(text) {
				entry.OverrideColor(0, nil)
			} else {
				entry.OverrideColor(0, &error_color)
			}
		})

		// A combobox with string IDs
		frame = gtk.NewFrame("Editable")
		vbox.PackStart(frame, false, false, 0)

		box = gtk.NewBox(gtk.OrientationVertical, 0)
		box.SetBorderWidth(5)
		frame.Add(box)

		combotext = gtk.NewComboBoxText()
		combotext.Append("never", "Not visible")
		combotext.Append("when-active", "Visible when active")
		combotext.Append("always", "Always visible")
		box.Add(combotext)

		entry = gtk.NewEntry()
		gobject.ObjectBindProperty(
			combotext, "active-id",
			entry, "text",
			gobject.BindingFlagsBidirectional)
		box.Add(entry)

	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}