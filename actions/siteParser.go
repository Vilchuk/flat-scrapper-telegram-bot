package actions

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"main.go/Models"
	"net/http"
)

func OtodomGetFlats(url string) []Models.Flat {
	doc, _ := getHtmlDocumentByUrl(url)

	flats := []Models.Flat{}

	//regex := regexp.MustCompile(`"totalPages":\d\d?\d?\d?`)
	//pages, _ := strconv.Atoi(strings.Split(regex.FindString(body), ":")[1])

	flats = append(flats, parseOtodomFlatsFromDoc(doc)...)

	//if pages > 1 {
	//	for i := 2; i <= pages; i++ {
	//		doc, _ := getHtmlDocumentByUrl(uppendParamToQuery(url, "page", strconv.Itoa(i)))
	//		flats = append(flats, parseOtodomFlatsFromDoc(doc)...)
	//	}
	//}

	return removeDuplicateFlats(flats)
}

func OlxGetFlats(url string) []Models.Flat {
	doc, _ := getHtmlDocumentByUrl(url)

	flats := []Models.Flat{}

	flats = append(flats, parseOlxFlatsFromDoc(doc)...)

	return removeDuplicateFlats(flats)
}

func removeDuplicateFlats(flats []Models.Flat) []Models.Flat {
	allKeys := make(map[string]bool)
	list := []Models.Flat{}
	for _, item := range flats {
		if _, value := allKeys[item.Hash]; !value {
			allKeys[item.Hash] = true
			list = append(list, item)
		}
	}
	return list
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(bs)
}

func parseOtodomFlatsFromDoc(doc *goquery.Document) []Models.Flat {
	result := []Models.Flat{}

	doc.Find("div[role='main'] > div[data-cy='search.listing']").Find("[data-cy='listing-item-link']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title, _ := s.Find("[data-cy='listing-item-title']").Attr("title")

		result = append(result, Models.Flat{Hash: hash(href), Title: title, Href: fmt.Sprintf("https://www.otodom.pl%s", href)})
	})

	return result
}

func parseOlxFlatsFromDoc(doc *goquery.Document) []Models.Flat {
	result := []Models.Flat{}

	doc.Find("#offers_table").Find("[data-cy='listing-ad-title']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := s.Find("strong").Text()

		result = append(result, Models.Flat{Hash: hash(href), Title: title, Href: href})
	})

	return result
}

func getHtmlDocumentByUrl(url string) (*goquery.Document, string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bodyStr := string(data)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	return doc, bodyStr
}
