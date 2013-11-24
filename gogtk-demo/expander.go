// Expander
//
// GtkExpander allows to provide additional content that is initially hidden.
// This is also known as "disclosure triangle".
package expander

import "gobject/gtk-3.0"

var dialog *gtk.Dialog

func Do(mainwin *gtk.Window) *gtk.Window {
	if dialog == nil {
		dialog = gtk.NewDialogWithButtons("GtkExpander", mainwin, 0,
			gtk.StockClose, gtk.ResponseTypeNone)
		dialog.SetResizable(false)
		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.Connect("destroy", func() { dialog = nil })

		content_area := gtk.ToBox(dialog.GetContentArea())

		vbox := gtk.NewBox(gtk.OrientationVertical, 5)
		content_area.PackStart(vbox, true, true, 0)
		vbox.SetBorderWidth(5)

		label := gtk.NewLabel("Expander demo. Click on the triangle for details.")
		vbox.PackStart(label, false, false, 0)

		// Create the expander
		expander := gtk.NewExpander("Details")
		vbox.PackStart(expander, false, false, 0)
		label = gtk.NewLabel("Details can be shown or hidden.")
		expander.Add(label)
	}

	if !dialog.GetVisible() {
		dialog.ShowAll()
	} else {
		dialog.Destroy()
	}

	return gtk.ToWindow(dialog)
}
