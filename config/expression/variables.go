package expression

import "github.com/wailsapp/wails/v2/pkg/runtime"

type Variables struct {
	PrimaryScreen    VariableScreen   `json:"PrimaryScreen"`
	CurrentScreen    VariableScreen   `json:"CurrentScreen"`
	Screens          []VariableScreen `json:"Screens"`
	SecondaryScreens []VariableScreen `json:"SecondaryScreens"`
}

type VariableScreen struct {
	Dimension VariableScreenDimension `json:"Dimension"`
}

type VariableScreenDimension struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

func FromScreens(screens []runtime.Screen) Variables {
	var variables Variables
	for _, screen := range screens {
		variables.Screens = append(variables.Screens, VariableScreen{
			Dimension: VariableScreenDimension{
				Width:  screen.PhysicalSize.Width,
				Height: screen.PhysicalSize.Height,
			},
		})
		if screen.IsPrimary {
			variables.PrimaryScreen = variables.Screens[len(variables.Screens)-1]
		} else {
			variables.SecondaryScreens = append(variables.SecondaryScreens, variables.Screens[len(variables.Screens)-1])
		}
		if screen.IsCurrent {
			variables.CurrentScreen = variables.Screens[len(variables.Screens)-1]
		}
	}

	return variables
}
