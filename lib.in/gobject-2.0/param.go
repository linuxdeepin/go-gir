package gobject

import (
	"runtime"
	"unsafe"
)

/*
#include "gobject.gen.h"
#include "gobject.h"
#include <string.h>

extern GParamSpec *g_param_spec_ref_sink(GParamSpec*);
extern void g_param_spec_unref(GParamSpec*);
*/
import "C"

//--------------------------------------------------------------
// ParamSpec utils
//--------------------------------------------------------------

// Let's implement these manually (not Object based and small amount of things
// to implement).

// First some utils
func param_spec_finalizer(pspec *ParamSpec) {
	if FQueue.Push(unsafe.Pointer(pspec), param_spec_finalizer2) {
		return
	}
	C.g_param_spec_unref((*C.GParamSpec)(pspec.C))
}

func param_spec_finalizer2(pspec_un unsafe.Pointer) {
	pspec := (*ParamSpec)(pspec_un)
	C.g_param_spec_unref((*C.GParamSpec)(pspec.C))
}

func set_param_spec_finalizer(pspec *ParamSpec) {
	runtime.SetFinalizer(pspec, param_spec_finalizer)
}

func ParamSpecGrabIfType(c unsafe.Pointer, t Type) unsafe.Pointer {
	if c == nil {
		return nil
	}
	obj := &ParamSpec{c}
	if obj.GetType().IsA(t) {
		C.g_param_spec_ref_sink((*C.GParamSpec)(obj.C))
		set_param_spec_finalizer(obj)
		return unsafe.Pointer(obj)
	}
	return nil
}

func ParamSpecWrap(c unsafe.Pointer, grab bool) unsafe.Pointer {
	if c == nil {
		return nil
	}
	obj := &ParamSpec{c}
	if grab {
		C.g_param_spec_ref_sink((*C.GParamSpec)(obj.C))
	}
	set_param_spec_finalizer(obj)
	return unsafe.Pointer(obj)
}

//--------------------------------------------------------------
// ParamSpec
//--------------------------------------------------------------

type ParamSpecLike interface {
	InheritedFromGParamSpec() *C.GParamSpec
}

type ParamSpec struct {
	C unsafe.Pointer
}

func ToParamSpec(pspeclike ParamSpecLike) *ParamSpec {
	t := (*ParamSpec)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpec()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpec)(obj)
	}
	panic("cannot cast to ParamSpec")
}

func (this *ParamSpec) InheritedFromGParamSpec() *C.GParamSpec {
	return (*C.GParamSpec)(this.C)
}

func (this *ParamSpec) GetStaticType() Type {
	return Type(C._g_type_param())
}

func (this *ParamSpec) GetType() Type {
	return Type(C._g_param_spec_type(this.InheritedFromGParamSpec()))
}

func (this *ParamSpec) GetValueType() Type {
	return Type(C._g_param_spec_value_type(this.InheritedFromGParamSpec()))
}

//--------------------------------------------------------------
// ParamSpecBoolean
//--------------------------------------------------------------

type ParamSpecBooleanLike interface {
	InheritedFromGParamSpecBoolean() *C.GParamSpecBoolean
}

type ParamSpecBoolean struct {
	ParamSpec
}

func ToParamSpecBoolean(pspeclike ParamSpecBooleanLike) *ParamSpecBoolean {
	t := (*ParamSpecBoolean)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecBoolean()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecBoolean)(obj)
	}
	panic("cannot cast to ParamSpecBoolean")
}

func (this *ParamSpecBoolean) InheritedFromGParamSpecBoolean() *C.GParamSpecBoolean {
	return (*C.GParamSpecBoolean)(this.C)
}

func (this *ParamSpecBoolean) GetStaticType() Type {
	return Type(C._g_type_param_boolean())
}

//--------------------------------------------------------------
// ParamSpecBoxed
//--------------------------------------------------------------

type ParamSpecBoxedLike interface {
	InheritedFromGParamSpecBoxed() *C.GParamSpecBoxed
}

type ParamSpecBoxed struct {
	ParamSpec
}

func ToParamSpecBoxed(pspeclike ParamSpecBoxedLike) *ParamSpecBoxed {
	t := (*ParamSpecBoxed)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecBoxed()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecBoxed)(obj)
	}
	panic("cannot cast to ParamSpecBoxed")
}

func (this *ParamSpecBoxed) InheritedFromGParamSpecBoxed() *C.GParamSpecBoxed {
	return (*C.GParamSpecBoxed)(this.C)
}

func (this *ParamSpecBoxed) GetStaticType() Type {
	return Type(C._g_type_param_boxed())
}

//--------------------------------------------------------------
// ParamSpecChar
//--------------------------------------------------------------

type ParamSpecCharLike interface {
	InheritedFromGParamSpecChar() *C.GParamSpecChar
}

type ParamSpecChar struct {
	ParamSpec
}

func ToParamSpecChar(pspeclike ParamSpecCharLike) *ParamSpecChar {
	t := (*ParamSpecChar)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecChar()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecChar)(obj)
	}
	panic("cannot cast to ParamSpecChar")
}

func (this *ParamSpecChar) InheritedFromGParamSpecChar() *C.GParamSpecChar {
	return (*C.GParamSpecChar)(this.C)
}

func (this *ParamSpecChar) GetStaticType() Type {
	return Type(C._g_type_param_char())
}

//--------------------------------------------------------------
// ParamSpecDouble
//--------------------------------------------------------------

type ParamSpecDoubleLike interface {
	InheritedFromGParamSpecDouble() *C.GParamSpecDouble
}

type ParamSpecDouble struct {
	ParamSpec
}

func ToParamSpecDouble(pspeclike ParamSpecDoubleLike) *ParamSpecDouble {
	t := (*ParamSpecDouble)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecDouble()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecDouble)(obj)
	}
	panic("cannot cast to ParamSpecDouble")
}

func (this *ParamSpecDouble) InheritedFromGParamSpecDouble() *C.GParamSpecDouble {
	return (*C.GParamSpecDouble)(this.C)
}

func (this *ParamSpecDouble) GetStaticType() Type {
	return Type(C._g_type_param_double())
}

//--------------------------------------------------------------
// ParamSpecEnum
//--------------------------------------------------------------

type ParamSpecEnumLike interface {
	InheritedFromGParamSpecEnum() *C.GParamSpecEnum
}

type ParamSpecEnum struct {
	ParamSpec
}

func ToParamSpecEnum(pspeclike ParamSpecEnumLike) *ParamSpecEnum {
	t := (*ParamSpecEnum)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecEnum()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecEnum)(obj)
	}
	panic("cannot cast to ParamSpecEnum")
}

func (this *ParamSpecEnum) InheritedFromGParamSpecEnum() *C.GParamSpecEnum {
	return (*C.GParamSpecEnum)(this.C)
}

func (this *ParamSpecEnum) GetStaticType() Type {
	return Type(C._g_type_param_enum())
}

//--------------------------------------------------------------
// ParamSpecFlags
//--------------------------------------------------------------

type ParamSpecFlagsLike interface {
	InheritedFromGParamSpecFlags() *C.GParamSpecFlags
}

type ParamSpecFlags struct {
	ParamSpec
}

func ToParamSpecFlags(pspeclike ParamSpecFlagsLike) *ParamSpecFlags {
	t := (*ParamSpecFlags)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecFlags()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecFlags)(obj)
	}
	panic("cannot cast to ParamSpecFlags")
}

func (this *ParamSpecFlags) InheritedFromGParamSpecFlags() *C.GParamSpecFlags {
	return (*C.GParamSpecFlags)(this.C)
}

func (this *ParamSpecFlags) GetStaticType() Type {
	return Type(C._g_type_param_flags())
}

//--------------------------------------------------------------
// ParamSpecFloat
//--------------------------------------------------------------

type ParamSpecFloatLike interface {
	InheritedFromGParamSpecFloat() *C.GParamSpecFloat
}

type ParamSpecFloat struct {
	ParamSpec
}

func ToParamSpecFloat(pspeclike ParamSpecFloatLike) *ParamSpecFloat {
	t := (*ParamSpecFloat)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecFloat()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecFloat)(obj)
	}
	panic("cannot cast to ParamSpecFloat")
}

func (this *ParamSpecFloat) InheritedFromGParamSpecFloat() *C.GParamSpecFloat {
	return (*C.GParamSpecFloat)(this.C)
}

func (this *ParamSpecFloat) GetStaticType() Type {
	return Type(C._g_type_param_float())
}

//--------------------------------------------------------------
// ParamSpecGType
//--------------------------------------------------------------

type ParamSpecGTypeLike interface {
	InheritedFromGParamSpecGType() *C.GParamSpecGType
}

type ParamSpecGType struct {
	ParamSpec
}

func ToParamSpecGType(pspeclike ParamSpecGTypeLike) *ParamSpecGType {
	t := (*ParamSpecGType)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecGType()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecGType)(obj)
	}
	panic("cannot cast to ParamSpecGType")
}

func (this *ParamSpecGType) InheritedFromGParamSpecGType() *C.GParamSpecGType {
	return (*C.GParamSpecGType)(this.C)
}

func (this *ParamSpecGType) GetStaticType() Type {
	return Type(C._g_type_param_gtype())
}

//--------------------------------------------------------------
// ParamSpecInt
//--------------------------------------------------------------

type ParamSpecIntLike interface {
	InheritedFromGParamSpecInt() *C.GParamSpecInt
}

type ParamSpecInt struct {
	ParamSpec
}

func ToParamSpecInt(pspeclike ParamSpecIntLike) *ParamSpecInt {
	t := (*ParamSpecInt)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecInt()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecInt)(obj)
	}
	panic("cannot cast to ParamSpecInt")
}

func (this *ParamSpecInt) InheritedFromGParamSpecInt() *C.GParamSpecInt {
	return (*C.GParamSpecInt)(this.C)
}

func (this *ParamSpecInt) GetStaticType() Type {
	return Type(C._g_type_param_int())
}

//--------------------------------------------------------------
// ParamSpecInt64
//--------------------------------------------------------------

type ParamSpecInt64Like interface {
	InheritedFromGParamSpecInt64() *C.GParamSpecInt64
}

type ParamSpecInt64 struct {
	ParamSpec
}

func ToParamSpecInt64(pspeclike ParamSpecInt64Like) *ParamSpecInt64 {
	t := (*ParamSpecInt64)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecInt64()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecInt64)(obj)
	}
	panic("cannot cast to ParamSpecInt64")
}

func (this *ParamSpecInt64) InheritedFromGParamSpecInt64() *C.GParamSpecInt64 {
	return (*C.GParamSpecInt64)(this.C)
}

func (this *ParamSpecInt64) GetStaticType() Type {
	return Type(C._g_type_param_int64())
}

//--------------------------------------------------------------
// ParamSpecLong
//--------------------------------------------------------------

type ParamSpecLongLike interface {
	InheritedFromGParamSpecLong() *C.GParamSpecLong
}

type ParamSpecLong struct {
	ParamSpec
}

func ToParamSpecLong(pspeclike ParamSpecLongLike) *ParamSpecLong {
	t := (*ParamSpecLong)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecLong()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecLong)(obj)
	}
	panic("cannot cast to ParamSpecLong")
}

func (this *ParamSpecLong) InheritedFromGParamSpecLong() *C.GParamSpecLong {
	return (*C.GParamSpecLong)(this.C)
}

func (this *ParamSpecLong) GetStaticType() Type {
	return Type(C._g_type_param_long())
}

//--------------------------------------------------------------
// ParamSpecObject
//--------------------------------------------------------------

type ParamSpecObjectLike interface {
	InheritedFromGParamSpecObject() *C.GParamSpecObject
}

type ParamSpecObject struct {
	ParamSpec
}

func ToParamSpecObject(pspeclike ParamSpecObjectLike) *ParamSpecObject {
	t := (*ParamSpecObject)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecObject()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecObject)(obj)
	}
	panic("cannot cast to ParamSpecObject")
}

func (this *ParamSpecObject) InheritedFromGParamSpecObject() *C.GParamSpecObject {
	return (*C.GParamSpecObject)(this.C)
}

func (this *ParamSpecObject) GetStaticType() Type {
	return Type(C._g_type_param_object())
}

//--------------------------------------------------------------
// ParamSpecOverride
//--------------------------------------------------------------

type ParamSpecOverrideLike interface {
	InheritedFromGParamSpecOverride() *C.GParamSpecOverride
}

type ParamSpecOverride struct {
	ParamSpec
}

func ToParamSpecOverride(pspeclike ParamSpecOverrideLike) *ParamSpecOverride {
	t := (*ParamSpecOverride)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecOverride()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecOverride)(obj)
	}
	panic("cannot cast to ParamSpecOverride")
}

func (this *ParamSpecOverride) InheritedFromGParamSpecOverride() *C.GParamSpecOverride {
	return (*C.GParamSpecOverride)(this.C)
}

func (this *ParamSpecOverride) GetStaticType() Type {
	return Type(C._g_type_param_override())
}

//--------------------------------------------------------------
// ParamSpecParam
//--------------------------------------------------------------

type ParamSpecParamLike interface {
	InheritedFromGParamSpecParam() *C.GParamSpecParam
}

type ParamSpecParam struct {
	ParamSpec
}

func ToParamSpecParam(pspeclike ParamSpecParamLike) *ParamSpecParam {
	t := (*ParamSpecParam)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecParam()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecParam)(obj)
	}
	panic("cannot cast to ParamSpecParam")
}

func (this *ParamSpecParam) InheritedFromGParamSpecParam() *C.GParamSpecParam {
	return (*C.GParamSpecParam)(this.C)
}

func (this *ParamSpecParam) GetStaticType() Type {
	return Type(C._g_type_param_param())
}

//--------------------------------------------------------------
// ParamSpecPointer
//--------------------------------------------------------------

type ParamSpecPointerLike interface {
	InheritedFromGParamSpecPointer() *C.GParamSpecPointer
}

type ParamSpecPointer struct {
	ParamSpec
}

func ToParamSpecPointer(pspeclike ParamSpecPointerLike) *ParamSpecPointer {
	t := (*ParamSpecPointer)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecPointer()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecPointer)(obj)
	}
	panic("cannot cast to ParamSpecPointer")
}

func (this *ParamSpecPointer) InheritedFromGParamSpecPointer() *C.GParamSpecPointer {
	return (*C.GParamSpecPointer)(this.C)
}

func (this *ParamSpecPointer) GetStaticType() Type {
	return Type(C._g_type_param_pointer())
}

//--------------------------------------------------------------
// ParamSpecString
//--------------------------------------------------------------

type ParamSpecStringLike interface {
	InheritedFromGParamSpecString() *C.GParamSpecString
}

type ParamSpecString struct {
	ParamSpec
}

func ToParamSpecString(pspeclike ParamSpecStringLike) *ParamSpecString {
	t := (*ParamSpecString)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecString()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecString)(obj)
	}
	panic("cannot cast to ParamSpecString")
}

func (this *ParamSpecString) InheritedFromGParamSpecString() *C.GParamSpecString {
	return (*C.GParamSpecString)(this.C)
}

func (this *ParamSpecString) GetStaticType() Type {
	return Type(C._g_type_param_string())
}

//--------------------------------------------------------------
// ParamSpecUChar
//--------------------------------------------------------------

type ParamSpecUCharLike interface {
	InheritedFromGParamSpecUChar() *C.GParamSpecUChar
}

type ParamSpecUChar struct {
	ParamSpec
}

func ToParamSpecUChar(pspeclike ParamSpecUCharLike) *ParamSpecUChar {
	t := (*ParamSpecUChar)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecUChar()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecUChar)(obj)
	}
	panic("cannot cast to ParamSpecUChar")
}

func (this *ParamSpecUChar) InheritedFromGParamSpecUChar() *C.GParamSpecUChar {
	return (*C.GParamSpecUChar)(this.C)
}

func (this *ParamSpecUChar) GetStaticType() Type {
	return Type(C._g_type_param_uchar())
}

//--------------------------------------------------------------
// ParamSpecUInt
//--------------------------------------------------------------

type ParamSpecUIntLike interface {
	InheritedFromGParamSpecUInt() *C.GParamSpecUInt
}

type ParamSpecUInt struct {
	ParamSpec
}

func ToParamSpecUInt(pspeclike ParamSpecUIntLike) *ParamSpecUInt {
	t := (*ParamSpecUInt)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecUInt()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecUInt)(obj)
	}
	panic("cannot cast to ParamSpecUInt")
}

func (this *ParamSpecUInt) InheritedFromGParamSpecUInt() *C.GParamSpecUInt {
	return (*C.GParamSpecUInt)(this.C)
}

func (this *ParamSpecUInt) GetStaticType() Type {
	return Type(C._g_type_param_uint())
}

//--------------------------------------------------------------
// ParamSpecUInt64
//--------------------------------------------------------------

type ParamSpecUInt64Like interface {
	InheritedFromGParamSpecUInt64() *C.GParamSpecUInt64
}

type ParamSpecUInt64 struct {
	ParamSpec
}

func ToParamSpecUInt64(pspeclike ParamSpecUInt64Like) *ParamSpecUInt64 {
	t := (*ParamSpecUInt64)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecUInt64()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecUInt64)(obj)
	}
	panic("cannot cast to ParamSpecUInt64")
}

func (this *ParamSpecUInt64) InheritedFromGParamSpecUInt64() *C.GParamSpecUInt64 {
	return (*C.GParamSpecUInt64)(this.C)
}

func (this *ParamSpecUInt64) GetStaticType() Type {
	return Type(C._g_type_param_uint64())
}

//--------------------------------------------------------------
// ParamSpecULong
//--------------------------------------------------------------

type ParamSpecULongLike interface {
	InheritedFromGParamSpecULong() *C.GParamSpecULong
}

type ParamSpecULong struct {
	ParamSpec
}

func ToParamSpecULong(pspeclike ParamSpecULongLike) *ParamSpecULong {
	t := (*ParamSpecULong)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecULong()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecULong)(obj)
	}
	panic("cannot cast to ParamSpecULong")
}

func (this *ParamSpecULong) InheritedFromGParamSpecULong() *C.GParamSpecULong {
	return (*C.GParamSpecULong)(this.C)
}

func (this *ParamSpecULong) GetStaticType() Type {
	return Type(C._g_type_param_ulong())
}

//--------------------------------------------------------------
// ParamSpecUnichar
//--------------------------------------------------------------

type ParamSpecUnicharLike interface {
	InheritedFromGParamSpecUnichar() *C.GParamSpecUnichar
}

type ParamSpecUnichar struct {
	ParamSpec
}

func ToParamSpecUnichar(pspeclike ParamSpecUnicharLike) *ParamSpecUnichar {
	t := (*ParamSpecUnichar)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecUnichar()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecUnichar)(obj)
	}
	panic("cannot cast to ParamSpecUnichar")
}

func (this *ParamSpecUnichar) InheritedFromGParamSpecUnichar() *C.GParamSpecUnichar {
	return (*C.GParamSpecUnichar)(this.C)
}

func (this *ParamSpecUnichar) GetStaticType() Type {
	return Type(C._g_type_param_unichar())
}

//--------------------------------------------------------------
// ParamSpecValueArray
//--------------------------------------------------------------

type ParamSpecValueArrayLike interface {
	InheritedFromGParamSpecValueArray() *C.GParamSpecValueArray
}

type ParamSpecValueArray struct {
	ParamSpec
}

func ToParamSpecValueArray(pspeclike ParamSpecValueArrayLike) *ParamSpecValueArray {
	t := (*ParamSpecValueArray)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecValueArray()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecValueArray)(obj)
	}
	panic("cannot cast to ParamSpecValueArray")
}

func (this *ParamSpecValueArray) InheritedFromGParamSpecValueArray() *C.GParamSpecValueArray {
	return (*C.GParamSpecValueArray)(this.C)
}

func (this *ParamSpecValueArray) GetStaticType() Type {
	return Type(C._g_type_param_value_array())
}

//--------------------------------------------------------------
// ParamSpecVariant
//--------------------------------------------------------------

type ParamSpecVariantLike interface {
	InheritedFromGParamSpecVariant() *C.GParamSpecVariant
}

type ParamSpecVariant struct {
	ParamSpec
}

func ToParamSpecVariant(pspeclike ParamSpecVariantLike) *ParamSpecVariant {
	t := (*ParamSpecVariant)(nil).GetStaticType()
	c := pspeclike.InheritedFromGParamSpecVariant()
	obj := ParamSpecGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*ParamSpecVariant)(obj)
	}
	panic("cannot cast to ParamSpecVariant")
}

func (this *ParamSpecVariant) InheritedFromGParamSpecVariant() *C.GParamSpecVariant {
	return (*C.GParamSpecVariant)(this.C)
}

func (this *ParamSpecVariant) GetStaticType() Type {
	return Type(C._g_type_param_variant())
}
