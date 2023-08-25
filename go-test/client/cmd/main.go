//go:build !wasm
// +build !wasm

package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"go.wasm.test/client"
)

func main() {
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

// func main() {  //it's ok
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		fmt.Println("test wg is oK")
// 	}()

// 	wg.Wait()
// }
