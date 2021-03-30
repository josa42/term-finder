package main

import (
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell/v2"
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

// Show a navigable tree view of the current directory.
func main() {

	pwd, _ := os.Getwd()

	root := tree.NewNode(pwd, ".", true)
	root.Expand()
	get(root).ReadChildren()

	treeView := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root).
		SetTopLevel(1)

	// treeView.SetBorder(true)
	treeView.SetBorderPadding(0, 0, 1, 1)
	treeView.SetGraphicsColor(tcell.NewHexColor(0x5c6370))
	treeView.SetBackgroundColor(tcell.NewHexColor(0x2c323c))
	// treeView.SetBorder(true)
	// treeView.SetBorderColor(tcell.NewHexColor(0x5c6370))

	contentView := tview.NewTextView()
	contentView.SetBorderPadding(0, 0, 1, 1)
	contentView.SetBackgroundColor(tcell.NewHexColor(0x282c34))
	// contentView.SetBorder(true)
	// contentView.SetBorderColor(tcell.NewHexColor(0x5c6370))

	// If a directory was selected, open it.
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

			if !fsnode.IsDir {
				// contentView.SetText(fsnode.Path)
				// contentView.SetTitle(fsnode.Path)
				// content, _ := ioutil.ReadFile(fsnode.Path)
				// contentView.SetText(string(content))
				// contentView.ScrollTo(0, 0)
			}
		}
	})

	grid := tview.NewGrid().
		SetBordersColor(tcell.NewHexColor(0x5c6370)).
		SetBorders(true).
		SetRows(0).
		SetColumns(50, 0)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 0, true)

	grid.
		AddItem(treeView, 0, 0, 1, 1, 0, 50, true).
		AddItem(contentView, 0, 1, 1, 1, 0, 50, false)

	if err := tview.NewApplication().SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
