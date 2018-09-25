package Algos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
)

const ContentTypeHeader = "Content-Type"
const JSONMediaType = "application/json"

type assertConfig struct {
	// true if we want to compare values.
	// false if we do not care about values.
	checkValue bool
	// true if we are asserting that actual contains expected.
	// false if we are asserting that actual should not contain anything in expected.
	onFoundReturns bool
}

// CheckResponse checks if a server response is equal to the expected response defined by the tester, returning true if it is and false if it is not.
// It takes in a pre-defined request, response, handler and a use case identifier.
// The handler refers to an existing http.HandlerFunc in the production code, and this allows CheckResponse to know which handler to test.
// The use case identifier is supplied by the tester to indicate what kind of checking he wants to carry out. The different identifiers and correspoding use cases are detailed below.
/*
KEY SPECIFIC:
0: I want response to CONTAIN AT LEAST this data. (check keys & values)
1: I want response to CONTAIN EXACTLY this data, nothing more. (check keys & values)
2: I want response to CONTAIN THESE keys, but I don't SPECIFY the values. (check keys)
RESPONSE SPECIFIC(?):
3: I want response to NOT CONTAIN any of these keys. (check keys)
4: I want response to NOT CONTAIN these values for these keys (check keys & values)
5: I want response to CONTAIN these values for these keys (check keys & values)
*/
func CheckResponse(req *http.Request, expected *http.Response, handler *http.HandlerFunc, identifier int) bool {
	//send req, receive response, pre-validate everything is ok.
	mediaType, expectedBod, actualBod, err := sendAndReceive(req, expected, handler)
	if err != nil {
		Fatal(err.Error())
	}

	switch mediaType {
	case JSONMediaType:
		//proceed on
	default:
		return bytes.Equal(expectedBod, actualBod)
	}

	//validate identifier
	switch identifier {
	case 0:
		return containsAtLeast(expectedBod, actualBod)
	case 1:
		return containsExactly(expectedBod, actualBod)
	case 2:
		return containsKeys(expectedBod, actualBod)
	case 3:
		return notContainKeys(expectedBod, actualBod)
	case 4:
		return notContainValues(expectedBod, actualBod)
	default: //invalid identifier
		return false
	}
}

// sendAndReceive sends the request, receives the response and extracts the response bodies(if any).
func sendAndReceive(req *http.Request, expected *http.Response, handler *http.HandlerFunc) (string, []byte, []byte, error) {
	//create new response recorder
	resRecorder := httptest.NewRecorder()
	handler.ServeHTTP(resRecorder, req)

	//1) Check if status codes are same.
	if expected.StatusCode != resRecorder.Result().StatusCode {
		return "", nil, nil, fmt.Errorf("Status codes of expected and actual were different. Expected %v but got %v", expected.StatusCode, resRecorder.Result().StatusCode)
	}

	//2) Check if content types are same.
	expectedType := expected.Header.Get("Content-Type")
	actualType := resRecorder.Result().Header.Get("Content-Type")
	if expectedType != "" && expectedType != actualType {
		return "", nil, nil, fmt.Errorf("Content types of expected and actual were different. Expected %v but got %v", expectedType, actualType)
	}

	//3) Parse the media type included in "Content-Type" header, without the additional parameters
	mediaType, _, err := mime.ParseMediaType(expectedType)
	if err != nil {
		Fatal(fmt.Sprintf("Error parsing media type from actual response in %v.", CurrentFunctionName()))
	}

	//extract the response bodies of expected and actual.
	expectedBod, err := ioutil.ReadAll(expected.Body)
	if err != nil {
		Fatal(err.Error())
	}
	actualBod, err := ioutil.ReadAll(resRecorder.Result().Body)
	if err != nil {
		Fatal(err.Error())
	}
	return mediaType, expectedBod, actualBod, nil
}

func containsAtLeast(expectedBod, actualBod []byte) bool {
	ac := assertConfig{
		checkValue:     true,
		onFoundReturns: true,
	}
	return AssertJSON(expectedBod, actualBod, ac)
}

func containsExactly(expectedBod, actualBod []byte) bool {
	ac := assertConfig{
		checkValue:     true,
		onFoundReturns: true,
	}
	return AssertJSON(expectedBod, actualBod, ac) && AssertJSON(actualBod, expectedBod, ac)
}

func containsKeys(expectedBod, actualBod []byte) bool {
	ac := assertConfig{
		checkValue:     false,
		onFoundReturns: true,
	}
	return AssertJSON(expectedBod, actualBod, ac)
}

func notContainKeys(expectedBod, actualBod []byte) bool {
	ac := assertConfig{
		checkValue:     false,
		onFoundReturns: false,
	}
	return AssertJSON(expectedBod, actualBod, ac)
}

func notContainValues(expectedBod, actualBod []byte) bool {
	ac := assertConfig{
		checkValue:     true,
		onFoundReturns: false,
	}
	return AssertJSON(expectedBod, actualBod, ac)
}
