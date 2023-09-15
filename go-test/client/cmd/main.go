package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	a := 15
	println(a % 8)
}

func readFile(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic("read file errror")
	}
	fmt.Printf("data:%x\n", data)
}

func writeFile(fileName string, data []byte) {
	fd, _ := os.Create(fileName)
	defer fd.Close()
	for i := 0; i < 8; i++ {
		var ii = data[i]
		err := binary.Write(fd, binary.BigEndian, ii)
		if err != nil {
			fmt.Println("err!", err)
		}
	}
}
