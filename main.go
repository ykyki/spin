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

func clearCurrentLine() {
	fmt.Printf("\033[2K")
}

func writeSpinnerLine(frame int, kind int) {
	switch kind {
	case 0:
		fmt.Printf("\rSpinner 1: %s\n", Spinner[frame%(len(Spinner))])
	case 1:
		fmt.Printf("\rSpinner 2: %s\n", Spinner2[frame%(len(Spinner2))])
	case 2:
		fmt.Printf("\rSpinner 3: %s\n", Spinner3[frame%(len(Spinner3))])
	case 3:
		fmt.Printf("\rSpinner 4: %s\n", Spinner4[frame%(len(Spinner4))])
	default:
		fmt.Printf("\rSpinners : %s%s%s%s\n", Spinner[frame%(len(Spinner))], Spinner[(frame+1)%(len(Spinner))], Spinner[(frame+2)%(len(Spinner))], Spinner[(frame+3)%(len(Spinner))])
	}
}

func main() {
	const linesInit = 5
	lines := linesInit
	for frame := 0; frame < 50; frame++ {
		if frame > 0 {
			for i := 0; i < lines; i++ {
				clearCurrentLine()
				moveCursorUp(1)
			}
			lines = linesInit + frame%10
		}

		for i := 0; i < lines; i++ {
			writeSpinnerLine(frame, i%5)
		}

		time.Sleep(100 * time.Millisecond)
	}

	for i := 0; i < lines-linesInit; i++ {
		clearCurrentLine()
		moveCursorUp(1)
	}
}
