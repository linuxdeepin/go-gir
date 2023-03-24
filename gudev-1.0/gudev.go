package gudev

/*
#include "gudev.gen.h"


GList* g_list_append(GList*, void*);
void g_list_free(GList*);


#cgo pkg-config: gudev-1.0
*/
import "C"
import "unsafe"

import (
    "github.com/linuxdeepin/go-gir/gobject-2.0"
)

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



type ClientLike interface {
	gobject.ObjectLike
	InheritedFromGUdevClient() *C.GUdevClient
}

type Client struct {
	gobject.Object
	
}

func ToClient(objlike gobject.ObjectLike) *Client {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*Client)(nil).GetStaticType()
	obj := gobject.ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*Client)(obj)
	}
	panic("cannot cast to Client")
}

func (this0 *Client) InheritedFromGUdevClient() *C.GUdevClient {
	if this0 == nil {
		return nil
	}
	return (*C.GUdevClient)(this0.C)
}

func (this0 *Client) GetStaticType() gobject.Type {
	return gobject.Type(C.g_udev_client_get_type())
}

func ClientGetType() gobject.Type {
	return (*Client)(nil).GetStaticType()
}
func NewClient(subsystems0 []string) *Client {
	var subsystems1 **C.char
	subsystems1 = (**C.char)(C.malloc(C.size_t(int(unsafe.Sizeof(*subsystems1)) * (len(subsystems0) + 1))))
	defer C.free(unsafe.Pointer(subsystems1))
	for i, e := range subsystems0 {
		(*(*[999999]*C.char)(unsafe.Pointer(subsystems1)))[i] = _GoStringToGString(e)
		defer C.free(unsafe.Pointer((*(*[999999]*C.char)(unsafe.Pointer(subsystems1)))[i]))
	}
	(*(*[999999]*C.char)(unsafe.Pointer(subsystems1)))[len(subsystems0)] = nil
	ret1 := C.g_udev_client_new(subsystems1)
	var ret2 *Client

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Client)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Client) QueryByDeviceFile(device_file0 string) *Device {
	var this1 *C.GUdevClient
	var device_file1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevClient)(this0.InheritedFromGUdevClient())
	}
	device_file1 = _GoStringToGString(device_file0)
	defer C.free(unsafe.Pointer(device_file1))
	ret1 := C.g_udev_client_query_by_device_file(this1, device_file1)
	var ret2 *Device

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Device)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Client) QueryByDeviceNumber(type0 DeviceType, number0 uint64) *Device {
	var this1 *C.GUdevClient
	var type1 C.GUdevDeviceType
	var number1 C.uint64_t
	if this0 != nil {
		this1 = (*C.GUdevClient)(this0.InheritedFromGUdevClient())
	}
	type1 = C.GUdevDeviceType(type0)
	number1 = C.uint64_t(number0)
	ret1 := C.g_udev_client_query_by_device_number(this1, type1, number1)
	var ret2 *Device

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Device)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Client) QueryBySubsystem(subsystem0 string) []*Device {
	var this1 *C.GUdevClient
	var subsystem1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevClient)(this0.InheritedFromGUdevClient())
	}
	subsystem1 = _GoStringToGString(subsystem0)
	defer C.free(unsafe.Pointer(subsystem1))
	ret1 := C.g_udev_client_query_by_subsystem(this1, subsystem1)
	var ret2 []*Device

//DEBUG: ret1(glist):flags = " conv_own_everything"
	for iter := (*_GList)(unsafe.Pointer(ret1)); iter != nil; iter = iter.next {
		var elt *Device
		elt = (*Device)(gobject.ObjectWrap(unsafe.Pointer((*C.GUdevDevice)(iter.data)), false))
		ret2 = append(ret2, elt)
	}
	C.g_list_free(ret1)
	return ret2
}
func (this0 *Client) QueryBySubsystemAndName(subsystem0 string, name0 string) *Device {
	var this1 *C.GUdevClient
	var subsystem1 *C.char
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevClient)(this0.InheritedFromGUdevClient())
	}
	subsystem1 = _GoStringToGString(subsystem0)
	defer C.free(unsafe.Pointer(subsystem1))
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_client_query_by_subsystem_and_name(this1, subsystem1, name1)
	var ret2 *Device

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Device)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Client) QueryBySysfsPath(sysfs_path0 string) *Device {
	var this1 *C.GUdevClient
	var sysfs_path1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevClient)(this0.InheritedFromGUdevClient())
	}
	sysfs_path1 = _GoStringToGString(sysfs_path0)
	defer C.free(unsafe.Pointer(sysfs_path1))
	ret1 := C.g_udev_client_query_by_sysfs_path(this1, sysfs_path1)
	var ret2 *Device

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Device)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
type DeviceLike interface {
	gobject.ObjectLike
	InheritedFromGUdevDevice() *C.GUdevDevice
}

type Device struct {
	gobject.Object
	
}

func ToDevice(objlike gobject.ObjectLike) *Device {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*Device)(nil).GetStaticType()
	obj := gobject.ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*Device)(obj)
	}
	panic("cannot cast to Device")
}

func (this0 *Device) InheritedFromGUdevDevice() *C.GUdevDevice {
	if this0 == nil {
		return nil
	}
	return (*C.GUdevDevice)(this0.C)
}

func (this0 *Device) GetStaticType() gobject.Type {
	return gobject.Type(C.g_udev_device_get_type())
}

func DeviceGetType() gobject.Type {
	return (*Device)(nil).GetStaticType()
}
func (this0 *Device) GetAction() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_action(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetDeviceFile() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_device_file(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetDeviceFileSymlinks() []string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_device_file_symlinks(this1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetDeviceNumber() uint64 {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_device_number(this1)
	var ret2 uint64

//DEBUG: ret1(guint64):flags = " conv_own_none"
	ret2 = uint64(ret1)
	return ret2
}
func (this0 *Device) GetDeviceType() DeviceType {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_device_type(this1)
	var ret2 DeviceType

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = DeviceType(ret1)
	return ret2
}
func (this0 *Device) GetDevtype() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_devtype(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetDriver() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_driver(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetIsInitialized() bool {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_is_initialized(this1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
func (this0 *Device) GetName() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_name(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetNumber() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_number(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetParent() *Device {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_parent(this1)
	var ret2 *Device

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Device)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Device) GetParentWithSubsystem(subsystem0 string, devtype0 string) *Device {
	var this1 *C.GUdevDevice
	var subsystem1 *C.char
	var devtype1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	subsystem1 = _GoStringToGString(subsystem0)
	defer C.free(unsafe.Pointer(subsystem1))
	devtype1 = _GoStringToGString(devtype0)
	defer C.free(unsafe.Pointer(devtype1))
	ret1 := C.g_udev_device_get_parent_with_subsystem(this1, subsystem1, devtype1)
	var ret2 *Device

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Device)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Device) GetProperty(key0 string) string {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_get_property(this1, key1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetPropertyAsBoolean(key0 string) bool {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_get_property_as_boolean(this1, key1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
func (this0 *Device) GetPropertyAsDouble(key0 string) float64 {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_get_property_as_double(this1, key1)
	var ret2 float64

//DEBUG: ret1(gdouble):flags = " conv_own_none"
	ret2 = float64(ret1)
	return ret2
}
func (this0 *Device) GetPropertyAsInt(key0 string) int32 {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_get_property_as_int(this1, key1)
	var ret2 int32

//DEBUG: ret1(gint32):flags = " conv_own_none"
	ret2 = int32(ret1)
	return ret2
}
func (this0 *Device) GetPropertyAsStrv(key0 string) []string {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_get_property_as_strv(this1, key1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetPropertyAsUint64(key0 string) uint64 {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_get_property_as_uint64(this1, key1)
	var ret2 uint64

//DEBUG: ret1(guint64):flags = " conv_own_none"
	ret2 = uint64(ret1)
	return ret2
}
func (this0 *Device) GetPropertyKeys() []string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_property_keys(this1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetSeqnum() uint64 {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_seqnum(this1)
	var ret2 uint64

//DEBUG: ret1(guint64):flags = " conv_own_none"
	ret2 = uint64(ret1)
	return ret2
}
func (this0 *Device) GetSubsystem() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_subsystem(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttr(name0 string) string {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr(this1, name1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrAsBoolean(name0 string) bool {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_boolean(this1, name1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
func (this0 *Device) GetSysfsAttrAsBooleanUncached(name0 string) bool {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_boolean_uncached(this1, name1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
func (this0 *Device) GetSysfsAttrAsDouble(name0 string) float64 {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_double(this1, name1)
	var ret2 float64

//DEBUG: ret1(gdouble):flags = " conv_own_none"
	ret2 = float64(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrAsDoubleUncached(name0 string) float64 {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_double_uncached(this1, name1)
	var ret2 float64

//DEBUG: ret1(gdouble):flags = " conv_own_none"
	ret2 = float64(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrAsInt(name0 string) int32 {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_int(this1, name1)
	var ret2 int32

//DEBUG: ret1(gint32):flags = " conv_own_none"
	ret2 = int32(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrAsIntUncached(name0 string) int32 {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_int_uncached(this1, name1)
	var ret2 int32

//DEBUG: ret1(gint32):flags = " conv_own_none"
	ret2 = int32(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrAsStrv(name0 string) []string {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_strv(this1, name1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetSysfsAttrAsStrvUncached(name0 string) []string {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_strv_uncached(this1, name1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetSysfsAttrAsUint64(name0 string) uint64 {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_uint64(this1, name1)
	var ret2 uint64

//DEBUG: ret1(guint64):flags = " conv_own_none"
	ret2 = uint64(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrAsUint64Uncached(name0 string) uint64 {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_as_uint64_uncached(this1, name1)
	var ret2 uint64

//DEBUG: ret1(guint64):flags = " conv_own_none"
	ret2 = uint64(ret1)
	return ret2
}
func (this0 *Device) GetSysfsAttrKeys() []string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_sysfs_attr_keys(this1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetSysfsAttrUncached(name0 string) string {
	var this1 *C.GUdevDevice
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_device_get_sysfs_attr_uncached(this1, name1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetSysfsPath() string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_sysfs_path(this1)
	var ret2 string

//DEBUG: ret1(utf8):flags = " conv_own_none"
	ret2 = C.GoString(ret1)
	return ret2
}
func (this0 *Device) GetTags() []string {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_tags(this1)
	var ret2 []string

//DEBUG: ret1(array):flags = " conv_own_none"
	ret2 = make([]string, C._array_length(unsafe.Pointer(ret1)))
	for i0 := range ret2 {
		ret2[i0] = C.GoString((*(*[999999]*C.char)(unsafe.Pointer(ret1)))[i0])
	}
	return ret2
}
func (this0 *Device) GetUsecSinceInitialized() uint64 {
	var this1 *C.GUdevDevice
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	ret1 := C.g_udev_device_get_usec_since_initialized(this1)
	var ret2 uint64

//DEBUG: ret1(guint64):flags = " conv_own_none"
	ret2 = uint64(ret1)
	return ret2
}
func (this0 *Device) HasProperty(key0 string) bool {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_has_property(this1, key1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
func (this0 *Device) HasSysfsAttr(key0 string) bool {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_has_sysfs_attr(this1, key1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
func (this0 *Device) HasSysfsAttrUncached(key0 string) bool {
	var this1 *C.GUdevDevice
	var key1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevDevice)(this0.InheritedFromGUdevDevice())
	}
	key1 = _GoStringToGString(key0)
	defer C.free(unsafe.Pointer(key1))
	ret1 := C.g_udev_device_has_sysfs_attr_uncached(this1, key1)
	var ret2 bool

//DEBUG: ret1(gboolean):flags = " conv_own_none"
	ret2 = ret1 != 0
	return ret2
}
type DeviceType C.uint32_t
const (
	DeviceTypeNone DeviceType = 0
	DeviceTypeBlock DeviceType = 98
	DeviceTypeChar DeviceType = 99
)
type EnumeratorLike interface {
	gobject.ObjectLike
	InheritedFromGUdevEnumerator() *C.GUdevEnumerator
}

type Enumerator struct {
	gobject.Object
	
}

func ToEnumerator(objlike gobject.ObjectLike) *Enumerator {
	c := objlike.InheritedFromGObject()
	if c == nil {
		return nil
	}
	t := (*Enumerator)(nil).GetStaticType()
	obj := gobject.ObjectGrabIfType(unsafe.Pointer(c), t)
	if obj != nil {
		return (*Enumerator)(obj)
	}
	panic("cannot cast to Enumerator")
}

func (this0 *Enumerator) InheritedFromGUdevEnumerator() *C.GUdevEnumerator {
	if this0 == nil {
		return nil
	}
	return (*C.GUdevEnumerator)(this0.C)
}

func (this0 *Enumerator) GetStaticType() gobject.Type {
	return gobject.Type(C.g_udev_enumerator_get_type())
}

func EnumeratorGetType() gobject.Type {
	return (*Enumerator)(nil).GetStaticType()
}
func NewEnumerator(client0 ClientLike) *Enumerator {
	var client1 *C.GUdevClient
	if client0 != nil {
		client1 = (*C.GUdevClient)(client0.InheritedFromGUdevClient())
	}
	ret1 := C.g_udev_enumerator_new(client1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_everything"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), false))
	return ret2
}
func (this0 *Enumerator) AddMatchIsInitialized() *Enumerator {
	var this1 *C.GUdevEnumerator
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	ret1 := C.g_udev_enumerator_add_match_is_initialized(this1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddMatchName(name0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var name1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	ret1 := C.g_udev_enumerator_add_match_name(this1, name1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddMatchProperty(name0 string, value0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var name1 *C.char
	var value1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	value1 = _GoStringToGString(value0)
	defer C.free(unsafe.Pointer(value1))
	ret1 := C.g_udev_enumerator_add_match_property(this1, name1, value1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddMatchSubsystem(subsystem0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var subsystem1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	subsystem1 = _GoStringToGString(subsystem0)
	defer C.free(unsafe.Pointer(subsystem1))
	ret1 := C.g_udev_enumerator_add_match_subsystem(this1, subsystem1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddMatchSysfsAttr(name0 string, value0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var name1 *C.char
	var value1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	value1 = _GoStringToGString(value0)
	defer C.free(unsafe.Pointer(value1))
	ret1 := C.g_udev_enumerator_add_match_sysfs_attr(this1, name1, value1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddMatchTag(tag0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var tag1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	tag1 = _GoStringToGString(tag0)
	defer C.free(unsafe.Pointer(tag1))
	ret1 := C.g_udev_enumerator_add_match_tag(this1, tag1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddNomatchSubsystem(subsystem0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var subsystem1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	subsystem1 = _GoStringToGString(subsystem0)
	defer C.free(unsafe.Pointer(subsystem1))
	ret1 := C.g_udev_enumerator_add_nomatch_subsystem(this1, subsystem1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddNomatchSysfsAttr(name0 string, value0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var name1 *C.char
	var value1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	name1 = _GoStringToGString(name0)
	defer C.free(unsafe.Pointer(name1))
	value1 = _GoStringToGString(value0)
	defer C.free(unsafe.Pointer(value1))
	ret1 := C.g_udev_enumerator_add_nomatch_sysfs_attr(this1, name1, value1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) AddSysfsPath(sysfs_path0 string) *Enumerator {
	var this1 *C.GUdevEnumerator
	var sysfs_path1 *C.char
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	sysfs_path1 = _GoStringToGString(sysfs_path0)
	defer C.free(unsafe.Pointer(sysfs_path1))
	ret1 := C.g_udev_enumerator_add_sysfs_path(this1, sysfs_path1)
	var ret2 *Enumerator

//DEBUG: ret1(interface):flags = " conv_own_none"
	ret2 = (*Enumerator)(gobject.ObjectWrap(unsafe.Pointer(ret1), true))
	return ret2
}
func (this0 *Enumerator) Execute() []*Device {
	var this1 *C.GUdevEnumerator
	if this0 != nil {
		this1 = (*C.GUdevEnumerator)(this0.InheritedFromGUdevEnumerator())
	}
	ret1 := C.g_udev_enumerator_execute(this1)
	var ret2 []*Device

//DEBUG: ret1(glist):flags = " conv_own_everything"
	for iter := (*_GList)(unsafe.Pointer(ret1)); iter != nil; iter = iter.next {
		var elt *Device
		elt = (*Device)(gobject.ObjectWrap(unsafe.Pointer((*C.GUdevDevice)(iter.data)), false))
		ret2 = append(ret2, elt)
	}
	C.g_list_free(ret1)
	return ret2
}
