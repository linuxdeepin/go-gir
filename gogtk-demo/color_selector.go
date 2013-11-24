// Color Selector
//
// GtkColorSelection lets the user choose a color. GtkColorSelectionDialog is
// a prebuilt dialog containing a GtkColorSelection.
package color_selector

import "gobject/gtk-3.0"
import "gobject/gdk-3.0"
import "gobject/cairo-1.0"
import "gobject/gobject-2.0"

var color gdk.RGBA
var window *gtk.Window

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		color = gdk.RGBA{0, 0, 1, 1}
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Color Selection")
		window.Connect("destroy", func() { window = nil })

		window.SetBorderWidth(8)

		vbox := gtk.NewBox(gtk.OrientationVertical, 8)
		vbox.SetBorderWidth(8)
		window.Add(vbox)

		// Create the color swatch area

		frame := gtk.NewFrame(gobject.NilString)
		frame.SetShadowType(gtk.ShadowTypeIn)
		vbox.PackStart(frame, true, true, 0)

		da := gtk.NewDrawingArea()
		da.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) bool {
			gdk.CairoSetSourceRGBA(cr, &color)
			cr.Paint()

			// unref explicitly, can't rely on GC here, leaks like crazy
			cr.Unref()
			return true
		})

		// set a minimum size
		da.SetSizeRequest(200, 200)
		frame.Add(da)

		button := gtk.NewButtonWithMnemonic("_Change the above color")
		button.SetHAlign(gtk.AlignEnd)
		button.SetVAlign(gtk.AlignCenter)

		vbox.PackStart(button, false, false, 0)
		button.Connect("clicked", func() {
			dialog := gtk.NewColorSelectionDialog("Changing color")
			dialog.SetTransientFor(window)
			colorsel := gtk.ToColorSelection(dialog.GetColorSelection())
			colorsel.SetPreviousRGBA(&color)
			colorsel.SetCurrentRGBA(&color)
			colorsel.SetHasPalette(true)
			response := dialog.Run()
			if response == gtk.ResponseTypeOk {
				color = colorsel.GetCurrentRGBA()
				da.QueueDraw()
			}
			dialog.Destroy()
		})
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
