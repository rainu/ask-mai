package expression

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"os"
)

var precompiledPrograms = map[string]*goja.Program{}

func Precompile(source string) error {
	file, err := os.Open(source)
	if err == nil && file != nil {
		defer file.Close()

		ast, err := parser.ParseFile(nil, file.Name(), file, 0)
		if err != nil {
			return fmt.Errorf("error parsing file: %w", err)
		}
		prog, err := goja.CompileAST(ast, false)
		if err != nil {
			return fmt.Errorf("error compiling file: %w", err)
		}
		precompiledPrograms[source] = prog
	} else {
		prog, err := goja.Compile("", source, false)
		if err != nil {
			return fmt.Errorf("error compiling file: %w", err)
		}
		precompiledPrograms[source] = prog
	}

	return nil
}
