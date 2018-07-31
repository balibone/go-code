package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

//JSONContainsJSON checks if baseJSON contains all json data defined in compareJSON.
func JSONContainsJSON(baseJSON interface{}, compareJSON interface{}, path string) bool {
	//recurse until you hit a pure key:val. At every recurse, append and pass in path. When pure key:val is hit, call CheckJSON method and pass in compareKey, compareValue, comparePath.
	switch compareJSON.(type) {
	case []interface{}:
		//type assert.
		jsonArray, ok := compareJSON.([]interface{})
		if !ok {
			log.Fatal("Error asserting json array in JSONContainsJSON")
		}
		//loop through array index and apply function to each index.
		for arrIndex, arrValue := range jsonArray {
			pathToArray := path
			//Every index in json array MUST be another array OR another object.
			//so just recurse
			if found := JSONContainsJSON(baseJSON, arrValue, path+"arr["+strconv.Itoa(arrIndex)+"]."); !found {
				return false
			}
			//reset path.
			path = pathToArray
		}
	case map[string]interface{}:
		//type assert.
		jsonObject, ok := compareJSON.(map[string]interface{})
		if !ok {
			log.Fatal("Error asserting json object in JSONContainsJSON")
		}
		for key, value := range jsonObject {
			pathToObject := path
			switch value.(type) {
			case []interface{}, map[string]interface{}:
				if found := JSONContainsJSON(baseJSON, jsonObject[key], path+key+"."); !found {
					return false
				}
			default:
				// fmt.Printf(test_utilities.ColorYellow("Searching for %v (%T) : %v (%T)\n"), key, key, value, value)
				//pure key:val hit. call CheckJSON and return bool
				if found := CheckJSON(baseJSON, "", key, value, path+key+"."); !found {
					// fmt.Printf(test_utilities.ColorRed("The key-value pair %v (%T): %v (%T) was not found in response json.\n"), key, key, value, value)
					return false
				} else {
					// fmt.Printf(test_utilities.ColorCyan("The key-value pair %v (%T): %v (%T) was found in response json.\n"), key, key, value, value)
				}
			}
			//if found, continue for loop
			path = pathToObject
		}
	}
	//if this is reached, means everything in compareJSON was found.
	return true
}

// CheckJSON takes a valid json and recursively searches it to validate hierarchical structure and return true on an exact match of "compareKey:compareValue" pair.
func CheckJSON(validJSON interface{}, path string, compareKey string, compareValue interface{}, comparePath string) bool {
	//basic idea: keep recursing until you hit a json object. When you hit an object, loop through that object and check each object attribute. If attribute is array or object, recurse again. Else, compare key:value pairs.  If equal, return true all the way up. If not equal, continue searching. If search exhausted, return false.

	//check type of validJSON
	switch validJSON.(type) {
	case []interface{}:
		//type assert.
		jsonArray, ok := validJSON.([]interface{})
		if !ok {
			log.Fatal("Error asserting json array in CheckJSON")
		}
		// fmt.Printf(test_utilities.ColorYellow("Searching Array --> %v (%T)\n"), jsonArray, jsonArray)
		//loop through array index and apply function to each index.
		for arrIndex, arrValue := range jsonArray {
			pathToArray := path
			//Every index in json array MUST be another array OR another object.
			//so just recurse
			if found := CheckJSON(arrValue, path+"arr["+strconv.Itoa(arrIndex)+"].", compareKey, compareValue, comparePath); found {
				return true
			}
			//reset path.
			path = pathToArray
		}
		return false
	case map[string]interface{}:
		//type assert.
		jsonObject, ok := validJSON.(map[string]interface{})
		if !ok {
			log.Fatal("Error asserting json object in CheckJSON")
		}
		// fmt.Printf(test_utilities.ColorYellow("Searching Object --> %v (%T)\n"), jsonObject, jsonObject)
		//loop through key value pairs and see if any of them match compareKey: compareValue
		for key, value := range jsonObject {
			pathToObject := path
			switch value.(type) {
			case []interface{}, map[string]interface{}:
				if found := CheckJSON(value, path+key+".", compareKey, compareValue, comparePath); found {
					return true
				}
			default:
				path = path + key + "."
				if key == compareKey && value == compareValue && path == comparePath {
					return true
				}
			}
			//reset path.
			path = pathToObject
			//continue looking.
		}
		return false
	}
	return false
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
