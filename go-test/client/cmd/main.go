package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
)

func main() {

	data := make([]byte, 10)
	data[0] = 32
	data[1] = 79
	data[2] = 129
	data[3] = 87
	data[4] = 144
	data[5] = 202
	data[6] = 59
	data[7] = 180

	// writeFile("junk.bin", data)
	readFile("junk.bin")

	b := make([]byte, 2)
	b[0] = 32
	b[1] = 78
	fmt.Println(reflect.DeepEqual(data[0:2], b))
	fmt.Printf("data:%x\n", data)
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
