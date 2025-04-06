package models

type Event struct {
	Id           string `json:"UID"`
	Dtstamp      string `json:"DTSTAMP"`
	Dtstart      string `json:"DTSTART"`
	Dtend        string `json:"DTEND"`
	Description  string `json:"DESCRIPTION"`
	Location     string `json:"LOCATION"`
	Created      string `json:"CREATED"`
	LastModified string `json:"LAST-MODIFIED"`
	ResourceID   string `json:"RESOURCE-ID"`
}
