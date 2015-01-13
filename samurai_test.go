package samurai

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestTokenizeLineAsyncApacheLog(t *testing.T) {
	exampleWithoutDelimiter := "[( (ip,ident,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	data := TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)

	if data["ip"] != "127.0.0.1" {
		t.Error("Failed to extract IP correctly")
	}

	if data["ident"] != "-" {
		t.Error("Failed to extract ident correctly")
	}

	if data["date"] != "10/Oct/2000:13:55:36 -0700" {
		t.Error("Failed to extract date correctly")
	}

	if data["method"] != "GET" {
		t.Error("Failed to extract http method correctly")
	}

	if data["url"] != "/apache_pb.gif" {
		t.Error("Failed to extract url correctly")
	}

	if data["httpver"] != "HTTP/1.0" {
		t.Error("Failed to extract http version correctly")
	}

	if data["httpcode"] != "200" {
		t.Error("Failed to extract http code correctly")
	}

	if data["reqsize"] != "2326" {
		t.Error("Failed to extract request size correctly")
	}
}

func TestTokenizeLineApacheLog(t *testing.T) {
	exampleWithoutDelimiter := "[( (ip,ident,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	data := TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)

	if data["ip"] != "127.0.0.1" {
		t.Error("Failed to extract IP correctly")
	}

	if data["ident"] != "-" {
		t.Error("Failed to extract ident correctly")
	}

	if data["date"] != "10/Oct/2000:13:55:36 -0700" {
		t.Error("Failed to extract date correctly")
	}

	if data["method"] != "GET" {
		t.Error("Failed to extract http method correctly")
	}

	if data["url"] != "/apache_pb.gif" {
		t.Error("Failed to extract url correctly")
	}

	if data["httpver"] != "HTTP/1.0" {
		t.Error("Failed to extract http version correctly")
	}

	if data["httpcode"] != "200" {
		t.Error("Failed to extract http code correctly")
	}

	if data["reqsize"] != "2326" {
		t.Error("Failed to extract request size correctly")
	}
}

func TestTokenizeLineApacheLogAltPattern(t *testing.T) {
	exampleWithoutDelimiter := " (ip,ident,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver),httpcode,reqsize)"
	data := TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)

	if data["ip"] != "127.0.0.1" {
		t.Error("Failed to extract IP correctly")
	}

	if data["ident"] != "-" {
		t.Error("Failed to extract ident correctly")
	}

	if data["date"] != "10/Oct/2000:13:55:36" {
		t.Error("Failed to extract date correctly")
	}

	if data["tz"] != "-0700" {
		t.Error("Failed to extract date correctly")
	}

	if data["method"] != "GET" {
		t.Error("Failed to extract http method correctly")
	}

	if data["url"] != "/apache_pb.gif" {
		t.Error("Failed to extract url correctly")
	}

	if data["httpver"] != "HTTP/1.0" {
		t.Error("Failed to extract http version correctly")
	}

	if data["httpcode"] != "200" {
		t.Log(data)
		t.Error("Failed to extract http code correctly")
	}

	if data["reqsize"] != "2326" {
		t.Error("Failed to extract request size correctly")
	}
}

func TestValidatePattern(t *testing.T) {
	exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	err := ValidatePattern(exampleWithoutDelimiter)
	t.Log(err)
}

func TestStripLayer(t *testing.T) {
	exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
	elements := stripLayer(exampleWithoutDelimiter)

	if elements[0] != " (ip,nil,user)" {
		t.Errorf("Failed to extract subpattern. Expected: \"%v\" Got: \"%v\"", " (ip,nil,user)", elements[0])
	}
	if elements[1] != "](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize)))" {
		t.Errorf("Failed to extract subpattern. Expected: \"%v\" Got: \"%v\"", "](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize)))", elements[1])
	}

	newElements := stripLayer(elements[1])
	if newElements[0] != "date" {
		t.Errorf("Failed to extract subpattern. Expected: \"%v\" Got: \"%v\"", "date", newElements[0])
	}
	if newElements[1] != "\"(nil, (method,url,httpver), (nil,httpcode,reqsize))" {
		t.Errorf("Failed to extract subpattern. Expected: \"%v\" Got: \"%v\"", "\"(nil, (method,url,httpver), (nil,httpcode,reqsize))", newElements[1])
	}
}

func BenchmarkTokenizeLine(b *testing.B) {
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver), code, size)"
	for i := 0; i < b.N; i++ {
		TokenizeBlock("127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326", exampleWithoutDelimiter)
	}
}

/* Use this to test on a local file
func TestTokenizeLine(t *testing.T) {
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver),code,size)"
	//exampleWithoutDelimiter := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"

	f, err := os.Open("./apache_access_log")
	if err != nil {
		t.Skip(err)
	}
	scanner := bufio.NewScanner(f)
	timeNow := time.Now()
	limit := 100
	c := 0
	for scanner.Scan() {
		if c == limit {
			break
		}
		c += 1
		TokenizeBlock(scanner.Text(), exampleWithoutDelimiter)
	}

	timeAfter := time.Now()

	fmt.Printf("Parsed %v lines in %v\n", c, timeAfter.Sub(timeNow))
	fmt.Printf("%v ns pr op\n", timeAfter.Sub(timeNow).Nanoseconds()/int64(c))
}
*/
