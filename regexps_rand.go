package crun

import (
	"log"
)

type Rand interface {
	Int() int
}

func rands(r Regexps, rand Rand, buf []rune) []rune {
	if len(r) == 0 {
		return buf
	}
	reg := r[0]

	switch reg.Op {
	case OpLiteral: // matches Runes sequence
		return rands(r[1:], rand, append(buf, reg.Rune...))
	case OpRepeat: // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
		size := rand.Int()%(1+reg.Max-reg.Min) + reg.Min
		if len(reg.Sub) == 0 {
			buf = randRepeatPossibilitie(reg.Rune, rand, buf, size)
		} else {
			buf = randPossibilitie2(reg.Sub, rand, buf, size)
		}
		return rands(r[1:], rand, buf)
	case OpAlternate: // matches alternation of Subs
		i := rand.Int() % len(reg.Sub)
		buf = rands(append(reg.Sub[i], r[1:]...), rand, buf)
		return buf
	default:
		log.Printf("crun: unsupported op %v", reg.Op)
		return nil
	}
}

func randPossibilitie(regs []Regexps, rand Rand, buf []rune) []rune {
	if len(regs) == 0 {
		return buf
	}
	return randPossibilitie(regs[1:], rand, rands(regs[0], rand, buf))
}

func randPossibilitie2(regs []Regexps, rand Rand, buf []rune, size int) []rune {
	if size == 0 {
		return buf
	}
	return randPossibilitie2(regs, rand, randPossibilitie(regs, rand, buf), size-1)
}

func randRepeatPossibilitie(runes []rune, rand Rand, buf []rune, size int) []rune {
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
