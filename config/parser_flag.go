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

	flag.Usage = func() {
		printUsage(os.Stderr, fields)
	}

	err := flag.CommandLine.Parse(arguments)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}
}

func (f *resolvedFieldInfo) setupFlag() {
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
