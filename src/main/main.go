package main

import (
	"fmt"
	"strings"
	"xmltest"
)

const path = "/Users/sreddy/space/HEAD/qa/scripts/tests/"
const file = "dbrep-ja.xml"

func main() {
	filepath := strings.Join([]string{path, file}, "")
	xmlTests := new(xmltest.XMLFile)
	tests := xmlTests.ReadFile(filepath)
	// for _, test := range tests {
	// 	fmt.Println(test.Name)
	// }
	results := xmlTests.ExecuteXMLFile(tests)
	for res := range results {
		fmt.Println(res.Result.QueryOutput)
	}

}
