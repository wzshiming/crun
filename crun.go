package crun

import (
	"fmt"
	"math"
	"regexp/syntax"
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

// Size The number of possibilities that can match regularity
func (r Regexps) Size() int {
	s := 0
	r.size(&s)
	return s
}

func (r Regexps) size(s *int) {
	if len(r) == 0 {
		return
	}
	if *s == 0 {
		*s += 1
	}

	reg := r[0]

	switch reg.Op {
	case OpLiteral: // matches Runes sequence
		r[1:].size(s)
		return
	case OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		ru := reg.Rune
		if len(reg.Sub) != 0 {
			ru = reg.Sub[0][0].Rune
		}
		*s *= SizePossibilities(ru, reg.Min, reg.Max)
		r[1:].size(s)
		return
	case OpAlternate: // matches alternation of Subs
		for _, v := range reg.Sub {
			append(v, r[1:]...).size(s)
		}
		return
	default:
		return
	}

}

// Range all possibilities
func (r Regexps) Range(f func([]rune) bool) bool {
	buf := []rune{}
	return r.ranges(buf, 0, f)
}

func (r Regexps) ranges(buf []rune, off int, f func([]rune) bool) bool {
	if len(r) == 0 {
		return f(buf)
	}
	reg := r[0]
	ff := func(s []rune) bool {
		l := len(s) - (len(buf) - off)
		if l > 0 {
			buf = append(buf, make([]rune, l)...)
		} else if l < 0 {
			buf = buf[:len(buf)+l]
		}
		copy(buf[off:], s)

		return r[1:].ranges(buf, off+len(s), f)
	}

	switch reg.Op {
	case OpLiteral: // matches Runes sequence
		return ff(reg.Rune)
	case OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		ru := reg.Rune
		if len(reg.Sub) != 0 {
			ru = reg.Sub[0][0].Rune
		}
		return RangePossibilities(ru, reg.Min, reg.Max, ff)
	case OpAlternate: // matches alternation of Subs
		for _, v := range reg.Sub {
			if !append(v, r[1:]...).ranges(buf, off, f) {
				return false
			}
		}
		return true
	default:
		fmt.Printf("Unsupported op %v", reg.Op)
		return false
	}
}

func NewSyntax(s string) Regexps {
	reg, _ := syntax.Parse(s, syntax.Perl)
	return NewSyntaxByRegexp(reg)
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

func rangePossibilities(runes []rune, buf []rune, ff func(r []rune) bool) bool {
	if len(buf) == cap(buf) {
		return ff(buf)
	}
	buf = append(buf, 0)
	for i := 0; i < len(runes); i += 2 {
		for j := runes[i]; j <= runes[i+1]; j++ {
			buf[len(buf)-1] = j
			if !rangePossibilities(runes, buf, ff) {
				return false
			}
		}
	}
	return true
}

// Range all possibilities
func RangePossibilities(runes []rune, min int, max int, ff func(r []rune) bool) bool {
	if len(runes) == 1 {
		runes = append(runes, runes[0])
	}
	for i := min; i <= max; i++ {
		buf := make([]rune, 0, i)
		if !rangePossibilities(runes, buf, ff) {
			return false
		}
	}
	return true
}

func SizePossibilities(runes []rune, min int, max int) int {
	if len(runes) == 1 {
		runes = append(runes, runes[0])
	}
	sum := 0
	for i := 0; i < len(runes); i += 2 {
		sum += int(runes[i+1]-runes[i]) + 1
	}

	r := 0
	for i := min; i <= max; i++ {
		r += int(math.Pow(float64(sum), float64(i)))
	}
	return r
}
