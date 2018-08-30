package crun

import (
	"fmt"
	"math/rand"
)

func rands(r Regexps, buf []rune) []rune {
	if len(r) == 0 {
		return buf
	}
	reg := r[0]

	switch reg.Op {
	case OpLiteral: // matches Runes sequence
		return rands(r[1:], append(buf, reg.Rune...))
	case OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		ru := reg.Rune
		if len(reg.Sub) != 0 {
			ru = reg.Sub[0][0].Rune
		}

		curr := rand.Int()%(1+reg.Max-reg.Min) + reg.Min
		for i := 0; i != curr; i++ {
			buf = randPossibilitie(ru, buf)
		}
		return rands(r[1:], buf)
	case OpAlternate: // matches alternation of Subs
		i := rand.Int() % len(reg.Sub)
		buf = rands(append(reg.Sub[i], r[1:]...), buf)
		return buf
	default:
		fmt.Printf("Unsupported op %v", reg.Op)
		return nil
	}
}

func randPossibilitie(runes []rune, buf []rune) []rune {
	if len(runes) == 1 {
		return append(runes, runes[0])
	}
	i := rand.Int() % (len(runes) / 2)
	curr := rune(rand.Int()%int(1+runes[i+1]-runes[i])) + runes[i]
	return append(buf, curr)
}
