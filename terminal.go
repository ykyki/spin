package main

import (
	"fmt"
	"io"

	"golang.org/x/term"
)

type Terminal struct {
	outStream     io.Writer
	outFd         int
	frame         int
	lineBuf       []string
	prevLineCount int
}

func NewTerminal(
	outStream io.Writer,
	errStream io.Writer,
	outFd int,
) (*Terminal, error) {
	if !term.IsTerminal(outFd) {
		return nil, fmt.Errorf("not a terminal")
	}

	return &Terminal{
		outStream:     outStream,
		outFd:         outFd,
		frame:         0,
		lineBuf:       []string{},
		prevLineCount: 0,
	}, nil
}

const maxBufLen = 100

func (t *Terminal) AddLine(line string) {
	t.lineBuf = append(t.lineBuf, line)

	if len(t.lineBuf) > maxBufLen {
		t.lineBuf = t.lineBuf[1:]
	}
}

func (t *Terminal) getHeight() (int, error) {
	_, height, err := term.GetSize(t.outFd)
	if err != nil {
		return -1, err
	}

	return height, nil
}

func (t *Terminal) Render() {
	for i := 0; i < t.prevLineCount; i++ {
		t.clearCurrentLine()
		t.moveCursorUp(1)
	}

	var maxLineCount int
	{
		height, err := t.getHeight()
		if err != nil {
			maxLineCount = 999
		} else {
			maxLineCount = height - 2
		}
	}

	t.writeSpinnerLine(t.frame, 3)

	lineCount := min(maxLineCount, len(t.lineBuf))

	for i := 0; i < lineCount; i++ {
		fmt.Fprintf(t.outStream, "%s\n", t.lineBuf[len(t.lineBuf)-lineCount+i])
	}

	t.frame++
	t.prevLineCount = lineCount + 1
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
		fmt.Fprintf(t.outStream,
			"\rSpinner 1: %s (%05d)\n", Spinner[frame%(len(Spinner))],
			frame,
		)
	case 1:
		fmt.Fprintf(t.outStream, "\rSpinner 2: %s\n", Spinner2[frame%(len(Spinner2))])
	case 2:
		fmt.Fprintf(t.outStream, "\rSpinner 3: %s\n", Spinner3[frame%(len(Spinner3))])
	case 3:
		fmt.Fprintf(t.outStream,
			"\rSpinner 4: %s (%05d)\n", Spinner4[frame%(len(Spinner4))],
			frame,
		)
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
