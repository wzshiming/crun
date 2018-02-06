package crun

import (
	"fmt"
	"regexp/syntax"
)

const MoreTimes = 18

type Regexp struct {
	Op       syntax.Op
	Sub      []Regexps
	Rune     []rune
	Min, Max int
}

type Regexps []*Regexp

func (r Regexps) Makes(f func([]rune)) {
	buf := []rune{}
	r.makes(buf, 0, f)
}

func (r Regexps) makes(buf []rune, off int, f func([]rune)) {
	if len(r) == 0 {
		f(buf)
		return
	}
	reg := r[0]
	ff := func(s []rune) {
		l := len(s) - (len(buf) - off)
		if l > 0 {
			buf = append(buf, make([]rune, l)...)
		} else if l < 0 {
			buf = buf[:len(buf)+l]
		}
		copy(buf[off:], s)

		r[1:].makes(buf, off+len(s), f)

	}
	switch reg.Op {

	case syntax.OpLiteral: // matches Runes sequence
		ff(reg.Rune)
	case syntax.OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		ru := reg.Rune
		if len(reg.Sub) != 0 {
			ru = reg.Sub[0][0].Rune
		}
		MakeCalendar(ru, reg.Min, reg.Max, ff)

	case syntax.OpAlternate: // matches alternation of Subs
		for _, v := range reg.Sub {
			append(v, r[1:]...).makes(buf, off, f)
		}
	default:
		fmt.Printf("Unsupported op %v", reg.Op)
	}
}

func NewSyntax(s string) Regexps {
	reg, _ := syntax.Parse(s, syntax.Perl)
	return newSyntax(reg)
}

func newSyntax(reg *syntax.Regexp) (out Regexps) {
	ff := func(rs ...*Regexp) {
		out = append(out, rs...)
	}
	switch reg.Op {
	case syntax.OpNoMatch: // matches no strings
	case syntax.OpEmptyMatch: // matches empty string
	case syntax.OpLiteral: // matches Runes sequence
		ff(&Regexp{
			Op:   syntax.OpLiteral,
			Rune: reg.Rune,
		})
	case syntax.OpCharClass: // matches Runes interpreted as range pair list
		ff(&Regexp{
			Op:   syntax.OpRepeat,
			Rune: reg.Rune,
			Max:  1,
			Min:  1,
		})
	case syntax.OpAnyCharNotNL: // matches any character except newline
		ff(&Regexp{
			Op:   syntax.OpRepeat,
			Rune: []rune{1, 127},
			Max:  1,
			Min:  1,
		})
	case syntax.OpAnyChar: // matches any character
		ff(&Regexp{
			Op:   syntax.OpRepeat,
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
			ff(newSyntax(v)...)
		}
	case syntax.OpStar: // matches Sub[0] zero or more times
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, newSyntax(v))
		}
		ff(&Regexp{
			Op:  syntax.OpRepeat,
			Sub: sub,
			Max: MoreTimes,
			Min: 0,
		})
	case syntax.OpPlus: // matches Sub[0] one or more times
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, newSyntax(v))
		}
		ff(&Regexp{
			Op:  syntax.OpRepeat,
			Sub: sub,
			Max: MoreTimes,
			Min: 1,
		})
	case syntax.OpQuest: // matches Sub[0] zero or one times
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, newSyntax(v))
		}
		ff(&Regexp{
			Op:  syntax.OpRepeat,
			Sub: sub,
			Max: 1,
			Min: 0,
		})
	case syntax.OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, newSyntax(v))
		}
		ff(&Regexp{
			Op:  syntax.OpRepeat,
			Sub: sub,
			Max: reg.Max,
			Min: reg.Min,
		})
	case syntax.OpConcat: // matches concatenation of Subs
		for _, v := range reg.Sub {
			ff(newSyntax(v)...)
		}
	case syntax.OpAlternate: // matches alternation of Subs
		sub := []Regexps{}
		for _, v := range reg.Sub {
			sub = append(sub, newSyntax(v))
		}
		ff(&Regexp{
			Op:  syntax.OpAlternate,
			Sub: sub,
		})
	default:
		fmt.Printf("Unsupported op %v", reg.Op)
	}
	return out
}
