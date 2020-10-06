package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// GetStringParam checks if give key exists and returns the value. if the key doesn't exist, retuns error
func GetStringParam(values url.Values, key string) (string, error) {
	val := values.Get(key)
	if val == "" {
		return "", errors.New("Couldn't find field '" + key + "'")
	}
	return val, nil
}

// GetIntegerParam checks if give key exists and returns a number if valid. if the key doesn't exist or not a valid number, retuns error
func GetIntegerParam(values url.Values, key string) (int64, error) {
	val := values.Get(key)
	if val == "" {
		return 0, errors.New("Couldn't find field '" + key + "'")
	}
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Expecting a valid integer value. Parser error: %w", err)
	}
	return num, nil
}
