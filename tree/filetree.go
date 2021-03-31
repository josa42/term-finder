package tree

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileTree struct {
	theme     *Theme
	view      *tview.TreeView
	root      *tview.TreeNode
	onSelect  func(node *FSNode)
	onChanged func(node *FSNode)
}

func get(node *tview.TreeNode) *FSNode {
	ref := node.GetReference()
	if ref == nil {
		return nil
	}
	return ref.(*FSNode)
}

func NewFileTree(theme *Theme) *FileTree {
	view := tview.NewTreeView().
		SetTopLevel(1)

	ft := &FileTree{
		theme: theme,
		view:  view,
	}

	view.SetBorderPadding(0, 0, 2, 2)
	view.SetGraphicsColor(theme.SidebarLines)
	view.SetBackgroundColor(theme.SidebarBackground)

	view.SetSelectedFunc(func(node *tview.TreeNode) {
		ft.selected(node)
	})

	view.SetChangedFunc(func(node *tview.TreeNode) {
		ft.changed(node)
	})

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return ft.inputCapture(event)
	})

	// Disable mouse scroll
	view.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		return ft.mouseCapture(action, event)
	})

	return ft
}

func (ft *FileTree) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	fsnode := get(ft.view.GetCurrentNode())

	switch event.Key() {
	case tcell.KeyLeft:

		parent := ft.GetParent(fsnode)

		if fsnode.IsDir && fsnode.IsExpanded() {
			fsnode.Collapse()

		} else if ft.IsRoot(parent) {
			ft.SetRoot(parent.CreateParent())

		} else {
			ft.SetCurrent(parent)
		}

		return nil

	case tcell.KeyRight:
		if fsnode.IsDir {
			if fsnode.IsExpanded() {
				ft.SetRoot(fsnode)

			} else {
				fsnode.Expand()
			}
		}
		return nil

	case tcell.KeyRune:
		switch event.Rune() {
		case 'K':
			return nil // noop

		case 'c':
			if fsnode.IsDir {
				ft.SetRoot(fsnode)
			}
			return nil
		case 'C':
			ft.SetRoot(ft.GetRoot().CreateParent())
			return nil

		default:
			return event
		}

	default:
		return event
	}
}

func (ft *FileTree) mouseCapture(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	switch action {
	case tview.MouseScrollUp:
		return action, nil
	case tview.MouseScrollDown:
		return action, nil
	default:
		return action, event
	}
}

func (ft *FileTree) selected(node *tview.TreeNode) {
	fsnode := get(node)
	if fsnode.IsExpanded() {
		fsnode.Collapse()

	} else if fsnode.IsDir {
		fsnode.Expand()

	} else {
		if ft.onSelect != nil {
			ft.onSelect(fsnode)
		}
	}
}

func (ft *FileTree) changed(node *tview.TreeNode) {
	if ft.onChanged != nil {
		ft.onChanged(get(node))
	}
}

func (ft *FileTree) GetParent(fsnode *FSNode) *FSNode {
	var currParent *tview.TreeNode
	ft.root.Walk(func(node, parent *tview.TreeNode) bool {
		if node == fsnode.Node {
			currParent = parent
			return false
		}
		return true
	})

	return get(currParent)
}

func (ft *FileTree) SetRoot(fsnode *FSNode) {
	if fsnode != nil {
		ft.root = fsnode.Node
		ft.view.SetRoot(fsnode.Node)
	}
}

func (ft *FileTree) GetRoot() *FSNode {
	return get(ft.root)
}

func (ft *FileTree) SetCurrent(fsnode *FSNode) {
	if fsnode != nil {
		ft.view.SetCurrentNode(fsnode.Node)
	}
}

func (ft *FileTree) IsRoot(fsnode *FSNode) bool {
	if !fsnode.IsExpanded() {
		fsnode.Expand()
	}

	return ft.root == fsnode.Node
}

func (ft *FileTree) Load(dir string) {
	ft.root = NewRootNode(dir)
	ft.view.SetRoot(ft.root).
		SetCurrentNode(ft.root)

}

func (ft *FileTree) OnSelect(fn func(node *FSNode)) {
	ft.onSelect = fn
}

func (ft *FileTree) OnChanged(fn func(node *FSNode)) {
	ft.onChanged = fn
}

// TODO remive these getters
func (ft *FileTree) GetView() *tview.TreeView {
	return ft.view
}

func (ft *FileTree) GetRootNode() *tview.TreeNode {
	return ft.root
}
