package crun

//go:generate stringer -type Op
type Op uint8

const (
	_ Op = iota
	OpLiteral
	OpRepeat
	OpAlternate
)
