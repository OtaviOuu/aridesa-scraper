package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

func parse(url string, header map[string]string) *goquery.Document {

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

func getCoursesLinks(doc *goquery.Document) []string {
	links := []string{}

	doc.Find("tbody tr").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Find("a").Attr("href")
		if strings.Contains(link, "/courses/") {
			links = append(links, link)
		}
	})

	return links
}

func getModulesData(moduleDoc *goquery.Selection) {
	moduleTitle := moduleDoc.Find(".name").Text()
	moduleDoc.Find(".ig-title.title.item_link").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		videoLectureLink := "https://aridesa.instructure.com" + link
		lectureDoc := parse(videoLectureLink, map[string]string{"Cookie": os.Getenv("COOKIES")})

		getLectureData(moduleTitle, lectureDoc)
	})

}

func getLectureData(moduleTitle string, lectureDoc *goquery.Document) {
	lectureTitle := lectureDoc.Find(".title").Text()
	lectureDoc.Find(".video-js").Each(func(index int, item *goquery.Selection) {
		videoLink, _ := item.Attr("src")
		log.Println(moduleTitle, lectureTitle, videoLink)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cookies := os.Getenv("COOKIES")
	doc := parse("https://aridesa.instructure.com/courses", map[string]string{"Cookie": cookies})

	coursesLinks := getCoursesLinks(doc)

	for _, link := range coursesLinks {
		url := "https://aridesa.instructure.com" + link
		doc := parse(url, map[string]string{"Cookie": cookies})

		doc.Find("#context_modules .item-group-condensed.context_module").Each(func(index int, item *goquery.Selection) {
			getModulesData(item)
		})
	}
}
