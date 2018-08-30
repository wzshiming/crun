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

var (
	e = flag.Bool("e", false, "exec")
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

	if *e {
		crun.NewSyntax(format).Range(func(s crun.String) bool {
			ss := strings.Split(s.String(), " ")
			out, err := exec.Command(ss[0], ss[1:]...).Output()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(out))
			}
			return true
		})
	} else {
		crun.NewSyntax(format).Range(func(s crun.String) bool {
			fmt.Println(s)
			return true
		})
	}

}
