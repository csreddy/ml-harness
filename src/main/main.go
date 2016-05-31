package main

import (
	"fmt"
	"strings"
	"xmltest"
)

const path = "/Users/sreddy/space/HEAD/qa/scripts/tests/"
const file = "bug3252.xml"

func main() {
	filepath := strings.Join([]string{path, file}, "")
	xmlTests := new(xmltest.XMLFile)
	// get tests
	tests := xmlTests.ReadFile(filepath)

	// execute tests
	results := xmlTests.ExecuteXMLFile(tests)
	for res := range results {
		fmt.Println(res.Result)
	}

}
