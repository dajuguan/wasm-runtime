//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		data := getKeyFromOracle()
		fmt.Println("Get resp from host", data)
	}()

	wg.Wait()
}

//export getKeyFromOracle
func getKeyFromOracle() uint32
