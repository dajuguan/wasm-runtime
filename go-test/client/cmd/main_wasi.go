package main

import (
	"encoding/binary"
	"encoding/hex"
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

		var key [32]byte

		str := "0100000000000000000000000000000000000000000000000000000000000001"
		// big hostio
		// str := "02f26283bdbd7992320d4d707dfa940095f638ad7c95d115150fb2a4417c3ad1"
		// str := "02dc1ac76a8d07580017d3d2120f6fae69df22c4709be4aeaa9aade4990a482e"
		b, _ := hex.DecodeString(str)
		copy(key[:], b)
		getPreimage(key)

		// for i := 0; i < 5; i++ {
		// 	str = "02a2b2ac180c92e0af73975c89261fb8e7fb0e10cb159a3c09883aaa812f3d68"
		// 	b, _ := hex.DecodeString(str)
		// 	copy(key[:], b)
		// 	getPreimage(key)
		// 	// println("received:", buf[0:8])
		// }

		// str = "02a2b2ac180c92e0af73975c89261fb8e7fb0e10cb159a3c09883aaa812f3d68"
		// b, _ = hex.DecodeString(str)
		// copy(key[:], b)
		// getPreimage(key)

		// hintHash := "l1-block-header 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab"
		// getHint(hintHash)
	}()

	wg.Wait()

}

//go:wasmimport _gotest sub
func testSub(uint32, uint32)

func getRandomString() []byte {
	var bufPtr *byte
	var bufSize uint32
	// println("&bufPtr", &bufPtr, &bufSize)
	// println("&bufPtr uintptr", uintptr(unsafe.Pointer(&bufPtr)))
	testSub(uint32(uintptr(unsafe.Pointer(&bufPtr))), uint32(uintptr(unsafe.Pointer(&bufSize))))
	res := unsafe.Slice(bufPtr, bufSize)
	return res
}

//go:wasmexport _gotest minn
func min(a, b uint32) uint32 {
	if a > b {
		return b
	}
	return a
}

func getPreimage(key [32]byte) []byte {
	size := getPreimageLenFromOracle(uint32(uintptr(unsafe.Pointer(&key[0]))))
	// println("len go", size)

	// size = min(size, uint32(65500))
	buf := make([]byte, size)
	readedLen := getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&key[0]))), uint32(uintptr(unsafe.Pointer(&buf[0]))), size)
	if readedLen < size {
		getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&key[0]))), uint32(uintptr(unsafe.Pointer(&buf[readedLen]))), size-readedLen)
	}
	// println("buf first", buf[0:8])
	// println("buf last", buf[size-8:size])
	return buf
}

//go:wasmimport _gotest get_preimage_len
func getPreimageLenFromOracle(keyPtr uint32) uint32

//go:wasmimport _gotest get_preimage_from_oracle
func getPreimageFromOracle(keyPtr uint32, retBufPtr uint32, size uint32) uint32

func getHint(hint string) {
	var hintBytes []byte
	hintBytes = binary.BigEndian.AppendUint32(hintBytes, uint32(len(hint)))
	hintBytes = append(hintBytes, []byte(hint)...)
	// println("hintbytes in go:", len(hintBytes))
	// println("hintbytes in go:", hintBytes)
	hintOracle(uint32(uintptr(unsafe.Pointer(&hintBytes[0]))), uint32(len(hintBytes)))
}

//go:wasmimport _gotest hint_oracle
func hintOracle(retBufPtr uint32, retBufSize uint32)
