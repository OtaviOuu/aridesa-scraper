package scrapingutils

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Parse(url string, header map[string]string) *goquery.Document {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Cookie", header["Cookie"])

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	return doc
}

func GetCoursesLinks(doc *goquery.Document) []string {
	links := []string{}

	doc.Find("tbody tr").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Find("a").Attr("href")
		if strings.Contains(link, "/courses/") {
			links = append(links, link)
		}
	})

	return links
}
