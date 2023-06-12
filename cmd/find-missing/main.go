package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/godfried/mcsa-member-sync/csv"
	"github.com/godfried/mcsa-member-sync/everlytic"
	"github.com/godfried/mcsa-member-sync/membaz"
	"github.com/godfried/mcsa-member-sync/types"
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
	missingMembaz := findMissing(membazContacts, everlyticContacts)
	missingEverlytic := findMissing(everlyticContacts, membazContacts)

	err = csv.WriteContacts(missingMembaz, destinationMembaz)
	if err != nil {
		return err
	}
	err = csv.WriteContacts(missingEverlytic, destinationEverlytic)
	return err
}

func findMissing[O types.Contact, C types.Contact](oracle []O, check []C) []types.Contact {
	result := make([]types.Contact, 0, len(check))
	for _, checking := range check {
		found := false
		for _, contact := range oracle {
			if strings.EqualFold(checking.EmailAddress(), contact.EmailAddress()) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, checking)
		}
	}
	return result
}
