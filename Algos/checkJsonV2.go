package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

var expected interface{}
var actual interface{}

type assertConfig struct {
	// true if we want to compare values.
	// false if we do not care about values.
	checkValues bool
	// true if we are asserting that actual contains expected.
	// false if we are asserting that actual should not contain anything in expected.
	checkExists bool
}

// NOTE: this function has been made to cater for both:
// 1) checking existence
// 2) checking for non-existence
// This is done through the use of checkExists field in assertConfig struct.
// if existence passes, returns true. else returns false.
// if non-existence passes, returns false. else returns true.
func assertJSON(expectedJSON, actualJSON []byte, ac assertConfig) (ok bool) {
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
	return isEqualJSON("actual", expected, actual, ac)
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

func isEqualElement(path string, expected, actual interface{}, ac assertConfig) bool {
	if jsonType, ok := isJSONBasicType(expected); ok {
		switch jsonType {
		case reflect.UnsafePointer:
			//compare types
			if reflect.TypeOf(actual) != nil && ac.checkExists {
				log.Fatal(fmt.Sprintf("Expected type `null` in %s, got `%v`\n", path, actual))
				return false
			}
			return true
		case reflect.Float64:
			//compare types
			if reflect.ValueOf(actual).Kind() != jsonType && ac.checkExists {
				log.Fatal(fmt.Sprintf("Expected type `float` in %s, got %v\n", path, actual))
				return false
			}
			//compare values
			if expected.(float64) != actual.(float64) && ac.checkExists && ac.checkValues {
				log.Fatal(fmt.Sprintf("Expected `%v` in %s, got `%v`\n", expected, path, actual))
				return false
			}
			return true
		case reflect.String:
			//compare types
			if reflect.ValueOf(actual).Kind() != jsonType && ac.checkExists {
				log.Fatal(fmt.Sprintf("Expected type `string` in %s, got %v\n", path, actual))
				return false
			}
			//compare values
			if expected.(string) != actual.(string) && ac.checkExists && ac.checkValues {
				log.Fatal(fmt.Sprintf("Expected `%v` in %s, got `%v`\n", expected, path, actual))
				return false
			}
			return true
		case reflect.Bool:
			//compare types
			if reflect.ValueOf(actual).Kind() != jsonType && ac.checkExists {
				log.Fatal(fmt.Sprintf("Expected type `bool` in %s, got %v\n", path, actual))
				return false
			}
			//compare values
			if expected.(bool) != actual.(bool) && ac.checkExists && ac.checkValues {
				log.Fatal(fmt.Sprintf("Expected `%v` in %s, got `%v`\n", expected, path, actual))
				return false
			}
			return true
		}
	}
	//not a json basic type, so its a json. call function on it again.
	return isEqualJSON(path, expected, actual, ac)
}

func isEqualJSON(path string, expected, actual interface{}, ac assertConfig) (ok bool) {
	//if expected is array,
	if expectedArray, ok := isArray(expected); ok {
		actualArray, ok := isArray(actual)
		//but actual is not array, return.
		if !ok && ac.checkExists { // if checking for existence, allow early return.
			return false
		}
		// To accomodate for notContains use cases,so that there are no false positives for this chunk below:
		/*`
		if !ac.checkExists{
			return true
		}
		*/
		// If expected and actual are both arrays,
		if ok {
			//but different lengths, return.
			if len(expectedArray) != len(actualArray) && ac.checkExists {
				log.Fatal(fmt.Sprintf("Expected array length of %d in %s, got %d\n",
					len(expectedArray), path, len(actualArray)))
				return false
			}
			//else start comparing each array index.
			for idx, element := range expectedArray {
				newPath := fmt.Sprintf("%s[%d]", path, idx)
				//if index elements are not equal, return
				if !isEqualElement(newPath, element, actualArray[idx], ac) && ac.checkExists {
					return false
				}
			}
		}
		//checking for existence, and all found, so return true.
		if ac.checkExists {
			return true
		}
		//not checking for existence, so all not found. return false.
		return false

	}
	//if expected is object,
	if expectedObj, ok := isObj(expected); ok {
		actualObj, ok := isObj(actual)
		//but actual is not object, return
		if !ok && ac.checkExists {
			return false
		}
		//start comparing each object element.
		for key, element := range expectedObj {
			actualElement, ok := actualObj[key]
			//if key does not exist in actual object, return.
			if !ok && ac.checkExists {
				log.Fatal(fmt.Sprintf("Expected key `%s` in %s\n", key, path))
				return false
			}

			newPath := fmt.Sprintf(`%s["%s"]`, path, key)
			// To accomodate for notContains use cases,so that there are no false positives for this chunk below:
			/*`
			if !ac.checkExists{
				return true
			}
			*/
			// If key exists in both actual and expected,
			if ok {
				//and their values are not equal, return.
				if !isEqualElement(newPath, element, actualElement, ac) && ac.checkExists {
					return false
				}
				//key exists in both actual and expected
				//not comparing values AND not checking for existence.
				//return true to indicate failure.
				if !ac.checkValues && !ac.checkExists {
					return true
				}
			}
		}
		//checking for existence and all found, so return true.
		if ac.checkExists {
			return true
		}
		//else, not checking for non-existence and all not found so return false.
		return false
	}
	//not a valid json, throw fatal.
	log.Fatal("Expected value at %s is not json: `%v`\n", path, expected)
	//unreachable, so safe to use any value
	return false
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
		"status": "fail"	
	}`
	fmt.Println("assertJSON =", assertJSON([]byte(expectedStr), []byte(actualStr), assertConfig{
		checkValues: false,
		checkExists: false,
	}))
}
