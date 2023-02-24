package gui

import (
	"fmt"
	"log"
	"strings"

	gocui "github.com/jroimartin/gocui"
)

func New(data []string) *Console {
	return &Console{
		Data: data,
	}
}

type Console struct {
	actual  int
	mark    int
	Data    []string
	message string
	v       *gocui.View
}

func (c *Console) Result() (a int, b string) {
	if c.mark > 0 {
		if c.mark+1 < len(c.Data) {
			return c.mark + 1, c.message
		}
	}
	return c.mark, c.message
}

func (c *Console) Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(c.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, c.up); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, c.down); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, c.markSelection); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, c.process); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
		return
	}
}

func (c *Console) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("tit", 0, 0, maxX, 2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Frame = false
	normal := "\033[0m"
	color := fmt.Sprintf("\033[3%d;%dm", 2, 2)
	color1 := fmt.Sprintf("\033[3%d;%dm", 2, 7)
	fmt.Fprintf(v, "%s git rebase%s - %s[Ctrl-C]%s Salir sin Procesar,  %s[Ctrl-Space]%s Marca los Commits a Unificar,  %sEnter%s Procesa el Rebase \n", color1, normal, color, normal, color, normal, color, normal)

	if v, err := g.SetView("git", 0, 3, maxX-1, maxY-1); err != nil {
		v.Frame = true
		v.Wrap = false
		v.Autoscroll = false
		if err != gocui.ErrUnknownView {
			return err
		}

		c.v = v
		c.print()
	}
	return nil
}

func (c *Console) print() {
	c.v.Clear()

	color := func(s string) string {
		c := "\033[33m"
		s = strings.ReplaceAll(s, "[", "["+c)
		s = strings.ReplaceAll(s, "]", "]\033[0m")
		return s
	}
	n := 0
	for _, s := range c.Data {
		if n == c.actual {
			fmt.Fprintln(c.v, "\033[0m\033[1;7m "+s+"\033[0m")
		} else {
			if n <= c.mark && c.mark > 0 {
				fmt.Fprintln(c.v, "\033[0m\033[33m "+s+"\033[0m")
			} else {
				fmt.Fprintln(c.v, " "+color(s)+"\033[0m")
			}
		}
		n++
	}
}

func (c *Console) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Console) markSelection(g *gocui.Gui, v *gocui.View) error {
	c.mark = c.actual
	c.print()
	return nil
}

func (c *Console) up(g *gocui.Gui, v *gocui.View) error {
	if c.actual-1 >= 0 {
		c.actual--
		c.print()
	}
	return nil
}

func (c *Console) down(g *gocui.Gui, v *gocui.View) error {
	if c.actual+1 < len(c.Data) {
		c.actual++
		c.print()

	}
	return nil
}

func (c *Console) process(g *gocui.Gui, v *gocui.View) error {

	if c.mark == 0 {
		return nil
	}

	iv, err := g.SetView("Input", 0, 0, 100, 2)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create input view:", err)
		return err
	}
	iv.Title = "Commit Message"
	iv.FgColor = gocui.ColorYellow
	iv.Editable = true
	err = iv.SetCursor(0, 0)
	if err != nil {
		log.Println("Failed to set cursor:", err)
		return err
	}

	err = g.SetKeybinding("Input", gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		c.message = iv.Buffer()
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	_, err = g.SetCurrentView("Input")
	if err != nil {
		log.Println("Cannot set focus to input view:", err)
		return err
	}

	return nil
}
