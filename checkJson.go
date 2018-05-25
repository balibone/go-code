package main

import (
	"log"
	"strconv"
)

//JsonContainsJson checks if baseJson contains all json data defined in compareJson.
func JsonContainsJson(baseJson interface{}, compareJson interface{}, path string) bool {
	//recurse until you hit a pure key:val. At every recurse, append and pass in path. When pure key:val is hit, call CheckJson method and pass in compareKey, compareValue, comparePath.
	switch compareJson.(type) {
	case []interface{}:
		//type assert.
		jsonArray, ok := compareJson.([]interface{})
		if !ok {
			log.Fatal("Error asserting json array in JsonContainsJson")
		}
		//loop through array index and apply function to each index.
		for arrIndex, arrValue := range jsonArray {
			pathToArray := path
			//Every index in json array MUST be another array OR another object.
			//so just recurse
			if found := JsonContainsJson(baseJson, arrValue, path+"arr["+strconv.Itoa(arrIndex)+"]."); !found {
				return false
			}
			//reset path.
			path = pathToArray
		}
	case map[string]interface{}:
		//type assert.
		jsonObject, ok := compareJson.(map[string]interface{})
		if !ok {
			log.Fatal("Error asserting json object in JsonContainsJson")
		}
		for key, value := range jsonObject {
			pathToObject := path
			switch value.(type) {
			case []interface{}, map[string]interface{}:
				if found := JsonContainsJson(baseJson, jsonObject[key], path+key+"."); !found {
					return false
				}
			default:
				// fmt.Printf(test_utilities.ColorYellow("Searching for %v (%T) : %v (%T)\n"), key, key, value, value)
				//pure key:val hit. call CheckJson and return bool
				if found := CheckJson(baseJson, "", key, value, path+key+"."); !found {
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
	//if this is reached, means everything in compareJson was found.
	return true
}

// CheckJson takes a valid json and recursively searches it to validate hierarchical structure and return true on an exact match of "compareKey:compareValue" pair.
func CheckJson(validJson interface{}, path string, compareKey string, compareValue interface{}, comparePath string) bool {
	//basic idea: keep recursing until you hit a json object. When you hit an object, loop through that object and check each object attribute. If attribute is array or object, recurse again. Else, compare key:value pairs.  If equal, return true all the way up. If not equal, continue searching. If search exhausted, return false.

	//check type of validJson
	switch validJson.(type) {
	case []interface{}:
		//type assert.
		jsonArray, ok := validJson.([]interface{})
		if !ok {
			log.Fatal("Error asserting json array in CheckJson")
		}
		// fmt.Printf(test_utilities.ColorYellow("Searching Array --> %v (%T)\n"), jsonArray, jsonArray)
		//loop through array index and apply function to each index.
		for arrIndex, arrValue := range jsonArray {
			pathToArray := path
			//Every index in json array MUST be another array OR another object.
			//so just recurse
			if found := CheckJson(arrValue, path+"arr["+strconv.Itoa(arrIndex)+"].", compareKey, compareValue, comparePath); found {
				return true
			}
			//reset path.
			path = pathToArray
		}
		return false
	case map[string]interface{}:
		//type assert.
		jsonObject, ok := validJson.(map[string]interface{})
		if !ok {
			log.Fatal("Error asserting json object in CheckJson")
		}
		// fmt.Printf(test_utilities.ColorYellow("Searching Object --> %v (%T)\n"), jsonObject, jsonObject)
		//loop through key value pairs and see if any of them match compareKey: compareValue
		for key, value := range jsonObject {
			pathToObject := path
			switch value.(type) {
			case []interface{}, map[string]interface{}:
				if found := CheckJson(value, path+key+".", compareKey, compareValue, comparePath); found {
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
