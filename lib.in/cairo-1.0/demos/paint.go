package main

import "gobject/cairo-1.0"

func main() {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, 120, 120)
	cr := cairo.NewContext(surface)

	// Examples are in 1.0 x 1.0 coordinate space
	cr.Scale(120, 120)

	// Drawing code goes here
	cr.SetSourceRGB(0, 0, 0)
	cr.PaintWithAlpha(0.5)

	// Write output
	surface.WriteToPNG("paint.png")
}