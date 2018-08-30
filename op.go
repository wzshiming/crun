package crun

//go:generate stringer -type Op

// An Op is a single regular expression operator.
type Op uint8

// Operators are listed in precedence order, tightest binding to weakest.
const (
	_ Op = iota
	OpLiteral
	OpRepeat
	OpAlternate
)
