package types

type Contact interface {
	IsEmpty() bool
	EmailAddress() string
	Record() []string
}
