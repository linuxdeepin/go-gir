// Entry/Entry Buffer
//
// GtkEntryBuffer provides the text content in a GtkEntry.
package entry_buffer

import "gobject/gtk-3.0"
import "gobject/gobject-2.0"

var dialog *gtk.Dialog

func Do(mainwin *gtk.Window) *gtk.Window {
	if dialog == nil {
		dialog = gtk.NewDialogWithButtons("GtkEntryBuffer", mainwin, 0,
			gtk.StockClose, gtk.ResponseTypeNone)
		dialog.SetResizable(false)
		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.Connect("destroy", func() { dialog = nil })

		content_area := gtk.ToBox(dialog.GetContentArea())

		vbox := gtk.NewBox(gtk.OrientationVertical, 5)
		content_area.PackStart(vbox, true, true, 0)
		vbox.SetBorderWidth(5)

		label := gtk.NewLabel("Entries share a buffer. Typing in one is reflected in the other.")
		vbox.PackStart(label, false, false, 0)

		// Create a buffer
		buffer := gtk.NewEntryBuffer(gobject.NilString, 0)

		// Create our first entry
		entry := gtk.NewEntryWithBuffer(buffer)
		vbox.PackStart(entry, false, false, 0)

		// Create the second entry
		entry = gtk.NewEntryWithBuffer(buffer)
		entry.SetVisibility(false)
		vbox.PackStart(entry, false, false, 0)
	}

	if !dialog.GetVisible() {
		dialog.ShowAll()
	} else {
		dialog.Destroy()
	}

	return gtk.ToWindow(dialog)
}
