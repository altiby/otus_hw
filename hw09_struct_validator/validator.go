package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNoStructType            = errors.New("validate support only struct type")
	ErrUnsupportedFieldType    = errors.New("unsupported field type")
	ErrUnsupportedConstraint   = errors.New("unsupported constraint")
	ErrConstraintCheckFailed   = errors.New("constraint check failed")
	ErrInvalidConstraintFormat = errors.New("invalid constraint format")
)

type ValidationError struct {
	Field string
	Err   error
}

type ConstraintType string

const (
	IntConstraintTypeMin ConstraintType = "min"
	IntConstraintTypeMax ConstraintType = "max"
	IntConstraintTypeIn  ConstraintType = "in"

	StringConstraintTypeLen    ConstraintType = "len"
	StringConstraintTypeRegexp ConstraintType = "regexp"
	StringConstraintTypeIn     ConstraintType = "in"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	ret := ""
	counter := 0
	for _, validationError := range v {
		counter++
		ret += fmt.Sprintf("#%d. field: %s, %s", counter, validationError.Field, validationError.Err)
	}
	return ret
}

func Validate(v interface{}) error {
	errors := ValidationErrors{}
	rValue := reflect.ValueOf(v)
	if rValue.Kind() != reflect.Struct {
		return ErrNoStructType
	}
	rType := rValue.Type()

	for i := 0; i < rType.NumField(); i++ {
		field := rValue.Field(i)
		fieldType := rType.Field(i)
		fieldTags := fieldType.Tag
		validateTag := fieldTags.Get("validate")
		if validateTag == "" {
			continue
		}
		if field.CanInterface() {
			constraints := getConstraints(validateTag)
			for _, constraint := range constraints {
				err := checkConstraint(constraint, field)
				if err != nil {
					errors = append(errors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func checkConstraint(constraint string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		return checkStringConstraint(constraint, field)
	case reflect.Int:
		return checkIntConstraint(constraint, field)
	case reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			err := checkConstraint(constraint, field.Index(i))
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%s, %w", field.Kind(), ErrUnsupportedFieldType)
	}
}

func checkIntConstraint(constraint string, field reflect.Value) error {
	constraintType, constraintValue, err := parseConstraint(constraint)
	if err != nil {
		return err
	}
	switch constraintType {
	case IntConstraintTypeMin:
		val, err := strconv.Atoi(constraintValue)
		if err != nil {
			return ErrInvalidConstraintFormat
		}
		if field.Int() < int64(val) {
			return ErrConstraintCheckFailed
		}
	case IntConstraintTypeMax:
		val, err := strconv.Atoi(constraintValue)
		if err != nil {
			return ErrInvalidConstraintFormat
		}
		if field.Int() > int64(val) {
			return ErrConstraintCheckFailed
		}
	case IntConstraintTypeIn:
		listIn := strings.Split(constraintValue, ",")
		found := false
		for _, s := range listIn {
			val, err := strconv.Atoi(s)
			if err != nil {
				return ErrInvalidConstraintFormat
			}
			if field.Int() == int64(val) {
				found = true
			}
		}
		if !found {
			return ErrConstraintCheckFailed
		}
	default:
		return ErrUnsupportedConstraint
	}
	return nil
}

func checkStringConstraint(constraint string, field reflect.Value) error {
	constraintType, constraintValue, err := parseConstraint(constraint)
	if err != nil {
		return err
	}
	switch constraintType {
	case StringConstraintTypeLen:
		val, err := strconv.Atoi(constraintValue)
		if err != nil {
			return ErrInvalidConstraintFormat
		}
		if len(field.String()) != val {
			return ErrConstraintCheckFailed
		}
	case StringConstraintTypeRegexp:
		mateched, err := regexp.MatchString(constraintValue, field.String())
		if err != nil {
			return ErrInvalidConstraintFormat
		}
		if !mateched {
			return ErrConstraintCheckFailed
		}
	case StringConstraintTypeIn:
		listIn := strings.Split(constraintValue, ",")
		found := false
		for _, s := range listIn {
			if field.String() == strings.TrimSpace(s) {
				found = true
			}
		}
		if !found {
			return ErrConstraintCheckFailed
		}
	default:
		return ErrUnsupportedConstraint
	}
	return nil
}

func parseConstraint(constraint string) (ConstraintType, string, error) {
	parts := strings.Split(constraint, ":")
	if len(parts) != 2 {
		return "", "", ErrInvalidConstraintFormat
	}
	return ConstraintType(parts[0]), parts[1], nil
}

func getConstraints(tag string) []string {
	ret := strings.Split(strings.TrimSpace(tag), "|")
	for i := 0; i < len(ret); i++ {
		ret[i] = strings.TrimSpace(ret[i])
	}
	return ret
}
