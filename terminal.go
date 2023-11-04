package main

import (
	"fmt"
	"io"

	"golang.org/x/term"
)

type Terminal struct {
	outStream io.Writer
	outFd     int
}

func NewTerminal(outStream io.Writer, outFd int) *Terminal {
	return &Terminal{
		outStream: outStream,
		outFd:     outFd,
	}
}

func (t *Terminal) getHeight() (int, error) {
	_, height, err := term.GetSize(t.outFd)
	if err != nil {
		return -1, err
	}
	return height, nil
}

// VT100 escape codes
// https://espterm.github.io/docs/VT100%20escape%20codes.html
func (t *Terminal) moveCursorUp(lines int) {
	fmt.Fprintf(t.outStream, "\033[%dA", lines)
}

func (t *Terminal) clearCurrentLine() {
	fmt.Fprintf(t.outStream, "\033[2K")
}

func (t *Terminal) writeSpinnerLine(frame int, kind int) {
	switch kind {
	case 0:
		fmt.Fprintf(t.outStream, "\rSpinner 1: %s\n", Spinner[frame%(len(Spinner))])
	case 1:
		fmt.Fprintf(t.outStream, "\rSpinner 2: %s\n", Spinner2[frame%(len(Spinner2))])
	case 2:
		fmt.Fprintf(t.outStream, "\rSpinner 3: %s\n", Spinner3[frame%(len(Spinner3))])
	case 3:
		fmt.Fprintf(t.outStream, "\rSpinner 4: %s\n", Spinner4[frame%(len(Spinner4))])
	default:
		fmt.Fprintf(
			t.outStream,
			"\rSpinners : %s%s%s%s\n",
			Spinner[frame%(len(Spinner))],
			Spinner[(frame+1)%(len(Spinner))],
			Spinner[(frame+2)%(len(Spinner))],
			Spinner[(frame+3)%(len(Spinner))],
		)
	}
}
