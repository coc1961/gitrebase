package gitutil

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func New(path string) *Git {
	return &Git{
		Path: path,
	}
}

type Git struct {
	Path string
}

func (g *Git) Log() ([]string, []*object.Commit, error) {
	var ret []string = nil
	var comm []*object.Commit = nil

	r, err := git.PlainOpen(g.Path)
	if err != nil {
		return ret, comm, err
	}

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return ret, comm, err
	}

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{
		From:  ref.Hash(),
		All:   false,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return ret, comm, err
	}

	noRL := func(s string) string {
		return strings.ReplaceAll(s, "\n", " ")
	}
	hash := func(s string) string {
		if len(s) < 8 {
			return s
		}
		return string([]byte(s)[0:7])
	}

	ret = make([]string, 0)
	comm = make([]*object.Commit, 0)
	err = cIter.ForEach(func(c *object.Commit) error {
		format := "commit %s Author: %s  Date: %s Message: %s"
		txt := fmt.Sprintf(format,
			hash(c.Hash.String()),
			c.Author.Name,
			c.Author.When.Format(object.DateFormat),
			noRL(c.Message),
		)
		ret = append(ret, txt)
		comm = append(comm, c)
		return nil
	})
	if err != nil {
		return ret, comm, err
	}

	return ret, comm, nil
}

func (g *Git) Rebase(commit *object.Commit, message string) error {
	r, err := git.PlainOpen(g.Path)
	if err != nil {
		return err
	}

	var w *git.Worktree
	w, err = r.Worktree()

	if err != nil {
		return err
	}

	err = w.Reset(&git.ResetOptions{Mode: git.SoftReset, Commit: commit.Hash})
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "user",
			Email: "mail@org.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
