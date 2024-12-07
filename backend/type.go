package backend

type Type int

const (
	TypeSingleShot Type = iota
	TypeMultiShot  Type = iota
)

type Role string

const (
	RoleUser Role = "user"
	RoleBot  Role = "bot"
)

type Message struct {
	Content string
	Role    Role
}

type Handle interface {
	AskSomething(chat []Message) (string, error)
	Close() error
}

type Builder struct {
	Build func() (Handle, error)
	Type  Type
}
