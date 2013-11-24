// Button Boxes
//
// The Button Box widgets are used to arrange buttons with padding.
package button_boxes

import "gobject/gtk-3.0"

var window *gtk.Window

func create_bbox(orient gtk.Orientation, title string, spacing int, layout gtk.ButtonBoxStyle) *gtk.Frame {
	frame := gtk.NewFrame(title)
	bbox := gtk.NewButtonBox(orient)

	bbox.SetBorderWidth(5)
	frame.Add(bbox)

	bbox.SetLayout(layout)
	bbox.SetSpacing(spacing)

	var button *gtk.Button
	button = gtk.NewButtonFromStock(gtk.StockOk)
	bbox.Add(button)

	button = gtk.NewButtonFromStock(gtk.StockCancel)
	bbox.Add(button)

	button = gtk.NewButtonFromStock(gtk.StockHelp)
	bbox.Add(button)

	return frame
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetTitle("Button Boxes")
		window.Connect("destroy", func() { window = nil })
		window.SetBorderWidth(10)

		main_vbox := gtk.NewBox(gtk.OrientationVertical, 0)
		window.Add(main_vbox)

		frame_horz := gtk.NewFrame("Horizontal Button Boxes")
		main_vbox.PackStart(frame_horz, true, true, 10)

		vbox := gtk.NewBox(gtk.OrientationVertical, 0)
		vbox.SetBorderWidth(10)
		frame_horz.Add(vbox)

		vbox.PackStart(
			create_bbox(gtk.OrientationHorizontal, "Spread",
				40, gtk.ButtonBoxStyleSpread), true, true, 0)
		vbox.PackStart(
			create_bbox(gtk.OrientationHorizontal, "Edge",
				40, gtk.ButtonBoxStyleEdge), true, true, 5)
		vbox.PackStart(
			create_bbox(gtk.OrientationHorizontal, "Start",
				40, gtk.ButtonBoxStyleStart), true, true, 5)
		vbox.PackStart(
			create_bbox(gtk.OrientationHorizontal, "End",
				40, gtk.ButtonBoxStyleEnd), true, true, 5)

		frame_vert := gtk.NewFrame("Vertical Button Boxes")
		main_vbox.PackStart(frame_vert, true, true, 10)

		hbox := gtk.NewBox(gtk.OrientationHorizontal, 0)
		hbox.SetBorderWidth(10)
		frame_vert.Add(hbox)

		hbox.PackStart(
			create_bbox(gtk.OrientationVertical, "Spread",
				30, gtk.ButtonBoxStyleSpread), true, true, 0)
		hbox.PackStart(
			create_bbox(gtk.OrientationVertical, "Edge",
				30, gtk.ButtonBoxStyleEdge), true, true, 5)
		hbox.PackStart(
			create_bbox(gtk.OrientationVertical, "Start",
				30, gtk.ButtonBoxStyleStart), true, true, 5)
		hbox.PackStart(
			create_bbox(gtk.OrientationVertical, "End",
				30, gtk.ButtonBoxStyleEnd), true, true, 5)
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}
