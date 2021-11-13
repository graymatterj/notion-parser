package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
)

const DatabaseType = "Database"
const BlockType = "Block"
const PageType = "Page"

const DatabasePath = "/databases/%s/query"
const BlockPath = "/blocks/%s/children"
const PagePath = "/pages/%s"

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
	} `json:"paragraph"`
	Properties Properties `json:"properties`
}

type RequestBody struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	Processed Processed `json:"Processed",omitempty`
}

type Processed struct {
	Checkbox bool `json:"checkbox"`
}

type QueryParams struct {
	Sorts []Sort `json:"sorts"`
}

type Sort struct {
	PropertyName string `json:"property"`
	Direction string `json:"direction"`
}

func (n Notion) updatePageStatus(id string) {
	requestPath, requestMethod, _ := buildRequestPath(PageType, n.path, id)

	properties := RequestBody{
		Properties: Properties{
			Processed: Processed{true},
		},
	}

	propertiesJSON, _ := json.Marshal(properties)

	response := n.sendRequest(requestPath, requestMethod, propertiesJSON)

	log.Println("Retrieved Response from API:")
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Println(string(bodyBytes))

}

func (n Notion) Fetch(Type, ObjectId string) {
	requestPath, requestMethod, _ := buildRequestPath(Type, n.path, ObjectId)

	var response *http.Response

	if Type == DatabaseType {
		query := QueryParams{
			Sorts: []Sort{
				Sort{"Lesson Date", "descending"},
			},
		}

		sortsJSON, _ := json.Marshal(query)
		response = n.sendRequest(requestPath, requestMethod, sortsJSON)
	} else {
		response = n.sendRequest(requestPath, requestMethod, nil)
	}


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

func (n Notion) sendRequest(RequestPath, RequestType string, RequestBody []byte) *http.Response {
	req, err := http.NewRequest(RequestType, RequestPath, bytes.NewBuffer(RequestBody))

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
	case PageType:
		println("Page Type")
		requestPath = Path + PagePath
		requestMethod = "PATCH"
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