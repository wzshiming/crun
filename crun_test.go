package crun

import (
	"regexp"
	"testing"
)

func TestCrun(t *testing.T) {

	tests := []struct {
		reg  string
		want int
	}{
		{``, 0},
		{` `, 1},
		{` ?`, 2},
		{`\d`, 10},
		{`\w`, 63},
		{`\s`, 5},
		{`(two){2}`, 1},
		{`\d?`, 11},
		{`\d{2}`, 100},
		{`\w{2}`, 3969},
		{`\d{0,2}`, 111},
		{`[A-Z0-9]`, 36},
		{`\d{2}\w{2}`, 396900},
		{`\d{2}|\w{2}`, 4069},
		{`hello|bye{1,2}`, 3},
		{`(hello|bye) (1|3|world)`, 6},
	}
	for _, tt := range tests {
		req, err := regexp.Compile(tt.reg)
		if err != nil {
			t.Error(err)
			continue
		}

		cr := NewSyntax(tt.reg)

		max := 10
		for i := 0; i != max && i != tt.want; i++ {
			if got := cr.Rand(); !req.MatchString(got) {
				t.Errorf("`%s` Rand() = %v", tt.reg, got)
			}
		}

		i := 0
		cr.Range(func(got string) bool {
			if !req.MatchString(got) {
				t.Errorf("`%s` Range() = %v", tt.reg, got)
			}
			i++
			return i != max && i != tt.want
		})

		if got := cr.Size(); got != tt.want {
			t.Errorf("`%s` Size() = %v, want %v", tt.reg, got, tt.want)
		}
	}
}
