package samurai

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTokenizeLine(t *testing.T) {
	InitSamurai(0, 0, 0)
	exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	fmt.Println(validatePattern(exampleWithoutDelimiter))
	data := TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)
	fmt.Println(data)
}

func TestTokenizeLine2(t *testing.T) {
	InitSamurai(0, 0, 0)
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver), code, size)"
	fmt.Println(validatePattern(exampleWithoutDelimiter))
	data := TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)
	fmt.Println(data)
}

func TestValidatePattern(t *testing.T) {
	InitSamurai(0, 0, 0)
	exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	err := validatePattern(exampleWithoutDelimiter)
	t.Log(err)
}

func TestStripLayer(t *testing.T) {
	InitSamurai(0, 0, 0)
	exampleWithoutDelimiter := "[(*(ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	elements := stripLayer(exampleWithoutDelimiter)
	for i := 0; i < len(elements); i++ {
		t.Log(elements[i])
	}

	newElements := stripLayer(elements[1])
	for i := 0; i < len(newElements); i++ {
		t.Log(newElements[i])
	}

	newElements2 := stripLayer(newElements[1])
	for i := 0; i < len(newElements2); i++ {
		t.Log(newElements2[i])
	}

	newelements3 := stripLayer(newElements2[1])
	t.Log(newelements3)
}

func BenchmarkTokenizeLine(b *testing.B) {
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver), code, size)"
	for i := 0; i < b.N; i++ {

	}
	TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)
}

func TestTokenizeLineNasaLog(t *testing.T) {
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver), code, size)"
	//exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"

	f, err := os.Open("./nasa_log.log")
	if err != nil {
		t.Log(err)
	}
	scanner := bufio.NewScanner(f)
	timeNow := time.Now()
	c := 0
	for scanner.Scan() {
		c += 1

		if c%100000 == 0 {
			TokenizeBlock(scanner.Text(), exampleWithoutDelimiter)
		} else {
			TokenizeBlock(scanner.Text(), exampleWithoutDelimiter)
		}
	}
	timeAfter := time.Now()

	fmt.Printf("Parsed %v lines in %v\n", c, timeAfter.Sub(timeNow))
	fmt.Printf("%v ns pr op\n", timeAfter.Sub(timeNow).Nanoseconds()/int64(c))
}

/*
func TestTokenizeLineNasaLogAsync(t *testing.T) {
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver), code, size)"
	//exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"

	f, err := os.Open("./nasa_log.log")
	if err != nil {
		t.Log(err)
	}
	scanner := bufio.NewScanner(f)
	timeNow := time.Now()
	c := 0
	for scanner.Scan() {
		c += 1

		if c%100000 == 0 {
			data := TokenizeBlockAsync(scanner.Text(), exampleWithoutDelimiter)
			fmt.Println(scanner.Text(), data)
		} else {
			TokenizeBlockAsync(scanner.Text(), exampleWithoutDelimiter)
		}
	}
	timeAfter := time.Now()

	fmt.Printf("Parsed %v lines in %v\n", c, timeAfter.Sub(timeNow))
	fmt.Printf("%v ns pr op\n", timeAfter.Sub(timeNow).Nanoseconds()/int64(c))
}
*/
