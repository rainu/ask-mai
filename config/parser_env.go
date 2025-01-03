package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func processEnvironment(env []string, fields resolvedFieldInfos) {
	parsedEnv := map[string]string{}

	for i := range env {
		parts := strings.SplitN(env[i], "=", 2)
		parsedEnv[parts[0]] = parts[1]
	}

	fields.setupEnvs(parsedEnv)
}

func (f *resolvedFieldInfo) setupEnv(env map[string]string) {
	esv, exists := env[f.Env]
	if !exists && !strings.HasPrefix(f.Value.Type().String(), "*[]") {
		return
	}

	switch f.Value.Type().String() {
	case "*string":
		sp := f.Value.Interface().(*string)
		*sp = esv
	case "*int":
		ip := f.Value.Interface().(*int)

		iv, err := strconv.Atoi(esv)
		if err != nil {
			panic(fmt.Errorf("invalid integer value '%s': %w", esv, err))
		}
		*ip = iv
	case "*uint":
		ip := f.Value.Interface().(*uint)

		iv, err := strconv.ParseUint(esv, 10, 32)
		if err != nil {
			panic(fmt.Errorf("invalid unsigned integer value '%s': %w", esv, err))
		}
		*ip = uint(iv)
	case "*bool":
		bp := f.Value.Interface().(*bool)

		switch strings.ToLower(esv) {
		case "1":
			fallthrough
		case "true":
			*bp = true
		default:
			*bp = false
		}
	case "*float64":
		fp := f.Value.Interface().(*float64)

		fv, err := strconv.ParseFloat(esv, 64)
		if err != nil {
			panic(fmt.Errorf("invalid float value '%s': %w", esv, err))
		}
		*fp = fv
	case "*[]string":
		sp := f.Value.Interface().(*[]string)

		var envKeys []string
		for key := range env {
			if strings.HasPrefix(key, f.Env+"_") {
				envKeys = append(envKeys, key)
			}
		}
		if len(envKeys) == 0 {
			return
		}
		slices.Sort(envKeys)

		var envValues []string
		for _, key := range envKeys {
			envValues = append(envValues, env[key])
		}

		*sp = envValues
	}
}

func (f resolvedFieldInfos) setupEnvs(env map[string]string) {
	for i := range f {
		f[i].setupEnv(env)
	}
}
