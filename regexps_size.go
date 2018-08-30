package crun

import (
	"math"
)

func size(r Regexps, s *int) {
	if len(r) == 0 {
		return
	}
	if *s == 0 {
		*s++
	}

	reg := r[0]

	switch reg.Op {
	case OpLiteral: // matches Runes sequence
		size(r[1:], s)
		return
	case OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		ru := reg.Rune
		if len(reg.Sub) != 0 {
			ru = reg.Sub[0][0].Rune
		}
		*s *= sizePossibilities(ru, reg.Min, reg.Max)
		size(r[1:], s)
		return
	case OpAlternate: // matches alternation of Subs
		for _, v := range reg.Sub {
			size(append(v, r[1:]...), s)
		}
		return
	default:
		return
	}

}

// sizePossibilities returns size of all possibilities.
func sizePossibilities(runes []rune, min int, max int) int {
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
