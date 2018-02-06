package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/wzshiming/crun"
)

func init() {
	flag.Usage = Usage
}

func Usage() {
	u := `
Usage of crun:
       crun [Options] [regexp]
    or crun "\d{3}"
    or crun "[0-9a-z]{2}"
    or crun "(root|admin) [0-9]{1}"

Options:
	-e # Execute the generated text
	`
	fmt.Print(u)
}

func main() {
	e := flag.Bool("e", false, "exec")

	flag.Parse()

	format := strings.Join(flag.Args(), " ")
	if format == "" {
		flag.Usage()
		return
	}

	if *e {
		crun.NewSyntax(format).Makes(func(s []rune) {
			ss := strings.Split(string(s), " ")
			out, err := exec.Command(ss[0], ss[1:]...).Output()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(out))
			}
		})
	} else {
		crun.NewSyntax(format).Makes(func(s []rune) {
			fmt.Println(string(s))
		})
	}

}
