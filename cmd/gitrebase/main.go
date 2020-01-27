package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/coc1961/gitrebase/internal/gitutil"

	"github.com/coc1961/gitrebase/internal/gui"
)

func main() {
	var path = flag.String("p", "", "git repository path (mandatory)")
	var list = flag.Bool("l", false, "commit list (optional)")
	var pcommit = flag.String("commit", "", "commit to squash (optional)")
	var pmsg = flag.String("msg", "", "commit message (optional)")
	flag.Parse()

	if path == nil || *path == "" {
		fmt.Println("gitrebase - param error\n ")
		flag.PrintDefaults()
		return
	}

	g := gitutil.New(*path)

	if list != nil && *list {
		arr, err := g.List()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, s := range arr {
			c := "\033[33m"
			s = strings.ReplaceAll(s, "[", "["+c)
			s = strings.ReplaceAll(s, "]", "]\033[0m")

			fmt.Println(s)
		}
		return
	}

	if pcommit != nil && *pcommit != "" {
		if pmsg == nil || *pmsg == "" {
			fmt.Println("if automatic squash is performed, the commit message must be indicated")
			return
		}
		fmt.Println("Processing, one moment please!...")
		err := g.Squash(*pcommit, *pmsg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Rebase Ok!")
		return

	}

	arr, commits, err := g.Log()
	if err != nil {
		fmt.Println(err)
		return
	}

	c := gui.New(arr)

	c.Run()

	i, m := c.Result()
	if i > 0 && m != "" {
		fmt.Println("Processing, one moment please!...")
		err := g.Rebase(commits[i], m)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Rebase Ok!")
		}
	} else {
		fmt.Println("Exit without processing")
	}
}
