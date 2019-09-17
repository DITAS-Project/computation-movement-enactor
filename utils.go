package main

import "fmt"

func GetParamsMap(getter func(string) string, params ...string) (map[string]string, error) {
	result := make(map[string]string)
	for _, param := range params {
		value := getter(param)
		if value == "" {
			return result, fmt.Errorf("Parameter %s is mandatory", param)
		}
		result[param] = value
	}
	return result, nil
}
