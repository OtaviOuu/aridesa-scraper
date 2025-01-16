package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	mongoClient "github.com/OtaviOuu/aridesa-scraper/internal/mongo"
	scrape "github.com/OtaviOuu/aridesa-scraper/internal/scraping_utils"
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

func main() {
	err := godotenv.Load()
	client = mongoClient.GetMongoClient()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cookies := os.Getenv("COOKIES")

	doc := scrape.Parse("https://aridesa.instructure.com/courses", map[string]string{"Cookie": cookies})

	coursesLinks := scrape.GetCoursesLinks(doc)

	for _, link := range coursesLinks {
		url := "https://aridesa.instructure.com" + link
		doc := scrape.Parse(url, map[string]string{"Cookie": cookies})

		doc.Find("#context_modules .item-group-condensed.context_module").Each(func(index int, item *goquery.Selection) {
			getModulesData(item)
		})
	}
}

func getModulesData(moduleDoc *goquery.Selection) {
	moduleTitle := moduleDoc.Find(".ig-header-title.collapse_module_link.ellipsis").Text()
	moduleDoc.Find(".ig-title.title.item_link").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		videoLectureLink := "https://aridesa.instructure.com" + link
		lectureDoc := scrape.Parse(videoLectureLink, map[string]string{"Cookie": os.Getenv("COOKIES")})
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
