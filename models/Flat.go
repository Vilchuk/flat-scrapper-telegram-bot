package models

import "time"

type Flat struct {
	FoundOnSite string
	DateCreated time.Time
	Title       string
	Href        string
}

func NewFlat(foundOnSite string, title string, href string) Flat {
	return Flat{
		FoundOnSite: foundOnSite,
		DateCreated: time.Now(),
		Title:       title,
		Href:        href}
}
