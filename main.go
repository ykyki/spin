package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
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

func getTerminalHeight() (int, error) {
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Failed to get terminal size", err)
		return -1, err
	}
	return height, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	height, err := getTerminalHeight()
	if err != nil {
		return
	}
	linesInit := min(5, height/2)

	lines := linesInit
	for frame := 0; frame < 100; frame++ {
		if frame > 0 {
			for i := 0; i < lines; i++ {
				clearCurrentLine()
				moveCursorUp(1)
			}

			height, err := getTerminalHeight()
			if err != nil {
				return
			}
			lines = min(linesInit+frame%200, height-2)
		}

		for i := 0; i < lines; i++ {
			writeSpinnerLine(frame, i%5)
		}

		time.Sleep(50 * time.Millisecond)
	}
}
