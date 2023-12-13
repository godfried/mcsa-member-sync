package csv

import (
	"encoding/csv"
	"os"

	"github.com/godfried/mcsa-member-sync/contacts"
)

func WriteContacts[T contacts.Contact](contacts []T, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	cr := csv.NewWriter(f)
	for _, contact := range contacts {
		err = cr.Write(contact.Record())
		if err != nil {
			return err
		}
	}
	cr.Flush()
	return nil
}
