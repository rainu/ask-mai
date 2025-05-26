package tools

import (
	"github.com/rainu/go-yacl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_GetTools(t *testing.T) {
	toTest := Config{}
	yacl.NewConfig(&toTest).ApplyDefaults()

	result, err := toTest.GetTools(t.Context())
	assert.NoError(t, err)

	_, contains := result[ServerNameBuiltin+"_getSystemInformation"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_getEnvironment"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_getSystemTime"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_getStats"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_changeMode"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_changeOwner"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_changeTimes"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_appendFile"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_createFile"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_deleteFile"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_createTempFile"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_createDirectory"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_deleteDirectory"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_createTempDirectory"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_readTextFile"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_executeCommand"]
	assert.True(t, contains)
	_, contains = result[ServerNameBuiltin+"_callHttp"]
	assert.True(t, contains)

	// deactivate builtin tool

	toTest.BuiltIns.SystemInfo.Disable = true
	toTest.BuiltIns.SystemTime.Disable = true
	toTest.BuiltIns.Environment.Disable = true
	toTest.BuiltIns.Stats.Disable = true
	toTest.BuiltIns.ChangeMode.Disable = true
	toTest.BuiltIns.ChangeOwner.Disable = true
	toTest.BuiltIns.ChangeTimes.Disable = true
	toTest.BuiltIns.FileAppending.Disable = true
	toTest.BuiltIns.FileCreation.Disable = true
	toTest.BuiltIns.FileTempCreation.Disable = true
	toTest.BuiltIns.FileReading.Disable = true
	toTest.BuiltIns.FileDeletion.Disable = true
	toTest.BuiltIns.DirectoryCreation.Disable = true
	toTest.BuiltIns.DirectoryTempCreation.Disable = true
	toTest.BuiltIns.DirectoryDeletion.Disable = true
	toTest.BuiltIns.CommandExec.Disable = true
	toTest.BuiltIns.Http.Disable = true

	result, err = toTest.GetTools(t.Context())
	assert.NoError(t, err)

	assert.Empty(t, result)
}
