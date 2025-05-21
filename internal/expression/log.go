package expression

import (
	"fmt"
	"os"
)

const FuncNameLog = "log"

var Log = func(args ...interface{}) {
	fmt.Fprint(os.Stderr, "EXPRESSION_LOG: ")
	fmt.Fprintln(os.Stderr, args...)
}
