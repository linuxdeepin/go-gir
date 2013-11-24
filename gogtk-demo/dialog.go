// Dialog and Message Boxes
//
// Dialog widgets are used to pop up a transient window for user feedback.
package dialog

import "gobject/gtk-3.0"

var window *gtk.Window
var entry1 *gtk.Entry
var entry2 *gtk.Entry
var i = 1

func message_dialog_clicked() {
	dialog := gtk.NewMessageDialog(window, gtk.DialogFlagsModal | gtk.DialogFlagsDestroyWithParent,
		gtk.MessageTypeInfo, gtk.ButtonsTypeOk,
		"This message box has been popped up the following\nnumber of times:")
	dialog.FormatSecondaryText("%d", i)
	dialog.Run()
	dialog.Destroy()
	i++
}

func interactive_dialog_clicked() {
	dialog := gtk.NewDialogWithButtons("Interactive Dialog", window,
		gtk.DialogFlagsModal | gtk.DialogFlagsDestroyWithParent,
		gtk.StockOk, gtk.ResponseTypeOk, "_Non-stock Button", gtk.ResponseTypeCancel)

	content_area := gtk.ToBox(dialog.GetContentArea())

	hbox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	hbox.SetBorderWidth(8)
	content_area.PackStart(hbox, false, false, 0)

	stock := gtk.NewImageFromStock(gtk.StockDialogQuestion, int(gtk.IconSizeDialog))
	hbox.PackStart(stock, false, false, 0)

	table := gtk.NewGrid()
	table.SetRowSpacing(4)
	table.SetColumnSpacing(4)
	hbox.PackStart(table, true, true, 0)
	label := gtk.NewLabelWithMnemonic("_Entry 1")
	table.Attach(label, 0, 0, 1, 1)
	local_entry1 := gtk.NewEntry()
	local_entry1.SetText(entry1.GetText())
	table.Attach(local_entry1, 1, 0, 1, 1)
	label.SetMnemonicWidget(local_entry1)

	label = gtk.NewLabelWithMnemonic("E_ntry 2")
	table.Attach(label, 0, 1, 1, 1)
	local_entry2 := gtk.NewEntry()
	local_entry2.SetText(entry2.GetText())
	table.Attach(local_entry2, 1, 1, 1, 1)
	label.SetMnemonicWidget(local_entry2)

	hbox.ShowAll()
	response := dialog.Run()

	if response == gtk.ResponseTypeOk {
		entry1.SetText(local_entry1.GetText())
		entry2.SetText(local_entry2.GetText())
	}

	dialog.Destroy()
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Dialogs")

		window.Connect("destroy", func() {
			window = nil
			entry1 = nil
			entry2 = nil
			i = 1
		})
		window.SetBorderWidth(8)

		frame := gtk.NewFrame("Dialogs")
		window.Add(frame)

		vbox := gtk.NewBox(gtk.OrientationVertical, 8)
		vbox.SetBorderWidth(8)
		frame.Add(vbox)

		// Standard message dialog
		hbox := gtk.NewBox(gtk.OrientationHorizontal, 8)
		vbox.PackStart(hbox, false, false, 0)
		button := gtk.NewButtonWithMnemonic("_Message Dialog")
		button.Connect("clicked", message_dialog_clicked)
		hbox.PackStart(button, false, false, 0)
		vbox.PackStart(gtk.NewSeparator(gtk.OrientationHorizontal), false, false, 0)

		// Interactive dialog
		hbox = gtk.NewBox(gtk.OrientationHorizontal, 8)
		vbox.PackStart(hbox, false, false, 0)
		vbox2 := gtk.NewBox(gtk.OrientationVertical, 0)

		button = gtk.NewButtonWithMnemonic("_Interactive Dialog")
		button.Connect("clicked", interactive_dialog_clicked)
		hbox.PackStart(vbox2, false, false, 0)
		vbox2.PackStart(button, false, false, 0)

		table := gtk.NewGrid()
		table.SetRowSpacing(4)
		table.SetColumnSpacing(4)
		hbox.PackStart(table, false, false, 0)

		label := gtk.NewLabelWithMnemonic("_Entry 1")
		table.Attach(label, 0, 0, 1, 1)

		entry1 = gtk.NewEntry()
		table.Attach(entry1, 1, 0, 1, 1)
		label.SetMnemonicWidget(entry1)

		label = gtk.NewLabelWithMnemonic("E_ntry 2")
		table.Attach(label, 0, 1, 1, 1)

		entry2 = gtk.NewEntry()
		table.Attach(entry2, 1, 1, 1, 1)
		label.SetMnemonicWidget(entry2)
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}