package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/wzshiming/crun"
)

func init() {
	flag.Usage = usage
}

func usage() {
	u := `
Usage of crun:
       crun [Options] [regexp]
    or crun "\d{3}"
    or crun "[0-9a-z]{2}"
    or crun "(root|admin) [0-9]{1}"

Options:
`
	fmt.Print(u)
	flag.PrintDefaults()
}

var (
	r = flag.Bool("r", false, "Random")
	l = flag.Int("l", 10, "Limit; If equal to -1 then unlimited")
)

func init() {
	flag.Parse()
}

func main() {

	format := strings.Join(flag.Args(), " ")
	if format == "" {
		flag.Usage()
		return
	}

	if *r {
		rand.Seed(time.Now().UnixNano())
		cs := crun.NewSyntax(format)
		for i := 0; i != *l; i++ {
			fmt.Println(cs.Rand())
		}
	} else {
		i := 0
		crun.NewSyntax(format).Range(func(s string) bool {
			if i == *l {
				return false
			}
			fmt.Println(s)
			i++
			return true
		})
	}
}
