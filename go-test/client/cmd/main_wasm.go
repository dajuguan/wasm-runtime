//go:build wasm
// +build wasm

package main

import (
	"encoding/hex"
	"fmt"
	"sync"
	"unsafe"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var key [32]byte
		str := "0100000000000000000000000000000000000000000000000000000000000001"
		b, _ := hex.DecodeString(str)
		copy(key[:], b)
		getPreimage(key)
		// data := getRandomString()
		// fmt.Println("res=====", string(data))
	}()

	wg.Wait()
}

//export allocate_buffer
func allocateBuffer(size uint32) *uint8 {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of get_random_string.
	buf := make([]uint8, size)
	return &buf[0]
}

func getPreimage(key [32]byte) {
	var bufPtr *byte
	var bufSize uint32
	getPreimageFromOracle(key, &bufPtr, &bufSize)
	res := unsafe.Slice(bufPtr, bufSize)
	fmt.Printf("received: %02x", res)
}

//export get_preimage_from_oracle
func getPreimageFromOracle(key [32]byte, retBufPtr **byte, retBufSize *uint32)

//export get_random_string
func getRandomStringRaw(retBufPtr **byte, retBufSize *uint32)

// Get random string from the hosts.
func getRandomString() []byte {
	var bufPtr *byte
	var bufSize uint32
	getRandomStringRaw(&bufPtr, &bufSize)
	println("bufPtr in go after", *bufPtr)
	res := unsafe.Slice(bufPtr, bufSize)
	return res
}
