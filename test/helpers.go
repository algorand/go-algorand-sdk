package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func loadMockJsons(pathToJsons, commaDelimitedFilenames string) ([][]byte, error) {
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

func confirmErrorContainsString(err error, desired string) error {
	if desired == "nil" {
		return err
	}
	if strings.Contains(err.Error(), desired) {
		return nil
	} else {
		return fmt.Errorf("validated error did not contain expected substring, expected substring: %s,"+
			"actual error string: %s", desired, err.Error())
	}
}
