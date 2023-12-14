package contacts

import (
	"fmt"
	"strings"
)

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

func LoadString(record []string) (StringContact, error) {
	if len(record) != 1 {
		return "", fmt.Errorf("Unexpected record length %d, expected one element. Record: %v", len(record), record)
	}
	return StringContact(record[0]), nil
}
