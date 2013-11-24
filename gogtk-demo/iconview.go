// Icon View/Icon View Basics
//
// The GtkIconView widget is used to display and manipulate icons.
// It uses a GtkTreeModel for data storage, so the list store
// example might be helpful.
package iconview

import "gobject/gtk-3.0"
import "gobject/gobject-2.0"
import "gobject/gdkpixbuf-2.0"
import "strings"
import "path/filepath"
import "os"
import "./gogtk-demo/common"

var window *gtk.Window
var file_pixbuf *gdkpixbuf.Pixbuf
var folder_pixbuf *gdkpixbuf.Pixbuf
var parent string
var up_button *gtk.ToolItem

const (
	col_path = iota
	col_display_name
	col_pixbuf
	col_is_directory
)

const folder_icon_filename = "gnome-fs-directory.png"
const file_icon_filename = "gnome-fs-regular.png"

// Loads the images for the demo and returns whether the operation succeeded
func load_pixbufs() error {
	if file_pixbuf != nil {
		return nil // already loaded earlier
	}

	var err error
	file_pixbuf, err = gdkpixbuf.NewPixbufFromFile(common.FindFile(file_icon_filename))
	if err != nil {
		return err
	}

	folder_pixbuf, err = gdkpixbuf.NewPixbufFromFile(common.FindFile(folder_icon_filename))
	if err != nil {
		return err
	}

	return nil
}

func sort_func(model *gtk.TreeModel, a, b *gtk.TreeIter) int {
	var is_dir_a, is_dir_b bool
	var name_a, name_b string

	model.Get(a,
		col_is_directory, &is_dir_a,
		col_display_name, &name_a)

	model.Get(b,
		col_is_directory, &is_dir_b,
		col_display_name, &name_b)

	name_a = strings.ToLower(name_a)
	name_b = strings.ToLower(name_b)

	if !is_dir_a && is_dir_b {
		return 1
	} else if is_dir_a && !is_dir_b {
		return -1
	} else {
		switch {
		case name_a < name_b:
			return -1
		case name_a == name_b:
			return 0
		case name_a > name_b:
			return 1
		}
	}

	return 0
}

func create_store() *gtk.ListStore {
	store := gtk.NewListStore(
		gobject.String,
		gobject.String,
		gdkpixbuf.PixbufGetType(),
		gobject.Boolean)

	// Set sort column and function
	store.SetDefaultSortFunc(sort_func)
	store.SetSortColumnID(-1, gtk.SortTypeAscending)

	return store
}

func fill_store(store *gtk.ListStore) {
	// temporarily disable sorting
	store.SetSortColumnID(-2, gtk.SortTypeAscending)
	defer store.SetSortColumnID(-1, gtk.SortTypeAscending)

	store.Clear()

	dir, err := os.Open(parent)
	if err != nil {
		println(err.Error())
		return
	}
	defer dir.Close()

	entries, err := dir.Readdir(0)
	if err != nil {
		return
	}

	for _, entry := range entries {
		display_name := entry.Name()
		if strings.HasPrefix(display_name, ".") {
			// We ignore hidden files that start with a '.'
			continue
		}

		path := filepath.Join(parent, display_name)
		is_dir := entry.IsDir()

		var pixbuf *gdkpixbuf.Pixbuf
		if is_dir {
			pixbuf = folder_pixbuf
		} else {
			pixbuf = file_pixbuf
		}

		store.Append(path, display_name, pixbuf, is_dir)
	}
}

func item_activated(icon_view *gtk.IconView, tree_path *gtk.TreePath) {
	store := gtk.ToListStore(icon_view.GetModel())
	iter, _ := store.GetIter(tree_path)

	var is_dir bool
	var path string
	store.Get(&iter,
		col_path, &path,
		col_is_directory, &is_dir)

	if !is_dir {
		return
	}

	// Replace parent with path and re-fill the model
	parent = path
	fill_store(store)

	// Sensitize the up button
	up_button.SetSensitive(true)
}

func up_clicked(store *gtk.ListStore) {
	parent = filepath.Dir(parent)
	fill_store(store)

	// Maybe de-sensitize the up button
	up_button.SetSensitive(parent != "/")
}

func home_clicked(store *gtk.ListStore) {
	parent = os.Getenv("HOME")
	fill_store(store)

	// Sensitize the up button
	up_button.SetSensitive(true)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetDefaultSize(650, 400)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("GtkIconView demo")
		window.Connect("destroy", func() {
			window = nil
			up_button = nil
		})

		if err := load_pixbufs(); err != nil {
			dialog := gtk.NewMessageDialog(window, gtk.DialogFlagsDestroyWithParent,
				gtk.MessageTypeError, gtk.ButtonsTypeClose,
				"Failed to load an image: %s", err)
			dialog.Connect("response", func() { dialog.Destroy() })
			dialog.Show()
			goto done
		}

		vbox := gtk.NewBox(gtk.OrientationVertical, 0)
		window.Add(vbox)

		tool_bar := gtk.NewToolbar()
		vbox.PackStart(tool_bar, false, false, 0)

		up_button = gtk.ToToolItem(gtk.NewToolButtonFromStock(gtk.StockGoUp))
		up_button.SetIsImportant(true)
		up_button.SetSensitive(false)
		tool_bar.Insert(up_button, -1)

		home_button := gtk.NewToolButtonFromStock(gtk.StockHome)
		home_button.SetIsImportant(true)
		tool_bar.Insert(home_button, -1)

		sw := gtk.NewScrolledWindow(nil, nil)
		sw.SetShadowType(gtk.ShadowTypeEtchedIn)
		sw.SetPolicy(gtk.PolicyTypeAutomatic, gtk.PolicyTypeAutomatic)
		vbox.PackStart(sw, true, true, 0)

		// Create the store and fill it with the contents of '/'
		parent = "/"
		store := create_store()
		fill_store(store)

		icon_view := gtk.NewIconViewWithModel(store)
		icon_view.SetSelectionMode(gtk.SelectionModeMultiple)

		// Connect to the "clicked" signal of the "Up" tool button
		up_button.Connect("clicked", func() { up_clicked(store) })

		// Connect to the "clicked" signal of the "Home" tool button
		home_button.Connect("clicked", func() { home_clicked(store) })

		// We now set which model columns that correspond to the text
		// and pixbuf of each item
		icon_view.SetTextColumn(col_display_name)
		icon_view.SetPixbufColumn(col_pixbuf)

		// Connect to the "item-activated" signal
		icon_view.Connect("item-activated", item_activated)
		sw.Add(icon_view)

		icon_view.GrabFocus()
	}

done:
	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}