package common

type NumberContainer struct {
	Expression string  `config:"" yaml:"expression"`
	Value      float64 `config:"-" yaml:"value"`
}

type StringContainer struct {
	Expression string `config:"" yaml:"expression"`
	Value      string `config:"-" yaml:"value"`
}
