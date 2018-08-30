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
		if len(reg.Sub) == 0 {
			*s *= sizeRepeatPossibilities(reg.Rune, reg.Min, reg.Max)
		} else {
			*s *= sizePossibilities(reg.Sub, reg.Min, reg.Max)
		}
		size(r[1:], s)
		return
	case OpAlternate: // matches alternation of Subs
		sum := 0
		for _, v := range reg.Sub {
			ns := 0
			size(append(v, r[1:]...), &ns)
			sum += ns
		}
		*s *= sum
		return
	default:
		return
	}
}

// sizePossibilities returns size of all possibilities.
func sizePossibilities(regs []Regexps, min int, max int) int {
	sum := 1
	for _, reg := range regs {
		i := 0
		size(reg, &i)
		sum *= i
	}

	// Geometric series formula
	// if q != 1 then x1*(1-q^n)/(1-q) or (x1-xn*q)/(1-q)
	// if q == 1 then x1*n
	q := float64(sum)
	x1 := math.Pow(q, float64(min))
	n := 1 + max - min
	if n == 1 {
		return int(x1)
	}
	if sum != 1 {
		sum = int(x1 * (1 - math.Pow(q, float64(n))) / (1 - q))
	} else {
		sum = int(x1 * float64(n))
	}
	return sum
}

// sizeRepeatPossibilities returns size of all possibilities.
func sizeRepeatPossibilities(runes []rune, min int, max int) int {
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
