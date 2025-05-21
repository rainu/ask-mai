package expression

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const VarNameScreens = "screens"

type Screen struct {
	PrimaryScreen    VariableScreen   `json:"PrimaryScreen"`
	CurrentScreen    VariableScreen   `json:"CurrentScreen"`
	Screens          []VariableScreen `json:"Screens"`
	SecondaryScreens []VariableScreen `json:"SecondaryScreens"`
}

type VariableScreen struct {
	Dimension         VariableScreenDimension `json:"Dimension"`
	PhysicalDimension VariableScreenDimension `json:"PhysicalDimension"`
}

type VariableScreenDimension struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

func SetScreens(screens []runtime.Screen) Screen {
	var s Screen
	for _, screen := range screens {
		s.Screens = append(s.Screens, VariableScreen{
			Dimension: VariableScreenDimension{
				Width:  screen.Size.Width,
				Height: screen.Size.Height,
			},
			PhysicalDimension: VariableScreenDimension{
				Width:  screen.PhysicalSize.Width,
				Height: screen.PhysicalSize.Height,
			},
		})
		if screen.IsPrimary {
			s.PrimaryScreen = s.Screens[len(s.Screens)-1]
		} else {
			s.SecondaryScreens = append(s.SecondaryScreens, s.Screens[len(s.Screens)-1])
		}
		if screen.IsCurrent {
			s.CurrentScreen = s.Screens[len(s.Screens)-1]
		}
	}

	globalVariables[VarNameScreens] = s

	return s
}
