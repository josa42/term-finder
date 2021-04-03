package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/josa42/term-finder/tree"
	"github.com/rivo/tview"
)

func setupLogging() func() error {
	f, _ := os.OpenFile("/tmp/term-finder.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	log.SetOutput(f)

	return f.Close
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
		SetColumns(50, 0)

	app.SetRoot(grid, true)

	treeView := tree.NewFileTree(theme)
	contentView := tree.NewContentView(theme)

	treeView.OnChanged(func(fsnode *tree.FSNode) {
		log.Printf("on changed: %s", fsnode.Name)
		contentView.SetPreview(fsnode)
		go app.Draw()
	})

	treeView.OnSelect(func(node *tree.FSNode) {
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

	treeView.OnOpen(func(node *tree.FSNode) {
		go func() {
			exec.Command("open", node.Path).Run()
		}()
	})

	treeView.Load(pwd)

	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		var x func()
		for len(treeView.AfterDraw) > 0 {
			x, treeView.AfterDraw = treeView.AfterDraw[0], treeView.AfterDraw[1:]
			x()
		}
	})

	grid.
		AddItem(treeView, 0, 0, 1, 2, 0, 0, true)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 50, true).
		AddItem(contentView, 0, 1, 1, 1, 0, 50, false)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
