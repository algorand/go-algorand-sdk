package logic

import "strings"

// SourceMap provides a mapping of the source to assembled program
type SourceMap struct {
	Version    int      `json:"version"`
	File       string   `json:"file,omitempty"`
	SourceRoot string   `json:"sourceRoot,omitempty"`
	Sources    []string `json:"sources"`
	Names      []string `json:"names"`
	Mapping    string   `json:"mapping"`
	// Decoded mapping results
	LineToPc map[int][]int
	PcToLine map[int]int
}

func NewSourceMap(ism map[string]interface{}) SourceMap {
	sm := SourceMap{}

	if v, ok := ism["version"]; ok {
		sm.Version = int(v.(float64))
	}

	if f, ok := ism["file"]; ok {
		sm.File = f.(string)
	}

	if sr, ok := ism["sourceRoot"]; ok {
		sm.SourceRoot = sr.(string)
	}

	if srcs, ok := ism["sources"]; ok {
		srcSlice := srcs.([]interface{})
		sm.Sources = make([]string, len(srcSlice))
		for idx, s := range srcSlice {
			sm.Sources[idx] = s.(string)
		}
	}

	if names, ok := ism["names"]; ok {
		nameSlice := names.([]interface{})
		sm.Names = make([]string, len(nameSlice))
		for idx, n := range nameSlice {
			sm.Names[idx] = n.(string)
		}
	}

	if m, ok := ism["mapping"]; ok {
		sm.Mapping = m.(string)
	}

	sm.PcToLine = map[int]int{0: 0}
	sm.LineToPc = map[int][]int{0: {0}}

	lastLine := 0
	for idx, chunk := range strings.Split(sm.Mapping, ";") {
		vals := DecodeSourceMapLine(chunk)
		if len(vals) > 3 {
			lineNum := vals[2]
			if _, ok := sm.LineToPc[lineNum]; !ok {
				sm.LineToPc[lineNum] = []int{}
			}
			sm.LineToPc[lineNum] = append(sm.LineToPc[lineNum], idx)

			lastLine = lineNum
		}
		sm.PcToLine[idx] = lastLine
	}

	return sm
}

func (s *SourceMap) GetLineForPc(pc int) int {
	return s.PcToLine[pc]
}

func (s *SourceMap) GetPCsForLine(line int) []int {
	return s.LineToPc[line]
}

const (
	// consts used for vlq encoding/decoding
	b64table     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	vlqShiftSize        = 5
	vlqFlag             = 1 << vlqShiftSize
	vlqMask             = vlqFlag - 1
)

func DecodeSourceMapLine(vlq string) []int {

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
