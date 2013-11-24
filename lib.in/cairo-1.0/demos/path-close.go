package main

import "gobject/cairo-1.0"
import "math"

func main() {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, 120, 120)
	cr := cairo.NewContext(surface)

	// Examples are in 1.0 x 1.0 coordinate space
	cr.Scale(120, 120)

	// Drawing code goes here
	cr.SetLineWidth(0.1)
	cr.SetSourceRGB(0, 0, 0)

	cr.MoveTo(0.25, 0.25)
	cr.LineTo(0.5, 0.375)
	cr.RelLineTo(0.25, -0.125)
	cr.Arc(0.5, 0.5, 0.25 * math.Sqrt2, -0.25 * math.Pi, 0.25 * math.Pi)
	cr.RelCurveTo(-0.25, -0.125, -0.25, 0.125, -0.5, 0)
	cr.ClosePath()

	cr.Stroke()

	// Write output
	surface.WriteToPNG("path-close.png")
}