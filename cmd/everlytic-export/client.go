package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EverlyticClient struct {
	http.Client
	username, apiKey string
}

func (ec EverlyticClient) DownloadEverlyticCSV() (contacts []EverlyticContact, err error) {
	listID := 138207
	url := fmt.Sprintf("https://live.everlytic.net/api/2.0/list_subscriptions/%d", listID)
	for {
		log.Printf("Processing %s", url)
		newContacts, next, err := ec.downloadEverlyticCSV(url)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, newContacts...)
		if next == "" {
			return contacts, nil
		}
		url = next
	}
}

func (ec EverlyticClient) makeRequest(url string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(ec.username, ec.apiKey))
	resp, err := ec.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// usefult for debugging
	/*
		buf := new(bytes.Buffer)
		tr := io.TeeReader(resp.Body, buf)
		out, _ := ioutil.ReadAll(tr)
		log.Printf("result for %s:\n%s", url, out)
		d := json.NewDecoder(buf)
	*/
	d := json.NewDecoder(resp.Body)
	return d.Decode(result)
}

func (ec EverlyticClient) downloadEverlyticCSV(url string) (contacts []EverlyticContact, next string, err error) {
	var everlyticResponse EverlyticResponse
	err = ec.makeRequest(url, &everlyticResponse)
	if err != nil {
		return nil, "", err
	}
	result := make([]EverlyticContact, 0, len(everlyticResponse.Collection))
	for _, item := range everlyticResponse.Collection {
		contact, err := ec.loadContact(item.Data.ContactID)
		if err != nil {
			return nil, "", err
		}
		result = append(result, contact)
	}
	for _, l := range everlyticResponse.Links {
		if l.Rel == "next" {
			next = l.Href
		}
	}
	// fmt.Printf("body: %+v\n", everlyticResponse)
	return result, next, nil
}

func (ec EverlyticClient) loadContact(contactID int) (EverlyticContact, error) {
	var resp EverlyticContactResponse
	url := fmt.Sprintf("https://live.everlytic.net/api/2.0/contacts/%d", contactID)
	err := ec.makeRequest(url, &resp)
	return resp.Item, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
