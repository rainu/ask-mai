package tools

import "reflect"

const BuiltInPrefix = "__"

type BuiltIns struct {
	SystemInfo   SystemInfo       `config:"system-info" yaml:"system-info" usage:"System information tool: "`
	SystemTime   SystemTime       `config:"system-time" yaml:"system-time" usage:"System time tool: "`
	FileCreation FileCreation     `config:"file-creation" yaml:"file-creation" usage:"File creation tool: "`
	FileReading  FileReading      `config:"file-reading" yaml:"file-reading" usage:"File reading tool: "`
	CommandExec  CommandExecution `config:"command-execution" yaml:"command-execution" usage:"Command execution tool: "`
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
