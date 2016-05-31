package mlerror

import (
	"encoding/json"
	//"fmt"
	//	"io"
	//"io/ioutil"
)

// Check check err
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// EvalResponse represents error response object
type EvalResponse struct {
	StatusCode  float64
	Status      string
	MessageCode string
	Message     string
}

// GetErrorResponse parses error response from eval endpoint
func GetErrorResponse(body []byte) (EvalResponse, error) {
	var evalResp EvalResponse
	var jsonError map[string]interface{}

	err := json.Unmarshal(body, &jsonError)
	if err != nil {
		return evalResp, err
	}

	//fmt.Println(jsonError["errorResponse"])

	status := jsonError["errorResponse"].(map[string]interface{})["status"]
	statusCode := jsonError["errorResponse"].(map[string]interface{})["statusCode"]
	message := jsonError["errorResponse"].(map[string]interface{})["message"]
	messageCode := jsonError["errorResponse"].(map[string]interface{})["messageCode"]
	//messageDetail := jsonError["errorResponse"].(map[string]interface{})["messageDetail"]

	//fmt.Println(status, statusCode, message, messageCode)

	evalResp = EvalResponse{
		StatusCode:  statusCode.(float64),
		Status:      status.(string),
		MessageCode: messageCode.(string),
		Message:     message.(string),
	}

	//	fmt.Printf("%v \n", evalResp)

	return evalResp, nil

}
