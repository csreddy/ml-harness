package xmltest

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
	"querylang"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}

const path = "/Users/sreddy/space/HEAD/qa/scripts/tests/"
const file = "dbrep-ja.xml"

// check err
func check(e error) {
	if e != nil {
		panic(e)
	}
}

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
	Result
	Auth

	//	Value    string `xml:",chardata"`
}

// Auth represents user credentials to run the test
type Auth struct {
	Username string `xml:"username,attr"`
	Password string `xml:"password,attr"`
}

//Script represents the xml test file composed of Test structs
type XMLFile struct {
	Filename   string
	Tests      []Test `xml:"test"`
	Status     string
	ResultFile string
}

// Result contains the output from test query and shell cmds (if present)
type Result struct {
	QueryOutput string
	ShellOutput string
}

// getSetup get all setup tests
func (s *XMLFile) GetSetup() []Test {
	var setups []Test
	for _, test := range s.Tests {
		if test.Type == "setup" {
			setups = append(setups, test)
		}
	}
	return setups
}

// getTeardown get all teardown tests
func (s *XMLFile) GetTeardown() []Test {
	var teardowns []Test
	for _, test := range s.Tests {
		if test.Type == "teardown" {
			teardowns = append(teardowns, test)
		}
	}
	return teardowns
}

// //
// func (s *Script) executeSetup() (bool, error) {
// 	setups := s.getSetup()
// 	setupChannel := make(chan Test, len(setups))
// 	for _, setup := range setups {
// 		setupChannel <- executeTest(setup)
// 	}
// 	close(setupChannel)
// }

// ReadFile reads the xml file into memory
func (s *XMLFile) ReadFile(filepath string) []Test {
	//filepath := strings.Join([]string{path, file}, "")
	data, err := ioutil.ReadFile(filepath)
	check(err)
	//fmt.Println(string(data))

	// unmarshal xml file
	xmldata := new(XMLFile)
	//xmldata.Filename = filepath
	xml.Unmarshal(data, xmldata)
	fmt.Println("xml test file : ", xmldata.Filename)
	for _, test := range xmldata.Tests {
		if test.QueryLang == "" {
			test.QueryLang = querylang.JS
		}
		//fmt.Println(test.QueryLang)
	}
	return xmldata.Tests
}

// ExecuteTest executes a given test
func (s *Test) Execute(t Test) Test {
	//api := "http://localhost:8000/LATEST/eval"
	// container
	api := "http://192.168.99.100:32809/LATEST/eval"
	query := "xdmp.version()"
	data := url.Values{}
	// set query lang
	t.QueryLang = querylang.JS

	if t.QueryLang == "js" {
		data.Set("javascript", query)
	} else {
		data.Set("xquery", query) // default query lang
	}

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
				return t
			}

			check(err)
			result, err := ioutil.ReadAll(p)
			check(err)
			t.Result = Result{QueryOutput: string(result), ShellOutput: ""}
			//fmt.Printf("Result  %q\n", t.Result)
		}
	}

	return t
}

// ExecuteXMLFile executes all tests inside the xml
func (s *XMLFile) ExecuteXMLFile(tests []Test) chan Test {
	wg := sync.WaitGroup{}
	resultChannel := make(chan Test, len(tests))

	for _, test := range tests {
		wg.Add(1)
		resultChannel <- test.Execute(test)
		wg.Done()
	}
	wg.Wait()
	close(resultChannel)
	return resultChannel
}
