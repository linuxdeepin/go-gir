// Links
//
// GtkLabel can show hyperlinks. The default action is to call
// gtk_show_uri() on their URI, but it is possible to override
// this with a custom handler.
package links

import "gobject/gtk-3.0"

var window *gtk.Window

const label_text = `Some <a href="http://en.wikipedia.org/wiki/Text" title="plain text">text</a> may be marked up
as hyperlinks, which can be clicked
or activated via <a href="keynav">keynav</a>`

const dialog_text = `The term <i>keynav</i> is a shorthand for ` +
	`keyboard navigation and refers to the process of using ` +
	`a program (exclusively) via keyboard input.`

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetTitle("Links")
		window.Connect("destroy", func() { window = nil })
		window.SetBorderWidth(12)

		label := gtk.NewLabel(label_text)
		label.SetUseMarkup(true)
		label.Connect("activate-link", activate_link)
		window.Add(label)
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}

func activate_link(label *gtk.Label, uri string) bool {
	if uri == "keynav" {
		parent := gtk.ToWindow(label.GetParent())
		dialog := gtk.NewMessageDialogWithMarkup(parent,
			gtk.DialogFlagsDestroyWithParent,
			gtk.MessageTypeInfo,
			gtk.ButtonsTypeOk,
			dialog_text)
		dialog.Present()
		dialog.Connect("response", func() {
			dialog.Destroy()
		})

		return true
	}

	return false
}
