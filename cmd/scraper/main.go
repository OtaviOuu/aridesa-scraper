package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

type Lecture struct {
	Subject string `json:"subject,omitempty"`
	Module  string `json:"module,omitempty"`
	Title   string `json:"title,omitempty"`
	Link    string `json:"link,omitempty"`
	Year    int    `json:"year,omitempty"`
	Type    string `json:"type,omitempty"`
	Pdf     string `json:"pdf,omitempty"`
}

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
	moduleTitle := moduleDoc.Find(".ig-header-title.collapse_module_link.ellipsis").Text()
	moduleDoc.Find(".ig-title.title.item_link").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		videoLectureLink := "https://aridesa.instructure.com" + link
		lectureDoc := parse(videoLectureLink, map[string]string{"Cookie": os.Getenv("COOKIES")})
		getLectureData(moduleTitle, lectureDoc)
	})
}

func insertStruct(data *Lecture) {
	file, err := os.OpenFile("lectures.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lectureJSON, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	file.Write(lectureJSON)
	_, err = file.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}
}

func getLectureData(moduleTitle string, lectureDoc *goquery.Document) {
	l := &Lecture{}
	re := regexp.MustCompile(`https?://(www\.)?(youtube\.com|youtu\.be)/[^\s]+`)
	matches := re.FindAllString(lectureDoc.Text(), -1)
	if len(matches) > 0 {
		match := matches[0]

		l.Link = strings.TrimSpace(strings.ReplaceAll(match, "\\\"", ""))
		l.Subject = strings.TrimSpace(strings.ReplaceAll(lectureDoc.Find(".mobile-header-title.expandable").First().Text(), "\n", ""))
		l.Module = strings.TrimSpace(moduleTitle)
		l.Title = strings.TrimSpace(lectureDoc.Find("title").Text())
		l.Type = "video"
		l.Year = 2023

	} else {
		// Todo:
	}
	insertStruct(l)
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
