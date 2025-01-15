# Aridesa Course Scraper

A Go-based web scraper designed to collect and store course content from Aridesa's Canvas LMS platform. This tool specifically targets course materials for ITA (Instituto Tecnológico de Aeronáutica) and IME (Instituto Militar de Engenharia) preparation courses.

## Features

- Scrapes course content from Aridesa's Canvas LMS
- Extracts video lectures and their metadata
- Stores data in MongoDB
- Organized by subjects, modules, and lecture types

## Prerequisites

- Go 1.x
- MongoDB Atlas account
- Access to Aridesa's course platform

## Environment Variables

Create a `.env` file in the root directory with the following variables (see `.env.example` for a template):

```
COOKIES=your_canvas_session_cookie
MONGO_USER=your_mongodb_username
MONGO_PASS=your_mongodb_password
```

## Data Structure

The scraper collects the following information for each lecture:

```go
type Lecture struct {
    Subject string `json:"subject,omitempty"`
    Module  string `json:"module,omitempty"`
    Title   string `json:"title,omitempty"`
    Link    string `json:"link,omitempty"`
    Year    int    `json:"year,omitempty"`
    Type    string `json:"type,omitempty"`
}
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/OtaviOuu/aridesa-scraper.git
cd aridesa-scraper
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your environment variables in `.env`

4. Run the scraper:
```bash
go run main.go
```

## Dependencies

- github.com/PuerkitoBio/goquery - HTML parsing
- github.com/joho/godotenv - Environment variable management
- go.mongodb.org/mongo-driver - MongoDB operations

## Project Structure

```
.
├── main.go                 # Main scraper logic
├── internal
│   └── mongo
│       └── client.go      # MongoDB connection handling
├── go.mod
├── go.sum
├── .env
└── .env.example          # Example environment variables template
```

## Warning

This tool is designed for educational purposes. Make sure you have the necessary permissions to access and scrape the content from the Aridesa platform.
