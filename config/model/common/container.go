package common

type NumberContainer struct {
	Expression *string  `yaml:"expression,omitempty"`
	Value      *float64 `yaml:"value,omitempty"`
}

type StringContainer struct {
	Expression *string `yaml:"expression,omitempty"`
	Value      *string ` yaml:"value,omitempty"`
}
