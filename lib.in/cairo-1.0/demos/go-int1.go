package main

import "gobject/cairo-1.0"
import "os"

func main() {
	// read
	r, err := os.Open("tux.png")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	surface := cairo.NewImageSurfaceFromPNGStream(r)
	if surface.Status() != cairo.StatusSuccess {
		panic(surface.Status().String())
	}

	// modify
	cr := cairo.NewContext(surface)
	cr.Scale(float64(surface.GetWidth()), float64(surface.GetHeight()))
	cr.SetSourceRGBA(0.2, 0, 0, 0.5)
	cr.Rectangle(0.05, 0.05, 0.9, 0.9)
	cr.Fill()

	// write
	w, err := os.Create("tux-2.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	status := surface.WriteToPNGStream(w)
	if status != cairo.StatusSuccess {
		panic(status.String())
	}
}