package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	fmt.Fprint(os.Stderr, u)
	flag.PrintDefaults()
}

var (
	r = flag.Bool("r", false, "Random")
	l = flag.Int("l", 10000, "Limit; If equal to -1 then unlimited")
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

	cs, err := crun.Compile(format)
	if err != nil {
		log.Println(err)
		return
	}

	if *r {
		for i := 0; i != *l; i++ {
			fmt.Println(cs.Rand())
		}
	} else {
		i := 0
		cs.Range(func(s string) bool {
			if i == *l {
				return false
			}
			fmt.Println(s)
			i++
			return true
		})
	}
}
