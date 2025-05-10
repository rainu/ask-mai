package tools

import "reflect"

const BuiltInPrefix = "_"

type BuiltIns struct {
	SystemInfo  SystemInfo  `yaml:"system-info,omitempty" usage:"System information tool: "`
	Environment Environment `yaml:"environment,omitempty" usage:"Environment tool: "`
	SystemTime  SystemTime  `yaml:"system-time,omitempty" usage:"System time tool: "`

	Stats Stats `yaml:"stats,omitempty" usage:"Stats tool: "`

	ChangeMode  ChangeMode  `yaml:"change-mode,omitempty" usage:"Change mode tool: "`
	ChangeOwner ChangeOwner `yaml:"change-owner,omitempty" usage:"Change owner tool: "`
	ChangeTimes ChangeTimes `yaml:"change-times,omitempty" usage:"Change times tool: "`

	FileCreation     FileCreation     `yaml:"file-creation,omitempty" usage:"File creation tool: "`
	FileTempCreation FileTempCreation `yaml:"temp-file-creation,omitempty" usage:"Temporary file creation tool: "`
	FileAppending    FileAppending    `yaml:"file-appending,omitempty" usage:"File appending tool: "`
	FileReading      FileReading      `yaml:"file-reading,omitempty" usage:"File reading tool: "`
	FileDeletion     FileDeletion     `yaml:"file-deletion,omitempty,omitempty" usage:"File deletion tool: "`

	DirectoryCreation     DirectoryCreation     `yaml:"dir-creation,omitempty" usage:"Directory creation tool: "`
	DirectoryTempCreation DirectoryTempCreation `yaml:"temp-dir-creation,omitempty" usage:"Temporary directory creation tool: "`
	DirectoryDeletion     DirectoryDeletion     `yaml:"dir-deletion,omitempty" usage:"Directory deletion tool: "`

	CommandExec CommandExecution `yaml:"command-execution,omitempty" usage:"Command execution tool: "`

	Http Http `yaml:"http,omitempty" usage:"HTTP tool: "`

	Disable bool `yaml:"disable,omitempty" usage:"Disable all builtin tools."`
}

func NewBuiltIns() *BuiltIns {
	return &BuiltIns{
		SystemInfo:            NewSystemInfo(),
		Environment:           NewEnvironment(),
		SystemTime:            NewSystemTime(),
		Stats:                 NewStats(),
		ChangeMode:            NewChangeMode(),
		ChangeOwner:           NewChangeOwner(),
		ChangeTimes:           NewChangeTimes(),
		FileCreation:          NewFileCreation(),
		FileTempCreation:      NewFileTempCreation(),
		FileAppending:         NewFileAppending(),
		FileReading:           NewFileReading(),
		FileDeletion:          NewFileDeletion(),
		DirectoryCreation:     NewDirectoryCreation(),
		DirectoryTempCreation: NewDirectoryTempCreation(),
		DirectoryDeletion:     NewDirectoryDeletion(),
		CommandExec:           NewCommandExecution(),
		Http:                  NewHttp(),
		Disable:               false,
	}
}

func (b BuiltIns) AsFunctionDefinitions() []FunctionDefinition {
	if b.Disable {
		return []FunctionDefinition{}
	}

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
