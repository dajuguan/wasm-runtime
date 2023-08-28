//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"sync"
	"unsafe"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		data := getRandomString()
		for _, v := range data {
			fmt.Println("res=====>", v)
		}
	}()

	wg.Wait()
}

//export allocate_buffer
func allocateBuffer(size uint32) *uint32 {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of get_random_string.
	buf := make([]uint32, size)
	buf[0] = 3
	buf[1] = 2
	// buf[3] = 2
	// buf[4] = 3
	// println("offset in go ", &buf[0])
	// println("offset in go ", &buf[1])
	// println("offset in go %v", buf)
	return &buf[0]
}

//export getKeyFromOracle
// func getKeyFromOracle() []byte

//export get_random_string
func getRandomStringRaw(retBufPtr **uint32, retBufSize *int)

// Get random string from the hosts.
func getRandomString() []uint32 {
	var bufPtr *uint32
	var bufSize int
	getRandomStringRaw(&bufPtr, &bufSize)
	println("bufPtr in go after", *bufPtr)
	return unsafe.Slice(bufPtr, bufSize)
	// return unsafe.String(bufPtr, bufSize)
}

// func getRandomString() []byte {
// 	var bufPtr *byte
// 	var bufSize int
// 	bufPtr = allocateBuffer(10)
// 	bufSize = 10

// 	// println("res...", *(bufPtr + 1) )
// 	var res []byte = unsafe.Slice(bufPtr, bufSize)
// 	// var  unsafe.String(bufPtr, bufSize)
// 	for v := range res {
// 		println("v===========> %", v)
// 	}
// 	fmt.Printf("res %v", string(res))
// 	return make([]byte, 1)
// }
