package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
		SetRows(0).
		SetColumns(50, 0)

	app.SetRoot(grid, true)

	root := tree.NewRootNode(pwd)
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
				fsnode.Expand()

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

		}
	})

	selectParent := func(curr *tview.TreeNode) {
		root.Walk(func(node, parent *tview.TreeNode) bool {
			if node == curr {
				if parent != nil && parent != root {
					treeView.SetCurrentNode(parent)
				}
				return false
			}
			return true
		})

	}

	treeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			curr := treeView.GetCurrentNode()
			if fsnode := get(curr); fsnode != nil {
				if fsnode.IsDir && curr.IsExpanded() {
					curr.Collapse()

				} else {
					selectParent(curr)
				}
			}
			return nil
		case tcell.KeyRight:
			curr := treeView.GetCurrentNode()
			if fsnode := get(curr); fsnode != nil {
				if fsnode.IsDir {
					fsnode.Expand()
				}
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'K':
				return nil // noop
			default:
				return event
			}
		default:
			return event
		}
	})

	treeView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseScrollUp:
			return action, nil
		case tview.MouseScrollDown:
			return action, nil
		default:
			return action, event
		}

	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
