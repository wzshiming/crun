package crun

func makeCalendar(runes []rune, buf []rune, ff func(r []rune)) {
	if len(buf) == cap(buf) {
		ff(buf)
		return
	}
	buf = append(buf, 0)
	for i := 0; i < len(runes); i += 2 {
		for j := runes[i]; j <= runes[i+1]; j++ {
			buf[len(buf)-1] = j
			makeCalendar(runes, buf, ff)
		}
	}
	return
}

// 产生历遍字符串
func MakeCalendar(runes []rune, min int, max int, ff func(r []rune)) {
	if len(runes) == 1 {
		runes = append(runes, runes[0])
	}
	for i := min; i <= max; i++ {
		buf := make([]rune, 0, i)
		makeCalendar(runes, buf, ff)
	}
}
