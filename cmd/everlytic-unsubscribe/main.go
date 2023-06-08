package main

import (
	"flag"
	"log"
	"os"

	"github.com/godfried/mcsa-member-sync/csv"
	"github.com/godfried/mcsa-member-sync/everlytic"
)

const format = "2006-01-02_15-04-05"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	_ = wd
	var username, apiKey, source string

	fs := flag.NewFlagSet("everlytic-unsubscribe", flag.ExitOnError)
	fs.StringVar(&apiKey, "api-key", "", "Everlytic API key.")
	fs.StringVar(&username, "username", "", "Everlytic username.")
	fs.StringVar(&source, "source", source, "File of emails to unsubscribe.")
	fs.Parse(os.Args[1:])

	if username == "" {
		log.Fatal("Everlytic username is not set.")
	}
	if apiKey == "" {
		log.Fatal("Everlytic API Key is not set.")
	}
	log.SetOutput(os.Stderr)
	contacts, err := csv.ReadContacts(source, everlytic.LoadContact)
	ec := everlytic.NewClient(username, apiKey)
	_ = ec
	err = ec.UnsubscribeAll(contacts[0].ListID, contacts)
	if err != nil {
		log.Fatal(err)
	}
}
