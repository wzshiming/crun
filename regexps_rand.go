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
		size := rand.Int()%(1+reg.Max-reg.Min) + reg.Min
		if len(reg.Sub) == 0 {
			buf = randRepeatPossibilitie(reg.Rune, buf, size)
		} else {
			buf = randPossibilitie2(reg.Sub, buf, size)
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

func randPossibilitie(regs []Regexps, buf []rune) []rune {
	if len(regs) == 0 {
		return buf
	}
	return randPossibilitie(regs[1:], rands(regs[0], buf))
}

func randPossibilitie2(regs []Regexps, buf []rune, size int) []rune {
	if size == 0 {
		return buf
	}
	return randPossibilitie2(regs, randPossibilitie(regs, buf), size-1)
}

func randRepeatPossibilitie(runes []rune, buf []rune, size int) []rune {
	if len(runes) == 1 {
		return append(buf, runes[0])
	}
	for i := 0; i != size; i++ {
		index := (rand.Int() % (len(runes) / 2)) << 1
		curr := rune(rand.Int()%int(1+runes[index+1]-runes[index])) + runes[index]
		buf = append(buf, curr)
	}
	return buf
}
