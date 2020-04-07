package main

import (
	"fmt"
	"io/ioutil"
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

var mockServer httptest.Server

func mockHttpResponsesInLoadedFrom(jsonfiles, directory string) error {
	jsons, err := loadMockJsons(jsonfiles, directory)
	if err != nil {
		return err
	}
	// TODO mockServer is constructed and will return each json in jsons on subsequent calls
	return godog.ErrPending
}

func expectThePathUsedToBe(expectedPath string) error {
	// TODO request recorder is used and records path used
	return godog.ErrPending
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
