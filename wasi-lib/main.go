package main

import "time"

func main() {
	for i := 0; i < 5; i++ {
		dataq := make([]byte, 10000)
		for j := 0; j < len(dataq); j++ {
			dataq[j] = 3
		}
		time.Sleep(time.Second)
	}
	println("go finished!")
}
