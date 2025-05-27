package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChat_SimpleTest(t *testing.T) {
	cfg, ctrl := initTest(t)
	deactivateAllTools(cfg)

	res, err := simpleAsk(ctrl, "Test")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.True(t, res.Consumption["prompt"] < 50, "Prompt consumption should be less than 50: %d", res.Consumption["prompt"])
	assert.True(t, res.Consumption["completion"] < 50, "Completion consumption should be less than 50: %d", res.Consumption["completion"])
}
