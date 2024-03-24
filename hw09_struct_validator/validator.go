package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrExpectedStruct      = errors.New("expected struct")
	ErrMinValueViolation   = errors.New("minimum constraint violation")
	ErrMaxValueViolation   = errors.New("maximum constraint violation")
	ErrRegexpViolation     = errors.New("regexp does not match")
	ErrLenViolation        = errors.New("length violation")
	ErrNotInRangeViolation = errors.New("value is not in range")
)

var validationTag = "validate"

var (
	ruleMax    = "max"
	ruleMin    = "min"
	ruleLen    = "len"
	ruleRegexp = "regexp"
	ruleRange  = "in"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStrs := strings.Builder{}

	for i, err := range v {
		errStrs.WriteString(fmt.Sprintf("%s %s", err.Field, err.Err))

		if i != len(v)-1 {
			errStrs.WriteString("\n ")
		}
	}

	return errStrs.String()
}

func Validate(v interface{}) error {
	errs := make(ValidationErrors, 0)

	iv := reflect.ValueOf(v)
	if iv.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}

	t := iv.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		valueOfField := iv.Field(i)

		validateTag, ok := field.Tag.Lookup(validationTag)

		if !ok {
			continue
		}

		errs = checkValue(errs, field.Name, validateTag, valueOfField)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func checkValue(
	valErrs ValidationErrors,
	fName string, validateTag string,
	valueOfField reflect.Value,
) ValidationErrors {
	var errs []error
	newValErrs := valErrs

	if !valueOfField.CanInterface() {
		return nil
	}

	switch valueOfField.Interface().(type) {
	case string:
		errs = validateValue(validateTag, valueOfField)
	case int:
		errs = validateValue(validateTag, valueOfField)
	case []string:
		for i := 0; i < valueOfField.Len(); i++ {
			newValErrs = checkValue(newValErrs, fName, validateTag, valueOfField.Index(i))
		}
	case []int:
		for i := 0; i < valueOfField.Len(); i++ {
			newValErrs = checkValue(newValErrs, fName, validateTag, valueOfField.Index(i))
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			newValErrs = append(newValErrs, ValidationError{fName, err})
		}
	}

	return newValErrs
}

func validateValue(validateTag string, rv reflect.Value) []error {
	rules := strings.Split(validateTag, "|")
	errs := make([]error, 0)

	for _, rule := range rules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			continue
		}

		ruleType, ruleValue := r[0], r[1]

		var err error
		tmpErrMessage := "%w: %s; "

		switch ruleType {
		case ruleLen:
			if !checkLen(rv, ruleValue) {
				err = fmt.Errorf(tmpErrMessage, ErrLenViolation, ruleValue)
			}
		case ruleRegexp:
			if !checkRegex(rv, ruleValue) {
				err = fmt.Errorf(tmpErrMessage, ErrRegexpViolation, ruleValue)
			}
		case ruleMin:
			if !checkMin(rv, ruleValue) {
				err = fmt.Errorf(tmpErrMessage, ErrMinValueViolation, ruleValue)
			}
		case ruleMax:
			if !checkMax(rv, ruleValue) {
				err = fmt.Errorf(tmpErrMessage, ErrMaxValueViolation, ruleValue)
			}
		case ruleRange:
			if !checkIn(rv, ruleValue) {
				err = fmt.Errorf(tmpErrMessage, ErrNotInRangeViolation, ruleValue)
			}
		default:
			continue
		}

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func checkLen(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		intValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}

		return rv.Len() == intValue
	}

	return false
}

func checkRegex(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		rx, err := regexp.Compile(ruleValue)
		if err != nil {
			return false
		}

		return rx.Match([]byte(rv.String()))
	}

	return false
}

func checkMin(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		min, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}

		return intValue > min
	}

	return false
}

func checkMax(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		max, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}

		return intValue < max
	}

	return false
}

func checkIn(rv reflect.Value, ruleValue string) bool {
	ins := strings.Split(ruleValue, ",")
	isValid := false

	//nolint:exhaustive
	switch rv.Kind() {
	case reflect.Int:
		intValue := int(rv.Int())

		for _, in := range ins {
			in, err := strconv.Atoi(in)
			if err != nil {
				continue
			}

			if in == intValue {
				isValid = true
			}
		}
	case reflect.String:
		strValue := rv.String()

		for _, in := range ins {
			if in == strValue {
				isValid = true
			}
		}
	}
	return isValid
}
