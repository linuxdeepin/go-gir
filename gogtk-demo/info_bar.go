// Info bar
//
// Info bar widgets are used to report important messages to the user.
package info_bar

import "gobject/gtk-3.0"

var window *gtk.Window

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Info Bars")

		window.Connect("destroy", func() { window = nil })
		window.SetBorderWidth(8)

		vbox := gtk.NewBox(gtk.OrientationVertical, 0)
		window.Add(vbox)

		bar := gtk.NewInfoBar()
		vbox.PackStart(bar, false, false, 0)
		bar.SetMessageType(gtk.MessageTypeInfo)
		label := gtk.NewLabel("This is an info bar with message type GTK_MESSAGE_INFO")
		gtk.ToBox(bar.GetContentArea()).PackStart(label, false, false, 0)

		bar = gtk.NewInfoBar()
		vbox.PackStart(bar, false, false, 0)
		bar.SetMessageType(gtk.MessageTypeWarning)
		label = gtk.NewLabel("This is an info bar with message type GTK_MESSAGE_WARNING")
		gtk.ToBox(bar.GetContentArea()).PackStart(label, false, false, 0)

		bar = gtk.NewInfoBarWithButtons(gtk.StockOk, gtk.ResponseTypeOk)
		bar.Connect("response", func(info_bar *gtk.InfoBar, response_id gtk.ResponseType) {
			dialog := gtk.NewMessageDialog(window, gtk.DialogFlagsModal|gtk.DialogFlagsDestroyWithParent,
				gtk.MessageTypeInfo, gtk.ButtonsTypeOk, "You clicked a button on an info bar")

			dialog.FormatSecondaryText("Your response has id %d", response_id)
			dialog.Run()
			dialog.Destroy()
		})
		vbox.PackStart(bar, false, false, 0)
		bar.SetMessageType(gtk.MessageTypeQuestion)
		label = gtk.NewLabel("This is an info bar with message type GTK_MESSAGE_QUESTION")
		gtk.ToBox(bar.GetContentArea()).PackStart(label, false, false, 0)

		bar = gtk.NewInfoBar()
		vbox.PackStart(bar, false, false, 0)
		bar.SetMessageType(gtk.MessageTypeError)
		label = gtk.NewLabel("This is an info bar with message type GTK_MESSAGE_ERROR")
		gtk.ToBox(bar.GetContentArea()).PackStart(label, false, false, 0)

		bar = gtk.NewInfoBar()
		vbox.PackStart(bar, false, false, 0)
		bar.SetMessageType(gtk.MessageTypeWarning)
		label = gtk.NewLabel("This is an info bar with message type GTK_MESSAGE_WARNING")
		gtk.ToBox(bar.GetContentArea()).PackStart(label, false, false, 0)

		frame := gtk.NewFrame("Info bars")
		vbox.PackStart(frame, false, false, 8)

		vbox2 := gtk.NewBox(gtk.OrientationVertical, 8)
		vbox2.SetBorderWidth(8)
		frame.Add(vbox2)

		label = gtk.NewLabel("An example of different info bars")
		vbox2.PackStart(label, false, false, 0)
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
