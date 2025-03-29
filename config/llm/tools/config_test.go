package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_GetTools(t *testing.T) {
	toTest := Config{}

	result := toTest.GetTools()

	fd, contains := result[BuiltInPrefix+"getSystemInformation"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	_, contains = result[BuiltInPrefix+"getEnvironment"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"getSystemTime"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"getStats"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"changeMode"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"changeOwner"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"changeTimes"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"appendFile"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"createFile"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"deleteFile"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"createTempFile"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"createDirectory"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"deleteDirectory"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"createTempDirectory"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"readTextFile"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())
	fd, contains = result[BuiltInPrefix+"executeCommand"]
	assert.True(t, contains)
	assert.True(t, fd.IsBuiltIn())

	// deactivate builtin tool

	toTest.BuiltInTools.SystemInfo.Disable = true
	toTest.BuiltInTools.SystemTime.Disable = true
	toTest.BuiltInTools.Environment.Disable = true
	toTest.BuiltInTools.Stats.Disable = true
	toTest.BuiltInTools.ChangeMode.Disable = true
	toTest.BuiltInTools.ChangeOwner.Disable = true
	toTest.BuiltInTools.ChangeTimes.Disable = true
	toTest.BuiltInTools.FileAppending.Disable = true
	toTest.BuiltInTools.FileCreation.Disable = true
	toTest.BuiltInTools.FileTempCreation.Disable = true
	toTest.BuiltInTools.FileReading.Disable = true
	toTest.BuiltInTools.FileDeletion.Disable = true
	toTest.BuiltInTools.DirectoryCreation.Disable = true
	toTest.BuiltInTools.DirectoryTempCreation.Disable = true
	toTest.BuiltInTools.DirectoryDeletion.Disable = true
	toTest.BuiltInTools.CommandExec.Disable = true
	assert.Empty(t, toTest.GetTools())
}
