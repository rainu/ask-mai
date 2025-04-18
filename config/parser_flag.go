package config

import (
	"errors"
	flag "github.com/spf13/pflag"
	"os"
)

type UsageProvider interface {
	GetUsage(field string) string
}

func processArguments(arguments []string, fields resolvedFieldInfos) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fields.setupFlags()

	var helpArgs, helpEnv, helpYaml, helpStyles, helpExpr, helpTool bool
	flag.BoolVarP(&helpArgs, "help-arg", "", false, "Show this help")
	flag.BoolVarP(&helpEnv, "help-env", "", false, "Show help for environment variables")
	flag.BoolVarP(&helpYaml, "help-config", "", false, "Show help for config file")
	flag.BoolVarP(&helpStyles, "help-styles", "", false, "Show help of code styles")
	flag.BoolVarP(&helpExpr, "help-expression", "", false, "Show help for expressions")
	flag.BoolVarP(&helpTool, "help-tool", "", false, "Show help for tools")

	flag.Usage = func() {
		printHelpArgs(os.Stderr, fields)
	}

	err := flag.CommandLine.Parse(arguments)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	} else if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if helpArgs {
		printHelpArgs(os.Stdout, fields)
		os.Exit(0)
	} else if helpEnv {
		printHelpEnv(os.Stdout, fields)
		os.Exit(0)
	} else if helpYaml {
		printHelpConfig(os.Stdout, fields)
		os.Exit(0)
	} else if helpStyles {
		printHelpStyles(os.Stdout)
		os.Exit(0)
	} else if helpExpr {
		printHelpExpression(os.Stdout)
		os.Exit(0)
	} else if helpTool {
		printHelpTool(os.Stdout)
		os.Exit(0)
	}
}

func (f *resolvedFieldInfo) setupFlag() {
	if f.Flag == "" {
		return
	}

	switch f.Value.Type().String() {
	case "*string":
		sp := f.Value.Interface().(*string)
		sv := *sp
		flag.StringVarP(sp, f.Flag, f.Short, sv, f.Usage)
	case "*int":
		ip := f.Value.Interface().(*int)
		iv := *ip
		flag.IntVarP(ip, f.Flag, f.Short, iv, f.Usage)
	case "*uint":
		ip := f.Value.Interface().(*uint)
		iv := *ip
		flag.UintVarP(ip, f.Flag, f.Short, iv, f.Usage)
	case "*bool":
		bp := f.Value.Interface().(*bool)
		bv := *bp
		flag.BoolVarP(bp, f.Flag, f.Short, bv, f.Usage)
	case "*float64":
		fp := f.Value.Interface().(*float64)
		fv := *fp
		flag.Float64VarP(fp, f.Flag, f.Short, fv, f.Usage)
	case "*[]string":
		sp := f.Value.Interface().(*[]string)
		sv := *sp
		if sv == nil {
			sv = []string{}
		}
		flag.StringSliceVarP(sp, f.Flag, f.Short, sv, f.Usage)
	}
}

func (f resolvedFieldInfos) setupFlags() {
	for i := range f {
		f[i].setupFlag()
	}
}
