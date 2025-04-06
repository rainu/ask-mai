package expression

import (
	"fmt"
	"github.com/dop251/goja"
)

type Result struct {
	result goja.Value
	err    error
}

func (r *Result) AsBoolean() (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	return r.result.ToBoolean(), nil
}

func (r *Result) AsFloat() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if r.result.ToNumber().SameAs(goja.NaN()) {
		return 0, fmt.Errorf("result is not a number")
	}

	return r.result.ToFloat(), nil
}

func (r *Result) AsString() (string, error) {
	if r.err != nil {
		return "", r.err
	}

	return r.result.String(), nil
}

func (r *Result) AsByteArray() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}

	return []byte(r.result.String()), nil
}
