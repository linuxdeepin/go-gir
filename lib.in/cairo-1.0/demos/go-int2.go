package main

import "gobject/cairo-1.0"

func main() {
	surface := cairo.NewImageSurface(cairo.FormatARGB32, 128, 128)
	cr := cairo.NewContext(surface)
	x := cr.GetTarget()
	y := cr.GetTarget()
	z := cr.GetTarget()
	println(x, y, z)
}