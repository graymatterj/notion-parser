package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var BlockPath = "/blocks/%s/children"

type Notion struct {
	path, key, version string
}

func (n Notion) Fetch(PageId string) {
	requestPath := fmt.Sprintf(n.path + BlockPath, PageId)
	req, err := http.NewRequest("GET", requestPath, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Default notion headers, required for every request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", n.key)
	req.Header.Add("Notion-Version", n.version)

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Retrieved Response from API:")
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Println(string(bodyBytes))

}
