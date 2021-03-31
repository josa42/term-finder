package tree

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FSNode struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
	node  *tview.TreeNode
}

func NewRootNode(path string) *tview.TreeNode {
	fsnode := &FSNode{
		Name:  ".",
		Path:  path,
		IsDir: true,
	}

	fsnode.node = createNode(fsnode)

	return fsnode.node
}

func NewNode(parentPath string, file fs.FileInfo) *tview.TreeNode {

	name := file.Name()
	fpath := filepath.Join(parentPath, name)

	fsnode := &FSNode{
		Name:  name,
		Path:  fpath,
		IsDir: file.IsDir(),
		Size:  file.Size(),
	}

	fsnode.node = createNode(fsnode)

	return fsnode.node
}

func (n *FSNode) Expand() {
	n.ReadChildren()
	n.node.Expand()
}

func (n *FSNode) ReadChildren() {
	if n.IsDir {
		n.node.ClearChildren()

		files, err := ioutil.ReadDir(n.Path)
		if err != nil {
			panic(err)
		}

		nodes := []*tview.TreeNode{}

		for _, file := range files {
			nodes = append(nodes, NewNode(n.Path, file))
		}

		sort.Slice(nodes, func(i, j int) bool {
			a := nodes[i].GetReference().(*FSNode)
			b := nodes[j].GetReference().(*FSNode)

			if a.IsDir == b.IsDir {
				return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) < 0
			}

			return a.IsDir
		})

		for _, node := range nodes {
			n.node.AddChild(node)
		}
	}
}

func (n *FSNode) Title() string {
	if n.Name == "." {
		return ".."
	}

	icon := ""
	if n.IsDir {
		icon = ""
	}
	return fmt.Sprintf("%s %s", icon, n.Name)
}

func createNode(n *FSNode) *tview.TreeNode {
	node := tview.NewTreeNode(n.Title()).
		SetReference(n).
		SetSelectable(true)

	if n.IsDir {
		node.SetColor(tcell.ColorBlue)
	}

	node.SetExpanded(false)

	return node
}
