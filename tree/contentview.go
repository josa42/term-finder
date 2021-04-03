package tree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ tview.Primitive = &ContentView{}

type ContentView struct {
	theme       *Theme
	view        tview.Primitive
	topbarView  *tview.TextView
	infoView    *tview.TextView
	contentView *tview.TextView
}

func formatSize(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func NewContentView(theme *Theme) *ContentView {
	view := tview.NewGrid()
	view.SetRows(3, 5, 0)

	topbar := tview.NewTextView()
	topbar.SetBorderPadding(0, 0, 1, 1)
	topbar.SetBorder(true)
	topbar.SetBorderColor(theme.TopbarBorder)
	topbar.SetDynamicColors(true)
	topbar.SetRegions(true)

	info := tview.NewTextView()
	info.SetDynamicColors(true)
	info.SetRegions(true)

	content := tview.NewTextView()
	content.SetBorderPadding(0, 0, 2, 2)
	content.SetBackgroundColor(theme.ContentBackground)

	view.AddItem(topbar, 0, 0, 1, 1, 0, 0, false)
	view.AddItem(info, 1, 0, 1, 1, 0, 0, false)
	view.AddItem(content, 2, 0, 1, 1, 0, 0, false)

	ft := &ContentView{
		theme:       theme,
		view:        view,
		topbarView:  topbar,
		infoView:    info,
		contentView: content,
	}

	return ft
}

// Primitive interface

func (ft *ContentView) Draw(screen tcell.Screen) {
	ft.view.Draw(screen)
}
func (ft *ContentView) GetRect() (int, int, int, int) {
	return ft.view.GetRect()
}
func (ft *ContentView) SetRect(x, y, width, height int) {
	ft.view.SetRect(x, y, width, height)
}
func (ft *ContentView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ft.view.InputHandler()
}
func (ft *ContentView) Focus(delegate func(p tview.Primitive)) {
	ft.view.Focus(delegate)
}

func (ft *ContentView) HasFocus() bool {
	return ft.view.HasFocus()
}
func (ft *ContentView) Blur() {
	ft.view.Blur()
}
func (ft *ContentView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return ft.view.MouseHandler()
}

func (v *ContentView) SetPreview(fsnode *FSNode) {
	v.topbarView.SetText(formatPath(fsnode.Path))

	if !fsnode.IsDir && fsnode.Size < 400_000 {
		v.contentView.SetText(fsnode.Path)
		v.contentView.SetTitle(fsnode.Path)
		content, _ := ioutil.ReadFile(fsnode.Path)
		v.contentView.SetText(string(content))
		v.contentView.ScrollTo(0, 0)
	} else {
		v.contentView.SetText("")
	}

	v.infoView.SetText(strings.Join([]string{
		fmt.Sprintf(" [#5c6370]│      Mode:[normal] %v", fsnode.Mode),
		fmt.Sprintf(" [#5c6370]│  Modified:[normal] %v", fsnode.ModTime),
		fmt.Sprintf(" [#5c6370]│      Size:[normal] %v", formatSize(fsnode.Size)),
		fmt.Sprintf(" [#5c6370]│ Mime Type:[normal] %v", fsnode.MimeType),
	}, "\n"))
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
