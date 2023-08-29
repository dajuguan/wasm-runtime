//go:build js && !wasm
// +build js,!wasm

package main

import (
	"encoding/binary"
	"fmt"
	"sync"
	"syscall/js"
	"unsafe"
)

func main() {

	js.Global().Set("allocate_buffer", js.FuncOf(allocateBufferFunc))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// res := getRandomString()
		// println("getRandomString:", string(res))
		// fmt.Printf("getRandomString:%02x\n", res)

		// var key [32]byte
		// str := "0100000000000000000000000000000000000000000000000000000000000001"
		// b, _ := hex.DecodeString(str)
		// copy(key[:], b)
		// getPreimage(key)

		hintHash := "l1-block-header 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab"
		getHint(hintHash)
	}()

	wg.Wait()

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

func getPreimage(key [32]byte) {
	var bufPtr *byte
	var bufSize uint32
	getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&key[0]))), uint32(uintptr(unsafe.Pointer(&bufPtr))), uint32(uintptr(unsafe.Pointer(&bufSize))))
	res := unsafe.Slice(bufPtr, bufSize)
	fmt.Printf("received: %02x", res)
}

//go:wasmimport _gotest get_preimage_from_oracle
func getPreimageFromOracle(keyPtr uint32, retBufPtr uint32, retBufSize uint32)

func getHint(hint string) {
	var hintBytes []byte
	hintBytes = binary.BigEndian.AppendUint32(hintBytes, uint32(len(hint)))
	hintBytes = append(hintBytes, []byte(hint)...)
	println("hintbytes in go:", len(hintBytes))
	fmt.Printf("hintbytes in go:%02x", hintBytes)
	hintOracle(uint32(uintptr(unsafe.Pointer(&hintBytes[0]))), uint32(len(hintBytes)))
}

//go:wasmimport _gotest hint_oracle
func hintOracle(retBufPtr uint32, retBufSize uint32)
