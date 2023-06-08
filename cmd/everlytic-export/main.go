package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"encoding/csv"

	"github.com/godfried/mcsa-member-sync/everlytic"
)

const format = "2006-01-02_15-04-05"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var apiKey, username string
	destination := filepath.Join(wd, fmt.Sprintf("everlytic-export-%s.csv", time.Now().Format(format)))

	fs := flag.NewFlagSet("everlytic-export", flag.ExitOnError)
	fs.StringVar(&apiKey, "api-key", "", "Everlytic API key.")
	fs.StringVar(&username, "username", "", "Everlytic username.")
	fs.StringVar(&destination, "destination", destination, "Everlytic file destination.")
	fs.Parse(os.Args[1:])

	if username == "" {
		log.Fatal("Everlytic username is not set.")
	}
	if apiKey == "" {
		log.Fatal("Everlytic API Key is not set.")
	}
	log.SetOutput(os.Stderr)
	ec := everlytic.NewClient(username, apiKey)
	contacts, err := ec.DownloadEverlyticCSV()
	if err != nil {
		log.Fatal(err)
	}
	err = writeEverlyticCSV(contacts, destination)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote Everlytic export to %s.", destination)
}

func writeEverlyticCSV(contacts []everlytic.Contact, csvOutput string) error {
	f, err := os.Create(csvOutput)
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
