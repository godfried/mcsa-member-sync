package csv

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/godfried/mcsa-member-sync/contacts"
)

func ReadContacts[T contacts.Contact](source string, loadContact func(record []string) (T, error)) ([]T, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cr := csv.NewReader(f)
	contacts := make([]T, 0, 4096)
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		contact, err := loadContact(record)
		if err != nil {
			return nil, err
		}
		if !contact.IsEmpty() {
			contacts = append(contacts, contact)
		}
	}
	return contacts, nil
}
