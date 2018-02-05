package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wzshiming/crun"
)

func init() {
	flag.Usage = Usage
}

func Usage() {
	s := filepath.Base(os.Args[0])
	fmt.Printf("Usage of %s:\n", s)
	fmt.Printf("\t   %s [Options] <regexp>\n", s)
	fmt.Printf("\t    -e Execute the generated text\n")
	fmt.Printf("\tor %s \"\\d{3}\"\n", s)
	fmt.Printf("\tor %s \"[0-9a-z]{2}\"\n", s)
	fmt.Printf("\tor %s \"(root|admin) [0-9]{1}\"\n", s)
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
