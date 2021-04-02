package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/josa42/term-finder/tree"
	"github.com/rivo/tview"
)

func setupLogging() func() error {
	f, _ := os.OpenFile("/tmp/term-finder.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	log.SetOutput(f)

	return f.Close
}

func formatPath(p string) string {
	dir := filepath.Dir(p)
	base := filepath.Base(p)

	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, home) {
		dir = strings.Replace(dir, home, "~", 1)
	}

	if dir == "/" {
		return fmt.Sprintf("[blue]/[normal]%s", base)
	}

	return fmt.Sprintf("[blue]%s/[normal]%s", dir, base)
}

// Show a navigable tree view of the current directory.
func main() {
	defer setupLogging()()

	pwd, _ := os.Getwd()
	log.Printf("open: %s", pwd)

	theme := tree.GetTheme()

	app := tview.NewApplication()
	app.EnableMouse(true)

	grid := tview.NewGrid().
		SetBordersColor(theme.Border).
		SetBorders(theme.Border != 0).
		SetRows(3, 0).
		SetColumns(50, 0)

	app.SetRoot(grid, true)

	topbar := tview.NewTextView()
	topbar.SetBorderPadding(0, 0, 1, 1)
	topbar.SetBorder(true)
	topbar.SetBorderColor(theme.TopbarBorder)
	topbar.SetDynamicColors(true)
	topbar.SetRegions(true)
	topbar.SetText(formatPath(pwd))

	contentView := tview.NewTextView()
	contentView.SetBorderPadding(0, 0, 2, 2)
	contentView.SetBackgroundColor(theme.ContentBackground)

	ft := tree.NewFileTree(theme)
	ft.OnChanged(func(fsnode *tree.FSNode) {
		log.Printf("on changed: %s", fsnode.Name)
		if !fsnode.IsDir && fsnode.Size < 400_000 {
			contentView.SetText(fsnode.Path)
			contentView.SetTitle(fsnode.Path)
			content, _ := ioutil.ReadFile(fsnode.Path)
			contentView.SetText(string(content))

			topbar.SetText(formatPath(fsnode.Path))

		} else {
			contentView.SetText("")
		}

		topbar.SetText(formatPath(fsnode.Path))
		contentView.ScrollTo(0, 0)
		go app.Draw()
	})

	ft.OnSelect(func(node *tree.FSNode) {
		if !node.IsDir {
			app.Suspend(func() {
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vim"
				}
				cmd := exec.Command(editor, node.Path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Run()
			})
		}
	})

	ft.OnOpen(func(node *tree.FSNode) {
		go func() {
			exec.Command("open", node.Path).Run()
			// app.Suspend(func() {
			// 	cmd := exec.Command("open", node.Path)
			// 	cmd.Stdin = os.Stdin
			// 	cmd.Stdout = os.Stdout
			// 	cmd.Run()
			// })
		}()
	})

	ft.Load(pwd)

	treeView := ft.GetView()

	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		var x func()
		for len(ft.AfterDraw) > 0 {
			x, ft.AfterDraw = ft.AfterDraw[0], ft.AfterDraw[1:]
			x()
		}
	})

	grid.
		AddItem(treeView, 0, 0, 2, 2, 0, 0, true)

	grid.
		AddItem(treeView, 0, 0, 2, 1, 0, 50, true).
		AddItem(topbar, 0, 1, 1, 1, 0, 50, false).
		AddItem(contentView, 1, 1, 1, 1, 0, 50, false)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
