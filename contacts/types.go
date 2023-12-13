package contacts

import "strings"

type Contact interface {
	IsEmpty() bool
	EmailAddress() string
	Record() []string
}

type StringContact string

func (s StringContact) IsEmpty() bool {
	return strings.TrimSpace(string(s)) == ""
}

func (s StringContact) EmailAddress() string {
	return string(s)
}

func (s StringContact) Record() []string {
	return []string{s.EmailAddress()}
}
