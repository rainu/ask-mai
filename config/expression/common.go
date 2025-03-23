package expression

import (
	"fmt"
	"github.com/dop251/goja"
	"os"
)

const VarNameVariables = "v"
const FuncNameLog = "log"

func SetupLog(vm *goja.Runtime) error {
	return vm.Set(FuncNameLog, func(args ...interface{}) {
		fmt.Fprint(os.Stderr, "EXPRESSION_LOG: ")
		fmt.Fprintln(os.Stderr, args...)
	})
}
