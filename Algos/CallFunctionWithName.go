package Algos

import (
	"fmt"
	"reflect"
)

// CallFunctionWithName calls a function that is stored in a map with its name as key. It takes in a slice containing parameters for that function.
// Returns an error (if any).
func CallFunctionWithName(funcName string, params []interface{}) error {
	// Attempt to extract function using this name as key
	function := reflect.ValueOf(SendingFunctionsMap[funcName])
	// function does not exist (zero value returned from map)
	// OR function is not a function (shouldn't happen if map only contains functions as elements)
	if !function.IsValid() || function.Kind() != reflect.Func {
		return fmt.Errorf(errFunctionNotFound, funcName)
	}
	if len(params) != function.Type().NumIn() {
		return fmt.Errorf(errParamsMismatch, funcName)
	}

	// initialise parameters (slice of reflect.Value)
	parameters := make([]reflect.Value, len(params))
	for idx, param := range params {
		parameters[idx] = reflect.ValueOf(param)
	}
	// call function with reflect.Call
	results := function.Call(parameters)

	funcErr := reflect.Value{}
	// store error value depending on which position it is in the list of return values
	switch funcName {
	case "RequestVerify":
		funcErr = results[1] // error is 2nd return value
	default:
		funcErr = results[0] // for the rest, error is 1st return value
	}

	// equivalent to 'if err != nil'
	if funcErr.Kind() != reflect.UnsafePointer {
		// type assert to error.
		err, ok := funcErr.Interface().(error)
		if !ok {
			return fmt.Errorf(errFunctionResultNotError, function)
		}
		return err
	}
	return nil
}
