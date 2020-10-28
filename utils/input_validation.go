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
}

// Validate checks if the given value is valid or not.
func (v *TimeRule) Validate(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	var l time.Time
	if t, ok := value.(time.Time); ok {
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

// Float64Range returns a validation rule that checks if a value's timestamp is within the specified range.
// If max is nil, it means there is no upper bound for the timestamp.
func Float64Range(min, max float64, hasMinBound, hasMaxBound bool) *Float64Rule {
	message := "the value must be empty"
	if !(hasMinBound || hasMaxBound) {
		if !hasMinBound {
			message = fmt.Sprintf("the timestamp must be no more than %v", max)
		} else if !hasMaxBound {
			message = fmt.Sprintf("the timestamp must be no less than %v", min)
		} else if min == max {
			message = fmt.Sprintf("the timestamp must be exactly %v", min)
		} else {
			message = fmt.Sprintf("the timestamp must be between %v and %v", min, max)
		}

	}
	return &Float64Rule{
		min:         min,
		max:         max,
		message:     message,
		hasMinBound: hasMinBound,
		hasMaxBound: hasMaxBound,
	}
}

// Float64Rule implements a validation rule with min and max range for float64
type Float64Rule struct {
	min, max                 float64
	message                  string
	hasMinBound, hasMaxBound bool
}

// Validate checks if the given value is valid or not.
func (v *Float64Rule) Validate(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	var l float64
	if t, ok := value.(float64); ok {
		l = t
	} else {
		return errors.New("Invalid float64stmap value")
	}

	if v.hasMinBound && (l < v.min) || v.hasMaxBound && (l > v.max) {
		return errors.New(v.message)
	}
	return nil
}

// Error sets the error message for the rule.
func (v *Float64Rule) Error(message string) *Float64Rule {
	v.message = message
	return v
}

// Uint32Range returns a validation rule that checks if a value is within the specified range.
// If max is nil, it means there is no upper bound for the uint32 value.
func Uint32Range(min, max uint32, hasMinBound, hasMaxBound bool) *Uint32Rule {
	message := "the value must be empty"
	if !(hasMinBound || hasMaxBound) {
		if !hasMinBound {
			message = fmt.Sprintf("the timestamp must be no more than %v", max)
		} else if !hasMaxBound {
			message = fmt.Sprintf("the timestamp must be no less than %v", min)
		} else if min == max {
			message = fmt.Sprintf("the timestamp must be exactly %v", min)
		} else {
			message = fmt.Sprintf("the timestamp must be between %v and %v", min, max)
		}

	}
	return &Uint32Rule{
		min:         min,
		max:         max,
		message:     message,
		hasMinBound: hasMinBound,
		hasMaxBound: hasMaxBound,
	}
}

// Uint32Rule implements a validation rule with min and max range for uint32
type Uint32Rule struct {
	min, max                 uint32
	message                  string
	hasMinBound, hasMaxBound bool
}

// Validate checks if the given value is valid or not.
func (v *Uint32Rule) Validate(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	var l uint32
	if t, ok := value.(uint32); ok {
		l = t
	} else {
		return errors.New("Invalid uint32 value")
	}

	if v.hasMinBound && (l < v.min) || v.hasMaxBound && (l > v.max) {
		return errors.New(v.message)
	}
	return nil
}

// Error sets the error message for the rule.
func (v *Uint32Rule) Error(message string) *Uint32Rule {
	v.message = message
	return v
}
