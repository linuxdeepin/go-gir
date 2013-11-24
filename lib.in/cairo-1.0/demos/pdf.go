package main

import "gobject/cairo-1.0"
import "math"

func main() {
	surface := cairo.NewPDFSurface("pdf.pdf", 256, 256)
	cr := cairo.NewContext(surface)

	lin := cairo.NewLinearGradient(0,0,0,256)
	lin.AddColorStopRGBA(1, 0, 0, 0, 1)
	lin.AddColorStopRGBA(0, 1, 1, 1, 1)
	cr.Rectangle(0, 0, 256, 256)
	cr.SetSource(lin)
	cr.Fill()

	rad := cairo.NewRadialGradient(115.2, 102.4, 25.6,
		102.4, 102.4, 128.0)
	rad.AddColorStopRGBA(0, 1, 1, 1, 1)
	rad.AddColorStopRGBA(1, 0, 0, 0, 1)
	cr.SetSource(rad)
	cr.Arc(128, 128, 76.8, 0, 2 * math.Pi)
	cr.Fill()

	surface.Finish()
}