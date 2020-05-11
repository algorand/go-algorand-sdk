package main

import (
	"container/ring"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

func loadMockJsons(commaDelimitedFilenames, pathToJsons string) ([][]byte, error) {
	jsonFilenames := strings.Split(commaDelimitedFilenames, ",")
	var bytesArray [][]byte
	for _, jsonFilename := range jsonFilenames {
		fullPath := path.Join("./features/unit/", pathToJsons, jsonFilename)
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

func mockHttpResponsesInLoadedFrom(jsonfiles, directory string) error {
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
		_, err = w.Write(json)
		responseRing = responseRing.Next()
	}))
	return err
}

var receivedPath string

func mockServerRecordingRequestPaths() error {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath, _ = url.PathUnescape(r.URL.String())
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

func comparisonCheck(varname string, expected, actual interface{}) error {
	if expected != actual {
		return fmt.Errorf("expected %s value %v did not match actual value %v", varname, expected, actual)
	}
	return nil
}

func findAssetInHoldingsList(list []models.AssetHolding, desiredId uint64) (models.AssetHolding, error) {
	for _, holding := range list {
		if holding.AssetId == desiredId {
			return holding, nil
		}
	}
	return models.AssetHolding{}, fmt.Errorf("could not find asset ID %d in passed list of asset holdings", desiredId)
}

func getSigtypeFromTransaction(transaction models.Transaction) string {
	if len(transaction.Signature.Sig) != 0 {
		return "sig"
	}
	if len(transaction.Signature.Multisig.Subsignature) != 0 {
		return "msig"
	}
	if len(transaction.Signature.Logicsig.MultisigSignature.Subsignature) != 0 ||
		len(transaction.Signature.Logicsig.Logic) != 0 {
		return "lsig"
	}
	return "unknown sigtype"
}
