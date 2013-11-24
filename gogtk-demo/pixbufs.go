// Pixbufs
//
// A GdkPixbuf represents an image, normally in RGB or RGBA format.
// Pixbufs are normally used to load files from disk and perform
// image scaling.
//
// This demo is not all that educational, but looks cool. It was written
// by Extreme Pixbuf Hacker Federico Mena Quintero. It also shows
// off how to use GtkDrawingArea to do a simple animation.
//
// Look at the Image demo for additional pixbuf usage examples.
package pixbufs

import "gobject/gtk-3.0"
import "gobject/gdk-3.0"
import "gobject/gdkpixbuf-2.0"
import "gobject/cairo-1.0"
import "math"
import "time"
import "sync"
import "./gogtk-demo/common"

const frame_delay = 50
const background_name = "background.jpg"

var image_names = []string{
	"apple-red.png",
	"gnome-applets.png",
	"gnome-calendar.png",
	"gnome-foot.png",
	"gnome-gmush.png",
	"gnome-gimp.png",
	"gnome-gsame.png",
	"gnu-keys.png",
}

var window *gtk.Window
var background *gdkpixbuf.Pixbuf
var images []*gdkpixbuf.Pixbuf
var back_width int
var back_height int

var frame_lock sync.Mutex
var frame *gdkpixbuf.Pixbuf
var da *gtk.DrawingArea

func load_pixbufs() error {
	if background != nil {
		return nil // already loaded earlier
	}

	var err error
	background, err = gdkpixbuf.NewPixbufFromFile(common.FindFile(background_name))
	if err != nil {
		return err
	}

	back_width = background.GetWidth()
	back_height = background.GetHeight()

	images = make([]*gdkpixbuf.Pixbuf, len(image_names))
	for i, name := range image_names {
		images[i], err = gdkpixbuf.NewPixbufFromFile(common.FindFile(name))
		if err != nil {
			return err
		}
	}

	return nil
}

// Expose callback for the drawing area
func draw_cb(widget *gtk.Widget, cr *cairo.Context) bool {
	frame_lock.Lock()
	gdk.CairoSetSourcePixbuf(cr, frame, 0, 0)
	cr.Paint()
	frame_lock.Unlock()

	// unref explicitly, can't rely on GC here, leaks like crazy
	cr.Unref()
	return true
}

const cycle_len = 60
var frame_num int

func draw_one_frame() {
	background.CopyArea(0, 0, back_width, back_height, frame, 0, 0)
	f := float64(frame_num % cycle_len) / cycle_len
	xmid := float64(back_width) / 2
	ymid := float64(back_height) / 2
	radius := math.Min(xmid, ymid) / 2

	for i, image := range images {
		ang := 2 * math.Pi * float64(i) / float64(len(images)) - f * 2 * math.Pi

		iw := image.GetWidth()
		ih := image.GetHeight()

		r := radius + (radius / 3) * math.Sin(f * 2 * math.Pi)

		xpos := math.Floor(xmid + r * math.Cos(ang) - float64(iw) / 2 + 0.5)
		ypos := math.Floor(ymid + r * math.Sin(ang) - float64(ih) / 2 + 0.5)

		var k, alpha float64
		if i & 1 != 0 {
			k = math.Sin(f * 2 * math.Pi)
			alpha = math.Max(127, math.Abs(255 * math.Sin(f * 2 * math.Pi)))
		} else {
			k = math.Cos(f * 2 * math.Pi)
			alpha = math.Max(127, math.Abs(255 * math.Cos(f * 2 * math.Pi)))
		}
		k = 2 * k * k
		k = math.Max(0.25, k)

		r1 := cairo.RectangleInt{
			int32(xpos), int32(ypos),
			int32(float64(iw) * k), int32(float64(ih) * k),
		}
		r2 := cairo.RectangleInt{
			0, 0,
			int32(back_width), int32(back_height),
		}

		if dest, ok := gdk.RectangleIntersect(&r1, &r2); ok {
			frame_lock.Lock()
			image.Composite(frame,
				int(dest.X), int(dest.Y), int(dest.Width), int(dest.Height),
				xpos, ypos, k, k,
				gdkpixbuf.InterpTypeNearest, int(alpha))
			frame_lock.Unlock()
		}
	}

}

func drawing_loop(cancel chan int) {
	ticker := time.NewTicker(frame_delay * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			draw_one_frame()

			gdk.ThreadsEnter()
			da.QueueDraw()
			gdk.ThreadsLeave()

			frame_num++
		case <-cancel:
			ticker.Stop()
			cancel <- 1
			return
		}
	}
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		cancel := make(chan int)
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Pixbufs")
		window.SetResizable(false)
		window.Connect("destroy", func() {
			window = nil
			frame = nil
			da = nil

		})

		if err := load_pixbufs(); err != nil {
			dialog := gtk.NewMessageDialog(mainwin, gtk.DialogFlagsDestroyWithParent,
				gtk.MessageTypeError, gtk.ButtonsTypeClose,
				"Failed to load an image: %s", err)

			dialog.Connect("response", func() { dialog.Destroy() })
			dialog.Show()
			goto done
		}

		window.SetSizeRequest(back_width, back_height)
		frame = gdkpixbuf.NewPixbuf(gdkpixbuf.ColorspaceRGB, false, 8, back_width, back_height)
		da = gtk.NewDrawingArea()
		da.Connect("draw", draw_cb)
		window.Add(da)

		go drawing_loop(cancel)
		window.Connect("destroy", func() {
			// stop drawing loop
			cancel <- 1
			<-cancel
		})
	}

done:
	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}