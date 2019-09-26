package main

import "fmt"

// GetParamsMap gets a map of values from a list of keys that are mandatory.
// It returns an error if a parameter is not found.
// The function as argument must return the value of the parameter or an empty string if it's not found.
func GetParamsMap(getter func(string) string, params ...string) (map[string]string, error) {
	result := make(map[string]string)
	if len(params) == 0 {
		return result, fmt.Errorf("No parameter passed to get")
	}
	for _, param := range params {
		value := getter(param)
		if value == "" {
			return result, fmt.Errorf("Parameter %s is mandatory", param)
		}
		result[param] = value
	}
	return result, nil
}
