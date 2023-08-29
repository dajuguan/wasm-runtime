//go:build !wasm

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	cl "go.wasm.test/client"
	oppio "go.wasm.test/io"
)

func main() {

	var (
		pClientRW oppio.FileChannel
		hClientRW oppio.FileChannel
	)
	defer func() {
		if pClientRW != nil {
			_ = pClientRW.Close()
		}
		if hClientRW != nil {
			_ = hClientRW.Close()
		}
	}()

	// Setup client I/O for preimage oracle interaction
	pClientRW, pHostRW, err := oppio.CreateBidirectionalChannel()
	if err != nil {
		fmt.Errorf("failed to create preimage pipe: %w", err)
	}

	pHostRW.Write([]byte("Hello,world!!!"))

	ctx := context.Background()
	parts := strings.Fields("node ../wasm-lib/main.js /root/now/wasm-runtime/go-test/client/cmd/client.wasm")
	// parts := strings.Fields("wasmtime --mapdir=/tmp::/root/test /root/now/wasm-runtime/go-test/cmd/client.wasi")
	// parts := strings.Fields("/root/now/wasm-runtime/go-test/cmd/client.exe")
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.ExtraFiles = make([]*os.File, cl.MaxFd-3) // not including stdin, stdout and stderr
	cmd.ExtraFiles[cl.PClientRFd-3] = pClientRW.Reader()
	cmd.ExtraFiles[cl.PClientWFd-3] = pClientRW.Writer()
	cmd.Stdout = os.Stdout // for debugging
	cmd.Stderr = os.Stderr // for debugging
	err = cmd.Start()
	if err != nil {
		fmt.Errorf("program cmd failed to start: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Errorf("failed to wait for child program: %w", err)
	}
	fmt.Println("Client program completed successfully")

}
