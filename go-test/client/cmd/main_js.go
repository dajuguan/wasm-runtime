//go:build js && !wasm
// +build js,!wasm

package main

import (
	"encoding/binary"
	"fmt"
	"sync"
	"unsafe"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// res := getRandomString()
		// println("getRandomString:", string(res))
		// fmt.Printf("getRandomString:%02x\n", res)

		// var key [32]byte
		// str := "0100000000000000000000000000000000000000000000000000000000000006"
		// b, _ := hex.DecodeString(str)
		// copy(key[:], b)
		// getPreimage(key)

		hintHash := "l1-block-header 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab"
		getHint(hintHash)
	}()

	wg.Wait()

}

//go:wasmimport _gotest sub
func testSub(uint32, uint32)

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
	size := getPreimageLenFromOracle(uint32(uintptr(unsafe.Pointer(&key[0]))))
	println("len go", size)

	buf := make([]byte, size)
	getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&key[0]))), uint32(uintptr(unsafe.Pointer(&buf[0]))), size)
	fmt.Printf("received: %02x \n", buf)
}

//go:wasmimport _gotest get_preimage_len
func getPreimageLenFromOracle(keyPtr uint32) uint32

//go:wasmimport _gotest get_preimage_from_oracle
func getPreimageFromOracle(keyPtr uint32, retBufPtr uint32, size uint32)

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
