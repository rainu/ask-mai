package backend

type Type int

const (
	TypeSingleShot Type = iota
	TypeMultiShot  Type = iota
)

type Handle interface {
	AskSomething(question string) (string, error)
	Close() error
}

type Builder struct {
	Build func() (Handle, error)
	Type  Type
}
