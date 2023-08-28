//go:build !wasm
// +build !wasm

package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"

	"go.wasm.test/client"
)

func main() {
	str := "0100000000000000000000000000000000000000000000000000000000000001"
	b, _ := hex.DecodeString(str)

	encodedStr := hex.EncodeToString(b)
	fmt.Println("@@@@--bytes-->", b)
	fmt.Printf("@@@@--bytes-->%02x \n", b)
	fmt.Printf("@@@@--string-->%s \n", encodedStr)
}

func test() {
	reader := os.NewFile(client.PClientRFd, "preimage-oracle-read")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		data := make([]byte, 1024)

		n, err := reader.Read(data)
		if n > 0 {
			fmt.Println(string(data[:n]))
		}

		if err != nil && err != io.EOF {
			fmt.Println(err)
		}
	}()

	wg.Wait()

}

//export get_random_string
// func getRandomStringRaw(retBufPtr **byte, retBufSize *int)

// Get random string from the hosts.
// func getRandomString() string {
// 	var bufPtr *byte
// 	var bufSize int
// 	getRandomStringRaw(&bufPtr, &bufSize)
// 	return unsafe.String(bufPtr, bufSize)
// }

// func main() {  //it's ok
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		fmt.Println("test wg is oK")
// 	}()

// 	wg.Wait()
// }
