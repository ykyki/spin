package main

import (
	"fmt"
	"io"
	"time"

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

func (t *Terminal) Render(duration time.Duration) {
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

	t.writeCarriageReturn()
	t.writeSpinnerLine(t.frame, PlainSpinner)
	t.writeString(fmt.Sprintf(" %.1fs", duration.Seconds()))
	t.writeNewLine()

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

func (t *Terminal) writeNewLine() {
	fmt.Fprintf(t.outStream, "\n")
}

func (t *Terminal) writeCarriageReturn() {
	fmt.Fprintf(t.outStream, "\r")
}

func (t *Terminal) writeSpinnerLine(frame int, kind SpinnerKind) {
	switch kind {
	case PlainSpinner:
		fmt.Fprintf(t.outStream, "%s", PainSpinnerSeq[frame%(len(PainSpinnerSeq))])
	case ColorfulSpinner:
		fmt.Fprintf(t.outStream, "%s", ColorfulSpinnerSeq[frame%(len(ColorfulSpinnerSeq))])
	case ArrowSpinner:
		fmt.Fprintf(t.outStream, "%s", ArrowSpinnerSeq[frame%(len(ArrowSpinnerSeq))])
	case EmojiArrowSpinner:
		fmt.Fprintf(t.outStream, "%s", EmojiArrowSpinnerSeq[frame%(len(EmojiArrowSpinnerSeq))])
	}
}

func (t *Terminal) writeString(s string) {
	fmt.Fprintf(t.outStream, "%s", s)
}
