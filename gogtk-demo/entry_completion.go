// Entry/Entry Completion
//
// GtkEntryCompletion provides a mechanism for adding support for
// completion in GtkEntry.
package entry_completion

import "gobject/gtk-3.0"
import "gobject/gobject-2.0"

func create_completion_model() *gtk.TreeModel {
	store := gtk.NewListStore(gobject.String)

	// Append few words
	store.Append("GNOME")
	store.Append("total")
	store.Append("totally")
	return gtk.ToTreeModel(store)
}

var dialog *gtk.Dialog

func Do(mainwin *gtk.Window) *gtk.Window {
	if dialog == nil {
		dialog = gtk.NewDialogWithButtons("GtkEntryCompletion",
			mainwin, 0, gtk.StockClose, gtk.ResponseTypeNone)
		dialog.SetResizable(false)
		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.Connect("destroy", func() { dialog = nil })

		content_area := gtk.ToBox(dialog.GetContentArea())

		vbox := gtk.NewBox(gtk.OrientationVertical, 5)
		content_area.PackStart(vbox, true, true, 0)
		vbox.SetBorderWidth(5)

		label := gtk.NewLabel(gobject.NilString)
		label.SetMarkup("Completion demo, try writing <b>total</b> or <b>gnome</b> for example.")
		vbox.PackStart(label, false, false, 0)

		// Create our entry
		entry := gtk.NewEntry()
		vbox.PackStart(entry, false, false, 0)

		// Create the completion object
		completion := gtk.NewEntryCompletion()

		// Assign the completion to the entry
		entry.SetCompletion(completion)

		// Create a tree model and use it as the completion model
		completion_model := create_completion_model()
		completion.SetModel(completion_model)

		// Use model column 0 as the text column
		completion.SetTextColumn(0)
	}

	if !dialog.GetVisible() {
		dialog.ShowAll()
	} else {
		dialog.Destroy()
	}
	return gtk.ToWindow(dialog)
}