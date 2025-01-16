package scrapingutils

import (
	"net/http"
	"os"
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

func GetModulesData(moduleDoc *goquery.Selection) {
	moduleTitle := moduleDoc.Find(".ig-header-title.collapse_module_link.ellipsis").Text()
	moduleDoc.Find(".ig-title.title.item_link").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		videoLectureLink := "https://aridesa.instructure.com" + link
		lectureDoc := scrape.Parse(videoLectureLink, map[string]string{"Cookie": os.Getenv("COOKIES")})
		getLectureData(moduleTitle, lectureDoc)
	})
}
