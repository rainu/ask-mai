package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThemes_Validate(t *testing.T) {
	toTest := Themes{
		Custom: map[string]Theme{
			"custom-theme": {
				Colors: map[string]string{
					"chat-system-message": "#FFFFFF",
				},
			},
		},
	}
	toTest.SetDefaults()

	toTest.Dark.Colors["chat-tool-call"] = "#FF00FF"

	assert.NoError(t, toTest.Validate())
	assert.Equal(t, Themes{
		Dark: &Theme{
			Dark: true,
			Colors: map[string]string{
				"chat-system-message":    "",
				"chat-assistant-message": "#E0E0E0",
				"chat-tool-call":         "#FF00FF",
				"chat-user-message":      "#69F0AE",
			},
		},
		Light: &Theme{
			Dark: false,
			Colors: map[string]string{
				"chat-system-message":    "",
				"chat-assistant-message": "#E0E0E0",
				"chat-tool-call":         "#E0E0E0",
				"chat-user-message":      "#69F0AE",
			},
		},
		Custom: map[string]Theme{
			"custom-theme": {
				Colors: map[string]string{
					"chat-system-message":    "#FFFFFF",
					"chat-assistant-message": "#E0E0E0",
					"chat-tool-call":         "#E0E0E0",
					"chat-user-message":      "#69F0AE",
				},
			},
		},
	}, toTest)
}
