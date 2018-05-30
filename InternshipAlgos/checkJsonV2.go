package main

import "fmt"
import "encoding/json"
import "reflect"
import "strconv"
import "time"

func jsonContains(expected, actual interface{}) (ok bool) {
	var isArray func(interface{}) ([]interface{}, bool)
	var isObj func(interface{}) (map[string]interface{}, bool)
	var isJSONBasicType func(interface{}) (reflect.Kind, bool)
	var isEqualElement func(string, interface{}, interface{}) bool
	var isEqualJSON func(string, interface{}, interface{}) bool

	isArray = func(input interface{}) (array []interface{}, ok bool) {
		if reflect.ValueOf(input).Kind() != reflect.Slice {
			return
		}
		return input.([]interface{}), true
	}

	isObj = func(input interface{}) (obj map[string]interface{}, ok bool) {
		if reflect.ValueOf(input).Kind() != reflect.Map {
			return
		}
		return input.(map[string]interface{}), true
	}

	isJSONBasicType = func(input interface{}) (JSONType reflect.Kind, ok bool) {
		if reflect.TypeOf(input) == nil {
			// JSON null
			return reflect.UnsafePointer, true
		}
		switch input.(type) {
		case int64:
			// JSON number - integer
			return reflect.Int64, true
		case float64:
			// JSON number - float
			return reflect.Float64, true
		case string:
			// JSON string
			return reflect.String, true
		case bool:
			// JSON bool
			return reflect.Bool, true
		}
		return
	}
	
	isEqualElement = func(path string, expected, actual interface{}) bool {
		if jsonType, ok := isJSONBasicType(expected); ok {
			switch jsonType {
			case reflect.UnsafePointer:
				if reflect.TypeOf(actual) != nil {
					fmt.Printf("Expected type `null` in %s, got `%v`\n", path, actual)
					return false
				}
				return true
			case reflect.Int64:
				if reflect.ValueOf(actual).Kind() != jsonType {
					fmt.Printf("Expected type `int` in %s, got %v\n", path, actual)
					return false
				}
				if expected.(int64) != actual.(int64) {
					fmt.Printf("Expected `%v` in %s, got `%v`\n", expected, path, actual)
					return false
				}
				return true
			case reflect.Float64:
				if reflect.ValueOf(actual).Kind() != jsonType {
					fmt.Printf("Expected type `float` in %s, got %v\n", path, actual)
					return false
				}
				if expected.(float64) != actual.(float64) {
					fmt.Printf("Expected `%v` in %s, got `%v`\n", expected, path, actual)
					return false
				}
				return true
			case reflect.String:
				if reflect.ValueOf(actual).Kind() != jsonType {
					fmt.Printf("Expected type `string` in %s, got %v\n", path, actual)
					return false
				}
				if expected.(string) != actual.(string) {
					fmt.Printf("Expected `%v` in %s, got `%v`\n", expected, path, actual)
					return false
				}
				return true
			case reflect.Bool:
				if reflect.ValueOf(actual).Kind() != jsonType {
					fmt.Printf("Expected type `bool` in %s, got %v\n", path, actual)
					return false
				}
				if expected.(bool) != actual.(bool) {
					fmt.Printf("Expected `%v` in %s, got `%v`\n", expected, path, actual)
					return false
				}
				return true
			}
		}

		return isEqualJSON(path, expected, actual)
	}

	isEqualJSON = func(path string, expected, actual interface{}) (ok bool) {
		if expectedArray, ok := isArray(expected); ok {
			actualArray, ok := isArray(actual)
			if !ok {
				return false
			}
			if len(expectedArray) != len(actualArray) {
				fmt.Printf("Expected array length of %d in %s, got %d\n",
					len(expectedArray), path, len(actualArray))
				return false
			}
			for idx, element := range expectedArray {
				newPath := fmt.Sprintf("%s[%d]", path, idx)
				if !isEqualElement(newPath, element, actualArray[idx]) {
					return false
				}
			}

			return true
		}

		if expectedObj, ok := isObj(expected); ok {
			actualObj, ok := isObj(actual)
			if !ok {
				return false
			}
			for key, element := range expectedObj {
				actualElement, ok := actualObj[key]
				if !ok {
					fmt.Printf("Expected key `%s` in %s\n", key, path)
					return false
				}
				newPath := fmt.Sprintf(`%s["%s"]`, path, key)
				if !isEqualElement(newPath, element, actualElement) {
					return false
				}
			}

			return true
		}

		fmt.Printf("Expected value at %s is not json: `%v`\n", path, expected)
		return false
	}

	return isEqualJSON("actual", expected, actual)
}

func jsonExists(input interface{}, path []string, cb func(interface{})) (ok bool) {
	if len(path) == 0 {
		fmt.Println("Empty path, nothing to check!")
		return
	}
	
	var isArray func(interface{}) ([]interface{}, bool)
	var isObj func(interface{}) (map[string]interface{}, bool)
	
	isArray = func(input interface{}) (array []interface{}, ok bool) {
		if reflect.ValueOf(input).Kind() != reflect.Slice {
			return
		}
		return input.([]interface{}), true
	}

	isObj = func(input interface{}) (obj map[string]interface{}, ok bool) {
		if reflect.ValueOf(input).Kind() != reflect.Map {
			return
		}
		return input.(map[string]interface{}), true
	}
	
	currentKey := path[0]
	
	if array, ok := isArray(input); ok {
		if idx, err := strconv.Atoi(currentKey); err != nil {
			fmt.Println("Cannot use non-integer as index for array")
			return false
		} else if idx >= len(array) {
			fmt.Println("Index out of range")
			return false
		} else {
			if len(path) == 1 {
				cb(array[idx])
				return true
			} else {
				return jsonExists(array[idx], path[1:], cb)
			}
		}
	}

	if obj, ok := isObj(input); ok {
		if value, ok := obj[currentKey]; !ok {
			fmt.Println("key not in object")
			return false
		} else {
			if len(path) == 1 {
				cb(value)
				return true
			} else {
				return jsonExists(value, path[1:], cb)
			}
		}
	}
	
	return false	
}

func main() {
	var expected interface{}
	var actual interface{}
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
            "share_referral": "ABCDEF",
            "currency": "SGD"
        },
        "balance": {
            "currency": "SGD",
            "amount": "0.00"
        },
        "rewards": null,
        "new": false,
        "verified": true,
        "email_verified": true,
        "phone_verified": true,
        "has_password": true,
        "currency": "SGD"
    }
	}`
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
	
	var err error
	err = json.Unmarshal([]byte(expectedStr), &expected)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	err = json.Unmarshal([]byte(actualStr), &actual)
	if err != nil {
		fmt.Println(err)
		return
	}
	start := time.Now()
	fmt.Println("jsonContains =", jsonContains(expected, actual))
	fmt.Println(time.Now().Sub(start))
}
