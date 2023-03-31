// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package gobject

/*
#include "gobject.gen.h"
#include <string.h>

extern void g_free(void*);

#include "gobject.h"

extern uint32_t g_quark_from_string(const char*);
extern void g_object_set_qdata(GObject*, uint32_t, void*);

extern void g_type_init();

#cgo pkg-config: gobject-2.0
*/
import "C"
import "unsafe"
import "runtime"
import "reflect"


const alot = 999999

type _GSList struct {
	data unsafe.Pointer
	next *_GSList
}

type _GList struct {
	data unsafe.Pointer
	next *_GList
	prev *_GList
}

type _GError struct {
	domain uint32
	code int32
	message *C.char
}
func (e _GError) ToGError() GError {
	return GError{e.domain, e.code, C.GoString(e.message)}
}

type GError struct {
	Domain uint32
	Code int32
	Message string
}
func (e GError) Error() string {
	return e.Message
}

func _GoStringToGString(x string) *C.char {
	if x == "\x00" {
		return nil
	}
	return C.CString(x)
}

func _GoBoolToCBool(x bool) C.int {
	if x { return 1 }
	return 0
}

func _CInterfaceToGoInterface(iface [2]unsafe.Pointer) interface{} {
	return *(*interface{})(unsafe.Pointer(&iface))
}

func _GoInterfaceToCInterface(iface interface{}) *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Pointer(&iface))
}


//export _GObject_go_callback_cleanup
func _GObject_go_callback_cleanup(gofunc unsafe.Pointer) {
	Holder.Release(gofunc)
}


// blacklisted: BaseFinalizeFunc (callback)
// blacklisted: BaseInitFunc (callback)
type BindingLike interface {
	ObjectLike
	InheritedFromGBinding() *C.GBinding
}

type Binding struct {
	Object
	
}

func ToBinding(objlike ObjectLike) *Binding {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*Binding)(nil).GetStaticType()
	obj := ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*Binding)(obj)
	}
	panic("cannot cast to Binding")
}

func (this0 *Binding) InheritedFromGBinding() *C.GBinding {
	if this0 == nil {
		return nil
	}
	return (*C.GBinding)(this0.C)
}

func (this0 *Binding) GetStaticType() Type {
	return Type(C.g_binding_get_type())
}

func BindingGetType() Type {
	return (*Binding)(nil).GetStaticType()
}
func (this0 *Binding) DupSource() *Object {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_dup_source(this1)
	var ret2 *Object

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Object)(ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Binding) DupTarget() *Object {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_dup_target(this1)
	var ret2 *Object

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Object)(ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Binding) GetFlags() BindingFlags {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_get_flags(this1)
	var ret2 BindingFlags

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = BindingFlags(ret1)
	return ret2
}
func (this0 *Binding) GetSource() *Object {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_get_source(this1)
	var ret2 *Object

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Object)(ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Binding) GetSourceProperty() string {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_get_source_property(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Binding) GetTarget() *Object {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_get_target(this1)
	var ret2 *Object

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Object)(ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Binding) GetTargetProperty() string {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	ret1 := C.g_binding_get_target_property(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Binding) Unbind() {
	var this1 *C.GBinding
	if this0 != nil {
		this1 = (*C.GBinding)(this0.InheritedFromGBinding())
	}
	C.g_binding_unbind(this1)
}
type BindingFlags C.uint32_t
const (
	BindingFlagsDefault BindingFlags = 0
	BindingFlagsBidirectional BindingFlags = 1
	BindingFlagsSyncCreate BindingFlags = 2
	BindingFlagsInvertBoolean BindingFlags = 4
)
// blacklisted: BindingGroup (object)
// blacklisted: BindingTransformFunc (callback)
// blacklisted: BoxedCopyFunc (callback)
// blacklisted: BoxedFreeFunc (callback)
type CClosure struct {
	Closure Closure
	Callback unsafe.Pointer
}
func CClosureMarshalBoolean_BoxedBoxed(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_BOOLEAN__BOXED_BOXED(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalBoolean_Flags(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_BOOLEAN__FLAGS(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalString_ObjectPointer(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_STRING__OBJECT_POINTER(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Boolean(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__BOOLEAN(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Boxed(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__BOXED(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Char(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__CHAR(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Double(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__DOUBLE(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Enum(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__ENUM(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Flags(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__FLAGS(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Float(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__FLOAT(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Int(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__INT(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Long(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__LONG(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Object(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__OBJECT(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Param(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__PARAM(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Pointer(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__POINTER(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_String(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__STRING(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Uchar(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__UCHAR(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Uint(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__UINT(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_UintPointer(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__UINT_POINTER(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Ulong(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__ULONG(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Variant(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__VARIANT(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalVoid_Void(closure0 *Closure, return_value0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_value1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_value1 = (*C.GValue)(unsafe.Pointer(return_value0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_VOID__VOID(closure1, return_value1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
func CClosureMarshalGeneric(closure0 *Closure, return_gvalue0 *Value, n_param_values0 uint32, param_values0 *Value, invocation_hint0 unsafe.Pointer, marshal_data0 unsafe.Pointer) {
	var closure1 *C.GClosure
	var return_gvalue1 *C.GValue
	var n_param_values1 C.uint32_t
	var param_values1 *C.GValue
	var invocation_hint1 unsafe.Pointer
	var marshal_data1 unsafe.Pointer
	closure1 = (*C.GClosure)(unsafe.Pointer(closure0))
	return_gvalue1 = (*C.GValue)(unsafe.Pointer(return_gvalue0))
	n_param_values1 = C.uint32_t(n_param_values0)
	param_values1 = (*C.GValue)(unsafe.Pointer(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	marshal_data1 = unsafe.Pointer(marshal_data0)
	C.g_cclosure_marshal_generic(closure1, return_gvalue1, n_param_values1, param_values1, invocation_hint1, marshal_data1)
}
// blacklisted: Callback (callback)
// blacklisted: ClassFinalizeFunc (callback)
// blacklisted: ClassInitFunc (callback)
type Closure struct {
	RefCount uint32
	MetaMarshalNouse uint32
	NGuards uint32
	NFnotifiers uint32
	NInotifiers uint32
	InInotify uint32
	Floating uint32
	DerivativeFlag uint32
	InMarshal uint32
	IsInvalid uint32
	Marshal unsafe.Pointer
	Data unsafe.Pointer
	Notifiers *ClosureNotifyData
}
func NewClosureObject(sizeof_closure0 uint32, object0 ObjectLike) *Closure {
	var sizeof_closure1 C.uint32_t
	var object1 *C.GObject
	sizeof_closure1 = C.uint32_t(sizeof_closure0)
	if object0 != nil {
		object1 = (*C.GObject)(object0.InheritedFromGObject())
	}
	ret1 := C.g_closure_new_object(sizeof_closure1, object1)
	var ret2 *Closure

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Closure)(unsafe.Pointer(ret1))
	return ret2
}
func NewClosureSimple(sizeof_closure0 uint32, data0 unsafe.Pointer) *Closure {
	var sizeof_closure1 C.uint32_t
	var data1 unsafe.Pointer
	sizeof_closure1 = C.uint32_t(sizeof_closure0)
	data1 = unsafe.Pointer(data0)
	ret1 := C.g_closure_new_simple(sizeof_closure1, data1)
	var ret2 *Closure

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Closure)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *Closure) Invalidate() {
	var this1 *C.GClosure
	this1 = (*C.GClosure)(unsafe.Pointer(this0))
	C.g_closure_invalidate(this1)
}
func (this0 *Closure) Invoke(param_values0 []Value, invocation_hint0 unsafe.Pointer) Value {
	var this1 *C.GClosure
	var param_values1 *C.GValue
	var n_param_values1 C.uint32_t
	var invocation_hint1 unsafe.Pointer
	var return_value1 C.GValue
	this1 = (*C.GClosure)(unsafe.Pointer(this0))
	param_values1 = (*C.GValue)(C.malloc(C.size_t(int(unsafe.Sizeof(*param_values1)) * len(param_values0))))
	defer C.free(unsafe.Pointer(param_values1))
	for i, e := range param_values0 {
		(*(*[999999]C.GValue)(unsafe.Pointer(param_values1)))[i] = *(*C.GValue)(unsafe.Pointer(&e))
	}
	n_param_values1 = C.uint32_t(len(param_values0))
	invocation_hint1 = unsafe.Pointer(invocation_hint0)
	C.g_closure_invoke(this1, &return_value1, n_param_values1, param_values1, invocation_hint1)
	var return_value2 Value

//DEBUG: return_value1(interface):flags = " conv_own_none"
	return_value2 = *(*Value)(unsafe.Pointer(&return_value1))
	return return_value2
}
func (this0 *Closure) Sink() {
	var this1 *C.GClosure
	this1 = (*C.GClosure)(unsafe.Pointer(this0))
	C.g_closure_sink(this1)
}
// blacklisted: ClosureMarshal (callback)
// blacklisted: ClosureNotify (callback)
type ClosureNotifyData struct {
	Data unsafe.Pointer
	Notify unsafe.Pointer
}
type ConnectFlags C.uint32_t
const (
	ConnectFlagsDefault ConnectFlags = 0
	ConnectFlagsAfter ConnectFlags = 1
	ConnectFlagsSwapped ConnectFlags = 2
)
type EnumClass struct {
	GTypeClass TypeClass
	Minimum int32
	Maximum int32
	NValues uint32
	_ [4]byte
	Values *EnumValue
}
type EnumValue struct {
	Value int32
	_ [4]byte
	value_name0 *C.char
	value_nick0 *C.char
}
func (this0 *EnumValue) ValueName() string {
	var value_name1 string
	value_name1 = C.GoString(this0.value_name0)
	return value_name1
}
func (this0 *EnumValue) ValueNick() string {
	var value_nick1 string
	value_nick1 = C.GoString(this0.value_nick0)
	return value_nick1
}
type FlagsClass struct {
	GTypeClass TypeClass
	Mask uint32
	NValues uint32
	Values *FlagsValue
}
type FlagsValue struct {
	Value uint32
	_ [4]byte
	value_name0 *C.char
	value_nick0 *C.char
}
func (this0 *FlagsValue) ValueName() string {
	var value_name1 string
	value_name1 = C.GoString(this0.value_name0)
	return value_name1
}
func (this0 *FlagsValue) ValueNick() string {
	var value_nick1 string
	value_nick1 = C.GoString(this0.value_nick0)
	return value_nick1
}
type InitiallyUnownedLike interface {
	ObjectLike
	InheritedFromGInitiallyUnowned() *C.GInitiallyUnowned
}

type InitiallyUnowned struct {
	Object
	
}

func ToInitiallyUnowned(objlike ObjectLike) *InitiallyUnowned {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*InitiallyUnowned)(nil).GetStaticType()
	obj := ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*InitiallyUnowned)(obj)
	}
	panic("cannot cast to InitiallyUnowned")
}

func (this0 *InitiallyUnowned) InheritedFromGInitiallyUnowned() *C.GInitiallyUnowned {
	if this0 == nil {
		return nil
	}
	return (*C.GInitiallyUnowned)(this0.C)
}

func (this0 *InitiallyUnowned) GetStaticType() Type {
	return Type(C.g_initially_unowned_get_type())
}

func InitiallyUnownedGetType() Type {
	return (*InitiallyUnowned)(nil).GetStaticType()
}
// blacklisted: InstanceInitFunc (callback)
// blacklisted: InterfaceFinalizeFunc (callback)
type InterfaceInfo struct {
	InterfaceInit unsafe.Pointer
	InterfaceFinalize unsafe.Pointer
	InterfaceData unsafe.Pointer
}
// blacklisted: InterfaceInitFunc (callback)
type ObjectLike interface {
	
	InheritedFromGObject() *C.GObject
}

type Object struct {
	C unsafe.Pointer
	
}

func ToObject(objlike ObjectLike) *Object {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*Object)(nil).GetStaticType()
	obj := ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*Object)(obj)
	}
	panic("cannot cast to Object")
}

func (this0 *Object) InheritedFromGObject() *C.GObject {
	if this0 == nil {
		return nil
	}
	return (*C.GObject)(this0.C)
}

func (this0 *Object) GetStaticType() Type {
	return Type(C.g_object_get_type())
}

func ObjectGetType() Type {
	return (*Object)(nil).GetStaticType()
}
// blacklisted: Object.newv (method)
// blacklisted: Object.compat_control (method)
// blacklisted: Object.interface_find_property (method)
// blacklisted: Object.interface_install_property (method)
// blacklisted: Object.interface_list_properties (method)
// blacklisted: Object.bind_property (method)
// blacklisted: Object.bind_property_full (method)
// blacklisted: Object.force_floating (method)
// blacklisted: Object.freeze_notify (method)
// blacklisted: Object.get_data (method)
// blacklisted: Object.get_property (method)
// blacklisted: Object.get_qdata (method)
// blacklisted: Object.getv (method)
// blacklisted: Object.is_floating (method)
// blacklisted: Object.notify (method)
// blacklisted: Object.notify_by_pspec (method)
// blacklisted: Object.ref (method)
// blacklisted: Object.ref_sink (method)
// blacklisted: Object.run_dispose (method)
// blacklisted: Object.set_data (method)
// blacklisted: Object.set_property (method)
// blacklisted: Object.steal_data (method)
// blacklisted: Object.steal_qdata (method)
// blacklisted: Object.thaw_notify (method)
// blacklisted: Object.unref (method)
// blacklisted: Object.watch_closure (method)
type ObjectConstructParam struct {
	pspec0 *C.GParamSpec
	Value *Value
}
func (this0 *ObjectConstructParam) Pspec() *ParamSpec {
	var pspec1 *ParamSpec
	pspec1 = (*ParamSpec)(ObjectWrap(unsafe.Pointer(this0.pspec0), true))
	return pspec1
}
// blacklisted: ObjectFinalizeFunc (callback)
// blacklisted: ObjectGetPropertyFunc (callback)
// blacklisted: ObjectSetPropertyFunc (callback)
const ParamMask = 255
const ParamStaticStrings = 224
const ParamUserShift = 8
type ParamFlags C.uint32_t
const (
	ParamFlagsReadable ParamFlags = 1
	ParamFlagsWritable ParamFlags = 2
	ParamFlagsReadwrite ParamFlags = 3
	ParamFlagsConstruct ParamFlags = 4
	ParamFlagsConstructOnly ParamFlags = 8
	ParamFlagsLaxValidation ParamFlags = 16
	ParamFlagsStaticName ParamFlags = 32
	ParamFlagsPrivate ParamFlags = 32
	ParamFlagsStaticNick ParamFlags = 64
	ParamFlagsStaticBlurb ParamFlags = 128
	ParamFlagsExplicitNotify ParamFlags = 1073741824
	ParamFlagsDeprecated ParamFlags = 2147483648
)
// blacklisted: ParamSpec (object)
// blacklisted: ParamSpecBoolean (object)
// blacklisted: ParamSpecBoxed (object)
// blacklisted: ParamSpecChar (object)
// blacklisted: ParamSpecDouble (object)
// blacklisted: ParamSpecEnum (object)
// blacklisted: ParamSpecFlags (object)
// blacklisted: ParamSpecFloat (object)
// blacklisted: ParamSpecGType (object)
// blacklisted: ParamSpecInt (object)
// blacklisted: ParamSpecInt64 (object)
// blacklisted: ParamSpecLong (object)
// blacklisted: ParamSpecObject (object)
// blacklisted: ParamSpecOverride (object)
// blacklisted: ParamSpecParam (object)
// blacklisted: ParamSpecPointer (object)
// blacklisted: ParamSpecPool (struct)
// blacklisted: ParamSpecString (object)
type ParamSpecTypeInfo struct {
	InstanceSize uint16
	NPreallocs uint16
	_ [4]byte
	InstanceInit unsafe.Pointer
	ValueType Type
	Finalize unsafe.Pointer
	ValueSetDefault unsafe.Pointer
	ValueValidate unsafe.Pointer
	ValuesCmp unsafe.Pointer
}
// blacklisted: ParamSpecUChar (object)
// blacklisted: ParamSpecUInt (object)
// blacklisted: ParamSpecUInt64 (object)
// blacklisted: ParamSpecULong (object)
// blacklisted: ParamSpecUnichar (object)
// blacklisted: ParamSpecValueArray (object)
// blacklisted: ParamSpecVariant (object)
type Parameter struct {
	name0 *C.char
	Value Value
}
func (this0 *Parameter) Name() string {
	var name1 string
	name1 = C.GoString(this0.name0)
	return name1
}
const SignalFlagsMask = 511
const SignalMatchMask = 63
// blacklisted: SignalAccumulator (callback)
// blacklisted: SignalEmissionHook (callback)
type SignalFlags C.uint32_t
const (
	SignalFlagsRunFirst SignalFlags = 1
	SignalFlagsRunLast SignalFlags = 2
	SignalFlagsRunCleanup SignalFlags = 4
	SignalFlagsNoRecurse SignalFlags = 8
	SignalFlagsDetailed SignalFlags = 16
	SignalFlagsAction SignalFlags = 32
	SignalFlagsNoHooks SignalFlags = 64
	SignalFlagsMustCollect SignalFlags = 128
	SignalFlagsDeprecated SignalFlags = 256
	SignalFlagsAccumulatorFirstRun SignalFlags = 131072
)
// blacklisted: SignalGroup (object)
type SignalInvocationHint struct {
	SignalId uint32
	Detail uint32
	RunType SignalFlags
}
type SignalMatchType C.uint32_t
const (
	SignalMatchTypeId SignalMatchType = 1
	SignalMatchTypeDetail SignalMatchType = 2
	SignalMatchTypeClosure SignalMatchType = 4
	SignalMatchTypeFunc SignalMatchType = 8
	SignalMatchTypeData SignalMatchType = 16
	SignalMatchTypeUnblocked SignalMatchType = 32
)
type SignalQuery struct {
	SignalId uint32
	_ [4]byte
	signal_name0 *C.char
	Itype Type
	SignalFlags SignalFlags
	_ [4]byte
	ReturnType Type
	NParams uint32
	_ [4]byte
	param_types0 *C.GType
}
func (this0 *SignalQuery) SignalName() string {
	var signal_name1 string
	signal_name1 = C.GoString(this0.signal_name0)
	return signal_name1
}
func (this0 *SignalQuery) ParamTypes() []Type {
	var param_types1 []Type
	for i0 := range param_types1 {
		param_types1[i0] = Type((*(*[999999]C.GType)(unsafe.Pointer(this0.param_types0)))[i0])
	}
	return param_types1
}
const TypeFlagReservedIdBit = 0x1
const TypeFundamentalMax = 255
const TypeFundamentalShift = 2
const TypeReservedBseFirst = 32
const TypeReservedBseLast = 48
const TypeReservedGlibFirst = 22
const TypeReservedGlibLast = 31
const TypeReservedUserFirst = 49
// blacklisted: ToggleNotify (callback)
type TypeCValue struct {
	_data [0]byte
}
type TypeClass struct {
	GType Type
}
func (this0 *TypeClass) AddPrivate(private_size0 uint64) {
	var this1 *C.GTypeClass
	var private_size1 C.uint64_t
	this1 = (*C.GTypeClass)(unsafe.Pointer(this0))
	private_size1 = C.uint64_t(private_size0)
	C.g_type_class_add_private(this1, private_size1)
}
func (this0 *TypeClass) GetPrivate(private_type0 Type) {
	var this1 *C.GTypeClass
	var private_type1 C.GType
	this1 = (*C.GTypeClass)(unsafe.Pointer(this0))
	private_type1 = C.GType(private_type0)
	C.g_type_class_get_private(this1, private_type1)
}
func (this0 *TypeClass) PeekParent() *TypeClass {
	var this1 *C.GTypeClass
	this1 = (*C.GTypeClass)(unsafe.Pointer(this0))
	ret1 := C.g_type_class_peek_parent(this1)
	var ret2 *TypeClass

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*TypeClass)(unsafe.Pointer(ret1))
	return ret2
}
func TypeClassAdjustPrivateOffset(g_class0 unsafe.Pointer, private_size_or_offset0 *int32) {
	var g_class1 unsafe.Pointer
	var private_size_or_offset1 *C.int32_t
	g_class1 = unsafe.Pointer(g_class0)
	private_size_or_offset1 = (*C.int32_t)(unsafe.Pointer(private_size_or_offset0))
	C.g_type_class_adjust_private_offset(g_class1, private_size_or_offset1)
}
func TypeClassPeek(type0 Type) *TypeClass {
	var type1 C.GType
	type1 = C.GType(type0)
	ret1 := C.g_type_class_peek(type1)
	var ret2 *TypeClass

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*TypeClass)(unsafe.Pointer(ret1))
	return ret2
}
func TypeClassPeekStatic(type0 Type) *TypeClass {
	var type1 C.GType
	type1 = C.GType(type0)
	ret1 := C.g_type_class_peek_static(type1)
	var ret2 *TypeClass

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*TypeClass)(unsafe.Pointer(ret1))
	return ret2
}
// blacklisted: TypeClassCacheFunc (callback)
type TypeDebugFlags C.uint32_t
const (
	TypeDebugFlagsNone TypeDebugFlags = 0
	TypeDebugFlagsObjects TypeDebugFlags = 1
	TypeDebugFlagsSignals TypeDebugFlags = 2
	TypeDebugFlagsInstanceCount TypeDebugFlags = 4
	TypeDebugFlagsMask TypeDebugFlags = 7
)
type TypeFlags C.uint32_t
const (
	TypeFlagsNone TypeFlags = 0
	TypeFlagsAbstract TypeFlags = 16
	TypeFlagsValueAbstract TypeFlags = 32
	TypeFlagsFinal TypeFlags = 64
	TypeFlagsDeprecated TypeFlags = 128
)
type TypeFundamentalFlags C.uint32_t
const (
	TypeFundamentalFlagsClassed TypeFundamentalFlags = 1
	TypeFundamentalFlagsInstantiatable TypeFundamentalFlags = 2
	TypeFundamentalFlagsDerivable TypeFundamentalFlags = 4
	TypeFundamentalFlagsDeepDerivable TypeFundamentalFlags = 8
)
type TypeFundamentalInfo struct {
	TypeFlags TypeFundamentalFlags
}
type TypeInfo struct {
	ClassSize uint16
	_ [6]byte
	BaseInit unsafe.Pointer
	BaseFinalize unsafe.Pointer
	ClassInit unsafe.Pointer
	ClassFinalize unsafe.Pointer
	ClassData unsafe.Pointer
	InstanceSize uint16
	NPreallocs uint16
	_ [4]byte
	InstanceInit unsafe.Pointer
	ValueTable *TypeValueTable
}
type TypeInstance struct {
	GClass *TypeClass
}
func (this0 *TypeInstance) GetPrivate(private_type0 Type) {
	var this1 *C.GTypeInstance
	var private_type1 C.GType
	this1 = (*C.GTypeInstance)(unsafe.Pointer(this0))
	private_type1 = C.GType(private_type0)
	C.g_type_instance_get_private(this1, private_type1)
}
type TypeInterface struct {
	GType Type
	GInstanceType Type
}
func (this0 *TypeInterface) PeekParent() *TypeInterface {
	var this1 *C.GTypeInterface
	this1 = (*C.GTypeInterface)(unsafe.Pointer(this0))
	ret1 := C.g_type_interface_peek_parent(this1)
	var ret2 *TypeInterface

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*TypeInterface)(unsafe.Pointer(ret1))
	return ret2
}
func TypeInterfaceAddPrerequisite(interface_type0 Type, prerequisite_type0 Type) {
	var interface_type1 C.GType
	var prerequisite_type1 C.GType
	interface_type1 = C.GType(interface_type0)
	prerequisite_type1 = C.GType(prerequisite_type0)
	C.g_type_interface_add_prerequisite(interface_type1, prerequisite_type1)
}
func TypeInterfaceGetPlugin(instance_type0 Type, interface_type0 Type) *TypePlugin {
	var instance_type1 C.GType
	var interface_type1 C.GType
	instance_type1 = C.GType(instance_type0)
	interface_type1 = C.GType(interface_type0)
	ret1 := C.g_type_interface_get_plugin(instance_type1, interface_type1)
	var ret2 *TypePlugin

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*TypePlugin)(ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func TypeInterfaceInstantiatablePrerequisite(interface_type0 Type) Type {
	var interface_type1 C.GType
	interface_type1 = C.GType(interface_type0)
	ret1 := C.g_type_interface_instantiatable_prerequisite(interface_type1)
	var ret2 Type

//DEBUG: ret1(GType):flags = " conv_own_none"
	ret2 = Type(ret1)
	return ret2
}
func TypeInterfacePeek(instance_class0 *TypeClass, iface_type0 Type) *TypeInterface {
	var instance_class1 *C.GTypeClass
	var iface_type1 C.GType
	instance_class1 = (*C.GTypeClass)(unsafe.Pointer(instance_class0))
	iface_type1 = C.GType(iface_type0)
	ret1 := C.g_type_interface_peek(instance_class1, iface_type1)
	var ret2 *TypeInterface

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*TypeInterface)(unsafe.Pointer(ret1))
	return ret2
}
func TypeInterfacePrerequisites(interface_type0 Type) (uint32, []Type) {
	var interface_type1 C.GType
	var n_prerequisites1 C.uint32_t
	interface_type1 = C.GType(interface_type0)
	ret1 := C.g_type_interface_prerequisites(interface_type1, &n_prerequisites1)
	var n_prerequisites2 uint32
	var ret2 []Type

//DEBUG: n_prerequisites1(guint32):flags = " conv_own_everything"
	n_prerequisites2 = uint32(n_prerequisites1)
	ret2 = make([]Type, n_prerequisites1)

//DEBUG: ret1(array):flags = " conv_own_everything"
	for i0 := range ret2 {
		ret2[i0] = Type((*(*[999999]C.GType)(unsafe.Pointer(ret1)))[i0])
	}
	C.g_free(unsafe.Pointer(ret1))
	return n_prerequisites2, ret2
}
// blacklisted: TypeInterfaceCheckFunc (callback)
// blacklisted: TypeModule (object)
type TypePluginLike interface {
	ImplementsGTypePlugin() *C.GTypePlugin
}

type TypePlugin struct {
	Object
	TypePluginImpl
}

func (*TypePlugin) GetStaticType() Type {
	return Type(C.g_type_plugin_get_type())
}


type TypePluginImpl struct {}

func ToTypePlugin(objlike ObjectLike) *TypePlugin {
	c := objlike.InheritedFromGObject()
	obj := ObjectGrabIfType(unsafe.Pointer(c), Type(C.g_type_plugin_get_type()))
	if obj != nil {
		return (*TypePlugin)(obj)
	}
	panic("cannot cast to TypePlugin")
}

func (this0 *TypePluginImpl) ImplementsGTypePlugin() *C.GTypePlugin {
	base := unsafe.Pointer(uintptr(unsafe.Pointer(this0)) - unsafe.Sizeof(uintptr(0)))
	return (*C.GTypePlugin)((*Object)(base).C)
}
func (this0 *TypePluginImpl) CompleteInterfaceInfo(instance_type0 Type, interface_type0 Type, info0 *InterfaceInfo) {
	var this1 *C.GTypePlugin
	var instance_type1 C.GType
	var interface_type1 C.GType
	var info1 *C.GInterfaceInfo
	if this0 != nil {
		this1 = this0.ImplementsGTypePlugin()
	}
	instance_type1 = C.GType(instance_type0)
	interface_type1 = C.GType(interface_type0)
	info1 = (*C.GInterfaceInfo)(unsafe.Pointer(info0))
	C.g_type_plugin_complete_interface_info(this1, instance_type1, interface_type1, info1)
}
func (this0 *TypePluginImpl) CompleteTypeInfo(g_type0 Type, info0 *TypeInfo, value_table0 *TypeValueTable) {
	var this1 *C.GTypePlugin
	var g_type1 C.GType
	var info1 *C.GTypeInfo
	var value_table1 *C.GTypeValueTable
	if this0 != nil {
		this1 = this0.ImplementsGTypePlugin()
	}
	g_type1 = C.GType(g_type0)
	info1 = (*C.GTypeInfo)(unsafe.Pointer(info0))
	value_table1 = (*C.GTypeValueTable)(unsafe.Pointer(value_table0))
	C.g_type_plugin_complete_type_info(this1, g_type1, info1, value_table1)
}
func (this0 *TypePluginImpl) Unuse() {
	var this1 *C.GTypePlugin
	if this0 != nil {
		this1 = this0.ImplementsGTypePlugin()
	}
	C.g_type_plugin_unuse(this1)
}
func (this0 *TypePluginImpl) Use() {
	var this1 *C.GTypePlugin
	if this0 != nil {
		this1 = this0.ImplementsGTypePlugin()
	}
	C.g_type_plugin_use(this1)
}
type TypePluginClass struct {
	BaseIface TypeInterface
	UsePlugin unsafe.Pointer
	UnusePlugin unsafe.Pointer
	CompleteTypeInfo unsafe.Pointer
	CompleteInterfaceInfo unsafe.Pointer
}
// blacklisted: TypePluginCompleteInterfaceInfo (callback)
// blacklisted: TypePluginCompleteTypeInfo (callback)
// blacklisted: TypePluginUnuse (callback)
// blacklisted: TypePluginUse (callback)
type TypeQuery struct {
	Type Type
	type_name0 *C.char
	ClassSize uint32
	InstanceSize uint32
}
func (this0 *TypeQuery) TypeName() string {
	var type_name1 string
	type_name1 = C.GoString(this0.type_name0)
	return type_name1
}
type TypeValueTable struct {
	ValueInit unsafe.Pointer
	ValueFree unsafe.Pointer
	ValueCopy unsafe.Pointer
	ValuePeekPointer unsafe.Pointer
	collect_format0 *C.char
	CollectValue unsafe.Pointer
	lcopy_format0 *C.char
	LcopyValue unsafe.Pointer
}
func (this0 *TypeValueTable) CollectFormat() string {
	var collect_format1 string
	collect_format1 = C.GoString(this0.collect_format0)
	return collect_format1
}
func (this0 *TypeValueTable) LcopyFormat() string {
	var lcopy_format1 string
	lcopy_format1 = C.GoString(this0.lcopy_format0)
	return lcopy_format1
}
const ValueInternedString = 268435456
const ValueNocopyContents = 134217728
type Value struct {
	GType Type
	Data [2]_Value__data__union
}
type ValueArray struct {
	NValues uint32
	_ [4]byte
	Values *Value
	NPrealloced uint32
	_ [4]byte
}
func NewValueArray(n_prealloced0 uint32) *ValueArray {
	var n_prealloced1 C.uint32_t
	n_prealloced1 = C.uint32_t(n_prealloced0)
	ret1 := C.g_value_array_new(n_prealloced1)
	var ret2 *ValueArray

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*ValueArray)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *ValueArray) Append(value0 *Value) *ValueArray {
	var this1 *C.GValueArray
	var value1 *C.GValue
	this1 = (*C.GValueArray)(unsafe.Pointer(this0))
	value1 = (*C.GValue)(unsafe.Pointer(value0))
	ret1 := C.g_value_array_append(this1, value1)
	var ret2 *ValueArray

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*ValueArray)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *ValueArray) Copy() *ValueArray {
	var this1 *C.GValueArray
	this1 = (*C.GValueArray)(unsafe.Pointer(this0))
	ret1 := C.g_value_array_copy(this1)
	var ret2 *ValueArray

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*ValueArray)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *ValueArray) GetNth(index_0 uint32) *Value {
	var this1 *C.GValueArray
	var index_1 C.uint32_t
	this1 = (*C.GValueArray)(unsafe.Pointer(this0))
	index_1 = C.uint32_t(index_0)
	ret1 := C.g_value_array_get_nth(this1, index_1)
	var ret2 *Value

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Value)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *ValueArray) Insert(index_0 uint32, value0 *Value) *ValueArray {
	var this1 *C.GValueArray
	var index_1 C.uint32_t
	var value1 *C.GValue
	this1 = (*C.GValueArray)(unsafe.Pointer(this0))
	index_1 = C.uint32_t(index_0)
	value1 = (*C.GValue)(unsafe.Pointer(value0))
	ret1 := C.g_value_array_insert(this1, index_1, value1)
	var ret2 *ValueArray

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*ValueArray)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *ValueArray) Prepend(value0 *Value) *ValueArray {
	var this1 *C.GValueArray
	var value1 *C.GValue
	this1 = (*C.GValueArray)(unsafe.Pointer(this0))
	value1 = (*C.GValue)(unsafe.Pointer(value0))
	ret1 := C.g_value_array_prepend(this1, value1)
	var ret2 *ValueArray

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*ValueArray)(unsafe.Pointer(ret1))
	return ret2
}
func (this0 *ValueArray) Remove(index_0 uint32) *ValueArray {
	var this1 *C.GValueArray
	var index_1 C.uint32_t
	this1 = (*C.GValueArray)(unsafe.Pointer(this0))
	index_1 = C.uint32_t(index_0)
	ret1 := C.g_value_array_remove(this1, index_1)
	var ret2 *ValueArray

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*ValueArray)(unsafe.Pointer(ret1))
	return ret2
}
// blacklisted: ValueTransform (callback)
// blacklisted: WeakNotify (callback)
type WeakRef struct {}
type _Value__data__union struct {
	_data [8]byte
}
// blacklisted: boxed_copy (function)
// blacklisted: boxed_free (function)
// blacklisted: cclosure_marshal_BOOLEAN__BOXED_BOXED (function)
// blacklisted: cclosure_marshal_BOOLEAN__FLAGS (function)
// blacklisted: cclosure_marshal_STRING__OBJECT_POINTER (function)
// blacklisted: cclosure_marshal_VOID__BOOLEAN (function)
// blacklisted: cclosure_marshal_VOID__BOXED (function)
// blacklisted: cclosure_marshal_VOID__CHAR (function)
// blacklisted: cclosure_marshal_VOID__DOUBLE (function)
// blacklisted: cclosure_marshal_VOID__ENUM (function)
// blacklisted: cclosure_marshal_VOID__FLAGS (function)
// blacklisted: cclosure_marshal_VOID__FLOAT (function)
// blacklisted: cclosure_marshal_VOID__INT (function)
// blacklisted: cclosure_marshal_VOID__LONG (function)
// blacklisted: cclosure_marshal_VOID__OBJECT (function)
// blacklisted: cclosure_marshal_VOID__PARAM (function)
// blacklisted: cclosure_marshal_VOID__POINTER (function)
// blacklisted: cclosure_marshal_VOID__STRING (function)
// blacklisted: cclosure_marshal_VOID__UCHAR (function)
// blacklisted: cclosure_marshal_VOID__UINT (function)
// blacklisted: cclosure_marshal_VOID__UINT_POINTER (function)
// blacklisted: cclosure_marshal_VOID__ULONG (function)
// blacklisted: cclosure_marshal_VOID__VARIANT (function)
// blacklisted: cclosure_marshal_VOID__VOID (function)
// blacklisted: cclosure_marshal_generic (function)
// blacklisted: clear_signal_handler (function)
// blacklisted: enum_complete_type_info (function)
// blacklisted: enum_get_value (function)
// blacklisted: enum_get_value_by_name (function)
// blacklisted: enum_get_value_by_nick (function)
// blacklisted: enum_register_static (function)
// blacklisted: enum_to_string (function)
// blacklisted: flags_complete_type_info (function)
// blacklisted: flags_get_first_value (function)
// blacklisted: flags_get_value_by_name (function)
// blacklisted: flags_get_value_by_nick (function)
// blacklisted: flags_register_static (function)
// blacklisted: flags_to_string (function)
// blacklisted: gtype_get_type (function)
// blacklisted: param_spec_boolean (function)
// blacklisted: param_spec_boxed (function)
// blacklisted: param_spec_char (function)
// blacklisted: param_spec_double (function)
// blacklisted: param_spec_enum (function)
// blacklisted: param_spec_flags (function)
// blacklisted: param_spec_float (function)
// blacklisted: param_spec_gtype (function)
// blacklisted: param_spec_int (function)
// blacklisted: param_spec_int64 (function)
// blacklisted: param_spec_long (function)
// blacklisted: param_spec_object (function)
// blacklisted: param_spec_param (function)
// blacklisted: param_spec_pointer (function)
// blacklisted: param_spec_string (function)
// blacklisted: param_spec_uchar (function)
// blacklisted: param_spec_uint (function)
// blacklisted: param_spec_uint64 (function)
// blacklisted: param_spec_ulong (function)
// blacklisted: param_spec_unichar (function)
// blacklisted: param_spec_variant (function)
// blacklisted: param_type_register_static (function)
// blacklisted: param_value_convert (function)
// blacklisted: param_value_defaults (function)
// blacklisted: param_value_is_valid (function)
// blacklisted: param_value_set_default (function)
// blacklisted: param_value_validate (function)
// blacklisted: param_values_cmp (function)
// blacklisted: pointer_type_register_static (function)
// blacklisted: signal_accumulator_first_wins (function)
// blacklisted: signal_accumulator_true_handled (function)
// blacklisted: signal_add_emission_hook (function)
// blacklisted: signal_chain_from_overridden (function)
// blacklisted: signal_connect_closure (function)
// blacklisted: signal_connect_closure_by_id (function)
// blacklisted: signal_emitv (function)
// blacklisted: signal_get_invocation_hint (function)
// blacklisted: signal_handler_block (function)
// blacklisted: signal_handler_disconnect (function)
// blacklisted: signal_handler_find (function)
// blacklisted: signal_handler_is_connected (function)
// blacklisted: signal_handler_unblock (function)
// blacklisted: signal_handlers_block_matched (function)
// blacklisted: signal_handlers_destroy (function)
// blacklisted: signal_handlers_disconnect_matched (function)
// blacklisted: signal_handlers_unblock_matched (function)
// blacklisted: signal_has_handler_pending (function)
// blacklisted: signal_is_valid_name (function)
// blacklisted: signal_list_ids (function)
// blacklisted: signal_lookup (function)
// blacklisted: signal_name (function)
// blacklisted: signal_override_class_closure (function)
// blacklisted: signal_parse_name (function)
// blacklisted: signal_query (function)
// blacklisted: signal_remove_emission_hook (function)
// blacklisted: signal_set_va_marshaller (function)
// blacklisted: signal_stop_emission (function)
// blacklisted: signal_stop_emission_by_name (function)
// blacklisted: signal_type_cclosure_new (function)
// blacklisted: source_set_closure (function)
// blacklisted: source_set_dummy_callback (function)
// blacklisted: strdup_value_contents (function)
// blacklisted: type_add_class_private (function)
// blacklisted: type_add_instance_private (function)
// blacklisted: type_add_interface_dynamic (function)
// blacklisted: type_add_interface_static (function)
// blacklisted: type_check_class_is_a (function)
// blacklisted: type_check_instance (function)
// blacklisted: type_check_instance_is_a (function)
// blacklisted: type_check_instance_is_fundamentally_a (function)
// blacklisted: type_check_is_value_type (function)
// blacklisted: type_check_value (function)
// blacklisted: type_check_value_holds (function)
// blacklisted: type_children (function)
// blacklisted: type_class_adjust_private_offset (function)
// blacklisted: type_class_peek (function)
// blacklisted: type_class_peek_static (function)
// blacklisted: type_class_ref (function)
// blacklisted: type_default_interface_peek (function)
// blacklisted: type_default_interface_ref (function)
// blacklisted: type_default_interface_unref (function)
// blacklisted: type_depth (function)
// blacklisted: type_ensure (function)
// blacklisted: type_free_instance (function)
// blacklisted: type_from_name (function)
// blacklisted: type_fundamental (function)
// blacklisted: type_fundamental_next (function)
// blacklisted: type_get_instance_count (function)
// blacklisted: type_get_plugin (function)
// blacklisted: type_get_qdata (function)
// blacklisted: type_get_type_registration_serial (function)
// blacklisted: type_init (function)
// blacklisted: type_init_with_debug_flags (function)
// blacklisted: type_interface_add_prerequisite (function)
// blacklisted: type_interface_get_plugin (function)
// blacklisted: type_interface_instantiatable_prerequisite (function)
// blacklisted: type_interface_peek (function)
// blacklisted: type_interface_prerequisites (function)
// blacklisted: type_interfaces (function)
// blacklisted: type_is_a (function)
// blacklisted: type_name (function)
// blacklisted: type_name_from_class (function)
// blacklisted: type_name_from_instance (function)
// blacklisted: type_next_base (function)
// blacklisted: type_parent (function)
// blacklisted: type_qname (function)
// blacklisted: type_query (function)
// blacklisted: type_register_dynamic (function)
// blacklisted: type_register_fundamental (function)
// blacklisted: type_register_static (function)
// blacklisted: type_set_qdata (function)
// blacklisted: type_test_flags (function)
// blacklisted: value_type_compatible (function)
// blacklisted: value_type_transformable (function)


//--------------------------------------------------------------
// NilString
//--------------------------------------------------------------

// its value will stay the same forever, use the value directly if you like
const NilString = "\x00"

//--------------------------------------------------------------
// Quark
//
// TODO: probably it's a temporary place for this, quarks are
// from glib
//--------------------------------------------------------------

type Quark uint32

func NewQuarkFromString(s string) Quark {
	cs := C.CString(s)
	quark := C.g_quark_from_string(cs)
	C.free(unsafe.Pointer(cs))
	return Quark(quark)
}

// we use this one to store Go's representation of the GObject
// as user data in that GObject once it was allocated. For the
// sake of avoiding allocations.
var go_repr Quark

func init() {
	go_repr = NewQuarkFromString("go-representation")
}


//--------------------------------------------------------------
// Object
//--------------------------------------------------------------

func object_finalizer(obj *Object) {
	if FQueue.Push(unsafe.Pointer(obj), object_finalizer2) {
		return
	}
	C.g_object_set_qdata((*C.GObject)(obj.C), C.uint32_t(go_repr), nil)
	C.g_object_unref((*C.GObject)(obj.C))
}

func object_finalizer2(obj_un unsafe.Pointer) {
	obj := (*Object)(obj_un)
	C.g_object_set_qdata((*C.GObject)(obj.C), C.uint32_t(go_repr), nil)
	C.g_object_unref((*C.GObject)(obj.C))
}

func set_object_finalizer(obj *Object) {
	runtime.SetFinalizer(obj, object_finalizer)
}

func ObjectWrap(c unsafe.Pointer, grab bool) unsafe.Pointer {
	if c == nil {
		return nil
	}
	obj := (*Object)(C.g_object_get_qdata((*C.GObject)(c), C.uint32_t(go_repr)))
	if obj != nil {
		return unsafe.Pointer(obj)
	}
	obj = &Object{c}
	if grab {
		C.g_object_ref_sink((*C.GObject)(obj.C))
	}
	set_object_finalizer(obj)
	C.g_object_set_qdata((*C.GObject)(obj.C),
		C.uint32_t(go_repr), unsafe.Pointer(obj))
	return unsafe.Pointer(obj)
}

func ObjectGrabIfType(c unsafe.Pointer, t Type) unsafe.Pointer {
	if c == nil {
		return nil
	}
	hasrepr := true
	obj := (*Object)(C.g_object_get_qdata((*C.GObject)(c), C.uint32_t(go_repr)))
	if obj == nil {
		obj = &Object{c}
		hasrepr = false
	}
	if obj.GetType().IsA(t) {
		if !hasrepr {
			C.g_object_ref_sink((*C.GObject)(obj.C))
			set_object_finalizer(obj)
			C.g_object_set_qdata((*C.GObject)(obj.C),
				C.uint32_t(go_repr), unsafe.Pointer(obj))
		}
		return unsafe.Pointer(obj)
	}
	return nil
}

func (this *Object) GetType() Type {
	return Type(C._g_object_type((*C.GObject)(this.C)))
}

func (this *Object) Connect(signal string, clo interface{}) {
	csignal := C.CString(signal)
	Holder.Grab(clo)
	if clo == nil {
		panic("Connect with nil")
	}
	goclosure := C.g_goclosure_new(toGoInterfaceHolder(clo), toGoInterfaceHolder(nil))
	C.g_signal_connect_closure((*C.GObject)(this.C), csignal, (*C.GClosure)(unsafe.Pointer(goclosure)), 0)
	C.free(unsafe.Pointer(csignal))
}

func (this *Object) ConnectMethod(signal string, clo interface{}, recv interface{}) {
	csignal := C.CString(signal)
	Holder.Grab(clo)
	Holder.Grab(recv)
	goclosure := C.g_goclosure_new(toGoInterfaceHolder(clo), toGoInterfaceHolder(recv))
	C.g_signal_connect_closure((*C.GObject)(this.C), csignal, (*C.GClosure)(unsafe.Pointer(goclosure)), 0)
	C.free(unsafe.Pointer(csignal))
}

func (this *Object) FindProperty(name string) *ParamSpec {
	cname := C.CString(name)
	ret := C._g_object_find_property(this.InheritedFromGObject(), cname)
	C.free(unsafe.Pointer(cname))
	return (*ParamSpec)(ParamSpecWrap(unsafe.Pointer(ret), true))
}

func (this *Object) SetProperty(name string, value interface{}) {
	cname := C.CString(name)
	pspec := this.FindProperty(name)
	if pspec == nil {
		panic("Object has no property with that name: " + name)
	}
	var gvalue Value
	gvalue.Init(pspec.GetValueType())
	gvalue.SetGoInterface(value)
	C.g_object_set_property(this.InheritedFromGObject(), cname,
		(*C.GValue)(unsafe.Pointer(&gvalue)))
	gvalue.Unset()
	C.free(unsafe.Pointer(cname))
}

func (this *Object) GetProperty(name string, value interface{}) {
	cname := C.CString(name)
	pspec := this.FindProperty(name)
	if pspec == nil {
		panic("Object has no property with that name: " + name)
	}
	var gvalue Value
	gvalue.Init(pspec.GetValueType())
	C.g_object_get_property(this.InheritedFromGObject(), cname,
		(*C.GValue)(unsafe.Pointer(&gvalue)))
	gvalue.GetGoInterface(value)
	gvalue.Unset()
	C.free(unsafe.Pointer(cname))
}

func ObjectBindProperty(source ObjectLike, source_property string, target ObjectLike, target_property string, flags BindingFlags) *Binding {
	csource_property := C.CString(source_property)
	ctarget_property := C.CString(target_property)
	obj := C.g_object_bind_property(
		source.InheritedFromGObject(), csource_property,
		target.InheritedFromGObject(), ctarget_property,
		C.GBindingFlags(flags))
	C.free(unsafe.Pointer(csource_property))
	C.free(unsafe.Pointer(ctarget_property))
	return (*Binding)(ObjectWrap(unsafe.Pointer(obj), true))
}

func (this *Object) Unref() {
	runtime.SetFinalizer(this, nil)
	C.g_object_set_qdata((*C.GObject)(this.C), C.uint32_t(go_repr), nil)
	C.g_object_unref((*C.GObject)(this.C))
	this.C = nil
}

//--------------------------------------------------------------
// Closures
//--------------------------------------------------------------

//export g_goclosure_finalize_go
func g_goclosure_finalize_go(goclosure_up unsafe.Pointer) {
	goclosure := (*C.GGoClosure)(goclosure_up)
	clo := fromGoInterfaceHolder(C.g_goclosure_get_func(goclosure))
	recv := fromGoInterfaceHolder(C.g_goclosure_get_recv(goclosure))
	Holder.Release(clo)
	Holder.Release(recv)
}

//export g_goclosure_marshal_go
func g_goclosure_marshal_go(goclosure_up, ret_up unsafe.Pointer, nargs int32, args_up unsafe.Pointer) {
	var callargs [20]reflect.Value
	var recv reflect.Value
	goclosure := (*C.GGoClosure)(goclosure_up)
	ret := (*Value)(ret_up)
	args := (*(*[alot]Value)(args_up))[:nargs]
	f := reflect.ValueOf(fromGoInterfaceHolder(C.g_goclosure_get_func(goclosure)))
	ft := f.Type()
	callargsn := ft.NumIn()

	recvi := fromGoInterfaceHolder(C.g_goclosure_get_recv(goclosure))
	if recvi != nil {
		recv = reflect.ValueOf(recvi)
	}

	if callargsn >= 20 {
		panic("too many arguments in a closure")
	}

	for i, n := 0, callargsn; i < n; i++ {
		idx := i
		if recvi != nil {
			idx--
			if i == 0 {
				callargs[i] = recv
				continue
			}
		}

		in := ft.In(i)

		// use default value, if there is not enough args
		if len(args) <= idx {
			callargs[i] = reflect.New(in).Elem()
			continue
		}

		v := args[idx].GetGoValue(in)
		callargs[i] = v
	}


	// The running thread of f need be the owner of global
	// default main content.
	// And g_goclosure_marshal_go is called by main loop,
	// So we can simply make this by using runtime.LockOSThread
	runtime.LockOSThread()
	out := f.Call(callargs[:callargsn])
	runtime.UnlockOSThread()
	if len(out) == 1 {
		ret.SetGoValue(out[0])
	}
}

//--------------------------------------------------------------
// Go Interface boxed type
//--------------------------------------------------------------

//export g_go_interface_copy_go
func g_go_interface_copy_go(boxed unsafe.Pointer) unsafe.Pointer {
	Holder.Grab(*(*interface{})(boxed))
	newboxed := C.malloc(C.size_t(unsafe.Sizeof([2]unsafe.Pointer{})))
	C.memcpy(newboxed, boxed, C.size_t(unsafe.Sizeof([2]unsafe.Pointer{})))
	return newboxed
}

//export g_go_interface_free_go
func g_go_interface_free_go(boxed unsafe.Pointer) {
	Holder.Release(*(*interface{})(boxed))
	C.free(boxed)
}

//--------------------------------------------------------------
// Type
//--------------------------------------------------------------

type Type C.GType

func (this Type) IsA(other Type) bool {
	return C.g_type_is_a(C.GType(this), C.GType(other)) != 0
}

func (this Type) String() string {
	cname := C.g_type_name(C.GType(this))
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

func (this Type) asC() C.GType {
	return C.GType(this)
}

var (
	Interface Type
	Char Type
	UChar Type
	Boolean Type
	Int Type
	UInt Type
	Long Type
	ULong Type
	Int64 Type
	UInt64 Type
	Enum Type
	Flags Type
	Float Type
	Double Type
	String Type
	Pointer Type
	Boxed Type
	Param Type
	GObject Type
	GType Type
	Variant Type
	GoInterface Type
)

func init() {
	C.g_type_init()

	Interface = Type(C._g_type_interface())
	Char = Type(C._g_type_char())
	UChar = Type(C._g_type_uchar())
	Boolean = Type(C._g_type_boolean())
	Int = Type(C._g_type_int())
	UInt = Type(C._g_type_uint())
	Long = Type(C._g_type_long())
	ULong = Type(C._g_type_ulong())
	Int64 = Type(C._g_type_int64())
	UInt64 = Type(C._g_type_uint64())
	Enum = Type(C._g_type_enum())
	Flags = Type(C._g_type_flags())
	Float = Type(C._g_type_float())
	Double = Type(C._g_type_double())
	String = Type(C._g_type_string())
	Pointer = Type(C._g_type_pointer())
	Boxed = Type(C._g_type_boxed())
	Param = Type(C._g_type_param())
	GObject = Type(C._g_type_object())
	GType = Type(C._g_type_gtype())
	Variant = Type(C._g_type_variant())
	GoInterface = Type(C._g_type_go_interface())
}

// Every GObject generated by this generator implements this interface
// and it must work even if the receiver is a nil value
type StaticTyper interface {
	GetStaticType() Type
}

//--------------------------------------------------------------
// Value
//--------------------------------------------------------------

func (this *Value) asC() *C.GValue {
	return (*C.GValue)(unsafe.Pointer(this))
}

// g_value_init
func (this *Value) Init(t Type) {
	C.g_value_init(this.asC(), t.asC())
}

// g_value_copy
func (this *Value) Set(src *Value) {
	C.g_value_copy(src.asC(), this.asC())
}

// g_value_reset
func (this *Value) Reset() {
	C.g_value_reset(this.asC())
}

// g_value_unset
func (this *Value) Unset() {
	C.g_value_unset(this.asC())
}

// G_VALUE_TYPE
func (this *Value) GetType() Type {
	return Type(C._g_value_type(this.asC()))
}

// g_value_type_compatible
func ValueTypeCompatible(src, dst Type) bool {
	return C.g_value_type_compatible(src.asC(), dst.asC()) != 0
}

// g_value_type_transformable
func ValueTypeTransformable(src, dst Type) bool {
	return C.g_value_type_transformable(src.asC(), dst.asC()) != 0
}

// g_value_transform
func (this *Value) Transform(src *Value) bool {
	return C.g_value_transform(src.asC(), this.asC()) != 0
}

// g_value_get_boolean
func (this *Value) GetBool() bool {
	return C.g_value_get_boolean(this.asC()) != 0
}

// g_value_set_boolean
func (this *Value) SetBool(v bool) {
	C.g_value_set_boolean(this.asC(), _GoBoolToCBool(v))
}

// g_value_get_int64
func (this *Value) GetInt() int64 {
	return int64(C.g_value_get_int64(this.asC()))
}

// g_value_set_int64
func (this *Value) SetInt(v int64) {
	C.g_value_set_int64(this.asC(), C.int64_t(v))
}

// g_value_get_uint64
func (this *Value) GetUint() uint64 {
	return uint64(C.g_value_get_uint64(this.asC()))
}

// g_value_set_uint64
func (this *Value) SetUint(v uint64) {
	C.g_value_set_uint64(this.asC(), C.uint64_t(v))
}

// g_value_get_double
func (this *Value) GetFloat() float64 {
	return float64(C.g_value_get_double(this.asC()))
}

// g_value_set_double
func (this *Value) SetFloat(v float64) {
	C.g_value_set_double(this.asC(), C.double(v))
}

// g_value_get_string
func (this *Value) GetString() string {
	return C.GoString(C.g_value_get_string(this.asC()))
}

// g_value_take_string
func (this *Value) SetString(v string) {
	cstr := C.CString(v)
	C.g_value_take_string(this.asC(), cstr)
	// not freeing, because GValue takes the ownership
}

// g_value_get_object
func (this *Value) GetObject() unsafe.Pointer {
	return unsafe.Pointer(C.g_value_get_object(this.asC()))
}

// g_value_set_object
func (this *Value) SetObject(x unsafe.Pointer) {
	C.g_value_set_object(this.asC(), (*C.GObject)(x))
}

// g_value_get_boxed
func (this *Value) GetBoxed() unsafe.Pointer {
	return C.g_value_get_boxed(this.asC())
}

// g_value_take_boxed
func (this *Value) SetBoxed(x unsafe.Pointer) {
	C.g_value_take_boxed(this.asC(), x)
}

func (this *Value) GetBoxedInterface() interface{} {
	return *(*interface{})(C.g_value_get_boxed(this.asC()))
}

func (this *Value) SetBoxedInterface(x interface{}) {
	Holder.Grab(x)
	newboxed := C.malloc(C.size_t(unsafe.Sizeof([2]unsafe.Pointer{})))
	C.memcpy(newboxed, unsafe.Pointer(&x), C.size_t(unsafe.Sizeof([2]unsafe.Pointer{})))
	C.g_value_take_boxed(this.asC(), newboxed)
}

//--------------------------------------------------------------
// A giant glue for connecting GType and Go's reflection
//--------------------------------------------------------------

var statictyper = reflect.TypeOf((*StaticTyper)(nil)).Elem()
var objectlike = reflect.TypeOf((*ObjectLike)(nil)).Elem()

func (this *Value) SetGoValue(v reflect.Value) {
	valuetype := this.GetType()
	var src Value

	if valuetype == GoInterface {
		// special case
		this.SetBoxedInterface(v.Interface())
		return
	}

	transform := func() {
		ok := this.Transform(&src)
		if !ok {
			panic("Go value (" + v.Type().String() + ") is not transformable to " + valuetype.String())
		}
	}

	switch v.Kind() {
	case reflect.Bool:
		src.Init(Boolean)
		src.SetBool(v.Bool())
		transform()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		src.Init(Int64)
		src.SetInt(v.Int())
		transform()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		src.Init(UInt64)
		src.SetUint(v.Uint())
		transform()
	case reflect.Float32, reflect.Float64:
		src.Init(Double)
		src.SetFloat(v.Float())
		transform()
	case reflect.String:
		src.Init(String)
		src.SetString(v.String())
		transform()
		src.Unset()
	case reflect.Ptr:
		gotype := v.Type()
		src.Init(GObject)
		if gotype.Implements(objectlike) {
			obj, ok := v.Interface().(ObjectLike)
			if !ok {
				panic(gotype.String() + " is not transformable to GValue")
			}

			src.SetObject(unsafe.Pointer(obj.InheritedFromGObject()))
			transform()
		}
		src.Unset()
	}
}

var CairoMarshaler func(*Value, reflect.Type) (reflect.Value, bool)

func (this *Value) GetGoValue(t reflect.Type) reflect.Value {
	var out reflect.Value
	var dst Value

	if (this.GetType() == GoInterface) {
		return reflect.ValueOf(this.GetBoxedInterface())
	}

	transform := func() {
		ok := dst.Transform(this)
		if !ok {
			panic("GValue is not transformable to " + t.String())
		}
	}

	switch t.Kind() {
	case reflect.Bool:
		dst.Init(Boolean)
		transform()
		out = reflect.New(t).Elem()
		out.SetBool(dst.GetBool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dst.Init(Int64)
		transform()
		out = reflect.New(t).Elem()
		out.SetInt(dst.GetInt())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dst.Init(UInt64)
		transform()
		out = reflect.New(t).Elem()
		out.SetUint(dst.GetUint())
	case reflect.Float32, reflect.Float64:
		dst.Init(Double)
		transform()
		out = reflect.New(t).Elem()
		out.SetFloat(dst.GetFloat())
	case reflect.String:
		dst.Init(String)
		transform()
		out = reflect.New(t).Elem()
		out.SetString(dst.GetString())
		dst.Unset() // need to clean up in this case
	case reflect.Ptr:
		if t.Implements(objectlike) {
			// at this point we're sure that this is a pointer to the ObjectLike
			out = reflect.New(t)
			st, ok := out.Elem().Interface().(StaticTyper)
			if !ok {
				panic("ObjectLike type must implement StaticTyper as well")
			}
			dst.Init(st.GetStaticType())
			transform()
			*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = ObjectWrap(dst.GetObject(), true)
			dst.Unset()
			out = out.Elem()
		} else {
			// cairo marshaler hook
			if CairoMarshaler != nil {
				var ok bool
				out, ok = CairoMarshaler(this, t)
				if ok {
					break
				}
			}

			// must be a struct then
			out = reflect.New(t)
			*(*unsafe.Pointer)(unsafe.Pointer(out.Pointer())) = this.GetBoxed()
			out = out.Elem()
		}
	}
	return out
}

func (this *Value) SetGoInterface(v interface{}) {
	this.SetGoValue(reflect.ValueOf(v))
}

func (this *Value) GetGoInterface(v interface{}) {
	vp := reflect.ValueOf(v)
	if vp.Kind() != reflect.Ptr {
		panic("a pointer to value is expected for Value.GetGoInterface")
	}
	vp.Elem().Set(this.GetGoValue(vp.Type().Elem()))
}