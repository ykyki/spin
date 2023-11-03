package main

import (
	"fmt"
	"time"
)

var Spinner = [4]string{"|", "/", "-", "\\"}
var Spinner2 = [4]string{"\u001b[31m|\u001b[0m", "\u001b[32m/\u001b[0m", "\u001b[33m-\u001b[0m", "\u001b[34m\\\u001b[0m"}
var Spinner3 = [4]string{"↑", "→", "↓", "←"}
var Spinner4 = [8]string{"⬆️", "↗️", "➡️", "↘️", "⬇️", "↙️", "⬅️", "↖️"}

// VT100 escape codes
// https://espterm.github.io/docs/VT100%20escape%20codes.html
func moveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

// func clearCurrentLine() {
// 	fmt.Printf("\033[2K")
// }

func main() {
	for i := 0; i < 50; i++ {
		if i > 0 {
			moveCursorUp(5)
		}
		fmt.Printf("\rSpinner 1: %s\n", Spinner[i%(len(Spinner))])
		fmt.Printf("\rSpinner 2: %s\n", Spinner2[i%(len(Spinner2))])
		fmt.Printf("\rSpinner 3: %s\n", Spinner3[i%(len(Spinner3))])
		fmt.Printf("\rSpinner 4: %s\n", Spinner4[i%(len(Spinner4))])
		fmt.Printf("\rSpinners : %s%s%s%s\n", Spinner[i%(len(Spinner))], Spinner[(i+1)%(len(Spinner))], Spinner[(i+2)%(len(Spinner))], Spinner[(i+3)%(len(Spinner))])
		time.Sleep(100 * time.Millisecond)
	}
}
