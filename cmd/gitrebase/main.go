package main

import (
	"flag"
	"fmt"

	"github.com/coc1961/gitrebase/internal/gitutil"

	"github.com/coc1961/gitrebase/internal/gui"
)

func main() {
	var path = flag.String("p", "", "git repository path (mandatory)")
	flag.Parse()

	if path == nil || *path == "" {
		fmt.Println("gitrebase - param error\n ")
		flag.PrintDefaults()
		return
	}

	g := gitutil.New(*path)

	arr, commits, err := g.Log()
	if err != nil {
		fmt.Println(err)
		return
	}

	c := gui.New(arr)

	c.Run()

	i, m := c.Result()
	if i > 0 {
		err := g.Rebase(commits[i], m)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Rebase Ok!")
		}
	} else {
		fmt.Println("Exitwithout processing")
	}
}
