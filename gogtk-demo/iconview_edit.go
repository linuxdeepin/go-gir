// Icon View/Editing and Drag-and-Drop
//
// The GtkIconView widget supports Editing and Drag-and-Drop.
// This example also demonstrates using the generic GtkCellLayout
// interface to set up cell renderers in an icon view.
package iconview_edit

import "gobject/gtk-3.0"
import "gobject/gdk-3.0"
import "gobject/gobject-2.0"
import "gobject/gdkpixbuf-2.0"

var window *gtk.Window

func create_store() *gtk.ListStore {
	store := gtk.NewListStore(gobject.String)
	for _, item := range []string{"Red", "Green", "Blue", "Yellow"} {
		store.Append(item)
	}
	return store
}

func set_cell_color(cell_layout *gtk.CellLayout, cell *gtk.CellRenderer, model *gtk.TreeModel, iter *gtk.TreeIter) {
	var text string
	var pixel int

	model.Get(iter, 0, &text)
	if color, ok := gdk.ColorParse(text); ok {
		pixel = int(color.Red   >> 8) << 24 |
			int(color.Green >> 8) << 16 |
			int(color.Blue  >> 8) << 8
	}

	pixbuf := gdkpixbuf.NewPixbuf(gdkpixbuf.ColorspaceRGB, false, 8, 24, 24)
	pixbuf.Fill(pixel)

	cell.SetProperty("pixbuf", pixbuf)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Editing and Drag-and-Drop")
		window.Connect("destroy", func() { window = nil })

		store := create_store()
		icon_view := gtk.NewIconViewWithModel(store)

		icon_view.SetSelectionMode(gtk.SelectionModeSingle)
		icon_view.SetItemOrientation(gtk.OrientationHorizontal)
		icon_view.SetColumns(2)
		icon_view.SetReorderable(true)

		r := gtk.ToCellRenderer(gtk.NewCellRendererPixbuf())
		icon_view.PackStart(r, true)
		icon_view.SetCellDataFunc(r, set_cell_color)

		r = gtk.ToCellRenderer(gtk.NewCellRendererText())
		icon_view.PackStart(r, true)
		r.SetProperty("editable", true)
		r.Connect("edited", func(cell *gtk.CellRendererText, path, text string) {
			p := gtk.NewTreePathFromString(path)
			iter, _ := store.GetIter(p)
			store.Set(&iter, 0, text)
			p.Free()
		})
		icon_view.SetAttributes(r, "text", 0)

		window.Add(icon_view)
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}