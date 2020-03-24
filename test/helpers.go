package main

import (
	"fmt"
	"strings"
)

func loadMockJson(filename string) ([]byte, error) {
	// TODO EJR where should the mockjsons be stored? or should the feature file just hold the full path somehow
	return nil, nil
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
