package gudev

//go:generate go-gir-generator -o . gudev.go.in

import "C"
import "unsafe"

//export _GUdev_go_callback_cleanup
func _GUdev_go_callback_cleanup(unsafe.Pointer) {
}
