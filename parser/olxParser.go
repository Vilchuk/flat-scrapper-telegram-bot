package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"main.go/models"
	"strings"
)

func ParseOlx(url string) ([]models.Flat, error) {
	doc, err := getHTMLDocumentByURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при парсинге HTML-документа")
	}

	flats, err := parseOlxFlatsFromDoc(doc)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при парсинге объявлений OLX")
	}

	return removeDuplicateFlats(flats), nil
}

func parseOlxFlatsFromDoc(doc *goquery.Document) ([]models.Flat, error) {
	var result []models.Flat
	const site = "olx.pl"

	doc.Find(".listing-grid-container").Find("div[data-cy='l-card']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}

		title := s.Find("h6").Text()

		if href == "" && title == "" {
			return
		}

		if !strings.Contains(href, "https://www.otodom.pl") {
			href = fmt.Sprintf("https://www.olx.pl%s", href)
		}

		result = append(result, models.NewFlat(site, title, href))
	})

	return result, nil
}
