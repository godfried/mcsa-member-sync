package everlytic

import (
	"strconv"
	"strings"
)

type ListResponse struct {
	Links      []Link                   `json:"links"`
	Collection []SubscriptionCollection `json:"collection"`
}

type Link struct {
	Title string `json:"title"`
	Rel   string `json:"rel"`
	Href  string `json:"href"`
}

type SubscriptionCollection struct {
	Data Subscription `json:"data"`
}

type Subscription struct {
	ListID                  int         `json:"list_id"`
	ContactID               int         `json:"contact_id"`
	DateCreate              int         `json:"date_create"`
	DateOptIn               int         `json:"date_opt_in"`
	DateOptInRequest        int         `json:"date_opt_in_request"`
	EmailStatus             string      `json:"email_status"`
	MobileStatus            string      `json:"mobile_status"`
	PushStatus              interface{} `json:"push_status"`
	EmailUnsubscribedDate   int         `json:"email_unsubscribed_date"`
	MobileUnsubscribedDate  int         `json:"mobile_unsubscribed_date"`
	EmailUnsubscribedReason interface{} `json:"email_unsubscribed_reason"`
	ContactUniqueID         string      `json:"contact_unique_id"`
}

type ContactResponse struct {
	Links []Link  `json:"links"`
	Item  Contact `json:"item"`
}

type Contact struct {
	ID                               int     `json:"id"`
	CountryID                        int     `json:"country_id"`
	CityID                           int     `json:"city_id"`
	Name                             string  `json:"name"`
	Lastname                         string  `json:"lastname"`
	Email                            string  `json:"email"`
	Mobile                           string  `json:"mobile"`
	DateCreate                       int     `json:"date_create"`
	DateModified                     int     `json:"date_modified"`
	Status                           string  `json:"status"`
	Score                            float64 `json:"score"`
	Rating                           int     `json:"rating"`
	EmailStatus                      string  `json:"email_status"`
	SmsStatus                        string  `json:"sms_status"`
	EmailBounceHardCount             int     `json:"email_bounce_hard_count"`
	EmailBounceSoftCount             int     `json:"email_bounce_soft_count"`
	BlockBounceCount                 int     `json:"block_bounce_count"`
	SmsBounceHardCount               int     `json:"sms_bounce_hard_count"`
	SmsBounceSoftCount               int     `json:"sms_bounce_soft_count"`
	SmsBounceConsecutiveCount        int     `json:"sms_bounce_consecutive_count"`
	ComplaintsCount                  int     `json:"complaints_count"`
	ForwardCount                     int     `json:"forward_count"`
	InviteCount                      int     `json:"invite_count"`
	UpdateCount                      int     `json:"update_count"`
	BounceUnidentifiedCount          int     `json:"bounce_unidentified_count"`
	AutoresponderCount               int     `json:"autoresponder_count"`
	EmailMessageIds                  string  `json:"email_message_ids"`
	EmailLastSendDate                int     `json:"email_last_send_date"`
	SmsLastSendDate                  int     `json:"sms_last_send_date"`
	EmailLastOpenDate                int     `json:"email_last_open_date"`
	EmailLastClickDate               int     `json:"email_last_click_date"`
	SmsMessageIds                    string  `json:"sms_message_ids"`
	EmailMessageCount                int     `json:"email_message_count"`
	SmsMessageCount                  int     `json:"sms_message_count"`
	MessageReads                     int     `json:"message_reads"`
	MessageReadsInferred             int     `json:"message_reads_inferred"`
	MessageReadsUnique               int     `json:"message_reads_unique"`
	MessageLinkClicks                int     `json:"message_link_clicks"`
	MessageLinkClicksUnique          int     `json:"message_link_clicks_unique"`
	ContactAttachmentDownloadsUnique int     `json:"contact_attachment_downloads_unique"`
	ContactAttachmentDownloadsTotal  int     `json:"contact_attachment_downloads_total"`
	PreferredEmailFormat             string  `json:"preferred_email_format"`
	Title                            string  `json:"title"`
	CompanyPosition                  string  `json:"company_position"`
	CompanyName                      string  `json:"company_name"`
	Department                       string  `json:"department"`
	Industry                         string  `json:"industry"`
	Address                          string  `json:"address"`
	Address2                         string  `json:"address_2"`
	City                             string  `json:"city"`
	Country                          string  `json:"country"`
	State                            string  `json:"state"`
	Zip                              string  `json:"zip"`
	TelephoneOffice                  string  `json:"telephone_office"`
	TelephoneHome                    string  `json:"telephone_home"`
	TelephoneFax                     string  `json:"telephone_fax"`
	DateOfBirth                      int     `json:"date_of_birth"`
	BirthDate                        string  `json:"birth_date"`
	Gender                           string  `json:"gender"`
	MaritalStatus                    string  `json:"marital_status"`
	EducationLevel                   string  `json:"education_level"`
	Hash                             string  `json:"hash"`
	UniqueID                         string  `json:"unique_id"`
	ListID                           int
}

func (c Contact) EmailAddress() string {
	return c.Email
}

func (c Contact) Record() []string {
	return []string{
		strconv.Itoa(c.ID),
		c.Email,
		c.Name,
		c.Lastname,
		strconv.Itoa(c.ListID),
	}
}

func (c Contact) IsEmpty() bool {
	return strings.EqualFold(c.Email, "email") || strings.EqualFold(c.Email, "")
}

func LoadContact(record []string) (Contact, error) {
	cid, err := strconv.Atoi(record[0])
	if err != nil {
		return Contact{}, err
	}
	lid, err := strconv.Atoi(record[4])
	if err != nil {
		return Contact{}, err
	}
	return Contact{
		ID:       cid,
		Email:    record[1],
		Name:     record[2],
		Lastname: record[3],
		ListID:   lid,
	}, nil
}
