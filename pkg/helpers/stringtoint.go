package helpers

import (
	"strconv"
)

// StringToInt convert string to int
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// IntToString convert int to string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}


// uint to int
func UintToInt(u uint64) int {
	return int(u)
}
