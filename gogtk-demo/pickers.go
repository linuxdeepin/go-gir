// Pickers
//
// These widgets are mainly intended for use in preference dialogs.
// They allow to select colors, fonts, files, directories and applications.
package pickers

import "gobject/gtk-3.0"

var window *gtk.Window

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetTitle("Pickers")
		window.Connect("destroy", func() { window = nil })
		window.SetBorderWidth(10)

		table := gtk.NewGrid()
		table.SetRowSpacing(3)
		table.SetColumnSpacing(10)
		window.Add(table)

		table.SetBorderWidth(10)

		make_label := func(label string) *gtk.Label {
			l := gtk.NewLabel(label)
			l.SetHAlign(gtk.AlignStart)
			l.SetVAlign(gtk.AlignCenter)
			l.SetHExpand(true)
			return l
		}

		table.Attach(make_label("Color:"), 0, 0, 1, 1)
		table.Attach(gtk.NewColorButton(), 1, 0, 1, 1)

		table.Attach(make_label("Font:"), 0, 1, 1, 1)
		table.Attach(gtk.NewFontButton(), 1, 1, 1, 1)

		table.Attach(make_label("File:"), 0, 2, 1, 1)
		table.Attach(gtk.NewFileChooserButton("Pick a File", gtk.FileChooserActionOpen), 1, 2, 1, 1)

		table.Attach(make_label("Folder:"), 0, 3, 1, 1)
		table.Attach(gtk.NewFileChooserButton("Pick a Folder", gtk.FileChooserActionSelectFolder), 1, 3, 1, 1)

		mail_picker := gtk.NewAppChooserButton("x-scheme-handler/mailto")
		mail_picker.SetShowDialogItem(true)
		table.Attach(make_label("Mail:"), 0, 4, 1, 1)
		table.Attach(mail_picker, 1, 4, 1, 1)
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
