package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"main.go/models"
)

func ParseOtodom(url string) ([]models.Flat, error) {
	doc, err := getHTMLDocumentByURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при парсинге HTML-документа")
	}

	flats, err := parseOtodomFlatsFromDoc(doc)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при парсинге объявлений Otodom")
	}

	return removeDuplicateFlats(flats), nil
}

func parseOtodomFlatsFromDoc(doc *goquery.Document) ([]models.Flat, error) {
	var result []models.Flat
	const site = "otodom.pl"

	liList := doc.Find("div[role='main'] > div[data-cy='search.listing.organic'] > ul > li")
	liList.Each(func(i int, li *goquery.Selection) {
		href, exists := li.Find("[data-cy='listing-item-link']").Attr("href")
		if !exists {
			return
		}

		title, exists := li.Find("[data-cy='listing-item-link']").Find("[data-cy='listing-item-title']").Attr("title")
		if !exists {
			return
		}

		result = append(result, models.NewFlat(site, title, fmt.Sprintf("https://www.otodom.pl%s", href)))
	})

	return result, nil
}
