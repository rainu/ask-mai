package config

import (
	"fmt"
	"github.com/rainu/ask-mai/config/model"
	"os"
	"reflect"
	"slices"
	"strings"
)

const EnvironmentPrefix = "ASK_MAI_"

func Parse(arguments []string, env []string) *model.Config {
	c := defaultConfig()

	fields := scanConfigTags(nil, c)

	// for possible config file path
	processArguments(arguments, fields)
	processEnvironment(env, fields)

	processYamlFiles(c)

	// again because otherwise the config content will override the command line arguments and environment variables
	processEnvironment(env, fields)
	processArguments(arguments, fields)

	c.Printer.Targets = nil
	for _, target := range c.Printer.TargetsRaw {
		target = strings.TrimSpace(target)

		if target == model.PrinterTargetOut {
			c.Printer.Targets = append(c.Printer.Targets, os.Stdout)
		} else if target == model.PrinterTargetErr {
			c.Printer.Targets = append(c.Printer.Targets, os.Stderr)
		} else {
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
			if err != nil {
				panic(fmt.Errorf("Error creating printer target file: %w", err))
			}
			c.Printer.Targets = append(c.Printer.Targets, file)
		}
	}

	return c
}

type fieldTagInfo struct {
	Name    string
	YamlKey string
	Short   string
	Usage   string
}

func scanConfigTags(parent []fieldTagInfo, v interface{}) (result resolvedFieldInfos) {
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
		path = append(path, fieldTagInfo{
			Name:    getName(field),
			YamlKey: getYamlKey(field),
			Short:   getShort(field),
			Usage:   getUsage(usageProvider, field),
		})

		fieldValue := val.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			result = append(result, scanConfigTags(path, fieldValue.Addr().Interface())...)
		} else {
			result = append(result, extractFieldInfo(path, fieldValue.Addr()))
		}
	}

	return
}

func shouldSkip(field reflect.StructField) bool {
	return field.Tag.Get("config") == "-"
}

func getName(field reflect.StructField) string {
	if ct, ok := field.Tag.Lookup("config"); ok {
		return ct
	}
	if ct, ok := field.Tag.Lookup("yaml"); ok {
		return ct
	}
	return strings.ToLower(field.Name)
}

func getYamlKey(field reflect.StructField) string {
	if ct, ok := field.Tag.Lookup("yaml"); ok {
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

type resolvedFieldInfo struct {
	YamlPath []string
	Flag     string
	Short    string
	Env      string
	Usage    string

	Value reflect.Value
}
type resolvedFieldInfos []resolvedFieldInfo

func extractFieldInfo(path []fieldTagInfo, val reflect.Value) resolvedFieldInfo {
	sbFlag := strings.Builder{}
	sbEnv := strings.Builder{}
	sbShort := strings.Builder{}
	sbUsage := strings.Builder{}

	var sPath []string
	for i, p := range path {
		sPath = append(sPath, p.YamlKey)
		if i > 0 && p.Name != "" {
			sbFlag.WriteString("-")
			sbEnv.WriteString("_")
		}
		sbFlag.WriteString(p.Name)
		sbEnv.WriteString(strings.Replace(p.Name, "-", "_", -1))
		sbShort.WriteString(p.Short)
		sbUsage.WriteString(p.Usage)
	}
	return resolvedFieldInfo{
		YamlPath: sPath,
		Flag:     strings.TrimLeft(sbFlag.String(), "-"),
		Short:    sbShort.String(),
		Env:      EnvironmentPrefix + strings.TrimLeft(strings.ToUpper(sbEnv.String()), "_"),
		Usage:    strings.Trim(sbUsage.String(), " "),
		Value:    val,
	}
}
