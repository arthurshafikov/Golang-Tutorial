package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type StructField struct {
	Key        string
	Value      interface{}
	ValueKind  reflect.Kind
	Validation string
}

type Rule struct {
	Signature string
	Value     string
}

type ValidateFunc func(value interface{}, rule Rule) error

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errors string
	for _, valErr := range v {
		errors += fmt.Sprintln(valErr.Field, valErr.Err)
	}

	return errors
}

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrUnsupportedType = errors.New("unsupported kind of input")
	ErrWrongRuleSyntax = errors.New("rule has wrong syntax")
	ErrUndefinedRule   = errors.New("rule is undefined")
)

const (
	maxWrongMes    = "maximum length is %v"
	minWrongMes    = "minimum length is %v"
	minIntWrongMes = "minimum value is %v"
	maxIntWrongMes = "maximum value is %v"
	lenWrongMes    = "length must be exactly %v"
	regexpWrongMes = "has wrong value"
	inWrongMes     = "must be in %v"
)

func validateByType(value interface{}, rules []string, valFun ValidateFunc) error {
	var validationError error

	for _, ruleString := range rules {
		ruleSplit := strings.Split(ruleString, ":")

		if len(ruleSplit) != 2 && ruleString != "nested" {
			return ErrWrongRuleSyntax
		}

		if ruleString == "nested" {
			ruleSplit = []string{ruleString, ""}
		}

		validationError = valFun(value, Rule{
			Signature: ruleSplit[0],
			Value:     ruleSplit[1],
		})
	}

	return validationError
}

func validateField(field StructField) ValidationError {
	var validateFieldResponse ValidationError

	var validationError error

	rules := strings.Split(field.Validation, "|")
	switch field.ValueKind {
	case reflect.String:
		validationError = validateByType(field.Value.(reflect.Value).String(), rules, validateString)
	case reflect.Int:
		validationError = validateByType(int(field.Value.(reflect.Value).Int()), rules, validateInt)
	case reflect.Slice:
		validationError = validateByType(field.Value, rules, validateSlice)
	case reflect.Struct:
		validationError = validateByType(field.Value, rules, validateStruct)
	default:
		validationError = ErrUnsupportedType
	}

	if validationError != nil {
		validateFieldResponse = ValidationError{
			Field: field.Key,
			Err:   validationError,
		}
	}

	return validateFieldResponse
}

func Validate(v interface{}) error {
	structReflect := reflect.ValueOf(v)

	result := ValidationErrors{}
	for i := 0; i < structReflect.Type().NumField(); i++ {
		fieldValue := structReflect.Field(i)
		fieldStruct := structReflect.Type().Field(i)

		if fieldStruct.PkgPath != "" { // IsExported go 1.17
			continue
		}

		field := StructField{
			Key:        fieldStruct.Name,
			ValueKind:  fieldStruct.Type.Kind(),
			Value:      fieldValue,
			Validation: fieldStruct.Tag.Get("validate"),
		}

		if field.Validation != "" && fieldValue.String() != "" {
			resultField := validateField(field)
			if resultField.Err != nil {
				result = append(result, resultField)
			}
		}
	}
	if len(result) > 0 {
		return result
	}

	return nil
}
