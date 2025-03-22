package tools

import "reflect"

const BuiltInPrefix = "__"

type BuiltIns struct {
	SystemInfo  SystemInfo  `config:"system-info" yaml:"system-info" usage:"System information tool: "`
	Environment Environment `config:"environment" yaml:"environment" usage:"Environment tool: "`
	SystemTime  SystemTime  `config:"system-time" yaml:"system-time" usage:"System time tool: "`

	Stats Stats `config:"stats" yaml:"stats" usage:"Stats tool: "`

	ChangeMode  ChangeMode  `config:"change-mode" yaml:"change-mode" usage:"Change mode tool: "`
	ChangeOwner ChangeOwner `config:"change-owner" yaml:"change-owner" usage:"Change owner tool: "`
	ChangeTimes ChangeTimes `config:"change-times" yaml:"change-times" usage:"Change times tool: "`

	FileCreation     FileCreation     `config:"file-creation" yaml:"file-creation" usage:"File creation tool: "`
	FileTempCreation FileTempCreation `config:"temp-file-creation" yaml:"temp-file-creation" usage:"Temporary file creation tool: "`
	FileAppending    FileAppending    `config:"file-appending" yaml:"file-appending" usage:"File appending tool: "`
	FileReading      FileReading      `config:"file-reading" yaml:"file-reading" usage:"File reading tool: "`
	FileDeletion     FileDeletion     `config:"file-deletion" yaml:"file-deletion" usage:"File deletion tool: "`

	DirectoryCreation     DirectoryCreation     `config:"dir-creation" yaml:"dir-creation" usage:"Directory creation tool: "`
	DirectoryTempCreation DirectoryTempCreation `config:"temp-dir-creation" yaml:"temp-dir-creation" usage:"Temporary directory creation tool: "`
	DirectoryDeletion     DirectoryDeletion     `config:"dir-deletion" yaml:"dir-deletion" usage:"Directory deletion tool: "`

	CommandExec CommandExecution `config:"command-execution" yaml:"command-execution" usage:"Command execution tool: "`
}

func (b BuiltIns) AsFunctionDefinitions() []FunctionDefinition {
	var functions []FunctionDefinition

	v := reflect.ValueOf(b)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check if the field has a method called AsFunctionDefinition
		method := field.MethodByName("AsFunctionDefinition")
		if !method.IsValid() {
			continue
		}

		// Call the method and get the result
		result := method.Call(nil)
		if len(result) == 0 || result[0].IsNil() {
			continue
		}

		// Convert the result to FunctionDefinition and add it to the list
		if fn, ok := result[0].Interface().(*FunctionDefinition); ok && fn != nil {
			fn.Name = BuiltInPrefix + fn.Name
			functions = append(functions, *fn)
		}
	}

	return functions
}
