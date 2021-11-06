package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

var EnvNotFound = "Could not load .env"

func loadDotEnv(fileName string) {
	err := godotenv.Load(fileName)

	if err != nil {
		log.Fatal(EnvNotFound)
	}
}

func main() {
	loadDotEnv(".env")

	path := os.Getenv("NOTION_API_PATH")
	key := os.Getenv("NOTION_API_KEY")
	version := os.Getenv("NOTION_API_VERSION")
	page := os.Getenv("NOTION_PAGE_ID")

	notion := Notion{path, key, version}
	notion.Fetch(page)

}