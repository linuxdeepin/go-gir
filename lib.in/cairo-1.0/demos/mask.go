package main

import "gobject/cairo-1.0"

func main() {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, 120, 120)
	cr := cairo.NewContext(surface)

	// Examples are in 1.0 x 1.0 coordinate space
	cr.Scale(120, 120)

	// Drawing code goes here
	linpat := cairo.NewLinearGradient(0, 0, 1, 1)
	linpat.AddColorStopRGB(0, 0, 0.3, 0.8)
	linpat.AddColorStopRGB(1, 0, 0.8, 0.3)

	radpat := cairo.NewRadialGradient(0.5, 0.5, 0.25, 0.5, 0.5, 0.75)
	radpat.AddColorStopRGBA(0, 0, 0, 0, 1)
	radpat.AddColorStopRGBA(0.5, 0, 0, 0, 0)

	cr.SetSource(linpat)
	cr.Mask(radpat)

	// Write output
	surface.WriteToPNG("mask.png")
}