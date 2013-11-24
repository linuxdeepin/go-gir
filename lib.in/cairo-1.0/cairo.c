#include "_cgo_export.h"

cairo_surface_t * _cairo_image_surface_create_from_png_stream(void *closure)
{
	return cairo_image_surface_create_from_png_stream(io_reader_wrapper, closure);
}

cairo_status_t _cairo_surface_write_to_png_stream(cairo_surface_t *surface, void *closure)
{
	return cairo_surface_write_to_png_stream(surface, (cairo_write_func_t)io_writer_wrapper, closure);
}

cairo_surface_t *_cairo_pdf_surface_create_for_stream(void *closure, double width_in_points, double height_in_points)
{
	return cairo_pdf_surface_create_for_stream((cairo_write_func_t)io_writer_wrapper,
						   closure, width_in_points, height_in_points);
}
