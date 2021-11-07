package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
)

const BlockType = "Block"

var BlockPath = "/blocks/%s/children"

var ErrTypeNotFound = errors.New("Error, could not find request type")

type Notion struct {
	path, key, version string
}

func (n Notion) Fetch(Type, ObjectId string) {
	requestPath, requestMethod, _ := buildRequestPath(Type, n.path, ObjectId)
	response := fetchData(n.key, n.version, requestPath, requestMethod)

	log.Println("Retrieved Response from API:")
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Println(string(bodyBytes))
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
