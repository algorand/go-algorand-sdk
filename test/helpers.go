package main

import (
	"container/ring"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"

	"github.com/cucumber/godog"
)

func loadMockJsons(commaDelimitedFilenames, pathToJsons string) ([][]byte, error) {
	jsonFilenames := strings.Split(commaDelimitedFilenames, ",")
	var bytesArray [][]byte
	for _, jsonFilename := range jsonFilenames {
		fullPath := path.Join(pathToJsons, jsonFilename)
		jsonfile, err := os.Open(fullPath)
		if err != nil {
			return nil, err
		}
		fileBytes, err := ioutil.ReadAll(jsonfile)
		if err != nil {
			return nil, err
		}
		bytesArray = append(bytesArray, fileBytes)
	}
	return bytesArray, nil
}

var mockServer *httptest.Server
var responseRing *ring.Ring

func mockHttpResponsesInLoadedFrom(jsonfiles, directory string) error {
	jsons, err := loadMockJsons(jsonfiles, directory)
	if err != nil {
		return err
	}
	responseRing = ring.New(len(jsons))
	for _, json := range jsons {
		responseRing.Value = json
	}
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json := responseRing.Value.([]byte)
		_, err = w.Write(json)
		responseRing = responseRing.Next()
	}))
	return err
}

var receivedPath string

func mockServerRecordingRequestPaths() error {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
	}))
	return nil
}
func expectThePathUsedToBe(expectedPath string) error {
	if receivedPath != expectedPath {
		return fmt.Errorf("path used to access mock server was %s but expected path %s", receivedPath, expectedPath)
	}
	return nil
}

var globalErrForExamination error

func expectErrorStringToContain(contains string) error {
	if contains == "nil" {
		if globalErrForExamination != nil {
			return fmt.Errorf("expected no error but error was found: %s", globalErrForExamination.Error())
		}
		return nil
	}
	if strings.Contains(globalErrForExamination.Error(), contains) {
		return nil
	}
	return fmt.Errorf("validated error did not contain expected substring, expected substring: %s,"+
		"actual error string: %s", contains, globalErrForExamination.Error())
}
