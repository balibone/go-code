package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"
)

var expected interface{}
var actual interface{}

type JSONPath string
type customCompare func(interface{}, interface{}) bool
type JSONConfig struct {
	checkPaths map[JSONPath]customCompare
	exactMatch bool
}

var (
	ignoreVal = func(expected, actual interface{}) bool {
		return true
	}
	equalType = func(expected, actual interface{}) bool {
		return reflect.TypeOf(actual).Kind() == reflect.TypeOf(expected).Kind()
	}
	//check that timestamp field is in RFC3339 format.
	isTimestamp = func(expected, actual interface{}) bool {
		switch actual.(type) {
		case string:
			_, err := time.Parse(time.RFC3339, actual.(string))
			return err == nil
		default:
			return false
		}
	}
	//validate field with regex supplied to matchRegex func.
	matchRegex = func(regex string) customCompare {
		return func(expected, actual interface{}) bool {
			switch actual.(type) {
			case string:
				match, _ := regexp.MatchString(regex, actual.(string))
				return match
			default:
				return false
			}
		}
	}
)

// AssertJSON compares 2 JSONs, in a way defined by the assertConfig struct passed to it.
func AssertJSON(expectedJSON, actualJSON []byte, config JSONConfig) (ok bool) {
	err := json.Unmarshal(expectedJSON, &expected)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(actualJSON, &actual)
	if err != nil {
		fmt.Println(err)
		return
	}
	return isEqualJSON("", expected, actual, config)
}

func isEqualElement(path JSONPath, expected, actual interface{}, config JSONConfig) bool {
	if customCompare, ok := config.checkPaths[path]; ok {
		return customCompare(expected, actual)
	}
	if jsonType, ok := isJSONBasicType(expected); ok {
		switch jsonType {
		case reflect.UnsafePointer:
			//compare types
			if reflect.TypeOf(actual) != nil {
				log.Fatal(fmt.Sprintf("Expected type `null` in %s, got `%v`\n", path, actual))
				return false
			}
			return true
		case reflect.Float64:
			//compare types
			if reflect.ValueOf(actual).Kind() != jsonType {
				log.Fatal(fmt.Sprintf("Expected type `float` in %s, got %v\n", path, actual))
				return false
			}
			if expected.(float64) != actual.(float64) {
				log.Fatal(fmt.Sprintf("Expected `%v` in %s, got `%v`\n", expected, path, actual))
				return false
			}
			return true
		case reflect.String:
			//compare types
			if reflect.ValueOf(actual).Kind() != jsonType {
				log.Fatal(fmt.Sprintf("Expected type `string` in %s, got %v\n", path, actual))
				return false
			}
			if expected.(string) != actual.(string) {
				log.Fatal(fmt.Sprintf("Expected `%v` in %s, got `%v`\n", expected, path, actual))
				return false
			}
			return true
		case reflect.Bool:
			//compare types
			if reflect.ValueOf(actual).Kind() != jsonType {
				log.Fatal(fmt.Sprintf("Expected type `bool` in %s, got %v\n", path, actual))
				return false
			}
			if expected.(bool) != actual.(bool) {
				log.Fatal(fmt.Sprintf("Expected `%v` in %s, got `%v`\n", expected, path, actual))
				return false
			}
			return true
		}
	}
	//not in comparePaths map and not a json basic type, so its a json. call function on it again.
	return isEqualJSON(path, expected, actual, config)
}

func isEqualJSON(path JSONPath, expected, actual interface{}, config JSONConfig) (ok bool) {
	//if expected is array,
	if expectedArray, ok := isArray(expected); ok {
		actualArray, ok := isArray(actual)
		//but actual is not array, return.
		if !ok {
			return false
		}
		//if expected and actual are both arrays but different lengths, return.
		if len(expectedArray) != len(actualArray) {
			log.Fatal(fmt.Sprintf("Expected array length of %d in %s, got %d\n",
				len(expectedArray), path, len(actualArray)))
			return false
		}
		//else start comparing each array index.
		for idx, element := range expectedArray {
			newPath := fmt.Sprintf("%s[%d]", path, idx)
			//if index elements are not equal, return
			if !isEqualElement(JSONPath(newPath), element, actualArray[idx], config) {
				return false
			}
		}
		//found
		return true
	}
	//if expected is object,
	if expectedObj, ok := isObj(expected); ok {
		actualObj, ok := isObj(actual)
		//but actual is not object, return
		if !ok {
			return false
		}
		//start comparing each object element.
		for key, element := range expectedObj {
			actualElement, ok := actualObj[key]
			//if key does not exist in actual object, return.
			if !ok {
				log.Fatal(fmt.Sprintf("Expected key `%s` in %s\n", key, path))
				return false
			}

			newPath := fmt.Sprintf(`%s["%s"]`, path, key)
			//if key values are not equal, return.
			if !isEqualElement(JSONPath(newPath), element, actualElement, config) {
				return false
			}
		}
		return true
	}

	fmt.Printf("Expected value at %s is not json: `%v`\n", path, expected)
	return false
}

func isArray(input interface{}) (array []interface{}, ok bool) {
	//if input type is not a slice (array), return zero values (nil, false)
	if reflect.ValueOf(input).Kind() != reflect.Slice {
		return
	}
	return input.([]interface{}), true
}

func isObj(input interface{}) (obj map[string]interface{}, ok bool) {
	//if input type is not a map (object), return zero values (nil, false)
	if reflect.ValueOf(input).Kind() != reflect.Map {
		return
	}
	return input.(map[string]interface{}), true
}

func isJSONBasicType(input interface{}) (JSONType reflect.Kind, ok bool) {
	if reflect.TypeOf(input) == nil {
		// JSON null
		return reflect.UnsafePointer, true
	}
	switch input.(type) {
	case float64:
		// JSON number
		return reflect.Float64, true
	case string:
		// JSON string
		return reflect.String, true
	case bool:
		// JSON bool
		return reflect.Bool, true
	}
	// not json basic type, return (0, false)
	return
}

func main() {
	actualStr := `{
    "status": "success",
    "message": "",
    "data": {
        "user": {
            "id": 1,
            "email": "admin@ratex.co",
            "name": "admin",
            "phone_no": 96787432,
            "account_type": 2,
            "share_referral": "ABCDEF",
            "currency": "SGD",
            "last_logged_in": "2018-05-25T10:50:37+08:00",
            "last_created": "2018-05-15T14:44:38.771142+08:00"
        },
        "balance": {
            "currency": "SGD",
            "amount": "0.00",
            "last_modified": "2018-05-15T15:38:27.781238+08:00"
        },
        "rewards": null,
        "new": false,
        "verified": true,
        "email_verified": true,
        "phone_verified": true,
        "has_password": true,
        "currency": "SGD",
        "csrf": "xmzu8vEjIkTawI0BvOkzyoxZEua01FfBVDlN4f9SsG8="
    }
	}`
	expectedStr := `{
	"status": "success",
    "message": "",
    "data": {
        "user": {
            "id": 1,
            "email": "admin@ratex.co",
            "name": "admin",
            "phone_no": 96787432,
            "account_type": 2,
            "currency": "SGD",
            "last_logged_in": "2018-05-25T10:50:37+08:00",
            "last_created": "2018-05-15T14:44:38.771142+08:00"
        },
        "balance": {
            "currency": "SGD",
            "amount": "0.00",
            "last_modified": "2018-05-15T15:38:27.781238+08:00"
        },
        "rewards": null,
        "new": false,
        "phone_verified": true,
        "has_password": true,
        "currency": "SGD"
    }	
	}`

	conf := JSONConfig{
		checkPaths: map[JSONPath]customCompare{
			`["data"]["user"]["share_referral"]`: matchRegex(`[[:upper:]]{5}`),
		},
	}

	fmt.Println("assertJSON =", AssertJSON([]byte(expectedStr), []byte(actualStr), conf))
}
