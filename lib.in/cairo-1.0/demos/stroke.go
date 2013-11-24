package main

import "gobject/cairo-1.0"

func main() {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, 120, 120)
	cr := cairo.NewContext(surface)

	// Examples are in 1.0 x 1.0 coordinate space
	cr.Scale(120, 120)

	// Drawing code goes here
	cr.SetLineWidth(0.1)
	cr.SetSourceRGB(0, 0, 0)
	cr.Rectangle(0.25, 0.25, 0.5, 0.5)
	cr.Stroke()

	// Write output
	surface.WriteToPNG("stroke.png")
}