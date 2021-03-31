package main

import (
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

// Show a navigable tree view of the current directory.
func main() {
	defer setupLogging()()

	pwd, _ := os.Getwd()
	log.Printf("open: %s", pwd)

	theme := tree.GetTheme()

	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetBordersColor(theme.Border).
		SetBorders(theme.Border != 0).
		SetRows(0).
		SetColumns(50, 0)

	app.SetRoot(grid, true)

	root := tree.NewNode(pwd, ".", true)
	root.Expand()
	get(root).ReadChildren()

	treeView := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root).
		SetTopLevel(1)

	treeView.SetBorderPadding(0, 0, 2, 2)
	treeView.SetGraphicsColor(theme.SidebarLines)
	treeView.SetBackgroundColor(theme.SidebarBackground)

	contentView := tview.NewTextView()
	contentView.SetBorderPadding(0, 0, 2, 2)
	contentView.SetBackgroundColor(theme.ContentBackground)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 0, true)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 50, true).
		AddItem(contentView, 0, 1, 1, 1, 0, 50, false)

	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		if fsnode := get(node); fsnode != nil {

			if fsnode.Name == "." {

			} else if node.IsExpanded() {
				node.Collapse()

			} else if fsnode.IsDir {
				node.Expand()
				fsnode.ReadChildren()

			} else {
				contentView.SetTitle(fsnode.Path)

				// https://github.com/alecthomas/chroma#try-it
				content, _ := ioutil.ReadFile(fsnode.Path)
				contentView.SetText(string(content))
				contentView.ScrollTo(0, 0)

			}
		}
	})

	// treeView.Set
	treeView.SetChangedFunc(func(pre *tview.TreeNode) {
		node := treeView.GetCurrentNode()
		if fsnode := get(node); fsnode != nil {

			log.Printf("changed: %s", fsnode.Path)
			if !fsnode.IsDir {
				// contentView.SetText(fsnode.Path)
				// contentView.SetTitle(fsnode.Path)
				// content, _ := ioutil.ReadFile(fsnode.Path)
				// contentView.SetText(string(content))
				// contentView.ScrollTo(0, 0)
			}
		}
	})

	app.SetFocus(treeView)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
