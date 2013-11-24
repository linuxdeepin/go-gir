// Spinner
//
// GtkSpinner allows to show that background activity is on-going.
package spinner

import "gobject/gtk-3.0"

var dialog *gtk.Dialog

func Do(mainwin *gtk.Window) *gtk.Window {
	if dialog == nil {
		dialog = gtk.NewDialogWithButtons("GtkSpinner", mainwin, 0,
			gtk.StockClose, gtk.ResponseTypeNone)

		dialog.SetResizable(false)
		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.Connect("destroy", func() { dialog = nil })

		content_area := gtk.ToBox(dialog.GetContentArea())
		vbox := gtk.NewBox(gtk.OrientationVertical, 5)
		content_area.PackStart(vbox, true, true, 0)
		vbox.SetBorderWidth(5)

		var hbox *gtk.Box

		// Sensitive
		hbox = gtk.NewBox(gtk.OrientationHorizontal, 5)
		spinner_sensitive := gtk.NewSpinner()
		hbox.Add(spinner_sensitive)
		hbox.Add(gtk.NewEntry())
		vbox.Add(hbox)

		// Disabled
		hbox = gtk.NewBox(gtk.OrientationHorizontal, 5)
		spinner_unsensitive := gtk.NewSpinner()
		hbox.Add(spinner_unsensitive)
		hbox.Add(gtk.NewEntry())
		vbox.Add(hbox)
		hbox.SetSensitive(false)

		button := gtk.NewButtonFromStock(gtk.StockMediaPlay)
		button.Connect("clicked", func() {
			spinner_sensitive.Start()
			spinner_unsensitive.Start()
		})
		vbox.Add(button)

		button = gtk.NewButtonFromStock(gtk.StockMediaStop)
		button.Connect("clicked", func() {
			spinner_sensitive.Stop()
			spinner_unsensitive.Stop()
		})
		vbox.Add(button)

		spinner_sensitive.Start()
		spinner_unsensitive.Start()
	}

	if !dialog.GetVisible() {
		dialog.ShowAll()
	} else {
		dialog.Destroy()
	}

	return gtk.ToWindow(dialog)
}
