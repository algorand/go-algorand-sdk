package test

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/nsf/jsondiff"
)

// EqualJson compares two json strings.
// returns true if considered equal, false otherwise.
// The error returns the difference.
// For reference: j1 is the baseline, j2 is the test
func EqualJson(j1, j2 string) (ans bool, err error) {

	// Sorting the elements inside json arrays is important to avoid false
	// negatives due to element order mismatch.
	j1, err = SortJsonArrays(j1)
	if err != nil {
		return false, err
	}

	j2, err = SortJsonArrays(j2)
	if err != nil {
		return false, err
	}

	// This func uses jsondiff.
	// options configures it.
	options := jsondiff.Options{
		Added:            jsondiff.Tag{Begin: "___ADDED___", End: ""},
		Removed:          jsondiff.Tag{Begin: "___REMED___", End: ""},
		Changed:          jsondiff.Tag{Begin: "___DIFFER___", End: ""},
		ChangedSeparator: " -> ",
	}
	options.PrintTypes = false
	d, str := jsondiff.Compare([]byte(j1), []byte(j2), &options)
	// If fully match, return true
	if d == jsondiff.FullMatch {
		return true, nil
	}

	// Check the difference, and decide if it is still accepted as "equal"
	// Scan line by line, and examine the difference
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		// If baseline has the property (j1), but the test (j2) has it removed
		if strings.Contains(line, "___REMED___") {
			// If the value of the property is false, then accept it
			// The default value is false. Missing property is considered as false
			propValue := strings.SplitN(line, ": ", 2)
			if len(propValue) == 1 {
				addedLine := strings.SplitN(propValue[0], "___REMED___", 2)
				if len(addedLine) == 2 && (addedLine[1] == "}," || addedLine[1] == "}") {
					continue
				}
			}
			if len(propValue) == 2 &&
				(propValue[1] == "null," ||
					propValue[1] == "[]," || propValue[1] == "[]" ||
					propValue[1] == "false," || propValue[1] == "false" ||
					propValue[1] == "0," || propValue[1] == "0" ||
					propValue[1] == "\"\"," || propValue[1] == "\"\"" ||
					propValue[1] == "{") {
				continue
			}
			// Any other ommision will not be considered as equal
			err = fmt.Errorf(line)
			return false, err
		}
		// If the test has properties not found in the baseline, then they are not equal
		if strings.Contains(line, "___ADDED___") {
			//			fmt.Printf("XXX j1:\n%s\n", j1)
			//			fmt.Printf("XXX j2:\n%s\n", j2)
			err = fmt.Errorf(line)
			return false, err
		}
		// If the properties are different, they they are not equal
		if strings.Contains(line, "___DIFFER___") {
			err = fmt.Errorf(line)
			return false, err
		}
	}
	// Extra layer of check. This is not expected to be true.
	// If only difference is removed with false value property, then the difference
	// will be qualified as SupersetMatch. If not, should be be considered equal.
	if d != jsondiff.SupersetMatch {
		err = fmt.Errorf(str)
		return false, err
	}
	return true, nil
}

// SortJsonArrays sorts the elements inside json arrays.
// Some new line organization may be broken.
// e.g. ...[  {a: "a3", b: "b3"},   {a: "a1", b: "b1"},    {a: "a2", b: "b2"}]...
// will be:
//      ...[{a: "a1", b: "b1"},{a: "a2", b: "b2"},{a: "a3", b: "b3"}]...
// Note: any characters between [ and {, as well as between elements will be lost: ("...},xxx{..." "xxx" is lost ) 
// This function will recursively sort arrays inside arrays as well.
func SortJsonArrays(json string) (sorted string, err error) {
	idx := 0
	for idx != len(json) {
		if json[idx] == '[' {
			c, err := sortElements(&json, idx+1, len(json), '[', ']')
			if err != nil {
				return "", err
			}
			idx = c
		}
		if json[idx] == '{' {
			c, err := matchingCloseIndex(&json, idx+1, len(json), '{', '}')
			if err != nil {
				return "", err
			}
			idx = c
		}			
		idx++
	}
	return json, nil
}

// return the index of the character that closes the bracket.
// If an array is detected, will call sortElements to sort the elements inside the array
func matchingCloseIndex(json *string, start, end int, ochar, cchar rune) (closeIndex int, err error) {

	for i := start; i < len(*json); i++ {
		if (*json)[i] == '[' {
			i, err = sortElements(json, i+1, end, '[', ']')
			if err != nil {
				return -1, err
			}
			continue
		}
		if (*json)[i] == '{' {
			c, err := matchingCloseIndex(json, i+1, end, ochar, cchar)
			if err != nil {
				return -1, err
			}
			i = c
			continue
		}
		if (*json)[i] == '}' {
			// found the match
			return i, nil
		}
		if (*json)[i] == ']' {
			return -1, fmt.Errorf("Could not find the matching %c", cchar)
		}
	}
	return -1, fmt.Errorf("Could not find the matching %c", cchar)	
	//	return -1, fmt.Errorf("Could not find the matching %c", cchar)
}

// sorts the elements inside the array
// if a nested arrray is detected, it will first sort the elements inside the nested array
func sortElements(json *string, start, end int, ochar, cchar rune) (closeIndex int, err error) {
	elements := make([]string, 0)
	for i := start; i < len(*json); i++ {
		if (*json)[i] == '[' {
			// nested array detected. sort it first
			i, err = sortElements(json, i+1, end, '[', ']')
			if err != nil {
				return -1, err
			}
			continue
		}
		if (*json)[i] == '{' {
			c, err := matchingCloseIndex(json, i+1, end, '{', '}')
			if err != nil {
				return -1, err
			}
			elements = append(elements, (*json)[i:c+1])
			i = c
			continue
		}
		if (*json)[i] == '}' {
			// expects ] but found }. Note: } should have been handled by matchingCloseIndex
			return -1, fmt.Errorf("Could not find the matching A %c", cchar)
		}
		if (*json)[i] == ']' {
			// end of the array is detected. Sort the elements in the array
			sort.Slice(elements, func(i, j int) bool { return strings.TrimSpace(elements[i]) < strings.TrimSpace(elements[j]) })
			// replace the sorted string with the unsorted one
			// note that, any characters between [ and { are lost, as well as between elements i.e. } and {
			processedStr := (*json)[0:start] + strings.Join(elements, ",")
			retPos := len(processedStr)
			*json = processedStr + (*json)[i:]
			return retPos, nil
		}
	}
	return -1, fmt.Errorf("Could not find the matching %c", cchar)
}
