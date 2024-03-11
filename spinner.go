package main

type SpinnerKind int

const (
	PlainSpinner SpinnerKind = iota
	ColorfulSpinner
	ArrowSpinner
	EmojiArrowSpinner
)

var (
	PainSpinnerSeq       = [4]string{"|", "/", "-", "\\"}
	ColorfulSpinnerSeq   = [4]string{"\u001b[31m|\u001b[0m", "\u001b[32m/\u001b[0m", "\u001b[33m-\u001b[0m", "\u001b[34m\\\u001b[0m"}
	ArrowSpinnerSeq      = [4]string{"↑", "→", "↓", "←"}
	EmojiArrowSpinnerSeq = [8]string{"⬆️", "↗️", "➡️", "↘️", "⬇️", "↙️", "⬅️", "↖️"}
)
