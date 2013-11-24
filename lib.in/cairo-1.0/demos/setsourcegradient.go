package main

import "gobject/cairo-1.0"

func main() {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, 120, 120)
	cr := cairo.NewContext(surface)

	// Examples are in 1.0 x 1.0 coordinate space
	cr.Scale(120, 120)

	// Drawing code goes here
	radpat := cairo.NewRadialGradient(0.25, 0.25, 0.1, 0.5, 0.5, 0.5)
	radpat.AddColorStopRGB(0, 1.0, 0.8, 0.8)
	radpat.AddColorStopRGB(1, 0.9, 0.0, 0.0)

	for i := 1; i < 10; i++ {
		fi := float64(i)
		for j := 1; j < 10; j++ {
			fj := float64(j)
			cr.Rectangle(fi/10.0 - 0.04, fj/10.0 - 0.04, 0.08, 0.08)
		}
	}
	cr.SetSource(radpat)
	cr.Fill()

	linpat := cairo.NewLinearGradient(0.25, 0.35, 0.75, 0.65)
	linpat.AddColorStopRGBA(0.00, 1, 1, 1, 0)
	linpat.AddColorStopRGBA(0.25, 0, 1, 0, 0.5)
	linpat.AddColorStopRGBA(0.50, 1, 1, 1, 0)
	linpat.AddColorStopRGBA(0.75, 0, 0, 1, 0.5)
	linpat.AddColorStopRGBA(1.00, 1, 1, 1, 0)

	cr.Rectangle(0, 0, 1, 1)
	cr.SetSource(linpat)
	cr.Fill()

	// Write output
	surface.WriteToPNG("setsourcegradient.png")
}