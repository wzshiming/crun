package crun

import (
	"fmt"
)

func ranges(r Regexps, buf []rune, off int, f func([]rune) bool) bool {
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

		return ranges(r[1:], buf, off+len(s), f)
	}

	switch reg.Op {
	case OpLiteral: // matches Runes sequence
		return ff(reg.Rune)
	case OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		ru := reg.Rune
		if len(reg.Sub) != 0 {
			ru = reg.Sub[0][0].Rune
		}
		return rangePossibilities(ru, reg.Min, reg.Max, ff)
	case OpAlternate: // matches alternation of Subs
		for _, v := range reg.Sub {
			if !ranges(append(v, r[1:]...), buf, off, f) {
				return false
			}
		}
		return true
	default:
		fmt.Printf("Unsupported op %v", reg.Op)
		return false
	}
}

func rangePossibilitie(runes []rune, buf []rune, ff func(r []rune) bool) bool {
	if len(buf) == cap(buf) {
		return ff(buf)
	}
	buf = append(buf, 0)
	for i := 0; i < len(runes); i += 2 {
		for j := runes[i]; j <= runes[i+1]; j++ {
			buf[len(buf)-1] = j
			if !rangePossibilitie(runes, buf, ff) {
				return false
			}
		}
	}
	return true
}

// rangePossibilities range all possibilities.
func rangePossibilities(runes []rune, min int, max int, ff func(r []rune) bool) bool {
	if len(runes) == 1 {
		runes = append(runes, runes[0])
	}
	for i := min; i <= max; i++ {
		buf := make([]rune, 0, i)
		if !rangePossibilitie(runes, buf, ff) {
			return false
		}
	}
	return true
}
