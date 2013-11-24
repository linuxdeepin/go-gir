// Drawing Area
//
// GtkDrawingArea is a blank area where you can draw custom displays
// of various kinds.
//
// This demo has two drawing areas. The checkerboard area shows
// how you can just draw something; all you have to do is write
// a signal handler for expose_event, as shown here.
//
// The "scribble" area is a bit more advanced, and shows how to handle
// events such as button presses and mouse motion. Click the mouse
// and drag in the scribble area to draw squiggles. Resize the window
// to clear the area.
package drawing_area

import "gobject/gdk-3.0"
import "gobject/gtk-3.0"
import "gobject/cairo-1.0"
import "gobject/gobject-2.0"

var window *gtk.Window
var surface *cairo.Surface
var cr *cairo.Context

func checkerboard_draw(da *gtk.DrawingArea, cr *cairo.Context) bool {
	const spacing = 2
	const check_size = 10

	// At the start of a draw handler, a clip region has been set on
	// the Cairo context, and the contents have been cleared to the
	// widget's background color. The docs for
	// gdk_window_begin_paint_region() give more details on how this
	// works.

	xcount := 0
	width := da.GetAllocatedWidth()
	height := da.GetAllocatedHeight()
	i := spacing

	for i < width {
		j := spacing
		ycount := xcount % 2 // start with even/odd depending on row
		for j < height {
			if ycount % 2 != 0 {
				cr.SetSourceRGB(0.45777, 0, 0.45777)
			} else {
				cr.SetSourceRGB(1, 1, 1)
			}

			// If we're outside the clip, this will do nothing.
			cr.Rectangle(float64(i), float64(j), check_size, check_size)
			cr.Fill()

			j += check_size + spacing
			ycount++
		}
		i += check_size + spacing
		xcount++
	}

	// return TRUE because we've handled this event, so no
	// further processing is required.
	return true
}

// Create a new surface of the appropriate size to store out scribbles
func scribble_configure_event(widget *gtk.Widget, event *gdk.EventConfigure) bool {
	allocation := widget.GetAllocation()
	surface = widget.GetWindow().CreateSimilarSurface(cairo.ContentColor,
		int(allocation.Width), int(allocation.Height))

	// Initialize the surface to white
	cr = cairo.NewContext(surface)
	cr.SetSourceRGB(1, 1, 1)
	cr.Paint()

	// We've handled the configure event, no need for further processing.
	return true
}

// Redraw the screen from the surface
func scribble_draw(widget *gtk.Widget, cr *cairo.Context) bool {
	cr.SetSourceSurface(surface, 0, 0)
	cr.Paint()

	return false
}

// Draw a rectangle on the screen
func draw_brush(widget *gtk.Widget, x, y float64) {
	update_rect := cairo.RectangleInt{int32(x) - 3, int32(y) - 3, 6, 6}
	gdk.CairoRectangle(cr, &update_rect)
	cr.SetSourceRGB(0, 0, 0)
	cr.Fill()

	widget.GetWindow().InvalidateRect(&update_rect, false)
}

func scribble_button_press_event(widget *gtk.Widget, event *gdk.EventButton) bool {
	if surface == nil {
		// paranoia check, in case we haven't gotten a configure event
		return false
	}

	if event.Button == 1 {
		draw_brush(widget, event.X, event.Y)
	}

	// We've handled the event, stop processing
	return true
}

func scribble_motion_notify_event(widget *gtk.Widget, event *gdk.EventMotion) bool {
	if surface == nil {
		// paranoia check, in case we haven't gotten a configure event
		return false
	}

	// This call is very important; it requests the next motion event.
	// If you don't call gdk_window_get_pointer() you'll only get
	// a single motion event. The reason is that we specified
	// GDK_POINTER_MOTION_HINT_MASK to gtk_widget_set_events().
	// If we hadn't specified that, we could just use event->x, event->y
	// as the pointer location. But we'd also get deluged in events.
	// By requesting the next event as we handle the current one,
	// we avoid getting a huge number of events faster than we
	// can cope.
	x, y, state, _ := event.Window().GetPointer()

	if state & gdk.ModifierTypeButton1Mask != 0 {
		draw_brush(widget, float64(x), float64(y))
	}

	// We've handled the event, stop processing
	return true
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Drawing Area")
		window.Connect("destroy", func() {
			window = nil
			surface = nil
			cr = nil
		})
		window.SetBorderWidth(8)

		vbox := gtk.NewBox(gtk.OrientationVertical, 8)
		vbox.SetBorderWidth(8)
		window.Add(vbox)

		// Create the checkerboard area
		label := gtk.NewLabel(gobject.NilString)
		label.SetMarkup("<u>Checkerboard patter</u>")
		vbox.PackStart(label, false, false, 0)

		frame := gtk.NewFrame(gobject.NilString)
		frame.SetShadowType(gtk.ShadowTypeIn)
		vbox.PackStart(frame, true, true, 0)

		da := gtk.NewDrawingArea()
		da.SetSizeRequest(100, 100)
		frame.Add(da)
		da.Connect("draw", checkerboard_draw)

		// Create the scribble area
		label = gtk.NewLabel(gobject.NilString)
		label.SetMarkup("<u>Scribble area</u>")
		vbox.PackStart(label, false, false, 0)

		frame = gtk.NewFrame(gobject.NilString)
		frame.SetShadowType(gtk.ShadowTypeIn)
		vbox.PackStart(frame, true, true, 0)

		da = gtk.NewDrawingArea()
		da.SetSizeRequest(100, 100)
		frame.Add(da)

		// Signals used to handle backing surface
		da.Connect("draw", scribble_draw)
		da.Connect("configure-event", scribble_configure_event)

		// Event signals
		da.Connect("motion-notify-event", scribble_motion_notify_event)
		da.Connect("button-press-event", scribble_button_press_event)

		da.SetEvents(int(gdk.EventMask(da.GetEvents()) |
			gdk.EventMaskLeaveNotifyMask |
			gdk.EventMaskButtonPressMask |
			gdk.EventMaskPointerMotionMask |
			gdk.EventMaskPointerMotionHintMask))
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}