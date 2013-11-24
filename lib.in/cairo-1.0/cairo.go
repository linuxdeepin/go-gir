// Full cairo bindings. Made in a stupid fashion (i.e. no sugar). But bits of
// memory management are here (GC finalizer hooks).

package cairo

/*
#include <stdlib.h>
#include <cairo.h>
#include <cairo-pdf.h>
#include <cairo-gobject.h>

extern cairo_surface_t * _cairo_image_surface_create_from_png_stream(void *closure);
extern cairo_status_t _cairo_surface_write_to_png_stream(cairo_surface_t *surface, void *closure);
extern cairo_surface_t *_cairo_pdf_surface_create_for_stream(void *closure, double width_in_points, double height_in_points);
#cgo pkg-config: cairo-gobject cairo-png
*/
import "C"
import "runtime"
import "strings"
import "reflect"
import "unsafe"
import "io"

import (
	"gobject/gobject-2.0"
)

//----------------------------------------------------------------------------
// GObject types
//----------------------------------------------------------------------------

// structs
// cairo_gobject_context_get_type (void);
func (*Context) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_context_get_type())
}

func ContextGetType() gobject.Type { return (*Context)(nil).GetStaticType() }

// cairo_gobject_device_get_type (void);

// cairo_gobject_pattern_get_type (void);
func (*Pattern) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_pattern_get_type())
}
func PatternGetType() gobject.Type { return (*Pattern)(nil).GetStaticType() }

// cairo_gobject_surface_get_type (void);
func (*Surface) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_surface_get_type())
}
func SurfaceGetType() gobject.Type { return (*Surface)(nil).GetStaticType() }

// cairo_gobject_rectangle_get_type (void);

// cairo_gobject_scaled_font_get_type (void);
func (*ScaledFont) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_scaled_font_get_type())
}
func ScaledFontGetType() gobject.Type { return (*ScaledFont)(nil).GetStaticType() }

// cairo_gobject_font_face_get_type (void);
func (*FontFace) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_font_face_get_type())
}
func FontFaceGetType() gobject.Type { return (*FontFace)(nil).GetStaticType() }

// cairo_gobject_font_options_get_type (void);
func (*FontOptions) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_font_options_get_type())
}
func FontOptionsGetType() gobject.Type { return (*FontOptions)(nil).GetStaticType() }

// cairo_gobject_rectangle_int_get_type (void);

// cairo_gobject_region_get_type (void);
func (*Region) GetStaticType() gobject.Type {
	return gobject.Type(C.cairo_gobject_region_get_type())
}
func RegionGetType() gobject.Type { return (*Region)(nil).GetStaticType() }

// TODO: No need for these?
// enums
// cairo_gobject_status_get_type (void);
// cairo_gobject_content_get_type (void);
// cairo_gobject_operator_get_type (void);
// cairo_gobject_antialias_get_type (void);
// cairo_gobject_fill_rule_get_type (void);
// cairo_gobject_line_cap_get_type (void);
// cairo_gobject_line_join_get_type (void);
// cairo_gobject_text_cluster_flags_get_type (void);
// cairo_gobject_font_slant_get_type (void);
// cairo_gobject_font_weight_get_type (void);
// cairo_gobject_subpixel_order_get_type (void);
// cairo_gobject_hint_style_get_type (void);
// cairo_gobject_hint_metrics_get_type (void);
// cairo_gobject_font_type_get_type (void);
// cairo_gobject_path_data_type_get_type (void);
// cairo_gobject_device_type_get_type (void);
// cairo_gobject_surface_type_get_type (void);
// cairo_gobject_format_get_type (void);
// cairo_gobject_pattern_type_get_type (void);
// cairo_gobject_extend_get_type (void);
// cairo_gobject_filter_get_type (void);
// cairo_gobject_region_overlap_get_type (void);

//----------------------------------------------------------------------------
// TODO MOVE
//----------------------------------------------------------------------------

// Cairo C to Go reflect.Value converter (for closure marshaling)
func cairo_marshaler(value *gobject.Value, t reflect.Type) (reflect.Value, bool) {
	if !strings.HasSuffix(t.Elem().PkgPath(), "cairo-1.0") {
		return reflect.Value{}, false
	}

	out := reflect.New(t)
	st, ok := out.Elem().Interface().(gobject.StaticTyper)
	if !ok {
		return reflect.Value{}, false
	}
	gtypedst := st.GetStaticType()

	var dst gobject.Value
	dst.Init(gtypedst)
	ok = dst.Transform(value)
	if !ok {
		panic("GValue is not transformable to " + t.String())
	}

	defer dst.Unset()

	// TODO: is it correct that I should always grab here?
	switch gtypedst {
	case ContextGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = ContextWrap(dst.GetBoxed(), true)
	case PatternGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = PatternWrap(dst.GetBoxed(), true)
	case SurfaceGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = SurfaceWrap(dst.GetBoxed(), true)
	case FontOptionsGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = FontOptionsWrap(dst.GetBoxed())
	case RegionGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = RegionWrap(dst.GetBoxed(), true)
	case FontFaceGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = FontFaceWrap(dst.GetBoxed(), true)
	case ScaledFontGetType():
		*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = ScaledFontWrap(dst.GetBoxed(), true)
	default:
		return reflect.Value{}, false
	}

	return out.Elem(), true
}

func init() {
	gobject.CairoMarshaler = cairo_marshaler
}

var go_repr_cookie C.cairo_user_data_key_t

type RectangleInt struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

func (this *RectangleInt) c() *C.cairo_rectangle_int_t {
	return (*C.cairo_rectangle_int_t)(unsafe.Pointer(this))
}

//----------------------------------------------------------------------------
// Error Handling
//----------------------------------------------------------------------------

type Status int

const (
	StatusSuccess                Status = C.CAIRO_STATUS_SUCCESS
	StatusNoMemory               Status = C.CAIRO_STATUS_NO_MEMORY
	StatusInvalidRestore         Status = C.CAIRO_STATUS_INVALID_RESTORE
	StatusInvalidPopGroup        Status = C.CAIRO_STATUS_INVALID_POP_GROUP
	StatusNoCurrentPoint         Status = C.CAIRO_STATUS_NO_CURRENT_POINT
	StatusInvalidMatrix          Status = C.CAIRO_STATUS_INVALID_MATRIX
	StatusInvalidStatus          Status = C.CAIRO_STATUS_INVALID_STATUS
	StatusNullPointer            Status = C.CAIRO_STATUS_NULL_POINTER
	StatusInvalidString          Status = C.CAIRO_STATUS_INVALID_STRING
	StatusInvalidPathData        Status = C.CAIRO_STATUS_INVALID_PATH_DATA
	StatusReadError              Status = C.CAIRO_STATUS_READ_ERROR
	StatusWriteError             Status = C.CAIRO_STATUS_WRITE_ERROR
	StatusSurfaceFinished        Status = C.CAIRO_STATUS_SURFACE_FINISHED
	StatusSurfaceTypeMismatch    Status = C.CAIRO_STATUS_SURFACE_TYPE_MISMATCH
	StatusPatternTypeMismatch    Status = C.CAIRO_STATUS_PATTERN_TYPE_MISMATCH
	StatusInvalidContent         Status = C.CAIRO_STATUS_INVALID_CONTENT
	StatusInvalidFormat          Status = C.CAIRO_STATUS_INVALID_FORMAT
	StatusInvalidVisual          Status = C.CAIRO_STATUS_INVALID_VISUAL
	StatusFileNotFound           Status = C.CAIRO_STATUS_FILE_NOT_FOUND
	StatusInvalidDash            Status = C.CAIRO_STATUS_INVALID_DASH
	StatusInvalidDscComment      Status = C.CAIRO_STATUS_INVALID_DSC_COMMENT
	StatusInvalidIndex           Status = C.CAIRO_STATUS_INVALID_INDEX
	StatusClipNotRepresentable   Status = C.CAIRO_STATUS_CLIP_NOT_REPRESENTABLE
	StatusTempFileError          Status = C.CAIRO_STATUS_TEMP_FILE_ERROR
	StatusInvalidStride          Status = C.CAIRO_STATUS_INVALID_STRIDE
	StatusFontTypeMismatch       Status = C.CAIRO_STATUS_FONT_TYPE_MISMATCH
	StatusUserFontImmutable      Status = C.CAIRO_STATUS_USER_FONT_IMMUTABLE
	StatusUserFontError          Status = C.CAIRO_STATUS_USER_FONT_ERROR
	StatusNegativeCount          Status = C.CAIRO_STATUS_NEGATIVE_COUNT
	StatusInvalidClusters        Status = C.CAIRO_STATUS_INVALID_CLUSTERS
	StatusInvalidSlant           Status = C.CAIRO_STATUS_INVALID_SLANT
	StatusInvalidWeight          Status = C.CAIRO_STATUS_INVALID_WEIGHT
	StatusInvalidSize            Status = C.CAIRO_STATUS_INVALID_SIZE
	StatusUserFontNotImplemented Status = C.CAIRO_STATUS_USER_FONT_NOT_IMPLEMENTED
	StatusDeviceTypeMismatch     Status = C.CAIRO_STATUS_DEVICE_TYPE_MISMATCH
	StatusDeviceError            Status = C.CAIRO_STATUS_DEVICE_ERROR
	StatusLastStatus             Status = C.CAIRO_STATUS_LAST_STATUS
)

func (this Status) c() C.cairo_status_t {
	return C.cairo_status_t(this)
}

func (this Status) String() string {
	cstr := C.cairo_status_to_string(this.c())
	return C.GoString(cstr)
}

//----------------------------------------------------------------------------
// Drawing Context
//----------------------------------------------------------------------------

// typedef             cairo_t;
type Context struct {
	C unsafe.Pointer
}

func (this *Context) c() *C.cairo_t {
	return (*C.cairo_t)(this.C)
}

func context_finalizer(this *Context) {
	if gobject.FQueue.Push(unsafe.Pointer(this), context_finalizer2) {
		return
	}
	C.cairo_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_destroy(this.c())
}

func context_finalizer2(this_un unsafe.Pointer) {
	this := (*Context)(this_un)
	C.cairo_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_destroy(this.c())
}

func ContextWrap(c_un unsafe.Pointer, grab bool) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_t)(c_un)
	go_repr := C.cairo_get_user_data(c, &go_repr_cookie)
	if go_repr != nil {
		return unsafe.Pointer(go_repr)
	}

	context := &Context{unsafe.Pointer(c)}
	if grab {
		C.cairo_reference(c)
	}
	runtime.SetFinalizer(context, context_finalizer)

	status := C.cairo_set_user_data(c, &go_repr_cookie, unsafe.Pointer(context), nil)
	if status != C.CAIRO_STATUS_SUCCESS {
		panic("failed to set user data, out of memory?")
	}
	return unsafe.Pointer(context)
}

// cairo_t *           cairo_create                        (cairo_surface_t *target);
func NewContext(target SurfaceLike) *Context {
	return (*Context)(ContextWrap(unsafe.Pointer(C.cairo_create(target.InheritedFromCairoSurface().c())), false))
}

// cairo_t *           cairo_reference                     (cairo_t *cr);
// void                cairo_destroy                       (cairo_t *cr);
func (this *Context) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_status                        (cairo_t *cr);
func (this *Context) Status() Status {
	return Status(C.cairo_status(this.c()))
}

// void                cairo_save                          (cairo_t *cr);
func (this *Context) Save() {
	C.cairo_save(this.c())
}

// void                cairo_restore                       (cairo_t *cr);
func (this *Context) Restore() {
	C.cairo_restore(this.c())
}

// cairo_surface_t *   cairo_get_target                    (cairo_t *cr);
func (this *Context) GetTarget() *Surface {
	return (*Surface)(SurfaceWrap(unsafe.Pointer(C.cairo_get_target(this.c())), true))
}

// void                cairo_push_group                    (cairo_t *cr);
func (this *Context) PushGroup() {
	C.cairo_push_group(this.c())
}

// void                cairo_push_group_with_content       (cairo_t *cr,
//                                                          cairo_content_t content);
func (this *Context) PushGroupWithContent(content Content) {
	C.cairo_push_group_with_content(this.c(), content.c())
}

// cairo_pattern_t *   cairo_pop_group                     (cairo_t *cr);
func (this *Context) PopGroup() *Pattern {
	return (*Pattern)(PatternWrap(unsafe.Pointer(C.cairo_pop_group(this.c())), false))
}

// void                cairo_pop_group_to_source           (cairo_t *cr);
func (this *Context) PopGroupToSource() {
	C.cairo_pop_group_to_source(this.c())
}

// cairo_surface_t *   cairo_get_group_target              (cairo_t *cr);
func (this *Context) GetGroupTarget() *Surface {
	return (*Surface)(SurfaceWrap(unsafe.Pointer(C.cairo_get_group_target(this.c())), true))
}

// void                cairo_set_source_rgb                (cairo_t *cr,
//                                                          double red,
//                                                          double green,
//                                                          double blue);
func (this *Context) SetSourceRGB(r, g, b float64) {
	C.cairo_set_source_rgb(this.c(), C.double(r), C.double(g), C.double(b))
}

// void                cairo_set_source_rgba               (cairo_t *cr,
//                                                          double red,
//                                                          double green,
//                                                          double blue,
//                                                          double alpha);
func (this *Context) SetSourceRGBA(r, g, b, a float64) {
	C.cairo_set_source_rgba(this.c(), C.double(r), C.double(g), C.double(b), C.double(a))
}

// void                cairo_set_source                    (cairo_t *cr,
//                                                          cairo_pattern_t *source);
func (this *Context) SetSource(source PatternLike) {
	C.cairo_set_source(this.c(), source.InheritedFromCairoPattern().c())
}

// void                cairo_set_source_surface            (cairo_t *cr,
//                                                          cairo_surface_t *surface,
//                                                          double x,
//                                                          double y);
func (this *Context) SetSourceSurface(surface SurfaceLike, x, y float64) {
	C.cairo_set_source_surface(this.c(), surface.InheritedFromCairoSurface().c(), C.double(x), C.double(y))
}

// cairo_pattern_t *   cairo_get_source                    (cairo_t *cr);
func (this *Context) GetSource() *Pattern {
	return (*Pattern)(PatternWrap(unsafe.Pointer(C.cairo_get_source(this.c())), true))
}

// enum                cairo_antialias_t;
type Antialias int

const (
	AntialiasDefault  Antialias = C.CAIRO_ANTIALIAS_DEFAULT
	AntialiasNone     Antialias = C.CAIRO_ANTIALIAS_NONE
	AntialiasGray     Antialias = C.CAIRO_ANTIALIAS_GRAY
	AntialiasSubpixel Antialias = C.CAIRO_ANTIALIAS_SUBPIXEL
)

func (this Antialias) c() C.cairo_antialias_t {
	return C.cairo_antialias_t(this)
}

// void                cairo_set_antialias                 (cairo_t *cr,
//                                                          cairo_antialias_t antialias);
func (this *Context) SetAntialias(antialias Antialias) {
	C.cairo_set_antialias(this.c(), antialias.c())
}

// cairo_antialias_t   cairo_get_antialias                 (cairo_t *cr);
func (this *Context) GetAntialias() Antialias {
	return Antialias(C.cairo_get_antialias(this.c()))
}

// void                cairo_set_dash                      (cairo_t *cr,
//                                                          const double *dashes,
//                                                          int num_dashes,
//                                                          double offset);
func (this *Context) SetDash(dashes []float64, offset float64) {
	var first *C.double
	if len(dashes) != 0 {
		first = (*C.double)(unsafe.Pointer(&dashes[0]))
	}
	C.cairo_set_dash(this.c(), first, C.int(len(dashes)), C.double(offset))
}

// int                 cairo_get_dash_count                (cairo_t *cr);
func (this *Context) GetDashCount() int {
	return int(C.cairo_get_dash_count(this.c()))
}

// void                cairo_get_dash                      (cairo_t *cr,
//                                                          double *dashes,
//                                                          double *offset);
func (this *Context) GetDash() ([]float64, float64) {
	var first *C.double
	var offset C.double
	var dashes []float64

	count := this.GetDashCount()
	if count != 0 {
		dashes = make([]float64, count)
		first = (*C.double)(unsafe.Pointer(&dashes[0]))
	}
	C.cairo_get_dash(this.c(), first, &offset)
	return dashes, float64(offset)
}

// enum                cairo_fill_rule_t;
type FillRule int

const (
	FillRuleWinding FillRule = C.CAIRO_FILL_RULE_WINDING
	FillRuleEvenOdd FillRule = C.CAIRO_FILL_RULE_EVEN_ODD
)

func (this FillRule) c() C.cairo_fill_rule_t {
	return C.cairo_fill_rule_t(this)
}

// void                cairo_set_fill_rule                 (cairo_t *cr,
//                                                          cairo_fill_rule_t fill_rule);
func (this *Context) SetFillRule(fill_rule FillRule) {
	C.cairo_set_fill_rule(this.c(), fill_rule.c())
}

// cairo_fill_rule_t   cairo_get_fill_rule                 (cairo_t *cr);
func (this *Context) GetFillRule() FillRule {
	return FillRule(C.cairo_get_fill_rule(this.c()))
}

// enum                cairo_line_cap_t;
type LineCap int

const (
	LineCapButt   LineCap = C.CAIRO_LINE_CAP_BUTT
	LineCapRound  LineCap = C.CAIRO_LINE_CAP_ROUND
	LineCapSquare LineCap = C.CAIRO_LINE_CAP_SQUARE
)

func (this LineCap) c() C.cairo_line_cap_t {
	return C.cairo_line_cap_t(this)
}

// void                cairo_set_line_cap                  (cairo_t *cr,
//                                                          cairo_line_cap_t line_cap);
func (this *Context) SetLineCap(line_cap LineCap) {
	C.cairo_set_line_cap(this.c(), line_cap.c())
}

// cairo_line_cap_t    cairo_get_line_cap                  (cairo_t *cr);
func (this *Context) GetLineCap() LineCap {
	return LineCap(C.cairo_get_line_cap(this.c()))
}

// enum                cairo_line_join_t;
type LineJoin int

const (
	LineJoinMiter LineJoin = C.CAIRO_LINE_JOIN_MITER
	LineJoinRound LineJoin = C.CAIRO_LINE_JOIN_ROUND
	LineJoinBevel LineJoin = C.CAIRO_LINE_JOIN_BEVEL
)

func (this LineJoin) c() C.cairo_line_join_t {
	return C.cairo_line_join_t(this)
}

// void                cairo_set_line_join                 (cairo_t *cr,
//                                                          cairo_line_join_t line_join);
func (this *Context) SetLineJoin(line_join LineJoin) {
	C.cairo_set_line_join(this.c(), line_join.c())
}

// cairo_line_join_t   cairo_get_line_join                 (cairo_t *cr);
func (this *Context) GetLineJoin() LineJoin {
	return LineJoin(C.cairo_get_line_join(this.c()))
}

// void                cairo_set_line_width                (cairo_t *cr,
//                                                          double width);
func (this *Context) SetLineWidth(width float64) {
	C.cairo_set_line_width(this.c(), C.double(width))
}

// double              cairo_get_line_width                (cairo_t *cr);
func (this *Context) GetLineWidth() float64 {
	return float64(C.cairo_get_line_width(this.c()))
}

// void                cairo_set_miter_limit               (cairo_t *cr,
//                                                          double limit);
func (this *Context) SetMiterLimit(limit float64) {
	C.cairo_set_miter_limit(this.c(), C.double(limit))
}

// double              cairo_get_miter_limit               (cairo_t *cr);
func (this *Context) GetMiterLimit() float64 {
	return float64(C.cairo_get_miter_limit(this.c()))
}

// enum                cairo_operator_t;
type Operator int

const (
	OperatorClear         Operator = C.CAIRO_OPERATOR_CLEAR
	OperatorSource        Operator = C.CAIRO_OPERATOR_SOURCE
	OperatorOver          Operator = C.CAIRO_OPERATOR_OVER
	OperatorIn            Operator = C.CAIRO_OPERATOR_IN
	OperatorOut           Operator = C.CAIRO_OPERATOR_OUT
	OperatorAtop          Operator = C.CAIRO_OPERATOR_ATOP
	OperatorDest          Operator = C.CAIRO_OPERATOR_DEST
	OperatorDestOver      Operator = C.CAIRO_OPERATOR_DEST_OVER
	OperatorDestIn        Operator = C.CAIRO_OPERATOR_DEST_IN
	OperatorDestOut       Operator = C.CAIRO_OPERATOR_DEST_OUT
	OperatorDestAtop      Operator = C.CAIRO_OPERATOR_DEST_ATOP
	OperatorXor           Operator = C.CAIRO_OPERATOR_XOR
	OperatorAdd           Operator = C.CAIRO_OPERATOR_ADD
	OperatorSaturate      Operator = C.CAIRO_OPERATOR_SATURATE
	OperatorMultiply      Operator = C.CAIRO_OPERATOR_MULTIPLY
	OperatorScreen        Operator = C.CAIRO_OPERATOR_SCREEN
	OperatorOverlay       Operator = C.CAIRO_OPERATOR_OVERLAY
	OperatorDarken        Operator = C.CAIRO_OPERATOR_DARKEN
	OperatorLighten       Operator = C.CAIRO_OPERATOR_LIGHTEN
	OperatorColorDodge    Operator = C.CAIRO_OPERATOR_COLOR_DODGE
	OperatorColorBurn     Operator = C.CAIRO_OPERATOR_COLOR_BURN
	OperatorHardLight     Operator = C.CAIRO_OPERATOR_HARD_LIGHT
	OperatorSoftLight     Operator = C.CAIRO_OPERATOR_SOFT_LIGHT
	OperatorDifference    Operator = C.CAIRO_OPERATOR_DIFFERENCE
	OperatorExclusion     Operator = C.CAIRO_OPERATOR_EXCLUSION
	OperatorHslHue        Operator = C.CAIRO_OPERATOR_HSL_HUE
	OperatorHslSaturation Operator = C.CAIRO_OPERATOR_HSL_SATURATION
	OperatorHslColor      Operator = C.CAIRO_OPERATOR_HSL_COLOR
	OperatorHslLuminosity Operator = C.CAIRO_OPERATOR_HSL_LUMINOSITY
)

func (this Operator) c() C.cairo_operator_t {
	return C.cairo_operator_t(this)
}

// void                cairo_set_operator                  (cairo_t *cr,
//                                                          cairo_operator_t op);
func (this *Context) SetOperator(op Operator) {
	C.cairo_set_operator(this.c(), op.c())
}

// cairo_operator_t    cairo_get_operator                  (cairo_t *cr);
func (this *Context) GetOperator() Operator {
	return Operator(C.cairo_get_operator(this.c()))
}

// void                cairo_set_tolerance                 (cairo_t *cr,
//                                                          double tolerance);
func (this *Context) SetTolerance(tolerance float64) {
	C.cairo_set_tolerance(this.c(), C.double(tolerance))
}

// double              cairo_get_tolerance                 (cairo_t *cr);
func (this *Context) GetTolerance() float64 {
	return float64(C.cairo_get_tolerance(this.c()))
}

// void                cairo_clip                          (cairo_t *cr);
func (this *Context) Clip() {
	C.cairo_clip(this.c())
}

// void                cairo_clip_preserve                 (cairo_t *cr);
func (this *Context) ClipPreserve() {
	C.cairo_clip_preserve(this.c())
}

// void                cairo_clip_extents                  (cairo_t *cr,
//                                                          double *x1,
//                                                          double *y1,
//                                                          double *x2,
//                                                          double *y2);
func (this *Context) ClipExtents() (x1, y1, x2, y2 float64) {
	C.cairo_clip_extents(this.c(),
		(*C.double)(unsafe.Pointer(&x1)),
		(*C.double)(unsafe.Pointer(&y1)),
		(*C.double)(unsafe.Pointer(&x2)),
		(*C.double)(unsafe.Pointer(&y2)))
	return
}

// cairo_bool_t        cairo_in_clip                       (cairo_t *cr,
//                                                          double x,
//                                                          double y);
func (this *Context) InClip(x, y float64) bool {
	return C.cairo_in_clip(this.c(), C.double(x), C.double(y)) != 0
}

// void                cairo_reset_clip                    (cairo_t *cr);
func (this *Context) ResetClip() {
	C.cairo_reset_clip(this.c())
}

//                     cairo_rectangle_t;
type Rectangle struct {
	X, Y, Width, Height float64
}

//                     cairo_rectangle_list_t;
// void                cairo_rectangle_list_destroy        (cairo_rectangle_list_t *rectangle_list);

// cairo_rectangle_list_t * cairo_copy_clip_rectangle_list (cairo_t *cr);
func (this *Context) CopyClipRectangleList() ([]Rectangle, Status) {
	var slice []Rectangle
	var status Status

	rl := C.cairo_copy_clip_rectangle_list(this.c())
	if rl.num_rectangles > 0 {
		var slice_header reflect.SliceHeader
		slice_header.Data = uintptr(unsafe.Pointer(rl.rectangles))
		slice_header.Len = int(rl.num_rectangles)
		slice_header.Cap = int(rl.num_rectangles)
		slice_src := *(*[]Rectangle)(unsafe.Pointer(&slice_header))
		slice = make([]Rectangle, rl.num_rectangles)
		copy(slice, slice_src)
	}

	status = Status(rl.status)

	C.cairo_rectangle_list_destroy(rl)
	return slice, status
}

// void                cairo_fill                          (cairo_t *cr);
func (this *Context) Fill() {
	C.cairo_fill(this.c())
}

// void                cairo_fill_preserve                 (cairo_t *cr);
func (this *Context) FillPreserve() {
	C.cairo_fill_preserve(this.c())
}

// void                cairo_fill_extents                  (cairo_t *cr,
//                                                          double *x1,
//                                                          double *y1,
//                                                          double *x2,
//                                                          double *y2);
func (this *Context) FillExtents() (x1, y1, x2, y2 float64) {
	C.cairo_fill_extents(this.c(),
		(*C.double)(unsafe.Pointer(&x1)),
		(*C.double)(unsafe.Pointer(&y1)),
		(*C.double)(unsafe.Pointer(&x2)),
		(*C.double)(unsafe.Pointer(&y2)))
	return
}

// cairo_bool_t        cairo_in_fill                       (cairo_t *cr,
//                                                          double x,
//                                                          double y);
func (this *Context) InFill(x, y float64) bool {
	return C.cairo_in_fill(this.c(), C.double(x), C.double(y)) != 0
}

// void                cairo_mask                          (cairo_t *cr,
//                                                          cairo_pattern_t *pattern);
func (this *Context) Mask(pattern PatternLike) {
	C.cairo_mask(this.c(), pattern.InheritedFromCairoPattern().c())
}

// void                cairo_mask_surface                  (cairo_t *cr,
//                                                          cairo_surface_t *surface,
//                                                          double surface_x,
//                                                          double surface_y);
func (this *Context) MaskSurface(surface SurfaceLike, surface_x, surface_y float64) {
	C.cairo_mask_surface(this.c(), surface.InheritedFromCairoSurface().c(), C.double(surface_x), C.double(surface_y))
}

// void                cairo_paint                         (cairo_t *cr);
func (this *Context) Paint() {
	C.cairo_paint(this.c())
}

// void                cairo_paint_with_alpha              (cairo_t *cr,
//                                                          double alpha);
func (this *Context) PaintWithAlpha(alpha float64) {
	C.cairo_paint_with_alpha(this.c(), C.double(alpha))
}

// void                cairo_stroke                        (cairo_t *cr);
func (this *Context) Stroke() {
	C.cairo_stroke(this.c())
}

// void                cairo_stroke_preserve               (cairo_t *cr);
func (this *Context) StrokePreserve() {
	C.cairo_stroke_preserve(this.c())
}

// void                cairo_stroke_extents                (cairo_t *cr,
//                                                          double *x1,
//                                                          double *y1,
//                                                          double *x2,
//                                                          double *y2);
func (this *Context) StrokeExtents() (x1, y1, x2, y2 float64) {
	C.cairo_stroke_extents(this.c(),
		(*C.double)(unsafe.Pointer(&x1)),
		(*C.double)(unsafe.Pointer(&y1)),
		(*C.double)(unsafe.Pointer(&x2)),
		(*C.double)(unsafe.Pointer(&y2)))
	return
}

// cairo_bool_t        cairo_in_stroke                     (cairo_t *cr,
//                                                          double x,
//                                                          double y);
func (this *Context) InStroke(x, y float64) bool {
	return C.cairo_in_stroke(this.c(), C.double(x), C.double(y)) != 0
}

// void                cairo_copy_page                     (cairo_t *cr);
func (this *Context) CopyPage() {
	C.cairo_copy_page(this.c())
}

// void                cairo_show_page                     (cairo_t *cr);
func (this *Context) ShowPage() {
	C.cairo_show_page(this.c())
}

// unsigned int        cairo_get_reference_count           (cairo_t *cr);
// cairo_status_t      cairo_set_user_data                 (cairo_t *cr,
//                                                          const cairo_user_data_key_t *key,
//                                                          void *user_data,
//                                                          cairo_destroy_func_t destroy);
// void *              cairo_get_user_data                 (cairo_t *cr,
//                                                          const cairo_user_data_key_t *key);

//----------------------------------------------------------------------------
// Paths
//----------------------------------------------------------------------------

//                     cairo_path_t;
type Path struct {
	C unsafe.Pointer
}

func (this *Path) c() *C.cairo_path_t {
	return (*C.cairo_path_t)(this.C)
}

func path_finalizer(this *Path) {
	if gobject.FQueue.Push(unsafe.Pointer(this), path_finalizer2) {
		return
	}
	C.cairo_path_destroy(this.c())
}

func path_finalizer2(this_un unsafe.Pointer) {
	this := (*Path)(this_un)
	C.cairo_path_destroy(this.c())
}

func PathWrap(c_un unsafe.Pointer) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_path_t)(c_un)
	path := &Path{unsafe.Pointer(c)}
	runtime.SetFinalizer(path, path_finalizer)
	return unsafe.Pointer(path)
}

// TODO: Implement?
// union               cairo_path_data_t;

// enum                cairo_path_data_type_t;
type PathDataType int

const (
	PathMoveTo    PathDataType = C.CAIRO_PATH_MOVE_TO
	PathLineTo    PathDataType = C.CAIRO_PATH_LINE_TO
	PathCurveTo   PathDataType = C.CAIRO_PATH_CURVE_TO
	PathClosePath PathDataType = C.CAIRO_PATH_CLOSE_PATH
)

// cairo_path_t *      cairo_copy_path                     (cairo_t *cr);
func (this *Context) CopyPath() *Path {
	return (*Path)(PathWrap(unsafe.Pointer(C.cairo_copy_path(this.c()))))
}

// cairo_path_t *      cairo_copy_path_flat                (cairo_t *cr);
func (this *Context) CopyPathFlat() *Path {
	return (*Path)(PathWrap(unsafe.Pointer(C.cairo_copy_path_flat(this.c()))))
}

// void                cairo_path_destroy                  (cairo_path_t *path);
func (this *Path) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_path_destroy(this.c())
	this.C = nil
}

// void                cairo_append_path                   (cairo_t *cr,
//                                                          const cairo_path_t *path);
func (this *Context) AppendPath(path *Path) {
	C.cairo_append_path(this.c(), path.c())
}

// cairo_bool_t        cairo_has_current_point             (cairo_t *cr);
func (this *Context) HasCurrentPoint() bool {
	return C.cairo_has_current_point(this.c()) != 0
}

// void                cairo_get_current_point             (cairo_t *cr,
//                                                          double *x,
//                                                          double *y);
func (this *Context) GetCurrentPoint() (x, y float64) {
	C.cairo_get_current_point(this.c(),
		(*C.double)(unsafe.Pointer(&x)),
		(*C.double)(unsafe.Pointer(&y)))
	return
}

// void                cairo_new_path                      (cairo_t *cr);
func (this *Context) NewPath() {
	C.cairo_new_path(this.c())
}

// void                cairo_new_sub_path                  (cairo_t *cr);
func (this *Context) NewSubPath() {
	C.cairo_new_sub_path(this.c())
}

// void                cairo_close_path                    (cairo_t *cr);
func (this *Context) ClosePath() {
	C.cairo_close_path(this.c())
}

// void                cairo_arc                           (cairo_t *cr,
//                                                          double xc,
//                                                          double yc,
//                                                          double radius,
//                                                          double angle1,
//                                                          double angle2);
func (this *Context) Arc(xc, yc, radius, angle1, angle2 float64) {
	C.cairo_arc(this.c(), C.double(xc), C.double(yc), C.double(radius), C.double(angle1), C.double(angle2))
}

// void                cairo_arc_negative                  (cairo_t *cr,
//                                                          double xc,
//                                                          double yc,
//                                                          double radius,
//                                                          double angle1,
//                                                          double angle2);
func (this *Context) ArcNegative(xc, yc, radius, angle1, angle2 float64) {
	C.cairo_arc_negative(this.c(), C.double(xc), C.double(yc), C.double(radius), C.double(angle1), C.double(angle2))
}

// void                cairo_curve_to                      (cairo_t *cr,
//                                                          double x1,
//                                                          double y1,
//                                                          double x2,
//                                                          double y2,
//                                                          double x3,
//                                                          double y3);
func (this *Context) CurveTo(x1, y1, x2, y2, x3, y3 float64) {
	C.cairo_curve_to(this.c(), C.double(x1), C.double(y1), C.double(x2), C.double(y2), C.double(x3), C.double(y3))
}

// void                cairo_line_to                       (cairo_t *cr,
//                                                          double x,
//                                                          double y);
func (this *Context) LineTo(x, y float64) {
	C.cairo_line_to(this.c(), C.double(x), C.double(y))
}

// void                cairo_move_to                       (cairo_t *cr,
//                                                          double x,
//                                                          double y);
func (this *Context) MoveTo(x, y float64) {
	C.cairo_move_to(this.c(), C.double(x), C.double(y))
}

// void                cairo_rectangle                     (cairo_t *cr,
//                                                          double x,
//                                                          double y,
//                                                          double width,
//                                                          double height);
func (this *Context) Rectangle(x, y, width, height float64) {
	C.cairo_rectangle(this.c(), C.double(x), C.double(y), C.double(width), C.double(height))
}

// void                cairo_glyph_path                    (cairo_t *cr,
//                                                          const cairo_glyph_t *glyphs,
//                                                          int num_glyphs);
func (this *Context) GlyphPath(glyphs []Glyph) {
	var first *C.cairo_glyph_t
	var n = C.int(len(glyphs))
	if n > 0 {
		first = glyphs[0].c()
	}

	C.cairo_glyph_path(this.c(), first, n)
}

// void                cairo_text_path                     (cairo_t *cr,
//                                                          const char *utf8);
func (this *Context) TextPath(utf8 string) {
	utf8c := C.CString(utf8)
	C.cairo_text_path(this.c(), utf8c)
	C.free(unsafe.Pointer(utf8c))
}

// void                cairo_rel_curve_to                  (cairo_t *cr,
//                                                          double dx1,
//                                                          double dy1,
//                                                          double dx2,
//                                                          double dy2,
//                                                          double dx3,
//                                                          double dy3);
func (this *Context) RelCurveTo(dx1, dy1, dx2, dy2, dx3, dy3 float64) {
	C.cairo_rel_curve_to(this.c(), C.double(dx1), C.double(dy1), C.double(dx2), C.double(dy2), C.double(dx3), C.double(dy3))
}

// void                cairo_rel_line_to                   (cairo_t *cr,
//                                                          double dx,
//                                                          double dy);
func (this *Context) RelLineTo(dx, dy float64) {
	C.cairo_rel_line_to(this.c(), C.double(dx), C.double(dy))
}

// void                cairo_rel_move_to                   (cairo_t *cr,
//                                                          double dx,
//                                                          double dy);
func (this *Context) RelMoveTo(dx, dy float64) {
	C.cairo_rel_move_to(this.c(), C.double(dx), C.double(dy))
}

// void                cairo_path_extents                  (cairo_t *cr,
//                                                          double *x1,
//                                                          double *y1,
//                                                          double *x2,
//                                                          double *y2);
func (this *Context) PathExtents() (x1, y1, x2, y2 float64) {
	C.cairo_path_extents(this.c(),
		(*C.double)(unsafe.Pointer(&x1)),
		(*C.double)(unsafe.Pointer(&y1)),
		(*C.double)(unsafe.Pointer(&x2)),
		(*C.double)(unsafe.Pointer(&y2)))
	return
}

//----------------------------------------------------------------------------
// Pattern
// 	SolidPattern
// 	SurfacePattern
// 	Gradient
// 		LinearGradient
// 		RadialGradient
//----------------------------------------------------------------------------

// typedef             cairo_pattern_t;
type PatternLike interface {
	InheritedFromCairoPattern() *Pattern
}

type Pattern struct{ C unsafe.Pointer }
type SolidPattern struct{ Pattern }
type SurfacePattern struct{ Pattern }
type Gradient struct{ Pattern }
type LinearGradient struct{ Gradient }
type RadialGradient struct{ Gradient }

func (this *Pattern) c() *C.cairo_pattern_t {
	return (*C.cairo_pattern_t)(this.C)
}

func pattern_finalizer(this *Pattern) {
	if gobject.FQueue.Push(unsafe.Pointer(this), pattern_finalizer2) {
		return
	}
	C.cairo_pattern_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_pattern_destroy(this.c())
}

func pattern_finalizer2(this_un unsafe.Pointer) {
	this := (*Pattern)(this_un)
	C.cairo_pattern_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_pattern_destroy(this.c())
}

func PatternWrap(c_un unsafe.Pointer, grab bool) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_pattern_t)(c_un)
	go_repr := C.cairo_pattern_get_user_data(c, &go_repr_cookie)
	if go_repr != nil {
		return unsafe.Pointer(go_repr)
	}

	pattern := &Pattern{unsafe.Pointer(c)}
	if grab {
		C.cairo_pattern_reference(c)
	}
	runtime.SetFinalizer(pattern, pattern_finalizer)

	status := C.cairo_pattern_set_user_data(c, &go_repr_cookie, unsafe.Pointer(pattern), nil)
	if status != C.CAIRO_STATUS_SUCCESS {
		panic("failed to set user data, out of memory?")
	}
	return unsafe.Pointer(pattern)
}

func ensure_pattern_type(c *C.cairo_pattern_t, types ...PatternType) {
	for _, t := range types {
		if C.cairo_pattern_get_type(c) == t.c() {
			return
		}
	}
	panic("unexpected pattern type")
}

func ToPattern(like PatternLike) *Pattern {
	return (*Pattern)(PatternWrap(unsafe.Pointer(like.InheritedFromCairoPattern().c()), true))
}

func ToSolidPattern(like PatternLike) *SolidPattern {
	c := like.InheritedFromCairoPattern().c()
	ensure_pattern_type(c, PatternTypeSolid)
	return (*SolidPattern)(PatternWrap(unsafe.Pointer(c), true))
}

func ToSurfacePattern(like PatternLike) *SurfacePattern {
	c := like.InheritedFromCairoPattern().c()
	ensure_pattern_type(c, PatternTypeSurface)
	return (*SurfacePattern)(PatternWrap(unsafe.Pointer(c), true))
}

func ToGradient(like PatternLike) *Gradient {
	c := like.InheritedFromCairoPattern().c()
	ensure_pattern_type(c, PatternTypeLinear, PatternTypeRadial)
	return (*Gradient)(PatternWrap(unsafe.Pointer(c), true))
}

func ToLinearGradient(like PatternLike) *LinearGradient {
	c := like.InheritedFromCairoPattern().c()
	ensure_pattern_type(c, PatternTypeLinear)
	return (*LinearGradient)(PatternWrap(unsafe.Pointer(c), true))
}

func ToRadialGradient(like PatternLike) *RadialGradient {
	c := like.InheritedFromCairoPattern().c()
	ensure_pattern_type(c, PatternTypeRadial)
	return (*RadialGradient)(PatternWrap(unsafe.Pointer(c), true))
}

func (this *Pattern) InheritedFromCairoPattern() *Pattern {
	return this
}

// void                cairo_pattern_add_color_stop_rgb    (cairo_pattern_t *pattern,
//                                                          double offset,
//                                                          double red,
//                                                          double green,
//                                                          double blue);
func (this *Gradient) AddColorStopRGB(offset, red, green, blue float64) {
	C.cairo_pattern_add_color_stop_rgb(this.c(),
		C.double(offset), C.double(red), C.double(green), C.double(blue))
}

// void                cairo_pattern_add_color_stop_rgba   (cairo_pattern_t *pattern,
//                                                          double offset,
//                                                          double red,
//                                                          double green,
//                                                          double blue,
//                                                          double alpha);
func (this *Gradient) AddColorStopRGBA(offset, red, green, blue, alpha float64) {
	C.cairo_pattern_add_color_stop_rgba(this.c(),
		C.double(offset), C.double(red), C.double(green), C.double(blue), C.double(alpha))
}

// cairo_status_t      cairo_pattern_get_color_stop_count  (cairo_pattern_t *pattern,
//                                                          int *count);
func (this *Gradient) GetColorStopCount() (count int) {
	C.cairo_pattern_get_color_stop_count(this.c(),
		(*C.int)(unsafe.Pointer(&count)))
	return
}

// cairo_status_t      cairo_pattern_get_color_stop_rgba   (cairo_pattern_t *pattern,
//                                                          int index,
//                                                          double *offset,
//                                                          double *red,
//                                                          double *green,
//                                                          double *blue,
//                                                          double *alpha);
func (this *Gradient) GetColorStopRGBA(index int) (offset, red, green, blue, alpha float64) {
	C.cairo_pattern_get_color_stop_rgba(this.c(), C.int(index),
		(*C.double)(unsafe.Pointer(&offset)),
		(*C.double)(unsafe.Pointer(&red)),
		(*C.double)(unsafe.Pointer(&green)),
		(*C.double)(unsafe.Pointer(&blue)),
		(*C.double)(unsafe.Pointer(&alpha)))
	return
}

// cairo_pattern_t *   cairo_pattern_create_rgb            (double red,
//                                                          double green,
//                                                          double blue);
func NewSolidPatternRGB(red, green, blue float64) *SolidPattern {
	return (*SolidPattern)(PatternWrap(unsafe.Pointer(C.cairo_pattern_create_rgb(C.double(red), C.double(green), C.double(blue))), false))
}

// cairo_pattern_t *   cairo_pattern_create_rgba           (double red,
//                                                          double green,
//                                                          double blue,
//                                                          double alpha);
func NewSolidPatternRGBA(red, green, blue, alpha float64) *SolidPattern {
	return (*SolidPattern)(PatternWrap(unsafe.Pointer(C.cairo_pattern_create_rgba(C.double(red), C.double(green), C.double(blue), C.double(alpha))), false))
}

// cairo_status_t      cairo_pattern_get_rgba              (cairo_pattern_t *pattern,
//                                                          double *red,
//                                                          double *green,
//                                                          double *blue,
//                                                          double *alpha);
func (this *SolidPattern) GetRGBA() (red, green, blue, alpha float64) {
	C.cairo_pattern_get_rgba(this.c(),
		(*C.double)(unsafe.Pointer(&red)),
		(*C.double)(unsafe.Pointer(&green)),
		(*C.double)(unsafe.Pointer(&blue)),
		(*C.double)(unsafe.Pointer(&alpha)))
	return
}

// cairo_pattern_t *   cairo_pattern_create_for_surface    (cairo_surface_t *surface);
func NewSurfacePattern(surface SurfaceLike) *SurfacePattern {
	return (*SurfacePattern)(PatternWrap(unsafe.Pointer(C.cairo_pattern_create_for_surface(surface.InheritedFromCairoSurface().c())), false))
}

// cairo_status_t      cairo_pattern_get_surface           (cairo_pattern_t *pattern,
//                                                          cairo_surface_t **surface);
func (this *SurfacePattern) GetSurface() *Surface {
	var surfacec *C.cairo_surface_t
	C.cairo_pattern_get_surface(this.c(), &surfacec)
	surface := (*Surface)(SurfaceWrap(unsafe.Pointer(surfacec), true))
	return surface
}

// cairo_pattern_t *   cairo_pattern_create_linear         (double x0,
//                                                          double y0,
//                                                          double x1,
//                                                          double y1);
func NewLinearGradient(x0, y0, x1, y1 float64) *LinearGradient {
	return (*LinearGradient)(PatternWrap(unsafe.Pointer(C.cairo_pattern_create_linear(
		C.double(x0), C.double(y0), C.double(x1), C.double(y1))), false))
}

// cairo_status_t      cairo_pattern_get_linear_points     (cairo_pattern_t *pattern,
//                                                          double *x0,
//                                                          double *y0,
//                                                          double *x1,
//                                                          double *y1);
func (this *LinearGradient) GetLinearPoints() (x0, y0, x1, y1 float64) {
	C.cairo_pattern_get_linear_points(this.c(),
		(*C.double)(unsafe.Pointer(&x0)),
		(*C.double)(unsafe.Pointer(&y0)),
		(*C.double)(unsafe.Pointer(&x1)),
		(*C.double)(unsafe.Pointer(&y1)))
	return
}

// cairo_pattern_t *   cairo_pattern_create_radial         (double cx0,
//                                                          double cy0,
//                                                          double radius0,
//                                                          double cx1,
//                                                          double cy1,
//                                                          double radius1);
func NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1 float64) *RadialGradient {
	return (*RadialGradient)(PatternWrap(unsafe.Pointer(C.cairo_pattern_create_radial(
		C.double(cx0), C.double(cy0), C.double(radius0), C.double(cx1), C.double(cy1), C.double(radius1))), false))
}

// cairo_status_t      cairo_pattern_get_radial_circles    (cairo_pattern_t *pattern,
//                                                          double *x0,
//                                                          double *y0,
//                                                          double *r0,
//                                                          double *x1,
//                                                          double *y1,
//                                                          double *r1);
func (this *RadialGradient) GetRadialCircles() (x0, y0, r0, x1, y1, r1 float64) {
	C.cairo_pattern_get_radial_circles(this.c(),
		(*C.double)(unsafe.Pointer(&x0)),
		(*C.double)(unsafe.Pointer(&y0)),
		(*C.double)(unsafe.Pointer(&r0)),
		(*C.double)(unsafe.Pointer(&x1)),
		(*C.double)(unsafe.Pointer(&y1)),
		(*C.double)(unsafe.Pointer(&r1)))
	return
}

// cairo_pattern_t *   cairo_pattern_reference             (cairo_pattern_t *pattern);
// void                cairo_pattern_destroy               (cairo_pattern_t *pattern);
func (this *Pattern) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_pattern_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_pattern_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_pattern_status                (cairo_pattern_t *pattern);
func (this *Pattern) Status() Status {
	return Status(C.cairo_pattern_status(this.c()))
}

// enum                cairo_extend_t;
type Extend int

const (
	ExtendNone    Extend = C.CAIRO_EXTEND_NONE
	ExtendRepeat  Extend = C.CAIRO_EXTEND_REPEAT
	ExtendReflect Extend = C.CAIRO_EXTEND_REFLECT
	ExtendPad     Extend = C.CAIRO_EXTEND_PAD
)

func (this Extend) c() C.cairo_extend_t {
	return C.cairo_extend_t(this)
}

// void                cairo_pattern_set_extend            (cairo_pattern_t *pattern,
//                                                          cairo_extend_t extend);
func (this *SurfacePattern) SetExtend(extend Extend) {
	C.cairo_pattern_set_extend(this.c(), extend.c())
}

// cairo_extend_t      cairo_pattern_get_extend            (cairo_pattern_t *pattern);
func (this *SurfacePattern) GetExtend() Extend {
	return Extend(C.cairo_pattern_get_extend(this.c()))
}

// enum                cairo_filter_t;
type Filter int

const (
	FilterFast     Filter = C.CAIRO_FILTER_FAST
	FilterGood     Filter = C.CAIRO_FILTER_GOOD
	FilterBest     Filter = C.CAIRO_FILTER_BEST
	FilterNearest  Filter = C.CAIRO_FILTER_NEAREST
	FilterBilinear Filter = C.CAIRO_FILTER_BILINEAR
	FilterGaussian Filter = C.CAIRO_FILTER_GAUSSIAN
)

func (this Filter) c() C.cairo_filter_t {
	return C.cairo_filter_t(this)
}

// void                cairo_pattern_set_filter            (cairo_pattern_t *pattern,
//                                                          cairo_filter_t filter);
func (this *SurfacePattern) SetFilter(filter Filter) {
	C.cairo_pattern_set_filter(this.c(), filter.c())
}

// cairo_filter_t      cairo_pattern_get_filter            (cairo_pattern_t *pattern);
func (this *SurfacePattern) GetFilter() Filter {
	return Filter(C.cairo_pattern_get_filter(this.c()))
}

// void                cairo_pattern_set_matrix            (cairo_pattern_t *pattern,
//                                                          const cairo_matrix_t *matrix);
func (this *Pattern) SetMatrix(matrix *Matrix) {
	C.cairo_pattern_set_matrix(this.c(), matrix.c())
}

// void                cairo_pattern_get_matrix            (cairo_pattern_t *pattern,
//                                                          cairo_matrix_t *matrix);
func (this *Pattern) GetMatrix() Matrix {
	var matrix C.cairo_matrix_t
	C.cairo_pattern_get_matrix(this.c(), &matrix)
	return *(*Matrix)(unsafe.Pointer(&matrix))
}

// enum                cairo_pattern_type_t;
type PatternType int

const (
	PatternTypeSolid   PatternType = C.CAIRO_PATTERN_TYPE_SOLID
	PatternTypeSurface PatternType = C.CAIRO_PATTERN_TYPE_SURFACE
	PatternTypeLinear  PatternType = C.CAIRO_PATTERN_TYPE_LINEAR
	PatternTypeRadial  PatternType = C.CAIRO_PATTERN_TYPE_RADIAL
)

func (this PatternType) c() C.cairo_pattern_type_t {
	return C.cairo_pattern_type_t(this)
}

// cairo_pattern_type_t  cairo_pattern_get_type            (cairo_pattern_t *pattern);
func (this *Pattern) GetType() PatternType {
	return PatternType(C.cairo_pattern_get_type(this.c()))
}

// TODO: Implement these
// unsigned int        cairo_pattern_get_reference_count   (cairo_pattern_t *pattern);
// cairo_status_t      cairo_pattern_set_user_data         (cairo_pattern_t *pattern,
//                                                          const cairo_user_data_key_t *key,
//                                                          void *user_data,
//                                                          cairo_destroy_func_t destroy);
// void *              cairo_pattern_get_user_data         (cairo_pattern_t *pattern,
//                                                          const cairo_user_data_key_t *key);

//----------------------------------------------------------------------------
// Regions
//----------------------------------------------------------------------------

// typedef             cairo_region_t;
type Region struct {
	C unsafe.Pointer
}

func (this *Region) c() *C.cairo_region_t {
	return (*C.cairo_region_t)(this.C)
}

func region_finalizer(this *Region) {
	if gobject.FQueue.Push(unsafe.Pointer(this), region_finalizer2) {
		return
	}
	C.cairo_region_destroy(this.c())
}

func region_finalizer2(this_un unsafe.Pointer) {
	this := (*Region)(this_un)
	C.cairo_region_destroy(this.c())
}

func RegionWrap(c_un unsafe.Pointer, grab bool) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_region_t)(c_un)
	if grab {
		C.cairo_region_reference(c)
	}
	region := &Region{unsafe.Pointer(c)}
	runtime.SetFinalizer(region, region_finalizer)
	return unsafe.Pointer(region)
}

// cairo_region_t *    cairo_region_create                 (void);
func NewRegion() *Region {
	return (*Region)(RegionWrap(unsafe.Pointer(C.cairo_region_create()), false))
}

// cairo_region_t *    cairo_region_create_rectangle       (const cairo_rectangle_int_t *rectangle);
func NewRegionRectangle(rectangle *RectangleInt) *Region {
	return (*Region)(RegionWrap(unsafe.Pointer(C.cairo_region_create_rectangle(rectangle.c())), false))
}

// cairo_region_t *    cairo_region_create_rectangles      (const cairo_rectangle_int_t *rects,
//                                                          int count);
func NewRegionRectangles(rects []RectangleInt) *Region {
	var first *C.cairo_rectangle_int_t
	count := C.int(len(rects))
	if count > 0 {
		first = rects[0].c()
	}
	return (*Region)(RegionWrap(unsafe.Pointer(C.cairo_region_create_rectangles(first, count)), false))
}

// cairo_region_t *    cairo_region_copy                   (const cairo_region_t *original);
func (this *Region) Copy() *Region {
	return (*Region)(RegionWrap(unsafe.Pointer(C.cairo_region_copy(this.c())), false))
}

// cairo_region_t *    cairo_region_reference              (cairo_region_t *region);
// void                cairo_region_destroy                (cairo_region_t *region);
func (this *Region) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_region_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_region_status                 (const cairo_region_t *region);
func (this *Region) Status() Status {
	return Status(C.cairo_region_status(this.c()))
}

// void                cairo_region_get_extents            (const cairo_region_t *region,
//                                                          cairo_rectangle_int_t *extents);
func (this *Region) GetExtents() (extents RectangleInt) {
	C.cairo_region_get_extents(this.c(), extents.c())
	return
}

// int                 cairo_region_num_rectangles         (const cairo_region_t *region);
func (this *Region) NumRectangles() int {
	return int(C.cairo_region_num_rectangles(this.c()))
}

// void                cairo_region_get_rectangle          (const cairo_region_t *region,
//                                                          int nth,
//                                                          cairo_rectangle_int_t *rectangle);
func (this *Region) GetRectangle(nth int) (rectangle RectangleInt) {
	C.cairo_region_get_rectangle(this.c(), C.int(nth), rectangle.c())
	return
}

// cairo_bool_t        cairo_region_is_empty               (const cairo_region_t *region);
func (this *Region) IsEmpty() bool {
	return C.cairo_region_is_empty(this.c()) != 0
}

// cairo_bool_t        cairo_region_contains_point         (const cairo_region_t *region,
//                                                          int x,
//                                                          int y);
func (this *Region) ContainsPoint(x, y int) bool {
	return C.cairo_region_contains_point(this.c(), C.int(x), C.int(y)) != 0
}

// enum                cairo_region_overlap_t;
type RegionOverlap int

const (
	RegionOverlapIn   RegionOverlap = C.CAIRO_REGION_OVERLAP_IN   /* completely inside region */
	RegionOverlapOut  RegionOverlap = C.CAIRO_REGION_OVERLAP_OUT  /* completely outside region */
	RegionOverlapPart RegionOverlap = C.CAIRO_REGION_OVERLAP_PART /* partly inside region */
)

// cairo_region_overlap_t  cairo_region_contains_rectangle (const cairo_region_t *region,
//                                                          const cairo_rectangle_int_t *rectangle);
func (this *Region) ContainsRectangle(rectangle *RectangleInt) RegionOverlap {
	return RegionOverlap(C.cairo_region_contains_rectangle(this.c(), rectangle.c()))
}

// cairo_bool_t        cairo_region_equal                  (const cairo_region_t *a,
//                                                          const cairo_region_t *b);
func (this *Region) Equal(b *Region) bool {
	return C.cairo_region_equal(this.c(), b.c()) != 0
}

// void                cairo_region_translate              (cairo_region_t *region,
//                                                          int dx,
//                                                          int dy);
func (this *Region) Translate(dx, dy int) {
	C.cairo_region_translate(this.c(), C.int(dx), C.int(dy))
}

// cairo_status_t      cairo_region_intersect              (cairo_region_t *dst,
//                                                          const cairo_region_t *other);
func (this *Region) Intersect(other *Region) Status {
	return Status(C.cairo_region_intersect(this.c(), other.c()))
}

// cairo_status_t      cairo_region_intersect_rectangle    (cairo_region_t *dst,
//                                                          const cairo_rectangle_int_t *rectangle);
func (this *Region) IntersectRectangle(rectangle *RectangleInt) Status {
	return Status(C.cairo_region_intersect_rectangle(this.c(), rectangle.c()))
}

// cairo_status_t      cairo_region_subtract               (cairo_region_t *dst,
//                                                          const cairo_region_t *other);
func (this *Region) Subtract(other *Region) Status {
	return Status(C.cairo_region_subtract(this.c(), other.c()))
}

// cairo_status_t      cairo_region_subtract_rectangle     (cairo_region_t *dst,
//                                                          const cairo_rectangle_int_t *rectangle);
func (this *Region) SubtractRectangle(rectangle *RectangleInt) Status {
	return Status(C.cairo_region_subtract_rectangle(this.c(), rectangle.c()))
}

// cairo_status_t      cairo_region_union                  (cairo_region_t *dst,
//                                                          const cairo_region_t *other);
func (this *Region) Union(other *Region) Status {
	return Status(C.cairo_region_union(this.c(), other.c()))
}

// cairo_status_t      cairo_region_union_rectangle        (cairo_region_t *dst,
//                                                          const cairo_rectangle_int_t *rectangle);
func (this *Region) UnionRectangle(rectangle *RectangleInt) Status {
	return Status(C.cairo_region_union_rectangle(this.c(), rectangle.c()))
}

// cairo_status_t      cairo_region_xor                    (cairo_region_t *dst,
//                                                          const cairo_region_t *other);
func (this *Region) Xor(other *Region) Status {
	return Status(C.cairo_region_xor(this.c(), other.c()))
}

// cairo_status_t      cairo_region_xor_rectangle          (cairo_region_t *dst,
//                                                          const cairo_rectangle_int_t *rectangle);
func (this *Region) XorRectangle(rectangle *RectangleInt) Status {
	return Status(C.cairo_region_xor_rectangle(this.c(), rectangle.c()))
}

//----------------------------------------------------------------------------
// Transformations
//----------------------------------------------------------------------------

// void                cairo_translate                     (cairo_t *cr,
//                                                          double tx,
//                                                          double ty);
func (this *Context) Translate(tx, ty float64) {
	C.cairo_translate(this.c(), C.double(tx), C.double(ty))
}

// void                cairo_scale                         (cairo_t *cr,
//                                                          double sx,
//                                                          double sy);
func (this *Context) Scale(sx, sy float64) {
	C.cairo_scale(this.c(), C.double(sx), C.double(sy))
}

// void                cairo_rotate                        (cairo_t *cr,
//                                                          double angle);
func (this *Context) Rotate(angle float64) {
	C.cairo_rotate(this.c(), C.double(angle))
}

// void                cairo_transform                     (cairo_t *cr,
//                                                          const cairo_matrix_t *matrix);
func (this *Context) Transform(matrix *Matrix) {
	C.cairo_transform(this.c(), matrix.c())
}

// void                cairo_set_matrix                    (cairo_t *cr,
//                                                          const cairo_matrix_t *matrix);
func (this *Context) SetMatrix(matrix *Matrix) {
	C.cairo_set_matrix(this.c(), matrix.c())
}

// void                cairo_get_matrix                    (cairo_t *cr,
//                                                          cairo_matrix_t *matrix);
func (this *Context) GetMatrix() Matrix {
	var matrix C.cairo_matrix_t
	C.cairo_get_matrix(this.c(), &matrix)
	return *(*Matrix)(unsafe.Pointer(&matrix))
}

// void                cairo_identity_matrix               (cairo_t *cr);
func (this *Context) IdentityMatrix() {
	C.cairo_identity_matrix(this.c())
}

// void                cairo_user_to_device                (cairo_t *cr,
//                                                          double *x,
//                                                          double *y);
func (this *Context) UserToDevice(x, y float64) (float64, float64) {
	C.cairo_user_to_device(this.c(),
		(*C.double)(unsafe.Pointer(&x)),
		(*C.double)(unsafe.Pointer(&y)))
	return x, y
}

// void                cairo_user_to_device_distance       (cairo_t *cr,
//                                                          double *dx,
//                                                          double *dy);
func (this *Context) UserToDeviceDistance(dx, dy float64) (float64, float64) {
	C.cairo_user_to_device_distance(this.c(),
		(*C.double)(unsafe.Pointer(&dx)),
		(*C.double)(unsafe.Pointer(&dy)))
	return dx, dy
}

// void                cairo_device_to_user                (cairo_t *cr,
//                                                          double *x,
//                                                          double *y);
func (this *Context) DeviceToUser(x, y float64) (float64, float64) {
	C.cairo_device_to_user(this.c(),
		(*C.double)(unsafe.Pointer(&x)),
		(*C.double)(unsafe.Pointer(&y)))
	return x, y
}

// void                cairo_device_to_user_distance       (cairo_t *cr,
//                                                          double *dx,
//                                                          double *dy);
func (this *Context) DeviceToUserDistance(dx, dy float64) (float64, float64) {
	C.cairo_device_to_user_distance(this.c(),
		(*C.double)(unsafe.Pointer(&dx)),
		(*C.double)(unsafe.Pointer(&dy)))
	return dx, dy
}

//----------------------------------------------------------------------------
// Text
//----------------------------------------------------------------------------

//                     cairo_glyph_t;
// type Glyph (see types_$GOARCH.go)

func (this *Glyph) c() *C.cairo_glyph_t {
	return (*C.cairo_glyph_t)(unsafe.Pointer(this))
}

// enum                cairo_font_slant_t;
type FontSlant int

const (
	FontSlantNormal  FontSlant = C.CAIRO_FONT_SLANT_NORMAL
	FontSlantItalic  FontSlant = C.CAIRO_FONT_SLANT_ITALIC
	FontSlantOblique FontSlant = C.CAIRO_FONT_SLANT_OBLIQUE
)

func (this FontSlant) c() C.cairo_font_slant_t {
	return C.cairo_font_slant_t(this)
}

// enum                cairo_font_weight_t;
type FontWeight int

const (
	FontWeightNormal FontWeight = C.CAIRO_FONT_WEIGHT_NORMAL
	FontWeightBold   FontWeight = C.CAIRO_FONT_WEIGHT_BOLD
)

func (this FontWeight) c() C.cairo_font_weight_t {
	return C.cairo_font_weight_t(this)
}

//                     cairo_text_cluster_t;
type TextCluster struct {
	NumBytes  int32
	NumGlyphs int32
}

// enum                cairo_text_cluster_flags_t;
type TextClusterFlags int

const (
	TextClusterFlagBackward TextClusterFlags = C.CAIRO_TEXT_CLUSTER_FLAG_BACKWARD
)

// void                cairo_select_font_face              (cairo_t *cr,
//                                                          const char *family,
//                                                          cairo_font_slant_t slant,
//                                                          cairo_font_weight_t weight);
func (this *Context) SelectFontFace(family string, slant FontSlant, weight FontWeight) {
	cfamily := C.CString(family)
	C.cairo_select_font_face(this.c(), cfamily, slant.c(), weight.c())
	C.free(unsafe.Pointer(cfamily))
}

// void                cairo_set_font_size                 (cairo_t *cr,
//                                                          double size);
func (this *Context) SetFontSize(size float64) {
	C.cairo_set_font_size(this.c(), C.double(size))
}

// void                cairo_set_font_matrix               (cairo_t *cr,
//                                                          const cairo_matrix_t *matrix);
func (this *Context) SetFontMatrix(matrix *Matrix) {
	C.cairo_set_font_matrix(this.c(), matrix.c())
}

// void                cairo_get_font_matrix               (cairo_t *cr,
//                                                          cairo_matrix_t *matrix);
func (this *Context) GetFontMatrix() Matrix {
	var matrix Matrix
	C.cairo_get_font_matrix(this.c(), matrix.c())
	return matrix
}

// TODO: Implement these
// void                cairo_set_font_options              (cairo_t *cr,
//                                                          const cairo_font_options_t *options);
// void                cairo_get_font_options              (cairo_t *cr,
//                                                          cairo_font_options_t *options);
// void                cairo_set_font_face                 (cairo_t *cr,
//                                                          cairo_font_face_t *font_face);
// cairo_font_face_t * cairo_get_font_face                 (cairo_t *cr);
// void                cairo_set_scaled_font               (cairo_t *cr,
//                                                          const cairo_scaled_font_t *scaled_font);
// cairo_scaled_font_t * cairo_get_scaled_font             (cairo_t *cr);

// void                cairo_show_text                     (cairo_t *cr,
//                                                          const char *utf8);
func (this *Context) ShowText(utf8 string) {
	cutf8 := C.CString(utf8)
	C.cairo_show_text(this.c(), cutf8)
	C.free(unsafe.Pointer(cutf8))
}

// void                cairo_show_glyphs                   (cairo_t *cr,
//                                                          const cairo_glyph_t *glyphs,
//                                                          int num_glyphs);
func (this *Context) ShowGlyphs(glyphs []Glyph) {
	var first *C.cairo_glyph_t
	var n = C.int(len(glyphs))
	if n > 0 {
		first = glyphs[0].c()
	}

	C.cairo_show_glyphs(this.c(), first, n)
}

// TODO: Implement these
// void                cairo_show_text_glyphs              (cairo_t *cr,
//                                                          const char *utf8,
//                                                          int utf8_len,
//                                                          const cairo_glyph_t *glyphs,
//                                                          int num_glyphs,
//                                                          const cairo_text_cluster_t *clusters,
//                                                          int num_clusters,
//                                                          cairo_text_cluster_flags_t cluster_flags);
// void                cairo_font_extents                  (cairo_t *cr,
//                                                          cairo_font_extents_t *extents);
// void                cairo_text_extents                  (cairo_t *cr,
//                                                          const char *utf8,
//                                                          cairo_text_extents_t *extents);
// void                cairo_glyph_extents                 (cairo_t *cr,
//                                                          const cairo_glyph_t *glyphs,
//                                                          int num_glyphs,
//                                                          cairo_text_extents_t *extents);
// cairo_font_face_t * cairo_toy_font_face_create          (const char *family,
//                                                          cairo_font_slant_t slant,
//                                                          cairo_font_weight_t weight);
// const char *        cairo_toy_font_face_get_family      (cairo_font_face_t *font_face);
// cairo_font_slant_t  cairo_toy_font_face_get_slant       (cairo_font_face_t *font_face);
// cairo_font_weight_t  cairo_toy_font_face_get_weight     (cairo_font_face_t *font_face);
// cairo_glyph_t *     cairo_glyph_allocate                (int num_glyphs);
// void                cairo_glyph_free                    (cairo_glyph_t *glyphs);
// cairo_text_cluster_t * cairo_text_cluster_allocate      (int num_clusters);
// void                cairo_text_cluster_free             (cairo_text_cluster_t *clusters);

//----------------------------------------------------------------------------
// FontFace
//----------------------------------------------------------------------------

// typedef             cairo_font_face_t;
type FontFace struct {
	C unsafe.Pointer
}

func (this *FontFace) c() *C.cairo_font_face_t {
	return (*C.cairo_font_face_t)(this.C)
}

func font_face_finalizer(this *FontFace) {
	if gobject.FQueue.Push(unsafe.Pointer(this), font_face_finalizer2) {
		return
	}
	C.cairo_font_face_destroy(this.c())
}

func font_face_finalizer2(this_un unsafe.Pointer) {
	this := (*FontFace)(this_un)
	C.cairo_font_face_destroy(this.c())
}

func FontFaceWrap(c_un unsafe.Pointer, grab bool) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_font_face_t)(c_un)
	go_repr := C.cairo_font_face_get_user_data(c, &go_repr_cookie)
	if go_repr != nil {
		return unsafe.Pointer(go_repr)
	}

	font_face := &FontFace{unsafe.Pointer(c)}
	if grab {
		C.cairo_font_face_reference(c)
	}
	runtime.SetFinalizer(font_face, font_face_finalizer)

	status := C.cairo_font_face_set_user_data(c, &go_repr_cookie, unsafe.Pointer(font_face), nil)
	if status != C.CAIRO_STATUS_SUCCESS {
		panic("failed to set user data, out of memory?")
	}
	return unsafe.Pointer(font_face)
}

// cairo_font_face_t * cairo_font_face_reference           (cairo_font_face_t *font_face);
// void                cairo_font_face_destroy             (cairo_font_face_t *font_face);
func (this *FontFace) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_font_face_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_font_face_status              (cairo_font_face_t *font_face);
func (this *FontFace) Status() Status {
	return Status(C.cairo_font_face_status(this.c()))
}

// enum                cairo_font_type_t;
type FontType int

const (
	FontTypeToy    FontType = C.CAIRO_FONT_TYPE_TOY
	FontTypeFT     FontType = C.CAIRO_FONT_TYPE_FT
	FontTypeWin32  FontType = C.CAIRO_FONT_TYPE_WIN32
	FontTypeQuartz FontType = C.CAIRO_FONT_TYPE_QUARTZ
	FontTypeUser   FontType = C.CAIRO_FONT_TYPE_USER
)

// cairo_font_type_t   cairo_font_face_get_type            (cairo_font_face_t *font_face);
func (this *FontFace) GetType() FontType {
	return FontType(C.cairo_font_face_get_type(this.c()))
}

// unsigned int        cairo_font_face_get_reference_count (cairo_font_face_t *font_face);
// cairo_status_t      cairo_font_face_set_user_data       (cairo_font_face_t *font_face,
//                                                          const cairo_user_data_key_t *key,
//                                                          void *user_data,
//                                                          cairo_destroy_func_t destroy);
// void *              cairo_font_face_get_user_data       (cairo_font_face_t *font_face,
//                                                          const cairo_user_data_key_t *key);

//----------------------------------------------------------------------------
// ScaledFont
//----------------------------------------------------------------------------

// typedef             cairo_scaled_font_t;
type ScaledFont struct {
	C unsafe.Pointer
}

func (this *ScaledFont) c() *C.cairo_scaled_font_t {
	return (*C.cairo_scaled_font_t)(this.C)
}

func scaled_font_finalizer(this *ScaledFont) {
	if gobject.FQueue.Push(unsafe.Pointer(this), scaled_font_finalizer2) {
		return
	}
	C.cairo_scaled_font_destroy(this.c())
}

func scaled_font_finalizer2(this_un unsafe.Pointer) {
	this := (*ScaledFont)(this_un)
	C.cairo_scaled_font_destroy(this.c())
}

func ScaledFontWrap(c_un unsafe.Pointer, grab bool) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_scaled_font_t)(c_un)
	go_repr := C.cairo_scaled_font_get_user_data(c, &go_repr_cookie)
	if go_repr != nil {
		return unsafe.Pointer(go_repr)
	}

	scaled_font := &ScaledFont{unsafe.Pointer(c)}
	if grab {
		C.cairo_scaled_font_reference(c)
	}
	runtime.SetFinalizer(scaled_font, scaled_font_finalizer)

	status := C.cairo_scaled_font_set_user_data(c, &go_repr_cookie, unsafe.Pointer(scaled_font), nil)
	if status != C.CAIRO_STATUS_SUCCESS {
		panic("failed to set user data, out of memory?")
	}
	return unsafe.Pointer(scaled_font)
}

// cairo_scaled_font_t * cairo_scaled_font_create          (cairo_font_face_t *font_face,
//                                                          const cairo_matrix_t *font_matrix,
//                                                          const cairo_matrix_t *ctm,
//                                                          const cairo_font_options_t *options);
func NewScaledFont(font_face *FontFace, font_matrix, ctm *Matrix, options *FontOptions) *ScaledFont {
	return (*ScaledFont)(ScaledFontWrap(unsafe.Pointer(C.cairo_scaled_font_create(
		font_face.c(), font_matrix.c(), ctm.c(), options.c())), false))

}

// cairo_scaled_font_t * cairo_scaled_font_reference       (cairo_scaled_font_t *scaled_font);
// void                cairo_scaled_font_destroy           (cairo_scaled_font_t *scaled_font);
func (this *ScaledFont) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_scaled_font_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_scaled_font_status            (cairo_scaled_font_t *scaled_font);
func (this *ScaledFont) Status() Status {
	return Status(C.cairo_scaled_font_status(this.c()))
}

//                     cairo_font_extents_t;
type FontExtents struct {
	Ascent      float64
	Descent     float64
	Height      float64
	MaxXAdvance float64
	MaxYAdvance float64
}

func (this *FontExtents) c() *C.cairo_font_extents_t {
	return (*C.cairo_font_extents_t)(unsafe.Pointer(this))
}

// void                cairo_scaled_font_extents           (cairo_scaled_font_t *scaled_font,
//                                                          cairo_font_extents_t *extents);
func (this *ScaledFont) Extents() (out FontExtents) {
	C.cairo_scaled_font_extents(this.c(), out.c())
	return
}

//                     cairo_text_extents_t;
type TextExtents struct {
	XBearing float64
	YBearing float64
	Width    float64
	Height   float64
	XAdvance float64
	YAdvance float64
}

func (this *TextExtents) c() *C.cairo_text_extents_t {
	return (*C.cairo_text_extents_t)(unsafe.Pointer(this))
}

// void                cairo_scaled_font_text_extents      (cairo_scaled_font_t *scaled_font,
//                                                          const char *utf8,
//                                                          cairo_text_extents_t *extents);
func (this *ScaledFont) TextExtents(utf8 string) (out TextExtents) {
	utf8c := C.CString(utf8)
	C.cairo_scaled_font_text_extents(this.c(), utf8c, out.c())
	C.free(unsafe.Pointer(utf8c))
	return
}

// void                cairo_scaled_font_glyph_extents     (cairo_scaled_font_t *scaled_font,
//                                                          const cairo_glyph_t *glyphs,
//                                                          int num_glyphs,
//                                                          cairo_text_extents_t *extents);
func (this *ScaledFont) GlyphExtents(glyphs []Glyph) (out TextExtents) {
	var first *C.cairo_glyph_t
	var n = C.int(len(glyphs))
	if n > 0 {
		first = glyphs[0].c()
	}

	C.cairo_scaled_font_glyph_extents(this.c(), first, n, out.c())
	return
}

// TODO
// cairo_status_t      cairo_scaled_font_text_to_glyphs    (cairo_scaled_font_t *scaled_font,
//                                                          double x,
//                                                          double y,
//                                                          const char *utf8,
//                                                          int utf8_len,
//                                                          cairo_glyph_t **glyphs,
//                                                          int *num_glyphs,
//                                                          cairo_text_cluster_t **clusters,
//                                                          int *num_clusters,
//                                                          cairo_text_cluster_flags_t *cluster_flags);

// cairo_font_face_t * cairo_scaled_font_get_font_face     (cairo_scaled_font_t *scaled_font);
func (this *ScaledFont) GetFontFace() *FontFace {
	return (*FontFace)(FontFaceWrap(unsafe.Pointer(C.cairo_scaled_font_get_font_face(this.c())), true))
}

// void                cairo_scaled_font_get_font_options  (cairo_scaled_font_t *scaled_font,
//                                                          cairo_font_options_t *options);
func (this *ScaledFont) GetFontOptions() *FontOptions {
	fo := NewFontOptions()
	C.cairo_scaled_font_get_font_options(this.c(), fo.c())
	return fo
}

// void                cairo_scaled_font_get_font_matrix   (cairo_scaled_font_t *scaled_font,
//                                                          cairo_matrix_t *font_matrix);
func (this *ScaledFont) GetFontMatrix() (out Matrix) {
	C.cairo_scaled_font_get_font_matrix(this.c(), out.c())
	return
}

// void                cairo_scaled_font_get_ctm           (cairo_scaled_font_t *scaled_font,
//                                                          cairo_matrix_t *ctm);
func (this *ScaledFont) GetCTM() (out Matrix) {
	C.cairo_scaled_font_get_ctm(this.c(), out.c())
	return
}

// void                cairo_scaled_font_get_scale_matrix  (cairo_scaled_font_t *scaled_font,
//                                                          cairo_matrix_t *scale_matrix);
func (this *ScaledFont) GetScaleMatrix() (out Matrix) {
	C.cairo_scaled_font_get_scale_matrix(this.c(), out.c())
	return
}

// cairo_font_type_t   cairo_scaled_font_get_type          (cairo_scaled_font_t *scaled_font);
func (this *ScaledFont) GetType() FontType {
	return FontType(C.cairo_scaled_font_get_type(this.c()))
}

// unsigned int        cairo_scaled_font_get_reference_count
//                                                         (cairo_scaled_font_t *scaled_font);
// cairo_status_t      cairo_scaled_font_set_user_data     (cairo_scaled_font_t *scaled_font,
//                                                          const cairo_user_data_key_t *key,
//                                                          void *user_data,
//                                                          cairo_destroy_func_t destroy);
// void *              cairo_scaled_font_get_user_data     (cairo_scaled_font_t *scaled_font,
//                                                          const cairo_user_data_key_t *key);

//----------------------------------------------------------------------------
// Font Options
//----------------------------------------------------------------------------

// typedef             cairo_font_options_t;
type FontOptions struct {
	C unsafe.Pointer
}

func (this *FontOptions) c() *C.cairo_font_options_t {
	return (*C.cairo_font_options_t)(this.C)
}

func font_options_finalizer(this *FontOptions) {
	if gobject.FQueue.Push(unsafe.Pointer(this), font_options_finalizer2) {
		return
	}
	C.cairo_font_options_destroy(this.c())
}

func font_options_finalizer2(this_un unsafe.Pointer) {
	this := (*FontOptions)(this_un)
	C.cairo_font_options_destroy(this.c())
}

func FontOptionsWrap(c_un unsafe.Pointer) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_font_options_t)(c_un)
	font_options := &FontOptions{unsafe.Pointer(c)}
	runtime.SetFinalizer(font_options, font_options_finalizer)
	return unsafe.Pointer(font_options)
}

// cairo_font_options_t * cairo_font_options_create        (void);
func NewFontOptions() *FontOptions {
	return (*FontOptions)(FontOptionsWrap(unsafe.Pointer(C.cairo_font_options_create())))
}

// cairo_font_options_t * cairo_font_options_copy          (const cairo_font_options_t *original);
func (this *FontOptions) Copy() *FontOptions {
	return (*FontOptions)(FontOptionsWrap(unsafe.Pointer(C.cairo_font_options_copy(this.c()))))
}

// void                cairo_font_options_destroy          (cairo_font_options_t *options);
func (this *FontOptions) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_font_options_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_font_options_status           (cairo_font_options_t *options);
func (this *FontOptions) Status() Status {
	return Status(C.cairo_font_options_status(this.c()))
}

// void                cairo_font_options_merge            (cairo_font_options_t *options,
//                                                          const cairo_font_options_t *other);
func (this *FontOptions) Merge(other *FontOptions) {
	C.cairo_font_options_merge(this.c(), other.c())
}

// unsigned long       cairo_font_options_hash             (const cairo_font_options_t *options);
func (this *FontOptions) Hash() uint64 {
	return uint64(C.cairo_font_options_hash(this.c()))
}

// cairo_bool_t        cairo_font_options_equal            (const cairo_font_options_t *options,
//                                                          const cairo_font_options_t *other);
func (this *FontOptions) Equal(other *FontOptions) bool {
	return C.cairo_font_options_equal(this.c(), other.c()) != 0
}

// void                cairo_font_options_set_antialias    (cairo_font_options_t *options,
//                                                          cairo_antialias_t antialias);
func (this *FontOptions) SetAntialias(antialias Antialias) {
	C.cairo_font_options_set_antialias(this.c(), antialias.c())
}

// cairo_antialias_t   cairo_font_options_get_antialias    (const cairo_font_options_t *options);
func (this *FontOptions) GetAntialias() Antialias {
	return Antialias(C.cairo_font_options_get_antialias(this.c()))
}

// enum                cairo_subpixel_order_t;
type SubpixelOrder int

const (
	SubpixelOrderDefault SubpixelOrder = C.CAIRO_SUBPIXEL_ORDER_DEFAULT
	SubpixelOrderRGB     SubpixelOrder = C.CAIRO_SUBPIXEL_ORDER_RGB
	SubpixelOrderBGR     SubpixelOrder = C.CAIRO_SUBPIXEL_ORDER_BGR
	SubpixelOrderVRGB    SubpixelOrder = C.CAIRO_SUBPIXEL_ORDER_VRGB
	SubpixelOrderVBGR    SubpixelOrder = C.CAIRO_SUBPIXEL_ORDER_VBGR
)

func (this SubpixelOrder) c() C.cairo_subpixel_order_t {
	return C.cairo_subpixel_order_t(this)
}

// void                cairo_font_options_set_subpixel_order
//                                                         (cairo_font_options_t *options,
//                                                          cairo_subpixel_order_t subpixel_order);
func (this *FontOptions) SetSubpixelOrder(subpixel_order SubpixelOrder) {
	C.cairo_font_options_set_subpixel_order(this.c(), subpixel_order.c())
}

// cairo_subpixel_order_t  cairo_font_options_get_subpixel_order
//                                                         (const cairo_font_options_t *options);
func (this *FontOptions) GetSubpixelOrder() SubpixelOrder {
	return SubpixelOrder(C.cairo_font_options_get_subpixel_order(this.c()))
}

// enum                cairo_hint_style_t;
type HintStyle int

const (
	HintStyleDefault HintStyle = C.CAIRO_HINT_STYLE_DEFAULT
	HintStyleNone    HintStyle = C.CAIRO_HINT_STYLE_NONE
	HintStyleSlight  HintStyle = C.CAIRO_HINT_STYLE_SLIGHT
	HintStyleMedium  HintStyle = C.CAIRO_HINT_STYLE_MEDIUM
	HintStyleFull    HintStyle = C.CAIRO_HINT_STYLE_FULL
)

func (this HintStyle) c() C.cairo_hint_style_t {
	return C.cairo_hint_style_t(this)
}

// void                cairo_font_options_set_hint_style   (cairo_font_options_t *options,
//                                                          cairo_hint_style_t hint_style);
func (this *FontOptions) SetHintStyle(hint_style HintStyle) {
	C.cairo_font_options_set_hint_style(this.c(), hint_style.c())
}

// cairo_hint_style_t  cairo_font_options_get_hint_style   (const cairo_font_options_t *options);
func (this *FontOptions) GetHintStyle() HintStyle {
	return HintStyle(C.cairo_font_options_get_hint_style(this.c()))
}

// enum                cairo_hint_metrics_t;
type HintMetrics int

const (
	HintMetricsDefault HintMetrics = C.CAIRO_HINT_METRICS_DEFAULT
	HintMetricsOff     HintMetrics = C.CAIRO_HINT_METRICS_OFF
	HintMetricsOn      HintMetrics = C.CAIRO_HINT_METRICS_ON
)

func (this HintMetrics) c() C.cairo_hint_metrics_t {
	return C.cairo_hint_metrics_t(this)
}

// void                cairo_font_options_set_hint_metrics (cairo_font_options_t *options,
//                                                          cairo_hint_metrics_t hint_metrics);
func (this *FontOptions) SetHintMetrics(hint_metrics HintMetrics) {
	C.cairo_font_options_set_hint_metrics(this.c(), hint_metrics.c())
}

// cairo_hint_metrics_t  cairo_font_options_get_hint_metrics
//                                                         (const cairo_font_options_t *options);
func (this *FontOptions) GetHintMetrics() HintMetrics {
	return HintMetrics(C.cairo_font_options_get_hint_metrics(this.c()))
}

//----------------------------------------------------------------------------
// Device (TODO)
//----------------------------------------------------------------------------

//----------------------------------------------------------------------------
// Surface
// 	ImageSurface
// 	PDFSurface
// 	... etc
//----------------------------------------------------------------------------

// #define             CAIRO_MIME_TYPE_JP2
// #define             CAIRO_MIME_TYPE_JPEG
// #define             CAIRO_MIME_TYPE_PNG
// #define             CAIRO_MIME_TYPE_URI
const (
	MimeTypeJp2  = "image/jp2"
	MimeTypeJpeg = "image/jpeg"
	MimeTypePng  = "image/png"
	MimeTypeUri  = "text/x-uri"
)

// typedef             cairo_surface_t;
type SurfaceLike interface {
	InheritedFromCairoSurface() *Surface
}

type Surface struct{ C unsafe.Pointer }
type ImageSurface struct{ Surface }
type PDFSurface struct{ Surface }

func (this *Surface) c() *C.cairo_surface_t {
	return (*C.cairo_surface_t)(this.C)
}

func surface_finalizer(this *Surface) {
	if gobject.FQueue.Push(unsafe.Pointer(this), surface_finalizer2) {
		return
	}
	C.cairo_surface_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_surface_destroy(this.c())
}

func surface_finalizer2(this_un unsafe.Pointer) {
	this := (*Surface)(this_un)
	C.cairo_surface_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_surface_destroy(this.c())
}

func SurfaceWrap(c_un unsafe.Pointer, grab bool) unsafe.Pointer {
	if c_un == nil {
		return nil
	}
	c := (*C.cairo_surface_t)(c_un)
	go_repr := C.cairo_surface_get_user_data(c, &go_repr_cookie)
	if go_repr != nil {
		return unsafe.Pointer(go_repr)
	}

	surface := &Surface{unsafe.Pointer(c)}
	if grab {
		C.cairo_pattern_reference(c)
	}
	runtime.SetFinalizer(surface, surface_finalizer)

	status := C.cairo_surface_set_user_data(c, &go_repr_cookie, unsafe.Pointer(surface), nil)
	if status != C.CAIRO_STATUS_SUCCESS {
		panic("failed to set user data, out of memory?")
	}
	return unsafe.Pointer(surface)
}

func ensure_surface_type(c *C.cairo_surface_t, types ...SurfaceType) {
	for _, t := range types {
		if C.cairo_surface_get_type(c) == t.c() {
			return
		}
	}
	panic("unexpected surface type")
}

func ToSurface(like SurfaceLike) *Surface {
	return (*Surface)(SurfaceWrap(unsafe.Pointer(like.InheritedFromCairoSurface().c()), true))
}

func ToImageSurface(like SurfaceLike) *ImageSurface {
	c := like.InheritedFromCairoSurface().c()
	ensure_surface_type(c, SurfaceTypeImage)
	return (*ImageSurface)(SurfaceWrap(unsafe.Pointer(c), true))
}

func ToPDFSurface(like SurfaceLike) *PDFSurface {
	c := like.InheritedFromCairoSurface().c()
	ensure_surface_type(c, SurfaceTypePDF)
	return (*PDFSurface)(SurfaceWrap(unsafe.Pointer(c), true))
}

func (this *Surface) InheritedFromCairoSurface() *Surface {
	return this
}

// enum                cairo_content_t;
type Content int

const (
	ContentColor      Content = C.CAIRO_CONTENT_COLOR
	ContentAlpha      Content = C.CAIRO_CONTENT_ALPHA
	ContentColorAlpha Content = C.CAIRO_CONTENT_COLOR_ALPHA
)

func (this Content) c() C.cairo_content_t {
	return C.cairo_content_t(this)
}

// cairo_surface_t *   cairo_surface_create_similar        (cairo_surface_t *other,
//                                                          cairo_content_t content,
//                                                          int width,
//                                                          int height);
func (this *Surface) CreateSimilar(content Content, width, height int) *Surface {
	return (*Surface)(SurfaceWrap(unsafe.Pointer(C.cairo_surface_create_similar(this.c(), content.c(), C.int(width), C.int(height))), false))
}

// cairo_surface_t *   cairo_surface_create_for_rectangle  (cairo_surface_t *target,
//                                                          double x,
//                                                          double y,
//                                                          double width,
//                                                          double height);
func (this *Surface) CreateForRectangle(x, y, width, height float64) *Surface {
	return (*Surface)(SurfaceWrap(unsafe.Pointer(C.cairo_surface_create_for_rectangle(this.c(),
		C.double(x), C.double(y), C.double(width), C.double(height))), false))
}

// cairo_surface_t *   cairo_surface_reference             (cairo_surface_t *surface);
// void                cairo_surface_destroy               (cairo_surface_t *surface);
func (this *Surface) Unref() {
	runtime.SetFinalizer(this, nil)
	C.cairo_surface_set_user_data(this.c(), &go_repr_cookie, nil, nil)
	C.cairo_surface_destroy(this.c())
	this.C = nil
}

// cairo_status_t      cairo_surface_status                (cairo_surface_t *surface);
func (this *Surface) Status() Status {
	return Status(C.cairo_surface_status(this.c()))
}

// void                cairo_surface_finish                (cairo_surface_t *surface);
func (this *Surface) Finish() {
	C.cairo_surface_finish(this.c())
}

// void                cairo_surface_flush                 (cairo_surface_t *surface);
func (this *Surface) Flush() {
	C.cairo_surface_flush(this.c())
}

// TODO: Implement these
// cairo_device_t *    cairo_surface_get_device            (cairo_surface_t *surface);
// void                cairo_surface_get_font_options      (cairo_surface_t *surface,
//                                                          cairo_font_options_t *options);

// cairo_content_t     cairo_surface_get_content           (cairo_surface_t *surface);
func (this *Surface) GetContent() Content {
	return Content(C.cairo_surface_get_content(this.c()))
}

// void                cairo_surface_mark_dirty            (cairo_surface_t *surface);
func (this *Surface) MarkDirty() {
	C.cairo_surface_mark_dirty(this.c())
}

// void                cairo_surface_mark_dirty_rectangle  (cairo_surface_t *surface,
//                                                          int x,
//                                                          int y,
//                                                          int width,
//                                                          int height);
func (this *Surface) MarkDirtyRectangle(x, y, width, height int) {
	C.cairo_surface_mark_dirty_rectangle(this.c(),
		C.int(x), C.int(y), C.int(width), C.int(height))
}

// void                cairo_surface_set_device_offset     (cairo_surface_t *surface,
//                                                          double x_offset,
//                                                          double y_offset);
func (this *Surface) SetDeviceOffset(x_offset, y_offset float64) {
	C.cairo_surface_set_device_offset(this.c(), C.double(x_offset), C.double(y_offset))
}

// void                cairo_surface_get_device_offset     (cairo_surface_t *surface,
//                                                          double *x_offset,
//                                                          double *y_offset);
func (this *Surface) GetDeviceOffset() (x_offset, y_offset float64) {
	C.cairo_surface_get_device_offset(this.c(),
		(*C.double)(unsafe.Pointer(&x_offset)),
		(*C.double)(unsafe.Pointer(&y_offset)))
	return
}

// void                cairo_surface_set_fallback_resolution
//                                                         (cairo_surface_t *surface,
//                                                          double x_pixels_per_inch,
//                                                          double y_pixels_per_inch);
func (this *Surface) SetFallbackResolution(x_pixels_per_inch, y_pixels_per_inch float64) {
	C.cairo_surface_set_fallback_resolution(this.c(), C.double(x_pixels_per_inch), C.double(y_pixels_per_inch))
}

// void                cairo_surface_get_fallback_resolution
//                                                         (cairo_surface_t *surface,
//                                                          double *x_pixels_per_inch,
//                                                          double *y_pixels_per_inch);
func (this *Surface) GetFallbackResolution() (x_pixels_per_inch, y_pixels_per_inch float64) {
	C.cairo_surface_get_fallback_resolution(this.c(),
		(*C.double)(unsafe.Pointer(&x_pixels_per_inch)),
		(*C.double)(unsafe.Pointer(&y_pixels_per_inch)))
	return
}

// enum                cairo_surface_type_t;
type SurfaceType int

const (
	SurfaceTypeImage         SurfaceType = C.CAIRO_SURFACE_TYPE_IMAGE
	SurfaceTypePDF           SurfaceType = C.CAIRO_SURFACE_TYPE_PDF
	SurfaceTypePS            SurfaceType = C.CAIRO_SURFACE_TYPE_PS
	SurfaceTypeXLib          SurfaceType = C.CAIRO_SURFACE_TYPE_XLIB
	SurfaceTypeXCB           SurfaceType = C.CAIRO_SURFACE_TYPE_XCB
	SurfaceTypeGlitz         SurfaceType = C.CAIRO_SURFACE_TYPE_GLITZ
	SurfaceTypeQuartz        SurfaceType = C.CAIRO_SURFACE_TYPE_QUARTZ
	SurfaceTypeWin32         SurfaceType = C.CAIRO_SURFACE_TYPE_WIN32
	SurfaceTypeBeOS          SurfaceType = C.CAIRO_SURFACE_TYPE_BEOS
	SurfaceTypeDirectFB      SurfaceType = C.CAIRO_SURFACE_TYPE_DIRECTFB
	SurfaceTypeSVG           SurfaceType = C.CAIRO_SURFACE_TYPE_SVG
	SurfaceTypeOs2           SurfaceType = C.CAIRO_SURFACE_TYPE_OS2
	SurfaceTypeWin32Printing SurfaceType = C.CAIRO_SURFACE_TYPE_WIN32_PRINTING
	SurfaceTypeQuartzImage   SurfaceType = C.CAIRO_SURFACE_TYPE_QUARTZ_IMAGE
	SurfaceTypeScript        SurfaceType = C.CAIRO_SURFACE_TYPE_SCRIPT
	SurfaceTypeQt            SurfaceType = C.CAIRO_SURFACE_TYPE_QT
	SurfaceTypeRecording     SurfaceType = C.CAIRO_SURFACE_TYPE_RECORDING
	SurfaceTypeVg            SurfaceType = C.CAIRO_SURFACE_TYPE_VG
	SurfaceTypeGL            SurfaceType = C.CAIRO_SURFACE_TYPE_GL
	SurfaceTypeDRM           SurfaceType = C.CAIRO_SURFACE_TYPE_DRM
	SurfaceTypeTee           SurfaceType = C.CAIRO_SURFACE_TYPE_TEE
	SurfaceTypeXML           SurfaceType = C.CAIRO_SURFACE_TYPE_XML
	SurfaceTypeSkia          SurfaceType = C.CAIRO_SURFACE_TYPE_SKIA
	SurfaceTypeSubsurface    SurfaceType = C.CAIRO_SURFACE_TYPE_SUBSURFACE
)

func (this SurfaceType) c() C.cairo_surface_type_t {
	return C.cairo_surface_type_t(this)
}

// cairo_surface_type_t  cairo_surface_get_type            (cairo_surface_t *surface);
func (this *Surface) GetType() SurfaceType {
	return SurfaceType(C.cairo_surface_get_type(this.c()))
}

// TODO: Implement these
// unsigned int        cairo_surface_get_reference_count   (cairo_surface_t *surface);
// cairo_status_t      cairo_surface_set_user_data         (cairo_surface_t *surface,
//                                                          const cairo_user_data_key_t *key,
//                                                          void *user_data,
//                                                          cairo_destroy_func_t destroy);
// void *              cairo_surface_get_user_data         (cairo_surface_t *surface,
//                                                          const cairo_user_data_key_t *key);

// void                cairo_surface_copy_page             (cairo_surface_t *surface);
func (this *Surface) CopyPage() {
	C.cairo_surface_copy_page(this.c())
}

// void                cairo_surface_show_page             (cairo_surface_t *surface);
func (this *Surface) ShowPage() {
	C.cairo_surface_show_page(this.c())
}

// cairo_bool_t        cairo_surface_has_show_text_glyphs  (cairo_surface_t *surface);
func (this *Surface) HasShowTextGlyphs() bool {
	return C.cairo_surface_has_show_text_glyphs(this.c()) != 0
}

// TODO: Implement these
// cairo_status_t      cairo_surface_set_mime_data         (cairo_surface_t *surface,
//                                                          const char *mime_type,
//                                                          unsigned char *data,
//                                                          unsigned long  length,
//                                                          cairo_destroy_func_t destroy,
//                                                          void *closure);
// void                cairo_surface_get_mime_data         (cairo_surface_t *surface,
//                                                          const char *mime_type,
//                                                          unsigned char **data,
//                                                          unsigned long *length);

//----------------------------------------------------------------------------
// Image Surfaces
//----------------------------------------------------------------------------

// enum                cairo_format_t;
type Format int

const (
	FormatInvalid   Format = C.CAIRO_FORMAT_INVALID
	FormatARGB32    Format = C.CAIRO_FORMAT_ARGB32
	FormatRGB24     Format = C.CAIRO_FORMAT_RGB24
	FormatA8        Format = C.CAIRO_FORMAT_A8
	FormatA1        Format = C.CAIRO_FORMAT_A1
	FormatRGB16_565 Format = C.CAIRO_FORMAT_RGB16_565
)

func (this Format) c() C.cairo_format_t {
	return C.cairo_format_t(this)
}

// int                 cairo_format_stride_for_width       (cairo_format_t format,
//                                                          int width);
func (this Format) StrideForWidth(width int) int {
	return int(C.cairo_format_stride_for_width(this.c(), C.int(width)))
}

// cairo_surface_t *   cairo_image_surface_create          (cairo_format_t format,
//                                                          int width,
//                                                          int height);
func NewImageSurface(format Format, width, height int) *ImageSurface {
	return (*ImageSurface)(SurfaceWrap(unsafe.Pointer(C.cairo_image_surface_create(format.c(), C.int(width), C.int(height))), false))
}

// TODO: Implement this (need a way to keep GC from freeing the 'data')
// cairo_surface_t *   cairo_image_surface_create_for_data (unsigned char *data,
//                                                          cairo_format_t format,
//                                                          int width,
//                                                          int height,
//                                                          int stride);

// TODO: Implement this (think about wrapping it into slice.. how?)
// unsigned char *     cairo_image_surface_get_data        (cairo_surface_t *surface);

// cairo_format_t      cairo_image_surface_get_format      (cairo_surface_t *surface);
func (this *ImageSurface) GetFormat() Format {
	return Format(C.cairo_image_surface_get_format(this.c()))
}

// int                 cairo_image_surface_get_width       (cairo_surface_t *surface);
func (this *ImageSurface) GetWidth() int {
	return int(C.cairo_image_surface_get_width(this.c()))
}

// int                 cairo_image_surface_get_height      (cairo_surface_t *surface);
func (this *ImageSurface) GetHeight() int {
	return int(C.cairo_image_surface_get_height(this.c()))
}

// int                 cairo_image_surface_get_stride      (cairo_surface_t *surface);
func (this *ImageSurface) GetStride() int {
	return int(C.cairo_image_surface_get_stride(this.c()))
}

//----------------------------------------------------------------------------
// PNG Support
//----------------------------------------------------------------------------

// cairo_surface_t *   cairo_image_surface_create_from_png (const char *filename);
func NewImageSurfaceFromPNG(filename string) *ImageSurface {
	cfilename := C.CString(filename)
	surface := (*ImageSurface)(SurfaceWrap(unsafe.Pointer(C.cairo_image_surface_create_from_png(cfilename)), false))
	C.free(unsafe.Pointer(cfilename))
	return surface
}

// cairo_status_t      (*cairo_read_func_t)                (void *closure,
//                                                          unsigned char *data,
//                                                          unsigned int length);

//export io_reader_wrapper
func io_reader_wrapper(reader_up unsafe.Pointer, data_uc *C.uchar, length uint32) uint32 {
	var reader io.Reader
	var data []byte
	var data_header reflect.SliceHeader
	reader = *(*io.Reader)(reader_up)
	data_header.Data = uintptr(unsafe.Pointer(data_uc))
	data_header.Len = int(length)
	data_header.Cap = int(length)
	data = *(*[]byte)(unsafe.Pointer(&data_header))

	_, err := reader.Read(data)
	if err != nil {
		return uint32(StatusReadError)
	}
	return uint32(StatusSuccess)
}

// cairo_surface_t *   cairo_image_surface_create_from_png_stream
//                                                         (cairo_read_func_t read_func,
//                                                          void *closure);
func NewImageSurfaceFromPNGStream(r io.Reader) *ImageSurface {
	return (*ImageSurface)(SurfaceWrap(unsafe.Pointer(C._cairo_image_surface_create_from_png_stream(unsafe.Pointer(&r))), false))
}

// cairo_status_t      cairo_surface_write_to_png          (cairo_surface_t *surface,
//                                                          const char *filename);
func (this *ImageSurface) WriteToPNG(filename string) Status {
	cfilename := C.CString(filename)
	status := C.cairo_surface_write_to_png(this.c(), cfilename)
	C.free(unsafe.Pointer(cfilename))
	return Status(status)
}

// cairo_status_t      (*cairo_write_func_t)               (void *closure,
//                                                          unsigned char *data,
//                                                          unsigned int length);

//export io_writer_wrapper
func io_writer_wrapper(writer_up unsafe.Pointer, data_uc *C.uchar, length uint32) uint32 {
	var writer io.Writer
	var data []byte
	var data_header reflect.SliceHeader
	writer = *(*io.Writer)(writer_up)
	data_header.Data = uintptr(unsafe.Pointer(data_uc))
	data_header.Len = int(length)
	data_header.Cap = int(length)
	data = *(*[]byte)(unsafe.Pointer(&data_header))

	_, err := writer.Write(data)
	if err != nil {
		return uint32(StatusWriteError)
	}
	return uint32(StatusSuccess)
}

// cairo_status_t      cairo_surface_write_to_png_stream   (cairo_surface_t *surface,
//                                                          cairo_write_func_t write_func,
//                                                          void *closure);
func (this *ImageSurface) WriteToPNGStream(w io.Writer) Status {
	return Status(C._cairo_surface_write_to_png_stream(this.c(), unsafe.Pointer(&w)))
}

//----------------------------------------------------------------------------
// PDF Surfaces
//----------------------------------------------------------------------------

// #define             CAIRO_HAS_PDF_SURFACE
// cairo_surface_t *   cairo_pdf_surface_create            (const char *filename,
//                                                          double width_in_points,
//                                                          double height_in_points);
func NewPDFSurface(filename string, width_in_points, height_in_points float64) *PDFSurface {
	cfilename := C.CString(filename)
	surface := (*PDFSurface)(SurfaceWrap(unsafe.Pointer(C.cairo_pdf_surface_create(cfilename, C.double(width_in_points), C.double(height_in_points))), false))
	C.free(unsafe.Pointer(cfilename))
	return surface
}

// cairo_surface_t *   cairo_pdf_surface_create_for_stream (cairo_write_func_t write_func,
//                                                          void *closure,
//                                                          double width_in_points,
//                                                          double height_in_points);
func NewPDFSurfaceForStream(w io.Writer, width_in_points, height_in_points float64) *PDFSurface {
	surface := (*PDFSurface)(SurfaceWrap(unsafe.Pointer(C._cairo_pdf_surface_create_for_stream(unsafe.Pointer(&w),
		C.double(width_in_points), C.double(height_in_points))), false))
	return surface
}

// void                cairo_pdf_surface_restrict_to_version
//                                                         (cairo_surface_t *surface,
//                                                          cairo_pdf_version_t version);
func (this *PDFSurface) RestrictToVersion(version PDFVersion) {
	C.cairo_pdf_surface_restrict_to_version(this.c(), version.c())
}

// enum                cairo_pdf_version_t;
type PDFVersion int

const (
	PDFVersion_1_4 PDFVersion = C.CAIRO_PDF_VERSION_1_4
	PDFVersion_1_5 PDFVersion = C.CAIRO_PDF_VERSION_1_5
)

func (this PDFVersion) c() C.cairo_pdf_version_t {
	return C.cairo_pdf_version_t(this)
}

// void                cairo_pdf_get_versions              (cairo_pdf_version_t const **versions,
//                                                          int *num_versions);
func PDFGetVersions() []PDFVersion {
	var versions *C.cairo_pdf_version_t
	var num_versions C.int
	C.cairo_pdf_get_versions(&versions, &num_versions)

	var out []PDFVersion
	if num_versions > 0 {
		out = make([]PDFVersion, num_versions)
		for i := range out {
			out[i] = (*(*[999999]PDFVersion)(unsafe.Pointer(&versions)))[i]
		}
	}
	return out
}

// const char *        cairo_pdf_version_to_string         (cairo_pdf_version_t version);
func (this PDFVersion) String() string {
	return C.GoString(C.cairo_pdf_version_to_string(this.c()))
}

// void                cairo_pdf_surface_set_size          (cairo_surface_t *surface,
//                                                          double width_in_points,
//                                                          double height_in_points);
func (this *PDFSurface) SetSize(width_in_points, height_in_points float64) {
	C.cairo_pdf_surface_set_size(this.c(), C.double(width_in_points), C.double(height_in_points))
}

//----------------------------------------------------------------------------
// Matrix
//----------------------------------------------------------------------------

//                     cairo_matrix_t;
type Matrix struct {
	xx, yx, xy, yy, x0, y0 float64
}

func (this *Matrix) c() *C.cairo_matrix_t {
	return (*C.cairo_matrix_t)(unsafe.Pointer(this))
}

// void                cairo_matrix_init                   (cairo_matrix_t *matrix,
//                                                          double xx,
//                                                          double yx,
//                                                          double xy,
//                                                          double yy,
//                                                          double x0,
//                                                          double y0);
func (this *Matrix) Init(xx, yx, xy, yy, x0, y0 float64) {
	C.cairo_matrix_init(this.c(), C.double(xx), C.double(yx), C.double(xy), C.double(yy), C.double(x0), C.double(y0))
}

// void                cairo_matrix_init_identity          (cairo_matrix_t *matrix);
func (this *Matrix) InitIdentity() {
	C.cairo_matrix_init_identity(this.c())
}

// void                cairo_matrix_init_translate         (cairo_matrix_t *matrix,
//                                                          double tx,
//                                                          double ty);
func (this *Matrix) InitTranslate(tx, ty float64) {
	C.cairo_matrix_init_translate(this.c(), C.double(tx), C.double(ty))
}

// void                cairo_matrix_init_scale             (cairo_matrix_t *matrix,
//                                                          double sx,
//                                                          double sy);
func (this *Matrix) InitScale(sx, sy float64) {
	C.cairo_matrix_init_scale(this.c(), C.double(sx), C.double(sy))
}

// void                cairo_matrix_init_rotate            (cairo_matrix_t *matrix,
//                                                          double radians);
func (this *Matrix) InitRotate(radians float64) {
	C.cairo_matrix_init_rotate(this.c(), C.double(radians))
}

// void                cairo_matrix_translate              (cairo_matrix_t *matrix,
//                                                          double tx,
//                                                          double ty);
func (this *Matrix) Translate(tx, ty float64) {
	C.cairo_matrix_translate(this.c(), C.double(tx), C.double(ty))
}

// void                cairo_matrix_scale                  (cairo_matrix_t *matrix,
//                                                          double sx,
//                                                          double sy);
func (this *Matrix) Scale(sx, sy float64) {
	C.cairo_matrix_scale(this.c(), C.double(sx), C.double(sy))
}

// void                cairo_matrix_rotate                 (cairo_matrix_t *matrix,
//                                                          double radians);
func (this *Matrix) Rotate(radians float64) {
	C.cairo_matrix_rotate(this.c(), C.double(radians))
}

// cairo_status_t      cairo_matrix_invert                 (cairo_matrix_t *matrix);
func (this *Matrix) Invert() Status {
	return Status(C.cairo_matrix_invert(this.c()))
}

// void                cairo_matrix_multiply               (cairo_matrix_t *result,
//                                                          const cairo_matrix_t *a,
//                                                          const cairo_matrix_t *b);
func (this *Matrix) Multiply(b *Matrix) Matrix {
	var result C.cairo_matrix_t
	C.cairo_matrix_multiply(&result, this.c(), b.c())
	return *(*Matrix)(unsafe.Pointer(&result))
}

// void                cairo_matrix_transform_distance     (const cairo_matrix_t *matrix,
//                                                          double *dx,
//                                                          double *dy);
func (this *Matrix) TransformDistance(dx, dy float64) (float64, float64) {
	C.cairo_matrix_transform_distance(this.c(),
		(*C.double)(unsafe.Pointer(&dx)),
		(*C.double)(unsafe.Pointer(&dy)))
	return dx, dy
}

// void                cairo_matrix_transform_point        (const cairo_matrix_t *matrix,
//                                                          double *x,
//                                                          double *y);
func (this *Matrix) TransformPoint(x, y float64) (float64, float64) {
	C.cairo_matrix_transform_point(this.c(),
		(*C.double)(unsafe.Pointer(&x)),
		(*C.double)(unsafe.Pointer(&y)))
	return x, y
}
