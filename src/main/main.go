package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	//	"log"
	"net/http"
	"net/url"
	//	"strconv"
	"mime"
	"mime/multipart"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}

const path = "/Users/sreddy/space/HEAD/qa/scripts/tests/"
const file = "dbrep-ja.xml"

//possibe values for test status
const (
	PASS    = "pass"
	Fail    = "fail"
	UNKNOWN = "unknown"
)

// QueryLang represents language options for query
type QueryLang struct {
	XQY string
	JS  string
}

// Status represents test status
type Status struct {
	PASS    string
	FAIL    string
	UNKNOWN string
}

// QUERY_LANG available query languages
var QUERY_LANG = &QueryLang{"xqy", "js"}

// Test represents data structure for <h:test> element
type Test struct {
	Name      string `xml:"name"`
	Expected  string `xml:"comment-expected-result"`
	Query     string `xml:"query"`
	QueryLang string `xml:"lang,attr"`
	ShellCmd  string `xml:"shell-cmd"`
	Type      string `xml:"type,attr"`
	Output    string `xml:"output-result,attr"`
	Status    string
	Auth

	//	Value    string `xml:",chardata"`
}

// Auth represents user credentials to run the test
type Auth struct {
	Username string `xml:"username,attr"`
	Password string `xml:"password,attr"`
}

//Script represents the xml test file composed of Test structs
type Script struct {
	Filename string
	Tests    []Test `xml:"test"`
}

// getSetup get all setup tests
func (s *Script) getSetup() []Test {
	var setups []Test
	for _, test := range s.Tests {
		if test.Type == "setup" {
			setups = append(setups, test)
		}
	}
	return setups
}

// getTeardown get all teardown tests
func (s *Script) getTeardown() []Test {
	var teardowns []Test
	for _, test := range s.Tests {
		if test.Type == "teardown" {
			teardowns = append(teardowns, test)
		}
	}
	return teardowns
}

// check err
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func executeTest(t Test) {
	api := "http://localhost:8000/LATEST/eval"
	query := "xdmp.version()"
	data := url.Values{}
	data.Set("javascript", query)

	body := strings.NewReader(data.Encode())
	//fmt.Println(body)
	req, err := http.NewRequest("POST", api, body)
	check(err)
	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "multipart/mixed; boundary=BOUNDARY")
	//req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	check(err)
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	check(err)

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			check(err)
			result, err := ioutil.ReadAll(p)
			check(err)
			fmt.Printf("Results  %q\n", result)
		}
	}
}

func sendResult(t Test, ch chan Test) {
	//fmt.Println(res)
	t.Status = PASS
	executeTest(t)
	ch <- t
	wg.Done()
}

func main() {
	filepath := strings.Join([]string{path, file}, "")
	data, err := ioutil.ReadFile(filepath)
	check(err)
	//fmt.Println(string(data))

	// unmarshal xml file
	xmldata := new(Script)
	//xmldata.Filename = filepath
	xml.Unmarshal(data, xmldata)
	fmt.Println("xml test file : ", xmldata.Filename)
	for _, test := range xmldata.Tests {
		if test.QueryLang == "" {
			test.QueryLang = QUERY_LANG.XQY
		}
		//fmt.Println(test.QueryLang)
	}

	// channel
	resultChannel := make(chan Test, len(xmldata.getTeardown()))

	for _, teardown := range xmldata.getTeardown() {
		wg.Add(1)
		go sendResult(teardown, resultChannel)
	}

	wg.Wait()
	close(resultChannel)

	// for elem := range resultChannel {
	// 	fmt.Println(elem.Status)
	// }

}
