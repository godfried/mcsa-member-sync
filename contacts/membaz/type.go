package membaz

import "strings"

type Contact struct {
	Name               string
	Surname            string
	Email              string
	MembershipNumber   string
	HutKeyNumber       string
	MembershipCategory string
}

func (c Contact) Record() []string {
	return []string{
		c.Name,
		c.Surname,
		c.Email,
		c.MembershipNumber,
		c.HutKeyNumber,
		c.MembershipCategory,
	}
}

func (c Contact) EmailAddress() string {
	return c.Email
}

func (c Contact) IsEmpty() bool {
	return strings.EqualFold(c.Email, "email") || strings.EqualFold(c.Email, "")
}

func LoadContact(record []string) (Contact, error) {
	return Contact{
		Name:               record[0],
		Surname:            record[1],
		Email:              record[2],
		MembershipNumber:   record[3],
		HutKeyNumber:       record[4],
		MembershipCategory: record[5],
	}, nil
}
