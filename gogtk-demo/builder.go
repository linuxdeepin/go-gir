// Builder
//
// Demonstrates an interface loaded from a XML description.
package builder

import "gobject/gtk-3.0"
import "./gogtk-demo/common"

var window *gtk.Window
var builder *gtk.Builder

type dummy struct{}

func (*dummy) QuitActivate(action *gtk.Action) {
	window := gtk.ToWidget(builder.GetObject("window1"))
	window.Destroy()
}

func (*dummy) AboutActivate(action *gtk.Action) {
	about_dlg := gtk.ToDialog(builder.GetObject("aboutdialog1"))
	about_dlg.Run()
	about_dlg.Hide()
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		builder = gtk.NewBuilder()
		_, err := builder.AddFromFile(common.FindFile("demo.ui"))
		if err != nil {
			println("ERROR: ", err.Error())
			return nil
		}

		builder.ConnectSignals((*dummy)(nil))
		window = gtk.ToWindow(builder.GetObject("window1"))
		window.SetScreen(mainwin.GetScreen())
		window.Connect("destroy", func() {
			window = nil
			builder = nil
		})
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
