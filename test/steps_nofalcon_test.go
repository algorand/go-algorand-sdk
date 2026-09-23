//go:build !falcon

package test

import (
	"fmt"
)

func genFalconKey() error {
	return fmt.Errorf("sdk built without falcon support")
}

func loadFalconKey() error {
	return fmt.Errorf("sdk built without falcon support")
}

func mnForFalcon(_ string) error {
	return fmt.Errorf("sdk built without falcon support")
}
