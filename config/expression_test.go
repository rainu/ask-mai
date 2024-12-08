package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpression_Variables(t *testing.T) {
	v := Variables{
		PrimaryScreen: VariableScreen{
			Dimension: VariableScreenDimension{
				Width:  1920,
				Height: 1080,
			},
		},
		Screens: []VariableScreen{
			{
				Dimension: VariableScreenDimension{
					Width:  1920,
					Height: 1080,
				},
			},
			{
				Dimension: VariableScreenDimension{
					Width:  3840,
					Height: 2160,
				},
			},
		},
	}

	m, e := v.ToFlatMap()
	require.NoError(t, e)

	assert.Equal(t, map[string]any{
		"PrimaryScreen.Dimension.Width":  float64(1920),
		"PrimaryScreen.Dimension.Height": float64(1080),
		"Screens[0].Dimension.Width":     float64(1920),
		"Screens[0].Dimension.Height":    float64(1080),
		"Screens[1].Dimension.Width":     float64(3840),
		"Screens[1].Dimension.Height":    float64(2160),
	}, m)
}

func TestExpression_Calculate(t *testing.T) {
	variables := Variables{
		PrimaryScreen: VariableScreen{
			Dimension: VariableScreenDimension{
				Width:  1920,
				Height: 1080,
			},
		},
		Screens: []VariableScreen{
			{
				Dimension: VariableScreenDimension{
					Width:  1920,
					Height: 1080,
				},
			},
			{
				Dimension: VariableScreenDimension{
					Width:  3840,
					Height: 2160,
				},
			},
		},
	}

	tests := []struct {
		expression string
		expected   float64
	}{
		{"PrimaryScreen.Dimension.Height * 2", float64(variables.PrimaryScreen.Dimension.Height * 2)},
		{"PrimaryScreen.Dimension.Width / 2", float64(variables.PrimaryScreen.Dimension.Width / 2)},
		{"Screens[1].Dimension.Width - Screens[0].Dimension.Width", float64(variables.Screens[1].Dimension.Width - variables.Screens[0].Dimension.Width)},
		{"Screens[1].Dimension.Height / Screens[0].Dimension.Height", float64(variables.Screens[1].Dimension.Height / variables.Screens[0].Dimension.Height)},
		{"3 + 2", float64(3 + 2)},
		{"(3 + 2) * 4", float64((3 + 2) * 4)},
		{"(3+Screens[1].Dimension.Height)*4", float64((3 + variables.Screens[1].Dimension.Height) * 4)},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := Expression(tt.expression).Calculate(variables)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
