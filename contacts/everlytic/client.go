package everlytic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	MCSAMembersList = 138207
	baseURL         = "https://live.everlytic.net/api/2.0"
)

func NewClient(username, apiKey string) Client {
	return Client{
		username: username,
		apiKey:   apiKey,
	}
}

type Client struct {
	http.Client
	username, apiKey string
}

func (ec Client) DownloadEverlyticCSV() (contacts []Contact, err error) {
	url := fmt.Sprintf("%s/list_subscriptions/%d", baseURL, MCSAMembersList)
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

func (ec Client) UnsubscribeAll(contacts []Contact) error {
	type contact struct {
		ListID      int    `json:"list_id"`
		ContactID   int    `json:"contact_id"`
		EmailStatus string `json:"email_status"`
	}
	//unsubscriptions := make(map[string]contact, len(contacts))
	for _, c := range contacts {
		//unsubscriptions[strconv.Itoa(c.ID)] = contact{
		unsub := contact{
			ContactID:   c.ID,
			ListID:      c.ListID,
			EmailStatus: "unsubscribed",
		}
		//}
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(unsub)
		if err != nil {
			return err
		}
		url := fmt.Sprintf("%s/list_subscriptions/%d", baseURL, MCSAMembersList)
		var result json.RawMessage
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		err = ec.makeRequest(http.MethodPost, url, body, &result, headers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ec Client) makeGetRequest(url string, result interface{}, headers map[string]string) error {
	return ec.makeRequest(http.MethodGet, url, nil, result, headers)
}

func (ec Client) makeRequest(method, url string, body io.Reader, result interface{}, headers map[string]string) error {
	req, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(ec.username, ec.apiKey))
	for k, v := range headers {
		req.Header.Add(k, v)
	}
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

func (ec Client) downloadEverlyticCSV(url string) (contacts []Contact, next string, err error) {
	var everlyticResponse ListResponse
	err = ec.makeGetRequest(url, &everlyticResponse, nil)
	if err != nil {
		return nil, "", err
	}
	result := make([]Contact, 0, len(everlyticResponse.Collection))
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

func (ec Client) loadContact(contactID int) (Contact, error) {
	var resp ContactResponse
	url := fmt.Sprintf("%s/contacts/%d", baseURL, contactID)
	err := ec.makeGetRequest(url, &resp, nil)
	return resp.Item, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
