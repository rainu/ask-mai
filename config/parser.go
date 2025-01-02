package config

import (
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"reflect"
	"slices"
	"strings"
)

type UsageProvider interface {
	GetUsage(field string) string
}

func Parse(arguments []string) *Config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	c := defaultConfig()
	scanConfigTags(nil, c)

	flag.Usage = func() {
		printUsage(os.Stderr)
	}

	err := flag.CommandLine.Parse(arguments)
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
		return nil
	}

	c.Printer.Targets = nil
	for _, target := range strings.Split(c.Printer.TargetsRaw, ",") {
		target = strings.TrimSpace(target)

		if target == PrinterTargetOut {
			c.Printer.Targets = append(c.Printer.Targets, os.Stdout)
		} else if target == PrinterTargetErr {
			c.Printer.Targets = append(c.Printer.Targets, os.Stderr)
		} else {
			file, err := os.Create(target)
			if err != nil {
				panic(fmt.Errorf("Error creating printer target file: %w", err))
			}
			c.Printer.Targets = append(c.Printer.Targets, file)
		}
	}

	return c
}

type fieldInfo struct {
	Name  string
	Short string
	Usage string
}

func scanConfigTags(parent []fieldInfo, v interface{}) {
	val := reflect.ValueOf(v).Elem()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	var usageProvider UsageProvider
	if up, ok := v.(UsageProvider); ok {
		usageProvider = up
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if shouldSkip(field) {
			continue
		}

		path := slices.Clone(parent)
		path = append(path, fieldInfo{
			Name:  getName(field),
			Short: getShort(field),
			Usage: getUsage(usageProvider, field),
		})

		fieldValue := val.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			scanConfigTags(path, fieldValue.Addr().Interface())
		} else {
			setupConfig(path, fieldValue.Addr())
		}
	}
}

func shouldSkip(field reflect.StructField) bool {
	return field.Tag.Get("config") == "-"
}

func getName(field reflect.StructField) string {
	if ct, ok := field.Tag.Lookup("config"); ok {
		return ct
	}
	return strings.ToLower(field.Name)
}

func getShort(field reflect.StructField) string {
	return field.Tag.Get("short")
}

func getUsage(up UsageProvider, field reflect.StructField) string {
	usage := field.Tag.Get("usage")
	if usage != "" {
		return usage
	}

	if up != nil {
		return up.GetUsage(field.Name)
	}
	return ""
}

func setupConfig(path []fieldInfo, val reflect.Value) {
	sbFlag := strings.Builder{}
	sbShort := strings.Builder{}
	sbUsage := strings.Builder{}

	for i, p := range path {
		if i > 0 && p.Name != "" {
			sbFlag.WriteString("-")
		}
		sbFlag.WriteString(p.Name)
		sbShort.WriteString(p.Short)
		sbUsage.WriteString(p.Usage)
	}
	sFlag := strings.TrimLeft(sbFlag.String(), "-")
	sShort := sbShort.String()
	sUsage := strings.Trim(sbUsage.String(), " ")

	switch val.Type().String() {
	case "*string":
		sp := val.Interface().(*string)
		sv := *sp
		flag.StringVarP(sp, sFlag, sShort, sv, sUsage)
	case "*int":
		ip := val.Interface().(*int)
		iv := *ip
		flag.IntVarP(ip, sFlag, sShort, iv, sUsage)
	case "*uint":
		ip := val.Interface().(*uint)
		iv := *ip
		flag.UintVarP(ip, sFlag, sShort, iv, sUsage)
	case "*bool":
		bp := val.Interface().(*bool)
		bv := *bp
		flag.BoolVarP(bp, sFlag, sShort, bv, sUsage)
	case "*float64":
		fp := val.Interface().(*float64)
		fv := *fp
		flag.Float64VarP(fp, sFlag, sShort, fv, sUsage)
	case "*[]string":
		sp := val.Interface().(*[]string)
		sv := *sp
		if sv == nil {
			sv = []string{}
		}
		flag.StringSliceVarP(sp, sFlag, sShort, sv, sUsage)
	}
}
