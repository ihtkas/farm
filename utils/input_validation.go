package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
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

// TimeRange returns a validation rule that checks if a value's timestamp is within the specified range.
// If max is nil, it means there is no upper bound for the timestamp.
// This rule should only be used for validating strings, slices, maps, and arrays.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func TimeRange(min, max time.Time) *TimeRule {
	message := "the value must be empty"
	if !(min.IsZero() && max.IsZero()) {
		if min.IsZero() {
			message = fmt.Sprintf("the timestamp must be no more than %v", max)
		} else if max.IsZero() {
			message = fmt.Sprintf("the timestamp must be no less than %v", min)
		} else if min.Equal(max) {
			message = fmt.Sprintf("the timestamp must be exactly %v", min)
		} else {
			message = fmt.Sprintf("the timestamp must be between %v and %v", min.String(), max.String())
		}

	}
	return &TimeRule{
		min:     min,
		max:     max,
		message: message,
	}
}

// TimeRule implements a validation rule with min and max range for time
type TimeRule struct {
	min, max time.Time
	message  string
	rune     bool
}

// Validate checks if the given value is valid or not.
func (v *TimeRule) Validate(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	var l time.Time
	if t, ok := value.(time.Time); ok && v.rune {
		l = t
	} else {
		return errors.New("Invalid timestmap value")
	}

	if !v.min.IsZero() && l.Before(v.min) || !v.max.IsZero() && l.After(v.max) {
		return errors.New(v.message)
	}
	return nil
}

// Error sets the error message for the rule.
func (v *TimeRule) Error(message string) *TimeRule {
	v.message = message
	return v
}
