package samurai

import (
	"errors"
	"fmt"
	s "strings"
)

var (
	spaceDelim = uint8(',')
	openDelim  = uint8('(')
	closeDelim = uint8(')')
)

func SetDelimiters(spaceDelimiter, openDelimiter, closeDelimiter uint8) {
	spaceDelim = spaceDelimiter
	openDelim = openDelimiter
	closeDelim = closeDelimiter
}

func TokenizeBlock(block string, patternBlock string) map[string]string {
	splitter := ""
	for i := 0; i < len(patternBlock); i++ {
		if patternBlock[i] == openDelim {
			splitter = patternBlock[0:i]
			break
		}
	}

	elements := stripLayer(patternBlock)
	str := s.Split(block, splitter)
	allData := make(map[string]string)

	for i := 0; i < len(elements); i++ {
		newElements := stripLayer(elements[i])
		if len(newElements) == 0 {
			if elements[i] == "nil" {
				continue
			}
			if i < len(str) {
				allData[elements[i]] = str[i]
			}
		} else {
			if i < len(str) && i < len(elements) {
				data := TokenizeBlock(str[i], elements[i])
				for k, v := range data {
					allData[k] = v
				}
			}
		}
	}
	return allData
}

func appendMaps(map1, map2 map[string]string) map[string]string {
	for k, v := range map2 {
		map1[k] = v
	}

	return map1
}

//Strips the outermost pattern group and returns the elements inside
func stripLayer(pattern string) (elements []string) {
	c := 0
	counting := false
	es := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == openDelim {
			c += 1

			if counting == false {
				counting = true
				es = i + 1
			}
		} else if (pattern[i] == spaceDelim || pattern[i] == closeDelim) && c == 1 {
			elements = append(elements, pattern[es:i])
			es = i + 1
		} else if pattern[i] == closeDelim {
			c -= 1
		}
	}

	return elements
}

//Should check if pattern has same amount of open and close delimiters, and check if all split values are present
//Could also check if number of keyswords > 0
func validatePattern(pattern string) error {
	delimCount := 0
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == openDelim {
			delimCount += 1
			//If open delim is present, there should also be a splitter just before
			if i-1 < 0 {
				return errors.New(fmt.Sprintf("Missing splitter string next to open delimiter (char %v)", 0))
			}
			if pattern[i-1] == openDelim || pattern[i-1] == closeDelim || pattern[i-1] == spaceDelim {
				return errors.New(fmt.Sprintf("Missing splitter string next to open delimiter (char %v)", i))
			}
		} else if pattern[i] == closeDelim {
			delimCount -= 1
		}
	}

	if delimCount != 0 {
		return errors.New("Unbalanced pattern")
	}

	return nil
}
