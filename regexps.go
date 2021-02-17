package crun

import (
	"fmt"
	"log"
	"math/rand"
	"regexp/syntax"
	"time"
)

// MoreTimes Maximum omitted default value
const MoreTimes = 18

// Regexp syntax tree translated from regexp/syntax
type Regexp struct {
	Op       Op
	Sub      []Regexps
	Rune     []rune
	Min, Max int
}

// Regexps syntax tree translated from regexp/syntax
type Regexps []*Regexp

// Compile parses a regular expression and returns.
func Compile(str string) (Regexps, error) {
	reg, err := syntax.Parse(str, syntax.Perl)
	if err != nil {
		return nil, fmt.Errorf("crun: Compile(%q): %w", str, err)
	}
	return NewRegexps(reg), nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular
// expressions.
func MustCompile(str string) Regexps {
	reg, err := Compile(str)
	if err != nil {
		panic(err)
	}
	return reg
}

// NewRegexps returns regexps translated from regexp/syntax
func NewRegexps(reg *syntax.Regexp) (out Regexps) {
	return std.NewRegexps(reg)
}

var std = &Optional{
	MoreTimes:    MoreTimes,
	AnyCharNotNL: []rune{33, 126},
	AnyChar:      []rune{33, 126},
}

// Optional is optional related option for regexps
type Optional struct {
	MoreTimes    int
	AnyCharNotNL []rune
	AnyChar      []rune
	Rand         Rand
}

// NewRegexps returns regexps translated from regexp/syntax
func (o *Optional) NewRegexps(reg *syntax.Regexp) (out Regexps) {
	ff := func(rs ...*Regexp) {
		out = append(out, rs...)
	}

	switch reg.Op {
	case syntax.OpNoMatch: // matches no strings
	case syntax.OpEmptyMatch: // matches empty string
	case syntax.OpLiteral: // matches Runes sequence
		ff(&Regexp{
			Op:   OpLiteral,
			Rune: reg.Rune,
		})
	case syntax.OpCharClass: // matches Runes interpreted as range pair list
		ff(&Regexp{
			Op:   OpRepeat,
			Rune: reg.Rune,
			Max:  1,
			Min:  1,
		})
	case syntax.OpAnyCharNotNL: // matches any character except newline
		ff(&Regexp{
			Op:   OpRepeat,
			Rune: o.AnyCharNotNL,
			Max:  1,
			Min:  1,
		})
	case syntax.OpAnyChar: // matches any character
		ff(&Regexp{
			Op:   OpRepeat,
			Rune: o.AnyChar,
			Max:  1,
			Min:  1,
		})
	case syntax.OpBeginLine: // matches empty string at beginning of line
	case syntax.OpEndLine: // matches empty string at end of line
	case syntax.OpBeginText: // matches empty string at beginning of text
	case syntax.OpEndText: // matches empty string at end of text
	case syntax.OpWordBoundary: // matches word boundary `\b`
	case syntax.OpNoWordBoundary: // matches word non-boundary `\B`
	case syntax.OpCapture: // capturing subexpression with index Cap, optional name Name
		for _, v := range reg.Sub {
			ff(o.NewRegexps(v)...)
		}
	case syntax.OpStar: // matches Sub[0] zero or more times
		sub := make([]Regexps, 0, len(reg.Sub))
		for _, v := range reg.Sub {
			sub = append(sub, o.NewRegexps(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: o.MoreTimes,
			Min: 0,
		})
	case syntax.OpPlus: // matches Sub[0] one or more times
		sub := make([]Regexps, 0, len(reg.Sub))
		for _, v := range reg.Sub {
			sub = append(sub, o.NewRegexps(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: o.MoreTimes,
			Min: 1,
		})
	case syntax.OpQuest: // matches Sub[0] zero or one times
		sub := make([]Regexps, 0, len(reg.Sub))
		for _, v := range reg.Sub {
			sub = append(sub, o.NewRegexps(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: 1,
			Min: 0,
		})
	case syntax.OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		sub := make([]Regexps, 0, len(reg.Sub))
		for _, v := range reg.Sub {
			sub = append(sub, o.NewRegexps(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: reg.Max,
			Min: reg.Min,
		})
	case syntax.OpConcat: // matches concatenation of Subs
		for _, v := range reg.Sub {
			ff(o.NewRegexps(v)...)
		}
	case syntax.OpAlternate: // matches alternation of Subs
		sub := make([]Regexps, 0, len(reg.Sub))
		for _, v := range reg.Sub {
			sub = append(sub, o.NewRegexps(v))
		}
		ff(&Regexp{
			Op:  OpAlternate,
			Sub: sub,
		})
	default:
		log.Printf("crun: unsupported op %v", reg.Op)
	}
	return out
}

// Size The number of possibilities that can match regularity
func (r Regexps) Size() int {
	s := 0
	size(r, &s)
	return s
}

// Range all possibilities
func (r Regexps) Range(f func(string) bool) bool {
	return r.RangeWithRunes(func(s []rune) bool {
		return f(string(s))
	})
}

// RangeWithRuns all possibilities
func (r Regexps) RangeWithRunes(f func([]rune) bool) bool {
	return ranges(r, []rune{}, 0, func(s []rune) bool {
		return f(s)
	})
}

// Rand possibilities
func (r Regexps) Rand() string {
	return string(r.RandWithRunes())
}

// RandWithRunes possibilities
func (r Regexps) RandWithRunes() []rune {
	return rands(r, stdRandSource, []rune{})
}

// RandSource possibilities
func (r Regexps) RandSource(rand Rand) string {
	return string(r.RandSourceWithRunes(rand))
}

// RandSourceWithRunes possibilities
func (r Regexps) RandSourceWithRunes(rand Rand) []rune {
	return rands(r, rand, []rune{})
}

var stdRandSource = rand.New(rand.NewSource(time.Now().UnixNano()))
