package test

import (
	"container/ring"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
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
		if strings.HasSuffix(fullPath, "base64") {
			fileBytesAsBase64String := string(fileBytes)
			decodedBytes, err := base64.StdEncoding.DecodeString(fileBytesAsBase64String)
			if err != nil {
				return nil, err
			}
			bytesArray = append(bytesArray, decodedBytes)
		} else {
			bytesArray = append(bytesArray, fileBytes)
		}
	}
	return bytesArray, nil
}

var mockServer *httptest.Server
var responseRing *ring.Ring

func mockHTTPResponsesInLoadedFromHelper(jsonfiles, directory string, status int) error {
	jsons, err := loadMockJsons(jsonfiles, directory)
	if err != nil {
		return err
	}
	responseRing = ring.New(len(jsons))
	for _, json := range jsons {
		responseRing.Value = json
		responseRing = responseRing.Next()
	}
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json := responseRing.Value.([]byte)
		if status > 0 {
			w.WriteHeader(status)
		}
		_, err = w.Write(json)
		responseRing = responseRing.Next()
	}))
	return err
}

var receivedMethod string
var receivedPath string

func mockServerRecordingRequestPaths() error {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.String()
	}))
	return nil
}

func expectTheRequestToBe(expectedMethod, expectedPath string) error {
	if !strings.EqualFold(expectedMethod, receivedMethod) {
		return fmt.Errorf("method used to access mock server was %s but expected %s", receivedMethod, expectedMethod)
	}
	return expectThePathUsedToBe(expectedPath)
}

func expectThePathUsedToBe(expectedPath string) error {
	if receivedPath != expectedPath {
		return fmt.Errorf("path used to access mock server was %s but expected path %s", receivedPath, expectedPath)
	}
	return nil
}

var globalErrForExamination error

func expectErrorStringToContain(contains string) error {
	if contains == "" {
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

func loadResource(filepath string) ([]byte, error) {
	return ioutil.ReadFile(path.Join("features", "resources", filepath))
}
