package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func validateStruct(value interface{}, rule Rule) error {
	validationError := Validate(value)

	return validationError
}

func validateInt(value interface{}, rule Rule) error {
	var validationError error

	ruleValueInt, err := getRuleValueInt(rule)
	if err != nil {
		return err
	}

	valueInt := value.(int)
	switch rule.Signature {
	case "max":
		if valueInt > ruleValueInt {
			validationError = fmt.Errorf(maxIntWrongMes, ruleValueInt)
		}
	case "min":
		if valueInt < ruleValueInt {
			validationError = fmt.Errorf(minIntWrongMes, ruleValueInt)
		}
	case "in":
		ruleValuesSlice := strings.Split(rule.Value, ",")
		maximum, err := strconv.Atoi(ruleValuesSlice[1])
		if err != nil {
			return err
		}
		minimum, err := strconv.Atoi(ruleValuesSlice[0])
		if err != nil {
			return err
		}
		if valueInt > maximum || valueInt < minimum {
			validationError = fmt.Errorf(inWrongMes, rule.Value)
		}
	default:
		return ErrUndefinedRule
	}

	return validationError
}

func validateString(value interface{}, rule Rule) error {
	var validationError error

	ruleValueInt, err := getRuleValueInt(rule)
	if err != nil {
		return err
	}

	valueString := value.(string)
	switch rule.Signature {
	case "len":
		if len(valueString) != ruleValueInt {
			validationError = fmt.Errorf(lenWrongMes, ruleValueInt)
		}
	case "regexp":
		match, err := regexp.MatchString(rule.Value, valueString)
		if !match || err != nil {
			validationError = fmt.Errorf(regexpWrongMes)
		}
	case "max":
		if len(valueString) > ruleValueInt {
			validationError = fmt.Errorf(maxWrongMes, ruleValueInt)
		}
	case "min":
		if len(valueString) < ruleValueInt {
			validationError = fmt.Errorf(minWrongMes, ruleValueInt)
		}
	case "in":
		validationError = fmt.Errorf(inWrongMes, rule.Value)
		for _, ruleValue := range strings.Split(rule.Value, ",") {
			if value == ruleValue {
				validationError = nil
				break
			}
		}
	default:
		return ErrUndefinedRule
	}

	return validationError
}

func validateSlice(value interface{}, rule Rule) error {
	var validationError error

	slice := value.(reflect.Value).Interface()

	switch reflect.TypeOf(slice).Elem().Kind() {
	case reflect.String:
		for _, el := range slice.([]string) {
			validationError = validateString(el, rule)
			if validationError != nil {
				break
			}
		}
	case reflect.Int:
		for _, el := range slice.([]int) {
			validationError = validateInt(el, rule)
			if validationError != nil {
				break
			}
		}
	default:
		return ErrUnsupportedType
	}

	return validationError
}

func getRuleValueInt(rule Rule) (ruleValueInt int, err error) {
	if rule.Signature != "in" && rule.Signature != "regexp" {
		ruleValueInt, err = strconv.Atoi(rule.Value)
	}
	return ruleValueInt, err
}
