package crun

import (
	"fmt"
	"regexp/syntax"
	"strconv"
)

// MoreTimes Maximum omitted default value
const MoreTimes = 18

type Regexp struct {
	Op       Op
	Sub      []Regexps
	Rune     []rune
	Min, Max int
}

type Regexps []*Regexp

func Compile(str string) (Regexps, error) {
	reg, err := syntax.Parse(str, syntax.Perl)
	if err != nil {
		return nil, fmt.Errorf("crun: Compile(`%s`): %s", strconv.Quote(str), err.Error())
	}
	return NewSyntaxByRegexp(reg), nil
}

func MustCompile(str string) Regexps {
	reg, err := Compile(str)
	if err != nil {
		panic(err)
	}
	return reg
}

func NewSyntax(str string) Regexps {
	reg, err := Compile(str)
	if err != nil {
		fmt.Println(err)
	}
	return reg
}

func NewSyntaxByRegexp(reg *syntax.Regexp) (out Regexps) {
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
			Rune: []rune{1, 127},
			Max:  1,
			Min:  1,
		})
	case syntax.OpAnyChar: // matches any character
		ff(&Regexp{
			Op:   OpRepeat,
			Rune: []rune{1, 127},
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
			ff(NewSyntaxByRegexp(v)...)
		}
	case syntax.OpStar: // matches Sub[0] zero or more times
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, NewSyntaxByRegexp(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: MoreTimes,
			Min: 0,
		})
	case syntax.OpPlus: // matches Sub[0] one or more times
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, NewSyntaxByRegexp(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: MoreTimes,
			Min: 1,
		})
	case syntax.OpQuest: // matches Sub[0] zero or one times
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, NewSyntaxByRegexp(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: 1,
			Min: 0,
		})
	case syntax.OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, NewSyntaxByRegexp(v))
		}
		ff(&Regexp{
			Op:  OpRepeat,
			Sub: sub,
			Max: reg.Max,
			Min: reg.Min,
		})
	case syntax.OpConcat: // matches concatenation of Subs
		for _, v := range reg.Sub {
			ff(NewSyntaxByRegexp(v)...)
		}
	case syntax.OpAlternate: // matches alternation of Subs
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, NewSyntaxByRegexp(v))
		}
		ff(&Regexp{
			Op:  OpAlternate,
			Sub: sub,
		})
	default:
		fmt.Printf("Unsupported op %v", reg.Op)
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
	return ranges(r, []rune{}, 0, func(s []rune) bool {
		return f(string(s))
	})
}

// Rand possibilities
func (r Regexps) Rand() string {
	return string(rands(r, []rune{}))
}
