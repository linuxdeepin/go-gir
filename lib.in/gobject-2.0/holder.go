package gobject

import (
	"sync"
	"unsafe"
)

/*
#include "gobject.gen.h"
#include "gobject.h"

typedef int32_t (*_GSourceFunc)(void*);
extern uint32_t g_timeout_add(uint32_t, _GSourceFunc, void*);
extern int32_t g_source_remove(uint32_t);

extern int32_t fqueue_dispatcher(void*);
static uint32_t _g_timeout_add_fqueue(uint32_t time) {
	return g_timeout_add(time, fqueue_dispatcher, 0);
}
*/
import "C"

//--------------------------------------------------------------
// Holder
//--------------------------------------------------------------
// holy crap, what am I doing here..

type holder_key [2]unsafe.Pointer
type holder_type map[holder_key]int

func init() {
	var i interface{} = 0
	if unsafe.Sizeof(i) != unsafe.Sizeof(holder_key{}) || unsafe.Sizeof(i) != unsafe.Sizeof(C.GoInterfaceHolder{}) {
		panic("The Go gir world be damaged because interface{} implement changed. :(")
	}
}

var Holder = holder_type(make(map[holder_key]int))

func toGoInterfaceHolder(x interface{}) C.GoInterfaceHolder {
	return *(*C.GoInterfaceHolder)(unsafe.Pointer(&x))
}
func fromGoInterfaceHolder(x C.GoInterfaceHolder) interface{} {
	return *(*interface{})(unsafe.Pointer(&x))
}

func (this holder_type) Grab(x interface{}) {
	if x == nil {
		return
	}

	key := *(*holder_key)(unsafe.Pointer(&x))
	count := this[key]
	this[key] = count + 1
}

func (this holder_type) Release(x interface{}) {
	if x == nil {
		return
	}

	key := *(*holder_key)(unsafe.Pointer(&x))
	count := this[key]
	if count <= 1 {
		delete(this, key)
	} else {
		this[key] = count - 1
	}
}

//--------------------------------------------------------------
// FinalizerQueue
//--------------------------------------------------------------

type finalizer_item struct {
	ptr       unsafe.Pointer
	finalizer func(unsafe.Pointer)
}

type fqueue_type struct {
	sync.Mutex
	queue      []finalizer_item
	exec_queue []finalizer_item
	tid        uint32
}

var FQueue fqueue_type

func (this *fqueue_type) Start(interval int) {
	this.Lock()
	this.queue = make([]finalizer_item, 0, 50)
	this.exec_queue = make([]finalizer_item, 50)
	this.tid = uint32(C._g_timeout_add_fqueue(C.uint32_t(interval)))
	this.Unlock()
}

func (this *fqueue_type) Stop() {
	this.Lock()
	// TODO: we'll discard few items here at Stop, is it ok?
	this.queue = nil
	C.g_source_remove(C.uint32_t(this.tid))
	this.Unlock()
}

// returns true if the item was enqueued, thread safe
func (this *fqueue_type) Push(ptr unsafe.Pointer, finalizer func(unsafe.Pointer)) bool {
	this.Lock()
	if this.queue != nil {
		this.queue = append(this.queue, finalizer_item{ptr, finalizer})
		this.Unlock()
		return true
	}
	this.Unlock()
	return false
}

// exec is only thread safe if executed by a single thread
func (this *fqueue_type) exec() {
	// exec_queue is used for not holding the lock a lot
	this.Lock()
	// common case
	if len(this.queue) == 0 {
		this.Unlock()
		return
	}

	// non-empty queue, copy everything to exec_queue
	if len(this.queue) > len(this.exec_queue) {
		this.exec_queue = make([]finalizer_item, len(this.queue))
	}
	nitems := copy(this.exec_queue, this.queue)
	this.queue = this.queue[:0]
	this.Unlock()

	// then do our work
	for i := 0; i < nitems; i++ {
		this.exec_queue[i].finalizer(this.exec_queue[i].ptr)
		this.exec_queue[i] = finalizer_item{}
	}
}

//export fqueue_dispatcher
func fqueue_dispatcher(unused unsafe.Pointer) int32 {
	FQueue.exec()
	return 1
}
