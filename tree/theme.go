package tree

import "github.com/gdamore/tcell/v2"

type Theme struct {
	SidebarBackground tcell.Color
	SidebarLines      tcell.Color
	ContentBackground tcell.Color
	Border            tcell.Color
}

var theme *Theme

func GetTheme() *Theme {
	if theme == nil {
		theme = &Theme{
			SidebarLines:      tcell.NewHexColor(0x5c6370),
			SidebarBackground: tcell.NewHexColor(0x21252B),
			ContentBackground: tcell.NewHexColor(0x282c34),
			// Border:            tcell.NewHexColor(0x5c6370),
		}
	}
	return theme
}
