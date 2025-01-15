package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	mongoClient "github.com/OtaviOuu/aridesa-scraper/internal/mongo"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type Lecture struct {
	Subject string `json:"subject,omitempty"`
	Module  string `json:"module,omitempty"`
	Title   string `json:"title,omitempty"`
	Link    string `json:"link,omitempty"`
	Year    int    `json:"year,omitempty"`
	Type    string `json:"type,omitempty"`
}

var client *mongo.Client

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

func insertStruct(lecture *Lecture) {
	collection := client.Database("scraper-ita").Collection("resourses")
	insertResult, err := collection.InsertOne(context.Background(), lecture)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(insertResult.InsertedID)
}

func getLectureData(moduleTitle string, lectureDoc *goquery.Document) {
	re := regexp.MustCompile(`https?://(www\.)?(youtube\.com|youtu\.be)/[^\s]+`)
	matches := re.FindAllString(lectureDoc.Text(), -1)
	if len(matches) > 0 {
		match := matches[0]

		subject := strings.TrimSpace(strings.ReplaceAll(lectureDoc.Find(".mobile-header-title.expandable").First().Text(), "\n", ""))
		formatedSubject := strings.Split(subject, "        ")[0]

		insertStruct(&Lecture{
			Module:  strings.TrimSpace(moduleTitle),
			Subject: formatedSubject,
			Title:   strings.TrimSpace(lectureDoc.Find("title").Text()),
			Link:    strings.TrimSpace(strings.ReplaceAll(match, "\\\"", "")),
			Year:    2023,
			Type:    "video",
		})

	} else {
		// Todo:
	}
}

func main() {
	err := godotenv.Load()
	client = mongoClient.GetMongoClient()
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
