package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/josa42/term-finder/tree"
	"github.com/rivo/tview"
)

func get(node *tview.TreeNode) *tree.FSNode {
	ref := node.GetReference()
	if ref == nil {
		return nil
	}
	return ref.(*tree.FSNode)
}

func setupLogging() func() error {
	f, _ := os.OpenFile("/tmp/term-finder.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	log.SetOutput(f)

	return f.Close
}

func findParent(root, curr *tview.TreeNode) *tview.TreeNode {
	var currParent *tview.TreeNode
	root.Walk(func(node, parent *tview.TreeNode) bool {
		if node == curr {
			currParent = parent
			return false
		}
		return true
	})

	return currParent
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
		SetRows(0).
		SetColumns(50, 0)

	app.SetRoot(grid, true)

	ft := tree.NewFileTree(theme)
	ft.Load(pwd)

	treeView := ft.GetView()

	contentView := tview.NewTextView()
	contentView.SetBorderPadding(0, 0, 2, 2)
	contentView.SetBackgroundColor(theme.ContentBackground)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 0, true)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 50, true).
		AddItem(contentView, 0, 1, 1, 1, 0, 50, false)

	ft.OnChanged(func(fsnode *tree.FSNode) {
		if !fsnode.IsDir && fsnode.Size < 400_000 {
			contentView.SetText(fsnode.Path)
			contentView.SetTitle(fsnode.Path)
			content, _ := ioutil.ReadFile(fsnode.Path)
			contentView.SetText(fmt.Sprintf("%s\n%d\n---\n%s", fsnode.Path, fsnode.Size, string(content)))

		} else {
			contentView.SetText(fmt.Sprintf("%s\n%d", fsnode.Path, fsnode.Size))
		}

		contentView.ScrollTo(0, 0)
		go app.Draw()
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
