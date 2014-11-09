package samurai

import (
	"errors"
	"fmt"
	s "strings"
	"sync"
)

var (
	spaceDelim, stringDelim, openDelim, closeDelim uint8
)

func InitSamurai(spaceDelimiter, openDelimiter, closeDelimiter uint8) {
	spaceDelimiter = uint8(',')
	openDelimiter = uint8('(')
	closeDelimiter = uint8(')')

	spaceDelim = spaceDelimiter
	openDelim = openDelimiter
	closeDelim = closeDelimiter
}

/*
func TokenizeLine(pattern, line string) {
	line = "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"
	//patternExampleWithStringDelimiter := "'['(' '(*ip*(), nil, *user*()),']'(*date*(),'\"'(nil, ' '(*method*(), *url*(), *httpver*()), ' '(nil, *httpcode*(), *reqsize*()))"

	exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"

	//Find pattern layer width
	patStep0 := s.Split(exampleWithoutDelimiter, string(openDelim))
	fmt.Println(patStep0)

	var data map[string]string
	if len(patStep0[0]) == 1 && patStep0[0] != string(closeDelim) && patStep0[0] != string(openDelim) && patStep0[0] != string(spaceDelim) && patStep0[0] != string(stringDelim) {
		data = TokenizeBlock(line, patStep0[1:])
	}

	fmt.Println(data)

		step0 := s.Split(line, "[")
		step1 := s.Split(step0[0], " ")
		ip := step1[0]
		user := step1[2]

		step2 := s.Split(step0[1], "]")
		date := step2[0]

		step3 := s.Split(step2[1], "\"")

		step4 := s.Split(step3[1], " ")
		method := step4[0]
		url := step4[1]
		httpver := step4[2]

		step5 := s.Split(step3[2], " ")

		httpcode := step5[1]
		reqsize := step5[2]

		fmt.Println(ip, user, date, method, url, httpver, httpcode, reqsize)
		return nil, nil

}
*/

func TokenizeBlockAsync(block string, patternBlock string) map[string]string {
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

	wg := sync.WaitGroup{}

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
				wg.Add(1)
				var data map[string]string
				go func() {
					data = TokenizeBlockAsync(block, patternBlock)
					for k, v := range data {
						allData[k] = v
					}
					wg.Done()
				}()

			}
		}
	}
	wg.Wait()

	return allData
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
