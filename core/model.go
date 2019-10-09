package core

type Model interface {
	ForeignKeys() []string
}
