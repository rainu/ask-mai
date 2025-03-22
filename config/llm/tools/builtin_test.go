package tools

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltIns_AsFunctionDefinitions(t *testing.T) {
	toTest := BuiltIns{
		SystemInfo: SystemInfo{
			Disable: false,
		},
	}

	functions := toTest.AsFunctionDefinitions()

	assert.NotEqual(t, 0, len(functions))

	expectedFD := toTest.SystemInfo.AsFunctionDefinition()
	expectedFD.Name = BuiltInPrefix + expectedFD.Name
	assert.Equal(t, fmt.Sprintf("%#v", *expectedFD), fmt.Sprintf("%#v", functions[0]))

	toTest = BuiltIns{
		SystemInfo: SystemInfo{
			Disable: true,
		},
		SystemTime: SystemTime{
			Disable: true,
		},
		FileCreation: FileCreation{
			Disable: true,
		},
		FileTempCreation: FileTempCreation{
			Disable: true,
		},
		FileAppending: FileAppending{
			Disable: true,
		},
		FileReading: FileReading{
			Disable: true,
		},
		FileDeletion: FileDeletion{
			Disable: true,
		},
		CommandExec: CommandExecution{
			Disable: true,
		},
	}
	assert.Equal(t, 0, len(toTest.AsFunctionDefinitions()))
}
