package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/godfried/mcsa-member-sync/contacts"
	"github.com/godfried/mcsa-member-sync/contacts/everlytic"
	"github.com/godfried/mcsa-member-sync/contacts/membaz"
	"github.com/godfried/mcsa-member-sync/csv"
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

	err = run(membazCSV, everlyticCSV, destinationMembaz, destinationEverlytic)
	if err != nil {
		log.Fatal(err)
	}
}

func run(sourceMembaz, sourceEverlytic, destinationMembaz, destinationEverlytic string) (err error) {
	membazContacts, err := csv.ReadContacts(sourceMembaz, membaz.LoadContact)
	if err != nil {
		return err
	}
	everlyticContacts, err := csv.ReadContacts(sourceEverlytic, everlytic.LoadContact)
	if err != nil {
		return err
	}
	missingMembaz := contacts.FindMissing(membazContacts, everlyticContacts)
	missingEverlytic := contacts.FindMissing(everlyticContacts, membazContacts)

	err = csv.WriteContacts(missingMembaz, destinationMembaz)
	if err != nil {
		return err
	}
	err = csv.WriteContacts(missingEverlytic, destinationEverlytic)
	return err
}
