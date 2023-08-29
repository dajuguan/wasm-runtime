//go:build js && !wasm
// +build js,!wasm

package main

import (
	"fmt"
	"syscall/js"
	"unsafe"
)

func main() {
	js.Global().Set("allocate_buffer", js.FuncOf(allocateBufferFunc))
	println("test....")
	res := getRandomString()
	println("getRandomString:", string(res))

	fmt.Printf("getRandomString:%02x\n", res)
}

func allocateBufferFunc(this js.Value, args []js.Value) interface{} {
	return int(allocateBuffer(uint32(args[0].Int())))
	// return js.ValueOf()
}

//go:wasmimport _gotest sub
func testSub(uint32, uint32)

//export allocate_buffer
func allocateBuffer(size uint32) uintptr {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of get_random_string.
	buf := make([]uint8, size)
	buf[0] = 3
	p := uintptr(unsafe.Pointer(&buf[0]))
	println("go uintptr:", &buf[0], p)

	return p
}

func getRandomString() []byte {
	var bufPtr *byte
	var bufSize uint32
	println("&bufPtr", &bufPtr, &bufSize)
	println("&bufPtr uintptr", uintptr(unsafe.Pointer(&bufPtr)))
	testSub(uint32(uintptr(unsafe.Pointer(&bufPtr))), uint32(uintptr(unsafe.Pointer(&bufSize))))
	res := unsafe.Slice(bufPtr, bufSize)
	return res
}
