/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package gi

/*
#include <stdlib.h>
#include <girepository.h>

static inline void free_string(char *p) { free(p); }
static inline void free_gstring(gchar *p) { if (p) g_free(p); }
static inline char *gpointer_to_charp(gpointer p) { return p; }
static inline gchar **next_gcharptr(gchar **s) { return s+1; }
#cgo pkg-config: gobject-introspection-1.0
*/
import "C"
import (
	"errors"
	"strings"
	"runtime"
	"unsafe"
	"fmt"
)

// utils

// Convert GSList containing strings to []string
func _GStringGSListToGoStringSlice(list *C.GSList) []string {
	var slice []string
	for list != nil {
		str := C.GoString(C.gpointer_to_charp(list.data))
		slice = append(slice, str)
		list = list.next
	}
	return slice
}

// Convert gchar** null-terminated glib string array to []string, frees "arr"
func _GStringArrayToGoStringSlice(arr **C.gchar) []string {
	var slice []string
	if arr == nil {
		return slice
	}

	iter := arr
	for *iter != nil {
		slice = append(slice, _GStringToGoString(*iter))
		iter = C.next_gcharptr(iter)
	}
	C.g_strfreev(arr)
	return slice
}

// Go string to glib C string, "" == NULL
func _GoStringToGString(s string) *C.gchar {
	if s == "" {
		return nil
	}
	return (*C.gchar)(unsafe.Pointer(C.CString(s)))
}

// glib C string to Go string, NULL == ""
func _GStringToGoString(s *C.gchar) string {
	if s == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

// C string to Go string, NULL == ""
func _CStringToGoString(s *C.char) string {
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

// GError to os.Error, frees "err"
func _GErrorToOSError(err *C.GError) (goerr error) {
	goerr = errors.New(_GStringToGoString(err.message))
	C.g_error_free(err)
	return
}

// Check for type
func _ExpectBaseInfoType(bil BaseInfoLike, types ...InfoType) {
	bi := bil.inheritedFromBaseInfo()
	for _, t := range types {
		if bi.Type() == t {
			return
		}
	}

	// error from here
	typeStrings := make([]string, len(types))
	for i, t := range types {
		typeStrings[i] = t.String()
	}
	panic(fmt.Sprintf("Type mismatch, expected: %s, got: %s",
		strings.Join(typeStrings, " or "), bi.Type()))
}

// Finalizer for BaseInfo structure, does the unref
func _BaseInfoFinalizer(bi *BaseInfo) {
	bi.Unref()
}

// Helper for initializing finalizer on BaseInfo
func _SetBaseInfoFinalizer(bi *BaseInfo) *BaseInfo {
	runtime.SetFinalizer(bi, _BaseInfoFinalizer)
	return bi
}

//------------------------------------------------------------------------------
// .types
//------------------------------------------------------------------------------

type BaseInfoLike interface {
	inheritedFromBaseInfo() *BaseInfo
}

func ToBaseInfo(bil BaseInfoLike) *BaseInfo {
	return bil.inheritedFromBaseInfo()
}

func ToArgInfo(bil BaseInfoLike) *ArgInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_ARG)
	return (*ArgInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToConstantInfo(bil BaseInfoLike) *ConstantInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_CONSTANT)
	return (*ConstantInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToFieldInfo(bil BaseInfoLike) *FieldInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_FIELD)
	return (*FieldInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToPropertyInfo(bil BaseInfoLike) *PropertyInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_PROPERTY)
	return (*PropertyInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToTypeInfo(bil BaseInfoLike) *TypeInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_TYPE)
	return (*TypeInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToCallableInfo(bil BaseInfoLike) *CallableInfo {
	_ExpectBaseInfoType(bil,
		INFO_TYPE_FUNCTION,
		INFO_TYPE_CALLBACK,
		INFO_TYPE_SIGNAL,
		INFO_TYPE_VFUNC)
	return (*CallableInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToFunctionInfo(bil BaseInfoLike) *FunctionInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_FUNCTION)
	return (*FunctionInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToSignalInfo(bil BaseInfoLike) *SignalInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_SIGNAL)
	return (*SignalInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToVFuncInfo(bil BaseInfoLike) *VFuncInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_VFUNC)
	return (*VFuncInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToRegisteredType(bil BaseInfoLike) *RegisteredType {
	_ExpectBaseInfoType(bil,
		INFO_TYPE_BOXED,
		INFO_TYPE_ENUM,
		INFO_TYPE_FLAGS,
		INFO_TYPE_INTERFACE,
		INFO_TYPE_OBJECT,
		INFO_TYPE_STRUCT,
		INFO_TYPE_UNION)
	return (*RegisteredType)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToEnumInfo(bil BaseInfoLike) *EnumInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_ENUM, INFO_TYPE_FLAGS)
	return (*EnumInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToInterfaceInfo(bil BaseInfoLike) *InterfaceInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_INTERFACE)
	return (*InterfaceInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToObjectInfo(bil BaseInfoLike) *ObjectInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_OBJECT)
	return (*ObjectInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToStructInfo(bil BaseInfoLike) *StructInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_STRUCT)
	return (*StructInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

func ToUnionInfo(bil BaseInfoLike) *UnionInfo {
	_ExpectBaseInfoType(bil, INFO_TYPE_UNION)
	return (*UnionInfo)(unsafe.Pointer(bil.inheritedFromBaseInfo()))
}

//------------------------------------------------------------------------------
// InfoType
//------------------------------------------------------------------------------

type InfoType int

const (
	INFO_TYPE_INVALID    InfoType = C.GI_INFO_TYPE_INVALID
	INFO_TYPE_FUNCTION   InfoType = C.GI_INFO_TYPE_FUNCTION
	INFO_TYPE_CALLBACK   InfoType = C.GI_INFO_TYPE_CALLBACK
	INFO_TYPE_STRUCT     InfoType = C.GI_INFO_TYPE_STRUCT
	INFO_TYPE_BOXED      InfoType = C.GI_INFO_TYPE_BOXED
	INFO_TYPE_ENUM       InfoType = C.GI_INFO_TYPE_ENUM
	INFO_TYPE_FLAGS      InfoType = C.GI_INFO_TYPE_FLAGS
	INFO_TYPE_OBJECT     InfoType = C.GI_INFO_TYPE_OBJECT
	INFO_TYPE_INTERFACE  InfoType = C.GI_INFO_TYPE_INTERFACE
	INFO_TYPE_CONSTANT   InfoType = C.GI_INFO_TYPE_CONSTANT
	INFO_TYPE_INVALID_0  InfoType = C.GI_INFO_TYPE_INVALID_0
	INFO_TYPE_UNION      InfoType = C.GI_INFO_TYPE_UNION
	INFO_TYPE_VALUE      InfoType = C.GI_INFO_TYPE_VALUE
	INFO_TYPE_SIGNAL     InfoType = C.GI_INFO_TYPE_SIGNAL
	INFO_TYPE_VFUNC      InfoType = C.GI_INFO_TYPE_VFUNC
	INFO_TYPE_PROPERTY   InfoType = C.GI_INFO_TYPE_PROPERTY
	INFO_TYPE_FIELD      InfoType = C.GI_INFO_TYPE_FIELD
	INFO_TYPE_ARG        InfoType = C.GI_INFO_TYPE_ARG
	INFO_TYPE_TYPE       InfoType = C.GI_INFO_TYPE_TYPE
	INFO_TYPE_UNRESOLVED InfoType = C.GI_INFO_TYPE_UNRESOLVED
)

// g_info_type_to_string
func (it InfoType) String() string {
	return _GStringToGoString(C.g_info_type_to_string(C.GIInfoType(it)))
}

//------------------------------------------------------------------------------
// Repository
//------------------------------------------------------------------------------

type Repository struct {
	c *C.GIRepository
}

type RepositoryLoadFlags int

const (
	REPOSITORY_LOAD_FLAG_LAZY RepositoryLoadFlags = C.G_IREPOSITORY_LOAD_FLAG_LAZY
)

// g_irepository_get_default
func DefaultRepository() *Repository {
	ret := C.g_irepository_get_default()
	if ret == nil {
		return nil
	}
	return &Repository{ret}
}

// g_irepository_prepend_search_path
func PreprendRepositorySearchPath(path string) {
	cpath := C.CString(path)
	C.g_irepository_prepend_search_path(cpath)
	C.free_string(cpath)
}

// g_irepository_get_search_path
func RepositorySearchPath() []string {
	return _GStringGSListToGoStringSlice(C.g_irepository_get_search_path())
}

//const char *        g_irepository_load_typelib          (GIRepository *repository,
//                                                         GITypelib *typelib,
//                                                         GIRepositoryLoadFlags flags,
//                                                         GError **error);

// g_irepository_is_registered
func (r *Repository) IsRegistered(namespace, version string) bool {
	gnamespace := _GoStringToGString(namespace)
	gversion := _GoStringToGString(version)
	ret := C.g_irepository_is_registered(r.c, gnamespace, gversion)
	C.free_gstring(gversion)
	C.free_gstring(gnamespace)
	return ret != 0
}

// g_irepository_find_by_name
func (r *Repository) FindByName(namespace, name string) *BaseInfo {
	gnamespace := _GoStringToGString(namespace)
	gname := _GoStringToGString(name)
	ret := C.g_irepository_find_by_name(r.c, gnamespace, gname)
	C.free_gstring(gname)
	C.free_gstring(gnamespace)
	return _SetBaseInfoFinalizer(&BaseInfo{ret})
}

// g_irepository_require
func (r *Repository) Require(namespace, version string, flags RepositoryLoadFlags) (*Typelib, error) {
	var err *C.GError
	gnamespace := _GoStringToGString(namespace)
	gversion := _GoStringToGString(version)
	tl := C.g_irepository_require(r.c, gnamespace, gversion, C.GIRepositoryLoadFlags(flags), &err)
	C.free_gstring(gversion)
	C.free_gstring(gnamespace)

	if err != nil {
		return nil, _GErrorToOSError(err)
	}

	var tlwrap *Typelib
	if tl != nil {
		tlwrap = &Typelib{tl}
	}

	return tlwrap, nil
}

//GITypelib *         g_irepository_require_private       (GIRepository *repository,
//                                                         const gchar *typelib_dir,
//                                                         const gchar *namespace_,
//                                                         const gchar *version,
//                                                         GIRepositoryLoadFlags flags,
//                                                         GError **error);

// g_irepository_get_dependencies
func (r *Repository) Dependencies(namespace string) []string {
	gnamespace := _GoStringToGString(namespace)
	arr := C.g_irepository_get_dependencies(r.c, gnamespace)
	C.free_gstring(gnamespace)
	return _GStringArrayToGoStringSlice(arr)
}

// g_irepository_get_loaded_namespaces
func (r *Repository) LoadedNamespaces() []string {
	arr := C.g_irepository_get_loaded_namespaces(r.c)
	return _GStringArrayToGoStringSlice(arr)
}

//GIBaseInfo *        g_irepository_find_by_gtype         (GIRepository *repository,
//                                                         GType gtype);

// g_irepository_get_n_infos
func (r *Repository) NumInfo(namespace string) int {
	gnamespace := _GoStringToGString(namespace)
	num := C.g_irepository_get_n_infos(r.c, gnamespace)
	C.free_gstring(gnamespace)
	return int(num)
}

// g_irepository_get_info
func (r *Repository) Info(namespace string, index int) *BaseInfo {
	gnamespace := _GoStringToGString(namespace)
	info := C.g_irepository_get_info(r.c, gnamespace, C.gint(index))
	C.free_gstring(gnamespace)
	return _SetBaseInfoFinalizer(&BaseInfo{info})
}

// g_irepository_get_typelib_path
func (r *Repository) TypelibPath(namespace string) string {
	gnamespace := _GoStringToGString(namespace)
	path := C.g_irepository_get_typelib_path(r.c, gnamespace)
	C.free_gstring(gnamespace)
	return _GStringToGoString(path)
}

// g_irepository_get_shared_library
func (r *Repository) SharedLibrary(namespace string) string {
	gnamespace := _GoStringToGString(namespace)
	shlib := C.g_irepository_get_shared_library(r.c, gnamespace)
	C.free_gstring(gnamespace)
	return _GStringToGoString(shlib)
}

// g_irepository_get_version
func (r *Repository) Version(namespace string) string {
	gnamespace := _GoStringToGString(namespace)
	ver := C.g_irepository_get_version(r.c, gnamespace)
	C.free_gstring(gnamespace)
	return _GStringToGoString(ver)
}

//GOptionGroup *      g_irepository_get_option_group      (void);

// g_irepository_get_c_prefix
func (r *Repository) CPrefix(namespace string) string {
	gnamespace := _GoStringToGString(namespace)
	prefix := C.g_irepository_get_c_prefix(r.c, gnamespace)
	C.free_gstring(gnamespace)
	return _GStringToGoString(prefix)
}

//gboolean            g_irepository_dump                  (const char *arg,
//                                                         GError **error);
//GList *             g_irepository_enumerate_versions    (GIRepository *repository,
//                                                         const gchar *namespace_);

//------------------------------------------------------------------------------
// Typelib
//------------------------------------------------------------------------------

type Typelib struct {
	c *C.GITypelib
}

//GITypelib *         g_typelib_new_from_memory           (guint8 *memory,
//                                                         gsize len,
//                                                         GError **error);
//GITypelib *         g_typelib_new_from_const_memory     (const guint8 *memory,
//                                                         gsize len,
//                                                         GError **error);
//GITypelib *         g_typelib_new_from_mapped_file      (GMappedFile *mfile,
//                                                         GError **error);
//void                g_typelib_free                      (GITypelib *typelib);
//gboolean            g_typelib_symbol                    (GITypelib *typelib,
//                                                         const gchar *symbol_name,
//                                                         gpointer *symbol);
//const gchar *       g_typelib_get_namespace             (GITypelib *typelib);

//------------------------------------------------------------------------------
// BaseInfo
//------------------------------------------------------------------------------

type BaseInfo struct {
	c *C.GIBaseInfo
}

func (bi *BaseInfo) inheritedFromBaseInfo() *BaseInfo {
	return bi
}

// g_base_info_ref
func (bi *BaseInfo) Ref() *BaseInfo {
	C.g_base_info_ref(bi.c)
	return bi
}

// g_base_info_unref
func (bi *BaseInfo) Unref() {
	C.g_base_info_unref(bi.c)
}

// g_base_info_get_type
func (bi *BaseInfo) Type() InfoType {
	return InfoType(C.g_base_info_get_type(bi.c))
}

// g_base_info_get_name
func (bi *BaseInfo) Name() string {
	return _GStringToGoString(C.g_base_info_get_name(bi.c))
}

// g_base_info_get_namespace
func (bi *BaseInfo) Namespace() string {
	return _GStringToGoString(C.g_base_info_get_namespace(bi.c))
}

// g_base_info_is_deprecated
func (bi *BaseInfo) IsDeprecated() bool {
	return C.g_base_info_is_deprecated(bi.c) != 0
}

// g_base_info_get_attribute
func (bi *BaseInfo) Attribute(name string) string {
	gname := _GoStringToGString(name)
	ret := _GStringToGoString(C.g_base_info_get_attribute(bi.c, gname))
	C.free_gstring(gname)
	return ret
}

// g_base_info_iterate_attributes
func (bi *BaseInfo) IterateAttributes(cb func(name, value string)) {
	var iter C.GIAttributeIter
	var cname, cvalue *C.char
	for C.g_base_info_iterate_attributes(bi.c, &iter, &cname, &cvalue) != 0 {
		name, value := C.GoString(cname), C.GoString(cvalue)
		cb(name, value)
	}
}

// g_base_info_get_container
func (bi *BaseInfo) Container() *BaseInfo {
	c := C.g_base_info_get_container(bi.c)
	if c == nil {
		return nil
	}
	return &BaseInfo{c}
}

// g_base_info_get_typelib
func (bi *BaseInfo) Typelib() *Typelib {
	return &Typelib{C.g_base_info_get_typelib(bi.c)}
}

//gboolean            g_base_info_equal                   (GIBaseInfo *info1,
//                                                         GIBaseInfo *info2);

//------------------------------------------------------------------------------
// ArgInfo
//------------------------------------------------------------------------------

type ArgInfo struct {
	BaseInfo
}

type Direction int
type ScopeType int
type Transfer int

const (
	DIRECTION_IN    Direction = C.GI_DIRECTION_IN
	DIRECTION_OUT   Direction = C.GI_DIRECTION_OUT
	DIRECTION_INOUT Direction = C.GI_DIRECTION_INOUT
)

const (
	SCOPE_TYPE_INVALID  ScopeType = C.GI_SCOPE_TYPE_INVALID
	SCOPE_TYPE_CALL     ScopeType = C.GI_SCOPE_TYPE_CALL
	SCOPE_TYPE_ASYNC    ScopeType = C.GI_SCOPE_TYPE_ASYNC
	SCOPE_TYPE_NOTIFIED ScopeType = C.GI_SCOPE_TYPE_NOTIFIED
)

const (
	TRANSFER_NOTHING    Transfer = C.GI_TRANSFER_NOTHING
	TRANSFER_CONTAINER  Transfer = C.GI_TRANSFER_CONTAINER
	TRANSFER_EVERYTHING Transfer = C.GI_TRANSFER_EVERYTHING
)

// g_arg_info_get_direction
func (ai *ArgInfo) Direction() Direction {
	return Direction(C.g_arg_info_get_direction((*C.GIArgInfo)(ai.c)))
}

// g_arg_info_is_caller_allocates
func (ai *ArgInfo) IsCallerAllocates() bool {
	return C.g_arg_info_is_caller_allocates((*C.GIArgInfo)(ai.c)) != 0
}

// g_arg_info_is_return_value
func (ai *ArgInfo) IsReturnValue() bool {
	return C.g_arg_info_is_return_value((*C.GIArgInfo)(ai.c)) != 0
}

// g_arg_info_is_optional
func (ai *ArgInfo) IsOptional() bool {
	return C.g_arg_info_is_optional((*C.GIArgInfo)(ai.c)) != 0
}

// g_arg_info_may_be_null
func (ai *ArgInfo) MayBeNil() bool {
	return C.g_arg_info_may_be_null((*C.GIArgInfo)(ai.c)) != 0
}

// g_arg_info_is_skip
func (ai *ArgInfo) IsSkip() bool {
	return C.g_arg_info_is_skip((*C.GIArgInfo)(ai.c)) != 0
}

// g_arg_info_get_ownership_transfer
func (ai *ArgInfo) OwnershipTransfer() Transfer {
	return Transfer(C.g_arg_info_get_ownership_transfer((*C.GIArgInfo)(ai.c)))
}

// g_arg_info_get_scope
func (ai *ArgInfo) Scope() ScopeType {
	return ScopeType(C.g_arg_info_get_scope((*C.GIArgInfo)(ai.c)))
}

// g_arg_info_get_closure
func (ai *ArgInfo) Closure() int {
	return int(C.g_arg_info_get_closure((*C.GIArgInfo)(ai.c)))
}

// g_arg_info_get_destroy
func (ai *ArgInfo) Destroy() int {
	return int(C.g_arg_info_get_destroy((*C.GIArgInfo)(ai.c)))
}

// g_arg_info_get_type
func (ai *ArgInfo) Type() *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_arg_info_get_type((*C.GIArgInfo)(ai.c)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//void                g_arg_info_load_type                (GIArgInfo *info,
//                                                         GITypeInfo *type)

//------------------------------------------------------------------------------
// ConstantInfo
//------------------------------------------------------------------------------

type ConstantInfo struct {
	BaseInfo
}

// g_constant_info_get_type
func (ci *ConstantInfo) Type() *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_constant_info_get_type((*C.GIConstantInfo)(ci.c)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_constant_info_get_value
func (ci *ConstantInfo) Value() interface{} {
	var arg C.GIArgument
	C.g_constant_info_get_value((*C.GIConstantInfo)(ci.c), &arg)

	ti := ci.Type()
	switch ti.Tag() {
	case TYPE_TAG_BOOLEAN:
		return *(*C.gboolean)(unsafe.Pointer(&arg)) != 0
	case TYPE_TAG_INT8:
		return *(*int8)(unsafe.Pointer(&arg))
	case TYPE_TAG_UINT8:
		return *(*uint8)(unsafe.Pointer(&arg))
	case TYPE_TAG_INT16:
		return *(*int16)(unsafe.Pointer(&arg))
	case TYPE_TAG_UINT16:
		return *(*uint16)(unsafe.Pointer(&arg))
	case TYPE_TAG_INT32:
		return *(*int32)(unsafe.Pointer(&arg))
	case TYPE_TAG_UINT32:
		return *(*uint32)(unsafe.Pointer(&arg))
	case TYPE_TAG_INT64:
		return *(*int64)(unsafe.Pointer(&arg))
	case TYPE_TAG_UINT64:
		return *(*uint64)(unsafe.Pointer(&arg))
	case TYPE_TAG_FLOAT:
		return *(*float32)(unsafe.Pointer(&arg))
	case TYPE_TAG_DOUBLE:
		return *(*float64)(unsafe.Pointer(&arg))
	case TYPE_TAG_UTF8, TYPE_TAG_FILENAME:
		return C.GoString(*(**C.char)(unsafe.Pointer(&arg)))
	}
	panic("unsupported constant value")
	return nil
}

//gint                g_constant_info_get_value           (GIConstantInfo *info,
//                                                         GIArgument *value);

//------------------------------------------------------------------------------
// FieldInfo
//------------------------------------------------------------------------------

type FieldInfo struct {
	BaseInfo
}

type FieldInfoFlags int

const (
	FIELD_IS_READABLE FieldInfoFlags = C.GI_FIELD_IS_READABLE
	FIELD_IS_WRITABLE FieldInfoFlags = C.GI_FIELD_IS_WRITABLE
)

// g_field_info_get_flags
func (fi *FieldInfo) Flags() FieldInfoFlags {
	return FieldInfoFlags(C.g_field_info_get_flags((*C.GIFieldInfo)(fi.c)))
}

// g_field_info_get_size
func (fi *FieldInfo) Size() int {
	return int(C.g_field_info_get_size((*C.GIFieldInfo)(fi.c)))
}

// g_field_info_get_offset
func (fi *FieldInfo) Offset() int {
	return int(C.g_field_info_get_offset((*C.GIFieldInfo)(fi.c)))
}

// g_field_info_get_type
func (fi *FieldInfo) Type() *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_field_info_get_type((*C.GIFieldInfo)(fi.c)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//gboolean            g_field_info_get_field              (GIFieldInfo *field_info,
//                                                         gpointer mem,
//                                                         GIArgument *value);
//gboolean            g_field_info_set_field              (GIFieldInfo *field_info,
//                                                         gpointer mem,
//                                                         const GIArgument *value);

//------------------------------------------------------------------------------
// PropertyInfo
//------------------------------------------------------------------------------

type PropertyInfo struct {
	BaseInfo
}

//GParamFlags         g_property_info_get_flags           (GIPropertyInfo *info);

// g_property_info_get_type
func (pi *PropertyInfo) Type() *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_property_info_get_type((*C.GIPropertyInfo)(pi.c)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_property_info_get_ownership_transfer
func (pi *PropertyInfo) OwnershipTransfer() Transfer {
	return Transfer(C.g_property_info_get_ownership_transfer((*C.GIPropertyInfo)(pi.c)))
}

//------------------------------------------------------------------------------
// TypeInfo
//------------------------------------------------------------------------------

type TypeInfo struct {
	BaseInfo
}

type ArrayType int
type TypeTag int

const (
	ARRAY_TYPE_C          ArrayType = C.GI_ARRAY_TYPE_C
	ARRAY_TYPE_ARRAY      ArrayType = C.GI_ARRAY_TYPE_ARRAY
	ARRAY_TYPE_PTR_ARRAY  ArrayType = C.GI_ARRAY_TYPE_PTR_ARRAY
	ARRAY_TYPE_BYTE_ARRAY ArrayType = C.GI_ARRAY_TYPE_BYTE_ARRAY
)

const (
	TYPE_TAG_VOID      TypeTag = C.GI_TYPE_TAG_VOID
	TYPE_TAG_BOOLEAN   TypeTag = C.GI_TYPE_TAG_BOOLEAN
	TYPE_TAG_INT8      TypeTag = C.GI_TYPE_TAG_INT8
	TYPE_TAG_UINT8     TypeTag = C.GI_TYPE_TAG_UINT8
	TYPE_TAG_INT16     TypeTag = C.GI_TYPE_TAG_INT16
	TYPE_TAG_UINT16    TypeTag = C.GI_TYPE_TAG_UINT16
	TYPE_TAG_INT32     TypeTag = C.GI_TYPE_TAG_INT32
	TYPE_TAG_UINT32    TypeTag = C.GI_TYPE_TAG_UINT32
	TYPE_TAG_INT64     TypeTag = C.GI_TYPE_TAG_INT64
	TYPE_TAG_UINT64    TypeTag = C.GI_TYPE_TAG_UINT64
	TYPE_TAG_FLOAT     TypeTag = C.GI_TYPE_TAG_FLOAT
	TYPE_TAG_DOUBLE    TypeTag = C.GI_TYPE_TAG_DOUBLE
	TYPE_TAG_GTYPE     TypeTag = C.GI_TYPE_TAG_GTYPE
	TYPE_TAG_UTF8      TypeTag = C.GI_TYPE_TAG_UTF8
	TYPE_TAG_FILENAME  TypeTag = C.GI_TYPE_TAG_FILENAME
	TYPE_TAG_ARRAY     TypeTag = C.GI_TYPE_TAG_ARRAY
	TYPE_TAG_INTERFACE TypeTag = C.GI_TYPE_TAG_INTERFACE
	TYPE_TAG_GLIST     TypeTag = C.GI_TYPE_TAG_GLIST
	TYPE_TAG_GSLIST    TypeTag = C.GI_TYPE_TAG_GSLIST
	TYPE_TAG_GHASH     TypeTag = C.GI_TYPE_TAG_GHASH
	TYPE_TAG_ERROR     TypeTag = C.GI_TYPE_TAG_ERROR
	TYPE_TAG_UNICHAR   TypeTag = C.GI_TYPE_TAG_UNICHAR
)

// g_type_tag_to_string
func (tt TypeTag) String() string {
	ret := C.g_type_tag_to_string(C.GITypeTag(tt))
	return _GStringToGoString(ret)
}

// g_type_info_is_pointer
func (ti *TypeInfo) IsPointer() bool {
	return C.g_type_info_is_pointer((*C.GITypeInfo)(ti.c)) != 0
}

// g_type_info_get_tag
func (ti *TypeInfo) Tag() TypeTag {
	return TypeTag(C.g_type_info_get_tag((*C.GITypeInfo)(ti.c)))
}

// g_type_info_get_param_type
func (ti *TypeInfo) ParamType(n int) *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_type_info_get_param_type((*C.GITypeInfo)(ti.c), C.gint(n)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_type_info_get_interface
func (ti *TypeInfo) Interface() *BaseInfo {
	cptr := C.g_type_info_get_interface((*C.GITypeInfo)(ti.c))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return _SetBaseInfoFinalizer(ptr)
}

// g_type_info_get_array_length
func (ti *TypeInfo) ArrayLength() int {
	return int(C.g_type_info_get_array_length((*C.GITypeInfo)(ti.c)))
}

// g_type_info_get_array_fixed_size
func (ti *TypeInfo) ArrayFixedSize() int {
	return int(C.g_type_info_get_array_fixed_size((*C.GITypeInfo)(ti.c)))
}

// g_type_info_is_zero_terminated
func (ti *TypeInfo) IsZeroTerminated() bool {
	return C.g_type_info_is_zero_terminated((*C.GITypeInfo)(ti.c)) != 0
}

// g_type_info_get_array_type
func (ti *TypeInfo) ArrayType() ArrayType {
	return ArrayType(C.g_type_info_get_array_type((*C.GITypeInfo)(ti.c)))
}

//------------------------------------------------------------------------------
// CallableInfo
//------------------------------------------------------------------------------

type CallableInfo struct {
	BaseInfo
}

// g_callable_info_get_return_type
func (ci *CallableInfo) ReturnType() *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_callable_info_get_return_type((*C.GICallableInfo)(ci.c)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_callable_info_get_caller_owns
func (ci *CallableInfo) CallerOwns() Transfer {
	return Transfer(C.g_callable_info_get_caller_owns((*C.GICallableInfo)(ci.c)))
}

// g_callable_info_may_return_null
func (ci *CallableInfo) MayReturnNil() bool {
	return C.g_callable_info_may_return_null((*C.GICallableInfo)(ci.c)) != 0
}

// g_callable_info_get_return_attribute
func (ci *CallableInfo) ReturnAttribute(name string) string {
	gname := _GoStringToGString(name)
	ret := C.g_callable_info_get_return_attribute((*C.GICallableInfo)(ci.c), gname)
	C.free_gstring(gname)
	return _GStringToGoString(ret)
}

// g_callable_info_iterate_return_attributes
func (ci *CallableInfo) IterateReturnAttributes(cb func(name, value string)) {
	var iter C.GIAttributeIter
	var cname, cvalue *C.char
	for C.g_callable_info_iterate_return_attributes((*C.GICallableInfo)(ci.c), &iter, &cname, &cvalue) != 0 {
		name, value := C.GoString(cname), C.GoString(cvalue)
		cb(name, value)
	}
}

// g_callable_info_get_n_args
func (ci *CallableInfo) NumArg() int {
	return int(C.g_callable_info_get_n_args((*C.GICallableInfo)(ci.c)))
}

// g_callable_info_get_arg
func (ci *CallableInfo) Arg(n int) *ArgInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_callable_info_get_arg((*C.GICallableInfo)(ci.c), C.gint(n)))}
	return (*ArgInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//void                g_callable_info_load_arg            (GICallableInfo *info,
//                                                         gint n,
//                                                         GIArgInfo *arg);
//void                g_callable_info_load_return_type    (GICallableInfo *info,
//                                                         GITypeInfo *type);

//------------------------------------------------------------------------------
// FunctionInfo
//------------------------------------------------------------------------------

type FunctionInfo struct {
	CallableInfo
}

type FunctionInfoFlags int

const (
	FUNCTION_IS_METHOD      FunctionInfoFlags = C.GI_FUNCTION_IS_METHOD
	FUNCTION_IS_CONSTRUCTOR FunctionInfoFlags = C.GI_FUNCTION_IS_CONSTRUCTOR
	FUNCTION_IS_GETTER      FunctionInfoFlags = C.GI_FUNCTION_IS_GETTER
	FUNCTION_IS_SETTER      FunctionInfoFlags = C.GI_FUNCTION_IS_SETTER
	FUNCTION_WRAPS_VFUNC    FunctionInfoFlags = C.GI_FUNCTION_WRAPS_VFUNC
	FUNCTION_THROWS         FunctionInfoFlags = C.GI_FUNCTION_THROWS
)

// g_function_info_get_symbol
func (fi *FunctionInfo) Symbol() string {
	ret := C.g_function_info_get_symbol((*C.GIFunctionInfo)(fi.c))
	return _GStringToGoString(ret)
}

// g_function_info_get_flags
func (fi *FunctionInfo) Flags() FunctionInfoFlags {
	return FunctionInfoFlags(C.g_function_info_get_flags((*C.GIFunctionInfo)(fi.c)))
}

// g_function_info_get_property
func (fi *FunctionInfo) Property() *PropertyInfo {
	cptr := (*C.GIBaseInfo)(C.g_function_info_get_property((*C.GIFunctionInfo)(fi.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*PropertyInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_function_info_get_vfunc
func (fi *FunctionInfo) VFunc() *VFuncInfo {
	cptr := (*C.GIBaseInfo)(C.g_function_info_get_vfunc((*C.GIFunctionInfo)(fi.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*VFuncInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//gboolean            g_function_info_invoke              (GIFunctionInfo *info,
//                                                         const GIArgument *in_args,
//                                                         int n_in_args,
//                                                         const GIArgument *out_args,
//                                                         int n_out_args,
//                                                         GIArgument *return_value,
//                                                         GError **error);

//------------------------------------------------------------------------------
// SignalInfo
//------------------------------------------------------------------------------

type SignalInfo struct {
	CallableInfo
}

//GSignalFlags        g_signal_info_get_flags             (GISignalInfo *info);

// g_signal_info_get_class_closure
func (si *SignalInfo) ClassClosure() *VFuncInfo {
	cptr := (*C.GIBaseInfo)(C.g_signal_info_get_class_closure((*C.GISignalInfo)(si.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*VFuncInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_signal_info_true_stops_emit
func (si *SignalInfo) TrueStopsEmit() bool {
	return C.g_signal_info_true_stops_emit((*C.GISignalInfo)(si.c)) != 0
}

//------------------------------------------------------------------------------
// VFuncInfo
//------------------------------------------------------------------------------

type VFuncInfo struct {
	CallableInfo
}

type VFuncInfoFlags int

const (
	VFUNC_MUST_CHAIN_UP     VFuncInfoFlags = C.GI_VFUNC_MUST_CHAIN_UP
	VFUNC_MUST_OVERRIDE     VFuncInfoFlags = C.GI_VFUNC_MUST_OVERRIDE
	VFUNC_MUST_NOT_OVERRIDE VFuncInfoFlags = C.GI_VFUNC_MUST_NOT_OVERRIDE
)

// g_vfunc_info_get_flags
func (vfi *VFuncInfo) Flags() VFuncInfoFlags {
	return VFuncInfoFlags(C.g_vfunc_info_get_flags((*C.GIVFuncInfo)(vfi.c)))
}

// g_vfunc_info_get_offset
func (vfi *VFuncInfo) Offset() int {
	return int(C.g_vfunc_info_get_offset((*C.GIVFuncInfo)(vfi.c)))
}

// g_vfunc_info_get_signal
func (vfi *VFuncInfo) Signal() *SignalInfo {
	cptr := (*C.GIBaseInfo)(C.g_vfunc_info_get_signal((*C.GIVFuncInfo)(vfi.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*SignalInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_vfunc_info_get_invoker
func (vfi *VFuncInfo) Invoker() *FunctionInfo {
	cptr := (*C.GIBaseInfo)(C.g_vfunc_info_get_invoker((*C.GIVFuncInfo)(vfi.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//------------------------------------------------------------------------------
// RegisteredType
//------------------------------------------------------------------------------

type RegisteredType struct {
	BaseInfo
}

// g_registered_type_info_get_type_name
func (rt *RegisteredType) TypeName() string {
	ret := C.g_registered_type_info_get_type_name((*C.GIRegisteredTypeInfo)(rt.c))
	return _GStringToGoString(ret)
}

// g_registered_type_info_get_type_init
func (rt *RegisteredType) TypeInit() string {
	ret := C.g_registered_type_info_get_type_init((*C.GIRegisteredTypeInfo)(rt.c))
	return _GStringToGoString(ret)
}

//GType               g_registered_type_info_get_g_type   (GIRegisteredTypeInfo *info);

//------------------------------------------------------------------------------
// EnumInfo
//------------------------------------------------------------------------------

type EnumInfo struct {
	RegisteredType
}

type ValueInfo struct {
	BaseInfo
}

// g_enum_info_get_n_values
func (ei *EnumInfo) NumValue() int {
	return int(C.g_enum_info_get_n_values((*C.GIEnumInfo)(ei.c)))
}

// g_enum_info_get_value
func (ei *EnumInfo) Value(n int) *ValueInfo {
	cptr := (*C.GIBaseInfo)(C.g_enum_info_get_value((*C.GIEnumInfo)(ei.c), C.gint(n)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*ValueInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_enum_info_get_n_methods
func (ei *EnumInfo) NumMethod() int {
	return int(C.g_enum_info_get_n_methods((*C.GIEnumInfo)(ei.c)))
}

// g_enum_info_get_method
func (ii *EnumInfo) Method(n int) *FunctionInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_enum_info_get_method((*C.GIEnumInfo)(ii.c), C.gint(n)))}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_enum_info_get_storage_type
func (ei *EnumInfo) StorageType() TypeTag {
	return TypeTag(C.g_enum_info_get_storage_type((*C.GIEnumInfo)(ei.c)))
}

// g_value_info_get_value
func (vi *ValueInfo) Value() int64 {
	return int64(C.g_value_info_get_value((*C.GIValueInfo)(vi.c)))
}

//------------------------------------------------------------------------------
// InterfaceInfo
//------------------------------------------------------------------------------

type InterfaceInfo struct {
	RegisteredType
}

// g_interface_info_get_n_prerequisites
func (ii *InterfaceInfo) NumPrerequisite() int {
	return int(C.g_interface_info_get_n_prerequisites((*C.GIInterfaceInfo)(ii.c)))
}

// g_interface_info_get_prerequisite
func (ii *InterfaceInfo) Prerequisite(n int) *BaseInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_interface_info_get_prerequisite((*C.GIInterfaceInfo)(ii.c), C.gint(n)))}
	return (*BaseInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_get_n_properties
func (ii *InterfaceInfo) NumProperty() int {
	return int(C.g_interface_info_get_n_properties((*C.GIInterfaceInfo)(ii.c)))
}

// g_interface_info_get_property
func (ii *InterfaceInfo) Property(n int) *PropertyInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_interface_info_get_property((*C.GIInterfaceInfo)(ii.c), C.gint(n)))}
	return (*PropertyInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_get_n_methods
func (ii *InterfaceInfo) NumMethod() int {
	return int(C.g_interface_info_get_n_methods((*C.GIInterfaceInfo)(ii.c)))
}

// g_interface_info_get_method
func (ii *InterfaceInfo) Method(n int) *FunctionInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_interface_info_get_method((*C.GIInterfaceInfo)(ii.c), C.gint(n)))}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_find_method
func (ii *InterfaceInfo) FindMethod(name string) *FunctionInfo {
	gname := _GoStringToGString(name)
	cptr := (*C.GIBaseInfo)(C.g_interface_info_find_method((*C.GIInterfaceInfo)(ii.c), gname))
	C.free_gstring(gname)

	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_get_n_signals
func (ii *InterfaceInfo) NumSignal() int {
	return int(C.g_interface_info_get_n_signals((*C.GIInterfaceInfo)(ii.c)))
}

// g_interface_info_get_signal
func (ii *InterfaceInfo) Signal(n int) *SignalInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_interface_info_get_signal((*C.GIInterfaceInfo)(ii.c), C.gint(n)))}
	return (*SignalInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_get_n_vfuncs
func (ii *InterfaceInfo) NumVFunc() int {
	return int(C.g_interface_info_get_n_vfuncs((*C.GIInterfaceInfo)(ii.c)))
}

// g_interface_info_get_vfunc
func (ii *InterfaceInfo) VFunc(n int) *VFuncInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_interface_info_get_vfunc((*C.GIInterfaceInfo)(ii.c), C.gint(n)))}
	return (*VFuncInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_get_n_constants
func (ii *InterfaceInfo) NumConstant() int {
	return int(C.g_interface_info_get_n_constants((*C.GIInterfaceInfo)(ii.c)))
}

// g_interface_info_get_constant
func (ii *InterfaceInfo) Constant(n int) *ConstantInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_interface_info_get_constant((*C.GIInterfaceInfo)(ii.c), C.gint(n)))}
	return (*ConstantInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//GIStructInfo *      g_interface_info_get_iface_struct   (GIInterfaceInfo *info);
// g_interface_info_get_iface_struct
func (ii *InterfaceInfo) InterfaceStruct() *StructInfo {
	cptr := (*C.GIBaseInfo)(C.g_interface_info_get_iface_struct((*C.GIInterfaceInfo)(ii.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*StructInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_interface_info_find_vfunc
func (ii *InterfaceInfo) FindVFunc(name string) *VFuncInfo {
	gname := _GoStringToGString(name)
	cptr := (*C.GIBaseInfo)(C.g_interface_info_find_vfunc((*C.GIInterfaceInfo)(ii.c), gname))
	C.free_gstring(gname)

	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*VFuncInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

//------------------------------------------------------------------------------
// ObjectInfo
//------------------------------------------------------------------------------

type ObjectInfo struct {
	RegisteredType
}

// g_object_info_get_type_name
func (oi *ObjectInfo) TypeName() string {
	ret := C.g_object_info_get_type_name((*C.GIObjectInfo)(oi.c))
	return _GStringToGoString(ret)
}

// g_object_info_get_type_init
func (oi *ObjectInfo) TypeInit() string {
	ret := C.g_object_info_get_type_init((*C.GIObjectInfo)(oi.c))
	return _GStringToGoString(ret)
}

// g_object_info_get_abstract
func (oi *ObjectInfo) Abstract() bool {
	return C.g_object_info_get_abstract((*C.GIObjectInfo)(oi.c)) != 0
}

// g_object_info_get_fundamental
func (oi *ObjectInfo) Fundamental() bool {
	return C.g_object_info_get_fundamental((*C.GIObjectInfo)(oi.c)) != 0
}

// g_object_info_get_parent
func (oi *ObjectInfo) Parent() *ObjectInfo {
	cptr := (*C.GIBaseInfo)(C.g_object_info_get_parent((*C.GIObjectInfo)(oi.c)))
	if cptr == nil {
		return nil
	}

	ptr := &BaseInfo{cptr}
	return (*ObjectInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_interfaces
func (oi *ObjectInfo) NumInterface() int {
	return int(C.g_object_info_get_n_interfaces((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_interface
func (oi *ObjectInfo) Interface(n int) *InterfaceInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_interface((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*InterfaceInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_fields
func (oi *ObjectInfo) NumField() int {
	return int(C.g_object_info_get_n_fields((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_field
func (oi *ObjectInfo) Field(n int) *FieldInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_field((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*FieldInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_properties
func (oi *ObjectInfo) NumProperty() int {
	return int(C.g_object_info_get_n_properties((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_field
func (oi *ObjectInfo) Property(n int) *PropertyInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_property((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*PropertyInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_methods
func (oi *ObjectInfo) NumMethod() int {
	return int(C.g_object_info_get_n_methods((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_method
func (oi *ObjectInfo) Method(n int) *FunctionInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_method((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_find_method
func (oi *ObjectInfo) FindMethod(name string) *FunctionInfo {
	gname := _GoStringToGString(name)
	cptr := (*C.GIBaseInfo)(C.g_object_info_find_method((*C.GIObjectInfo)(oi.c), gname))
	C.free_gstring(gname)

	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_signals
func (oi *ObjectInfo) NumSignal() int {
	return int(C.g_object_info_get_n_signals((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_signal
func (oi *ObjectInfo) Signal(n int) *SignalInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_signal((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*SignalInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_vfuncs
func (oi *ObjectInfo) NumVFunc() int {
	return int(C.g_object_info_get_n_vfuncs((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_vfunc
func (oi *ObjectInfo) VFunc(n int) *VFuncInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_vfunc((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*VFuncInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_n_constants
func (oi *ObjectInfo) NumConstant() int {
	return int(C.g_object_info_get_n_constants((*C.GIObjectInfo)(oi.c)))
}

// g_object_info_get_constant
func (oi *ObjectInfo) Constant(n int) *ConstantInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_object_info_get_constant((*C.GIObjectInfo)(oi.c), C.gint(n)))}
	return (*ConstantInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_class_struct
func (oi *ObjectInfo) ClassStruct() *StructInfo {
	cptr := (*C.GIBaseInfo)(C.g_object_info_get_class_struct((*C.GIObjectInfo)(oi.c)))
	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*StructInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_find_vfunc
func (oi *ObjectInfo) FindVFunc(name string) *VFuncInfo {
	gname := _GoStringToGString(name)
	cptr := (*C.GIBaseInfo)(C.g_object_info_find_vfunc((*C.GIObjectInfo)(oi.c), gname))
	C.free_gstring(gname)

	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*VFuncInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_object_info_get_unref_function
func (oi *ObjectInfo) UnrefFunction() string {
	ret := C.g_object_info_get_unref_function((*C.GIObjectInfo)(oi.c))
	return _CStringToGoString(ret)
}

//GIObjectInfoUnrefFunction  g_object_info_get_unref_function_pointer
//                                                        (GIObjectInfo *info);

// g_object_info_get_ref_function
func (oi *ObjectInfo) RefFunction() string {
	ret := C.g_object_info_get_ref_function((*C.GIObjectInfo)(oi.c))
	return _CStringToGoString(ret)
}

//GIObjectInfoRefFunction  g_object_info_get_ref_function_pointer
//                                                        (GIObjectInfo *info);

// g_object_info_get_set_value_function
func (oi *ObjectInfo) SetValueFunction() string {
	ret := C.g_object_info_get_set_value_function((*C.GIObjectInfo)(oi.c))
	return _CStringToGoString(ret)
}

//GIObjectInfoSetValueFunction  g_object_info_get_set_value_function_pointer
//                                                        (GIObjectInfo *info);

// g_object_info_get_get_value_function
func (oi *ObjectInfo) GetValueFunction() string {
	ret := C.g_object_info_get_get_value_function((*C.GIObjectInfo)(oi.c))
	return _CStringToGoString(ret)
}

//GIObjectInfoGetValueFunction  g_object_info_get_get_value_function_pointer
//                                                        (GIObjectInfo *info);

//------------------------------------------------------------------------------
// StructInfo
//------------------------------------------------------------------------------

type StructInfo struct {
	RegisteredType
}

// g_struct_info_get_n_fields
func (si *StructInfo) NumField() int {
	return int(C.g_struct_info_get_n_fields((*C.GIStructInfo)(si.c)))
}

// g_struct_info_get_field
func (si *StructInfo) Field(n int) *FieldInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_struct_info_get_field((*C.GIStructInfo)(si.c), C.gint(n)))}
	return (*FieldInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_struct_info_get_n_methods
func (si *StructInfo) NumMethod() int {
	return int(C.g_struct_info_get_n_methods((*C.GIStructInfo)(si.c)))
}

// g_struct_info_get_method
func (si *StructInfo) Method(n int) *FunctionInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_struct_info_get_method((*C.GIStructInfo)(si.c), C.gint(n)))}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_struct_info_find_method
func (si *StructInfo) FindMethod(name string) *FunctionInfo {
	gname := _GoStringToGString(name)
	cptr := (*C.GIBaseInfo)(C.g_struct_info_find_method((*C.GIStructInfo)(si.c), gname))
	C.free_gstring(gname)

	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_struct_info_get_size
func (si *StructInfo) Size() int {
	return int(C.g_struct_info_get_size((*C.GIStructInfo)(si.c)))
}

// g_struct_info_get_alignment
func (si *StructInfo) Alignment() int {
	return int(C.g_struct_info_get_alignment((*C.GIStructInfo)(si.c)))
}

// g_struct_info_is_gtype_struct
func (si *StructInfo) IsGTypeStruct() bool {
	return C.g_struct_info_is_gtype_struct((*C.GIStructInfo)(si.c)) != 0
}

// g_struct_info_is_foreign
func (si *StructInfo) IsForeign() bool {
	return C.g_struct_info_is_foreign((*C.GIStructInfo)(si.c)) != 0
}

//------------------------------------------------------------------------------
// UnionInfo
//------------------------------------------------------------------------------

type UnionInfo struct {
	RegisteredType
}

// g_union_info_get_n_fields
func (ui *UnionInfo) NumField() int {
	return int(C.g_union_info_get_n_fields((*C.GIUnionInfo)(ui.c)))
}

// g_union_info_get_field
func (ui *UnionInfo) Field(n int) FieldInfo {
	return FieldInfo{BaseInfo{
		(*C.GIBaseInfo)(C.g_union_info_get_field((*C.GIUnionInfo)(ui.c), C.gint(n)))}}
}

// g_union_info_get_n_methods
func (ui *UnionInfo) NumMethod() int {
	return int(C.g_union_info_get_n_methods((*C.GIUnionInfo)(ui.c)))
}

// g_union_info_get_method
func (ui *UnionInfo) Method(n int) *FunctionInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_union_info_get_method((*C.GIUnionInfo)(ui.c), C.gint(n)))}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_union_info_is_discriminated
func (ui *UnionInfo) IsDiscriminated() bool {
	return C.g_union_info_is_discriminated((*C.GIUnionInfo)(ui.c)) != 0
}

// g_union_info_get_discriminator_offset
func (ui *UnionInfo) DiscriminatorOffset() int {
	return int(C.g_union_info_get_discriminator_offset((*C.GIUnionInfo)(ui.c)))
}

// g_union_info_get_discriminator_type
func (ui *UnionInfo) DiscriminatorType() *TypeInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_union_info_get_discriminator_type((*C.GIUnionInfo)(ui.c)))}
	return (*TypeInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_union_info_get_discriminator
func (ui *UnionInfo) Discriminator(n int) *ConstantInfo {
	ptr := &BaseInfo{(*C.GIBaseInfo)(C.g_union_info_get_discriminator((*C.GIUnionInfo)(ui.c), C.gint(n)))}
	return (*ConstantInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_union_info_find_method
func (ui *UnionInfo) FindMethod(name string) *FunctionInfo {
	gname := _GoStringToGString(name)
	cptr := (*C.GIBaseInfo)(C.g_union_info_find_method((*C.GIUnionInfo)(ui.c), gname))
	C.free_gstring(gname)

	if cptr == nil {
		return nil
	}
	ptr := &BaseInfo{cptr}
	return (*FunctionInfo)(unsafe.Pointer(_SetBaseInfoFinalizer(ptr)))
}

// g_union_info_get_size
func (ui *UnionInfo) Size() int {
	return int(C.g_union_info_get_size((*C.GIUnionInfo)(ui.c)))
}

// g_union_info_get_alignment
func (ui *UnionInfo) Alignment() int {
	return int(C.g_union_info_get_alignment((*C.GIUnionInfo)(ui.c)))
}
