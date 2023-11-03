package main

import (
	"fmt"
	"time"
)

func main() {
	spinner := []string{"|", "/", "-", "\\"}
	for i := 0; i < 50; i++ {
		fmt.Printf("\r%s ", spinner[i%4])
		time.Sleep(100 * time.Millisecond)
	}
}
