package samurai

import (
	"fmt"
	s "strings"
)

func TokenizeLine(spaceDelimiter, stringDelimiter, openDelimiter, closeDelimiter, pattern, line string) {
	line = "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"
	//patternExampleWithStringDelimiter := "'['(' '(*ip*(), nil, *user*()),']'(*date*(),'\"'(nil, ' '(*method*(), *url*(), *httpver*()), ' '(nil, *httpcode*(), *reqsize*()))"

	exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	spaceDelimiter = ","
	stringDelimiter = ""
	openDelimiter = "("
	closeDelimiter = ")"

	//Find pattern layer width
	patStep0 := s.Split(exampleWithoutDelimiter, openDelimiter)
	fmt.Println(patStep0)

	var data map[string]string
	if len(patStep0[0]) == 1 && patStep0[0] != closeDelimiter && patStep0[0] != openDelimiter && patStep0[0] != spaceDelimiter && patStep0[0] != stringDelimiter {
		data = TokenizeBlock(patStep0[0], spaceDelimiter, stringDelimiter, openDelimiter, closeDelimiter, line, patStep0[1:])
	}

	fmt.Println(data)

	/*
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
	*/
}

func TokenizeBlock(splitter, spaceDelimiter, stringDelimiter, openDelimiter, closeDelimiter, block string, patternTokens []string) map[string]string {
	fmt.Println(splitter, spaceDelimiter, stringDelimiter, openDelimiter, closeDelimiter, block, patternTokens[0])
	str := s.Split(block, splitter)
	fmt.Println(str)
	allData := make(map[string]string)

	for i := 0; i < len(str); i++ {
		var data map[string]string
		if len(patternTokens[0]) == 1 {
			fmt.Println("Found stringDelimiter ", patternTokens[0])
			data = TokenizeBlock(patternTokens[0], spaceDelimiter, stringDelimiter, openDelimiter, closeDelimiter, str[i], patternTokens[1:])
		} else {
			fmt.Println("Found pattern ", patternTokens[0])
			if s.Contains(patternTokens[0], closeDelimiter) {
				fmt.Println("Found closeDelimiter ", patternTokens[0])
				data := make(map[string]string)
				s0 := s.Split(patternTokens[0], closeDelimiter)
				s1 := s.Split(s0[0], spaceDelimiter)

				for i := 0; i < len(s1); i++ {
					if s1[i] != "nil" {
						fmt.Println(s1[i], str[i])
						data[s1[i]] = str[i]
					}
				}

				moreData := make(map[string]string)
				if string(patternTokens[0][len(patternTokens[0])-1]) != closeDelimiter {
					fmt.Println("Found new stringDelimiter ", patternTokens[0])
					moreData = TokenizeBlock(string(patternTokens[0][len(patternTokens[0])-1]), spaceDelimiter, stringDelimiter, openDelimiter, closeDelimiter, str[i], patternTokens[1:])
				}

				data = appendMaps(data, moreData)
			}
		}

		allData = appendMaps(allData, data)
	}
	return allData
}

func appendMaps(map1, map2 map[string]string) map[string]string {
	for k, v := range map2 {
		map1[k] = v
	}

	return map1
}
