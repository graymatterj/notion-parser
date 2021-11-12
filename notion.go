package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
)

const DatabaseType = "Database"
const BlockType = "Block"

var DatabasePath = "/databases/%s/query"
var BlockPath = "/blocks/%s/children"

var ErrTypeNotFound = errors.New("Error, could not find request type")

type Notion struct {
	path, key, version string
}

type NotionResponse struct {
	Object string `json:"object"`
	Results []Result `json:"results"`
}

type Result struct {
	Object string `json:"object"`
	Id string `json:"id"`
	LastEditedTime string `json:"last_edited_time"`
	HasChildren bool `json:"has_children"`
	Paragraph struct {
		Text []struct {
			PlainText string `json:"plain_text"`
		} `json:"text"`
	} `json:"paragraph,omitempty"`
	Properties struct {
		Processed struct {
			Checkbox bool `json:"checkbox"`
		} `json:"Processed"`
	} `json:"properties,omitemtpy"`
}
}

func (n Notion) Fetch(Type, ObjectId string) {
	requestPath, requestMethod, _ := buildRequestPath(Type, n.path, ObjectId)
	response := fetchData(n.key, n.version, requestPath, requestMethod)

	log.Println("Retrieved Response from API:")
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Println(string(bodyBytes))

	var result Response
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		println("Unmarshal 1 failed")
		println(err.Error())
	}

	fmt.Println(PrettyPrint(result))
}

func fetchData(ApiKey, ApiVersion, RequestPath, RequestType string) *http.Response{
	req, err := http.NewRequest(RequestType, RequestPath, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Default notion headers, required for every request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ApiKey)
	req.Header.Add("Notion-Version", ApiVersion)

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err.Error())
	}

	return response
}

func buildRequestPath(Type, Path, ObjectId string) (requestPath, requestMethod string, err error) {
	switch Type {
	case DatabaseType:
		println("Database Type")
		requestPath = Path + DatabasePath
		requestMethod = "POST"
	case BlockType:
		println("Block Type")
		requestPath = Path + BlockPath
		requestMethod = "GET"
	default:
		println("Could not find request type")
		err = ErrTypeNotFound
		return
	}
	requestPath = fmt.Sprintf(requestPath, ObjectId)
	return
}

func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}