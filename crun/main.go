package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wzshiming/crun"
)

func init() {

	flag.Parse()
}

func main() {
	format := strings.Join(flag.Args(), " ")
	if format == "" {
		s := os.Args[0]
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", s)
		fmt.Fprintf(os.Stderr, "\t   %s <regexp>\n", s)
		fmt.Fprintf(os.Stderr, "\tor %s \"\\d{3}\"\n", s)
		fmt.Fprintf(os.Stderr, "\tor %s \"[0-9a-z]{2}\"\n", s)
		fmt.Fprintf(os.Stderr, "\tor %s \"(root|admin) [0-9]{1}\"\n", s)

		return
	}
	crun.NewSyntax(format).Makes(func(s []rune) {
		fmt.Println(string(s))
	})
}
