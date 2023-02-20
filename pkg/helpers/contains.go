package helpers

import "strings"


func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// check if a substring is part of string
func ContainsString(s string, substr string) bool {
	if strings.Contains(s, substr) {
		return true
	}
	return false
}
