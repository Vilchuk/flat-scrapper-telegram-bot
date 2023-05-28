package parser

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"main.go/models"
)

func removeDuplicateFlats(flats []models.Flat) []models.Flat {
	allKeys := make(map[string]bool)
	var list []models.Flat
	for _, item := range flats {
		if _, value := allKeys[item.Href]; !value {
			allKeys[item.Href] = true
			list = append(list, item)
		}
	}
	return list
}

func getHTMLDocumentByURL(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения HTML-документа по URL")
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка статуса кода: %d %s", res.StatusCode, res.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(res.Body)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка чтения тела ответа")
	}

	// Загрузка HTML-документа
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "ошибка создания HTML-документа из ридера")
	}

	return doc, nil
}
