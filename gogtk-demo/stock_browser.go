// Stock Item and Icon Browser
//
// This source code for this demo doesn't demonstrate anything
// particularly useful in applications. The purpose of the "demo" is
// just to provide a handy place to browse the available stock icons
// and stock items.
package stock_browser

import "gobject/gtk-3.0"
import "gobject/gobject-2.0"
import "gobject/gdkpixbuf-2.0"
import "fmt"
import "sort"
import "bytes"
import "strings"

var window *gtk.Window

type stock_item_info struct {
	id         string
	item       gtk.StockItem
	small_icon *gdkpixbuf.Pixbuf
	constant   string
	accel_str  string
}

type stock_item_display struct {
	type_label        *gtk.Label
	constant_label    *gtk.Label
	id_label          *gtk.Label
	label_accel_label *gtk.Label
	icon_image        *gtk.Image
}

func id_to_constant(id string) string {
	var out bytes.Buffer
	if strings.HasPrefix(id, "gtk-") {
		out.WriteString("gtk.Stock")
		id = id[4:]
	}

	for _, word := range strings.Split(id, "-") {
		word = strings.ToLower(word)
		if len(word) == 1 {
			out.WriteString(strings.ToUpper(word))
		} else {
			out.WriteString(strings.ToUpper(word[:1]))
			out.WriteString(word[1:])
		}
	}

	return out.String()
}

func create_model() *gtk.TreeModel {
	store := gtk.NewListStore(gobject.GoInterface, gobject.String)

	ids := gtk.StockListIDs()
	sort.Strings(ids)

	for _, id := range ids {
		var info stock_item_info

		info.id = id
		if item, ok := gtk.StockLookup(id); ok {
			info.item = item
		} else {
			info.item.Label = gobject.NilString
			info.item.StockID = gobject.NilString
			info.item.Modifier = 0
			info.item.Keyval = 0
			info.item.TranslationDomain = gobject.NilString
		}

		// only show icons for stock IDs that have default icons
		icon_set := gtk.IconFactoryLookupDefault(info.id)
		if icon_set != nil {
			// See what sizes this stock icon really exists at
			sizes := icon_set.GetSizes()

			// Use menu size if it exists, otherwise first size found
			size := sizes[0]
			for _, s := range sizes {
				if gtk.IconSize(s) == gtk.IconSizeMenu {
					size = int(gtk.IconSizeMenu)
					break
				}
			}

			info.small_icon = window.RenderIconPixbuf(info.id, size)
			if gtk.IconSize(size) != gtk.IconSizeMenu {
				// Make the result the proper size for our thumbnail
				w, h, _ := gtk.IconSizeLookup(int(gtk.IconSizeMenu))
				scaled := info.small_icon.ScaleSimple(w, h, gdkpixbuf.InterpTypeBilinear)
				info.small_icon = scaled
			}
		} else {
			info.small_icon = nil
		}

		if info.item.Keyval != 0 {
			info.accel_str = gtk.AcceleratorName(
				int(info.item.Keyval), info.item.Modifier)
		} else {
			info.accel_str = ""
		}

		info.constant = id_to_constant(info.id)
		store.Append(&info, info.id)
	}

	return gtk.ToTreeModel(store)
}

// Finds the largest size at which the given image stock id is
// available. This would not be useful for a normal application
func get_largest_size(id string) int {
	set := gtk.IconFactoryLookupDefault(id)
	sizes := set.GetSizes()
	best_size := int(gtk.IconSizeInvalid)
	best_pixels := 0

	for _, size := range sizes {
		w, h, _ := gtk.IconSizeLookup(size)
		if w*h > best_pixels {
			best_size = size
			best_pixels = w * h
		}
	}

	return best_size
}

func selection_changed(selection *gtk.TreeSelection, display *stock_item_display) {
	if model, iter, ok := selection.GetSelected(); ok {
		var infoi interface{}
		model.Get(&iter, 0, &infoi)
		info := infoi.(*stock_item_info)

		switch {
		case info.small_icon != nil && info.item.Label != gobject.NilString:
			display.type_label.SetText("Icon and Item")
		case info.small_icon != nil:
			display.type_label.SetText("Icon Only")
		case info.item.Label != gobject.NilString:
			display.type_label.SetText("Item Only")
		default:
			display.type_label.SetText("???????")
		}

		display.constant_label.SetText(info.constant)
		display.id_label.SetText(info.id)

		if info.item.Label != gobject.NilString {
			str := fmt.Sprintf("%s %s", info.item.Label, info.accel_str)
			display.label_accel_label.SetTextWithMnemonic(str)
		} else {
			display.label_accel_label.SetText("")
		}

		if info.small_icon != nil {
			display.icon_image.SetFromStock(info.id, get_largest_size(info.id))
		} else {
			display.icon_image.Clear()
		}
	} else {
		display.type_label.SetText("No selected item")
		display.constant_label.SetText("")
		display.id_label.SetText("")
		display.label_accel_label.SetText("")
		display.icon_image.Clear()
	}
}

func constant_set_func_text(tree_column *gtk.TreeViewColumn, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	var infoi interface{}
	model.Get(iter, 0, &infoi)

	info := infoi.(*stock_item_info)
	cell.SetProperty("text", info.constant)
}

func id_set_func(tree_column *gtk.TreeViewColumn, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	var infoi interface{}
	model.Get(iter, 0, &infoi)

	info := infoi.(*stock_item_info)
	cell.SetProperty("text", info.id)
}

func accel_set_func(tree_column *gtk.TreeViewColumn, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	var infoi interface{}
	model.Get(iter, 0, &infoi)

	info := infoi.(*stock_item_info)
	cell.SetProperty("text", info.accel_str)
}

func label_set_func(tree_column *gtk.TreeViewColumn, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	var infoi interface{}
	model.Get(iter, 0, &infoi)

	info := infoi.(*stock_item_info)
	cell.SetProperty("text", info.item.Label)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Stock Icons and Items")
		window.SetDefaultSize(-1, 500)
		window.Connect("destroy", func() { window = nil })
		window.SetBorderWidth(8)

		hbox := gtk.NewBox(gtk.OrientationHorizontal, 8)
		window.Add(hbox)

		sw := gtk.NewScrolledWindow(nil, nil)
		sw.SetPolicy(gtk.PolicyTypeNever, gtk.PolicyTypeAutomatic)
		hbox.PackStart(sw, false, false, 0)

		model := create_model()
		treeview := gtk.NewTreeViewWithModel(model)

		sw.Add(treeview)

		var r *gtk.CellRenderer
		c := gtk.NewTreeViewColumn()
		c.SetTitle("Constant")

		r = gtk.ToCellRenderer(gtk.NewCellRendererPixbuf())
		c.PackStart(r, false)
		c.SetAttributes(r, "stock-id", 1)

		r = gtk.ToCellRenderer(gtk.NewCellRendererText())
		c.PackStart(r, true)
		c.SetCellDataFunc(r, constant_set_func_text)

		treeview.AppendColumn(c)

		r = gtk.ToCellRenderer(gtk.NewCellRendererText())
		treeview.InsertColumnWithDataFunc(-1, "Label", r, label_set_func)

		r = gtk.ToCellRenderer(gtk.NewCellRendererText())
		treeview.InsertColumnWithDataFunc(-1, "Accel", r, accel_set_func)

		r = gtk.ToCellRenderer(gtk.NewCellRendererText())
		treeview.InsertColumnWithDataFunc(-1, "ID", r, id_set_func)

		frame := gtk.NewFrame("Selected Item")
		frame.SetVAlign(gtk.AlignStart)
		hbox.PackEnd(frame, false, false, 0)

		vbox := gtk.NewBox(gtk.OrientationVertical, 8)
		vbox.SetBorderWidth(4)
		frame.Add(vbox)

		display := &stock_item_display{
			type_label:        gtk.NewLabel(gobject.NilString),
			constant_label:    gtk.NewLabel(gobject.NilString),
			id_label:          gtk.NewLabel(gobject.NilString),
			label_accel_label: gtk.NewLabel(gobject.NilString),
			icon_image:        gtk.NewImage(), // empty image
		}

		vbox.PackStart(display.type_label, false, false, 0)
		vbox.PackStart(display.icon_image, false, false, 0)
		vbox.PackStart(display.label_accel_label, false, false, 0)
		vbox.PackStart(display.constant_label, false, false, 0)
		vbox.PackStart(display.id_label, false, false, 0)

		selection := treeview.GetSelection()
		selection.SetMode(gtk.SelectionModeSingle)

		selection.Connect("changed", func(sel *gtk.TreeSelection) {
			selection_changed(sel, display)
		})
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
