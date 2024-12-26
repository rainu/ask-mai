package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
		{"v.PrimaryScreen.Dimension.Height * 2", float64(variables.PrimaryScreen.Dimension.Height * 2)},
		{"v.PrimaryScreen.Dimension.Width / 2", float64(variables.PrimaryScreen.Dimension.Width / 2)},
		{"v.Screens[1].Dimension.Width - v.Screens[0].Dimension.Width", float64(variables.Screens[1].Dimension.Width - variables.Screens[0].Dimension.Width)},
		{"v.Screens[1].Dimension.Height / v.Screens[0].Dimension.Height", float64(variables.Screens[1].Dimension.Height / variables.Screens[0].Dimension.Height)},
		{"3 + 2", float64(3 + 2)},
		{"(3 + 2) * 4", float64((3 + 2) * 4)},
		{"(3+v.Screens[1].Dimension.Height)*4", float64((3 + variables.Screens[1].Dimension.Height) * 4)},
		{"if(v.PrimaryScreen.Dimension.Height >= 2160){ 1 } else { 2 }", float64(2)},
		{"let r=0; for(let i=0; i < 5; i++){ r+=i }; r", float64(10)},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			result, err := Expression(tt.expression).Calculate(variables)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpression_Calculate_NaN(t *testing.T) {
	_, err := Expression("log('test'); 'test'").Calculate(Variables{})
	assert.Error(t, err)
}
