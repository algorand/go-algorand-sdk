package test

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
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
			err = fmt.Errorf(str)
			return false, err
		}
		// If the test has properties not found in the baseline, then they are not equal
		if strings.Contains(line, "___ADDED___") {
			err = fmt.Errorf(str)
			return false, err
		}
		// If the properties are different, they they are not equal
		if strings.Contains(line, "___DIFFER___") {
			err = fmt.Errorf(str)
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
			return -1, fmt.Errorf("Could not find the matching %c", cchar)
		}
		if (*json)[i] == ']' {
			// end of the array is detected. Sort the elements in the array
			sort.Slice(elements, func(i, j int) bool {
				valA := removeDefaultValues(elements[i])
				valB := removeDefaultValues(elements[j])
				return strings.TrimSpace(valA) < strings.TrimSpace(valB)
			})
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

func removeDefaultValues(data string) string {
	data = removeSpaces(data)
	startingLen := len(data)
	for x := 0; x < len(data); x++ {
		// find key start
		for ; x < len(data) && data[x] != '"'; x++ {
		}
		if x == len(data) {
			break
		}

		// get the key
		column, err := getString(&data, x, true)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		dv, lastIndex, err := isDefaultValue(&data, column)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		if dv {
			// x where the key starts for the defalut value element
			// lastIndex is where the value ends of the defalut value
			// it will be the index of ',' if one exists
			data = data[0:x] + data[lastIndex+1:]
			// x should point to the start of the new element
			// since the for loop will increment x by 1, and now
			// it is already pointing at the start, decrement by 1
			x--
		} else {
			x = column
			// advance until:
			// 1. the value is a quoted value: stop at the start '"'
			// 2. the value is a unquoted premitive: stop at the ',' ending it
			// 3. the value is another '{' block. stop at '{'
			for ; x < len(data) &&
				data[x] != '"' &&
				data[x] != ',' &&
				data[x] != '{'; x++ {
			}
			if x == len(data) {
				break
			}
			if data[x] == '"' {
				// read the quote value
				x, err = getString(&data, x, false)
				if err != nil {
					log.Fatal(err)
					return ""
				}
				// now get to ','
				for ; x < len(data) && data[x] != ','; x++ {
				}
				if x == len(data) {
					break
				}
			}
		}
	}

	// all default values are removed
	// now remove the empty blocks
	if startingLen > len(data) {
		// clear the empty blocks
		return removeDefaultValues(data)
	}
	return data
}

func isWhiteSpace(char rune) bool {
	switch char {
	case ' ':
		return true
	case '\n':
		return true
	case '\t':
		return true
	default:
		return false

	}
}

func getString(json *string, start int, findKey bool) (keyEnd int, err error) {
	if (*json)[start] != '"' {
		return -1, fmt.Errorf("getKey should start with a quote")
	}
	var i int
	for i = start + 1; (*json)[i] != '"'; i++ {
		if (*json)[i] == '\\' && (*json)[i+1] == '"' {
			i++
		}
	}
	// now i is the index of the closeing '"'
	if !findKey {
		return i, nil
	}
	for i = i + 1; (*json)[i] != ':'; i++ {
		if !isWhiteSpace(rune((*json)[i])) {
			return -1, fmt.Errorf("Only spaces are expected here")
		}
	}
	return i, nil
}

func isDefaultValue(json *string, start int) (isDefaultVal bool, defaultValueLasChartIndex int, err error) {
	if (*json)[start] != ':' {
		return false, -1, fmt.Errorf("The falue string should start with ':'")
	}
	var i int
	end := len(*json)
	for i = start + 1; i < end && isWhiteSpace(rune((*json)[i])); i++ {
	}

	// check for 0
	if (*json)[i:i+2] == "0 " || (*json)[i:i+2] == "0\n" || (*json)[i:i+2] == "0," {
		return true, i + 1, nil
	}
	if (*json)[i:i+2] == "0}" {
		return true, i, nil
	}

	// check for []
	if (*json)[i:i+2] == "[]" {
		if (*json)[i+2] == ',' {
			return true, i + 2, nil
		}
		return true, i + 1, nil
	}

	// check for {}
	if (*json)[i] == '{' {
		for i = i + 1; i < end && isWhiteSpace(rune((*json)[i])); i++ {
		}
		if (*json)[i] == '}' {
			retIdx := i
			for i = i + 1; i < end && isWhiteSpace(rune((*json)[i])); i++ {
			}
			if (*json)[i] == ',' {
				return true, i, nil
			}
			return true, retIdx, nil
		}
	}

	return false, -1, nil
}

func removeSpaces(data string) string {
	var buff bytes.Buffer
	for i := 0; i < len(data); i++ {
		if data[i] == '"' {
			// false: we want it to stop at the closing '"'
			ei, err := getString(&data, i, false)
			if err != nil {
				log.Fatal(err)
				return ""
			}
			buff.Write([]byte(data[i : ei+1]))
			i = ei
			continue
		}
		if isWhiteSpace(rune(data[i])) {
			continue
		}
		buff.WriteByte(data[i])
	}
	return buff.String()
}
