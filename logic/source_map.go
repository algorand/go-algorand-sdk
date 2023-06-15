package logic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SourceMap provides a mapping of the source to assembled program
type SourceMap struct {
	Version    int      `json:"version"`
	File       string   `json:"file,omitempty"`
	SourceRoot string   `json:"sourceRoot,omitempty"`
	Sources    []string `json:"sources"`
	Names      []string `json:"names"`
	Mappings   string   `json:"mappings"`
	// Decoded mapping results
	LineToPc map[int][]int
	PcToLine map[int]int
}

// DecodeSourceMap decodes a source map
func DecodeSourceMap(ism map[string]interface{}) (SourceMap, error) {
	var sm SourceMap

	buff, err := json.Marshal(ism)
	if err != nil {
		return sm, err
	}

	err = json.Unmarshal(buff, &sm)
	if err != nil {
		return sm, err
	}

	if sm.Version != 3 {
		return sm, fmt.Errorf("only version 3 is supported")
	}

	if sm.Mappings == "" {
		return sm, fmt.Errorf("no mappings defined")
	}

	sm.PcToLine = map[int]int{}
	sm.LineToPc = map[int][]int{}

	lastLine := 0
	for idx, chunk := range strings.Split(sm.Mappings, ";") {
		vals := decodeSourceMapLine(chunk)
		// If the vals length >= 3 the lineDelta
		if len(vals) >= 3 {
			lastLine = lastLine + vals[2] // Add the line delta
		}

		if _, ok := sm.LineToPc[lastLine]; !ok {
			sm.LineToPc[lastLine] = []int{}
		}

		sm.LineToPc[lastLine] = append(sm.LineToPc[lastLine], idx)
		sm.PcToLine[idx] = lastLine
	}

	return sm, nil
}

// GetLineForPc returns the line number for the given pc
func (s *SourceMap) GetLineForPc(pc int) (int, bool) {
	line, ok := s.PcToLine[pc]
	return line, ok
}

// GetPcsForLine returns the program counter for the given line
func (s *SourceMap) GetPcsForLine(line int) []int {
	return s.LineToPc[line]
}

const (
	// consts used for vlq encoding/decoding
	b64table     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	vlqShiftSize        = 5
	vlqFlag             = 1 << vlqShiftSize
	vlqMask             = vlqFlag - 1
)

func decodeSourceMapLine(vlq string) []int {

	var (
		results      []int
		value, shift int
	)

	for i := 0; i < len(vlq); i++ {
		digit := strings.Index(b64table, string(vlq[i]))

		value |= (digit & int(vlqMask)) << shift

		if digit&vlqFlag > 0 {
			shift += vlqShiftSize
			continue
		}

		if value&1 > 0 {
			value = (value >> 1) * -1
		} else {
			value = value >> 1
		}

		results = append(results, value)

		// Reset
		value, shift = 0, 0
	}

	return results
}
