package tree

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
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
	Node  *tview.TreeNode
}

func newRootFsnode(path string) *FSNode {
	fsnode := &FSNode{
		Name:  filepath.Base(path),
		Path:  path,
		IsDir: true,
	}

	fsnode.Node = createNode(fsnode)

	return fsnode
}

func NewRootNode(path string) *tview.TreeNode {
	fsnode := newRootFsnode(path)

	if !fsnode.Node.IsExpanded() {
		fsnode.Node.Expand()
		fsnode.ReadChildren()
	}

	return fsnode.Node
}

func newFsnode(parentPath string, file fs.FileInfo) *FSNode {

	name := file.Name()
	fpath := filepath.Join(parentPath, name)

	fsnode := &FSNode{
		Name:  name,
		Path:  fpath,
		IsDir: file.IsDir(),
		Size:  file.Size(),
	}

	fsnode.Node = createNode(fsnode)

	return fsnode
}

func NewNode(parentPath string, file fs.FileInfo) *tview.TreeNode {
	fsnode := newFsnode(parentPath, file)
	return fsnode.Node
}

func (n *FSNode) Expand() {
	n.ReadChildren()
	n.Node.Expand()
}

func (n *FSNode) Collapse() {
	n.Node.ClearChildren()
	n.Node.Collapse()
}

func (n *FSNode) IsExpanded() bool {
	return n.Node.IsExpanded()
}

func (n *FSNode) readChildren(node *FSNode) {
	if n.IsDir {
		n.Node.ClearChildren()

		files, err := ioutil.ReadDir(n.Path)
		if err != nil {
			panic(err)
		}

		nodes := []*tview.TreeNode{}

		if node != nil {
			log.Printf("looking for node %s", node.Path)
		}

		for _, file := range files {
			fpath := filepath.Join(n.Path, file.Name())

			if node != nil && node.Path == fpath {
				log.Printf("reuse node %s", node.Path)
				nodes = append(nodes, node.Node)
			} else {
				if node != nil {
					log.Printf("new node %s", fpath)
				}
				nodes = append(nodes, NewNode(n.Path, file))
			}
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
			n.Node.AddChild(node)
		}
	}
}

func (n *FSNode) ReadChildren() {
	n.readChildren(nil)
}

func (n *FSNode) CreateParent() *FSNode {
	dir := filepath.Dir(n.Path)
	log.Printf("Create parent for: %s => %s", n.Path, dir)

	if n.Path == dir {
		return n
	}

	rnode := newRootFsnode(dir)

	rnode.readChildren(n)
	rnode.Node.SetExpanded(true)

	return rnode
}

func (n *FSNode) Title() string {
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
