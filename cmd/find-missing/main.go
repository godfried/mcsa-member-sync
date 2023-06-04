package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const format = "2006-01-02_15-04-05"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var membazCSV, everlyticCSV string
	destinationMembaz := filepath.Join(wd, fmt.Sprintf("missing-membaz-%s.csv", time.Now().Format(format)))
	destinationEverlytic := filepath.Join(wd, fmt.Sprintf("missing-everlytic-%s.csv", time.Now().Format(format)))

	fs := flag.NewFlagSet("member-sync", flag.ExitOnError)
	fs.StringVar(&destinationMembaz, "membaz-destination", destinationMembaz, "Missing members from Membaz file destination.")
	fs.StringVar(&destinationEverlytic, "everlytic-destination", destinationEverlytic, "Missing members from Everlytic file destination.")
	fs.StringVar(&membazCSV, "membaz-csv", membazCSV, "Membaz CSV source.")
	fs.StringVar(&everlyticCSV, "everlytic-csv", everlyticCSV, "Everlytic CSV source.")
	fs.Parse(os.Args[1:])

	membazContacts, err := loadMembaz(membazCSV)
	if err != nil {
		log.Fatal(err)
	}
	everlyticContacts, err := loadEverlytic(everlyticCSV)
	if err != nil {
		log.Fatal(err)
	}
	missingMembaz := findMissing(membazContacts, everlyticContacts)
	missingEverlytic := findMissing(everlyticContacts, membazContacts)
	err = writeContacts(missingEverlytic, destinationEverlytic)
	if err != nil {
		log.Fatal(err)
	}
	err = writeContacts(missingMembaz, destinationMembaz)
	if err != nil {
		log.Fatal(err)
	}
}

func writeContacts(contacts []contact, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	cr := csv.NewWriter(f)
	for _, contact := range contacts {
		err = cr.Write(contact.record())
		if err != nil {
			return err
		}
	}
	cr.Flush()
	return nil
}

func findMissing(oracle, check []contact) []contact {
	missing := make([]contact, 0, len(check))
	for _, checking := range check {
		found := false
		for _, contact := range oracle {
			if strings.EqualFold(checking.email, contact.email) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, checking)
		}
	}
	return missing
}

type contact struct {
	email       string
	first, last string
}

func (c contact) record() []string {
	return []string{
		c.email,
		c.first,
		c.last,
	}
}

func loadEverlytic(source string) ([]contact, error) {
	return loadCSV(source, func(record []string) contact {
		return contact{
			email: record[1],
			first: record[2],
			last:  record[3],
		}
	})
}

func loadMembaz(source string) ([]contact, error) {
	return loadCSV(source, func(record []string) contact {
		return contact{
			first: record[0],
			last:  record[1],
			email: record[2],
		}
	})

}

func loadCSV(source string, loadFunc func(record []string) contact) ([]contact, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cr := csv.NewReader(f)
	contacts := make([]contact, 0, 4096)
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		contact := loadFunc(record)
		if !strings.EqualFold(contact.email, "email") && !strings.EqualFold(contact.email, "") {
			contacts = append(contacts, contact)
		}
	}
	return contacts, nil
}
