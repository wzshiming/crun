package crun

type Op uint8

const (
	_ Op = 1 + iota
	OpLiteral
	OpRepeat
	OpAlternate
)
