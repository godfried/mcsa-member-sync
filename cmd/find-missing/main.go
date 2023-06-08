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

	membazContacts, err := csv.ReadContacts(membazCSV, membaz.LoadContact)
	if err != nil {
		log.Fatal(err)
	}
	everlyticContacts, err := csv.ReadContacts(everlyticCSV, everlytic.LoadContact)
	if err != nil {
		log.Fatal(err)
	}
	missingMembaz := findMissing(membazContacts, everlyticContacts, make([]everlytic.Contact, 0, len(everlyticContacts)))
	missingEverlytic := findMissing(everlyticContacts, membazContacts, make([]membaz.Contact, 0, len(membazContacts)))
	err = csv.WriteContacts(missingEverlytic, destinationEverlytic)
	if err != nil {
		log.Fatal(err)
	}
	err = csv.WriteContacts(missingMembaz, destinationMembaz)
	if err != nil {
		log.Fatal(err)
	}
}

func findMissing[O types.Contact, C types.Contact](oracle []O, check, result []C) []C {
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
