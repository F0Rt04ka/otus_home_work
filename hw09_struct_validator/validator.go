package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errors []string
	for _, err := range v {
		errors = append(errors, err.Field+": "+err.Err.Error())
	}

	return strings.Join(errors, ", ")
}

func Validate(v interface{}) error {

	reflectType := reflect.TypeOf(v)
	if reflectType.Kind() != reflect.Struct {
		return nil
	}

	valiadtionErrors := ValidationErrors{}

	reflectValue := reflect.ValueOf(v)
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		tagString := field.Tag.Get("validate")
		if tagString == "" {
			continue
		}

		tags := strings.Split(tagString, "|")

		val := reflectValue.Field(i)
		_ = val

		for _, tag := range tags {
			switch field.Type.Kind() {
			case reflect.String, reflect.Int:
				err := validateBaseType(val, tag)
				if err != nil {
					valiadtionErrors = append(valiadtionErrors, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				}
			case reflect.Slice:
				if val.Len() == 0 {
					continue
				}

				switch field.Type.Elem().Kind() {
				case reflect.String, reflect.Int:
					for i := 0; i < val.Len(); i++ {
						err := validateBaseType(val.Index(i), tag)
						if err != nil {
							valiadtionErrors = append(valiadtionErrors, ValidationError{
								Field: field.Name + "[" + strconv.Itoa(i) + "]",
								Err:   err,
							})
						}
					}
				}
			}
		}
	}

	if len(valiadtionErrors) == 0 {
		return nil
	}

	return valiadtionErrors
}

func validateBaseType(value reflect.Value, tag string) error {
	switch value.Kind() {
	case reflect.String:
		return validateString(value.String(), tag)
	case reflect.Int:
		return validateInt(int(value.Int()), tag)
	}

	panic(fmt.Sprintf("unknown type: %T", value))
}

func validateString(value string, tag string) error {
	var res [][]string

	res = regexp.MustCompile(`^len:(\d+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		val, err := strconv.Atoi(res[0][1])
		if err == nil {
			if len(value) != val {
				return errors.New("value length must be equal to " + strconv.Itoa(val))
			}

			return nil
		}
	}

	res = regexp.MustCompile(`^min:(\d+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		val, err := strconv.Atoi(res[0][1])
		if err == nil {
			if len(value) < val {
				return errors.New("value length must be greater than " + strconv.Itoa(val))
			}

			return nil
		}
	}

	res = regexp.MustCompile(`^max:(\d+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		val, err := strconv.Atoi(res[0][1])
		if err == nil {
			if len(value) > val {
				return errors.New("value length must be less than " + strconv.Itoa(val))
			}

			return nil
		}
	}

	res = regexp.MustCompile(`^regexp:(.+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		re, err := regexp.Compile(res[0][1])
		if err == nil {
			if !re.MatchString(value) {
				return errors.New("value must match the regular expression " + res[0][1])
			}

			return nil
		}
	}

	res = regexp.MustCompile(`^in:(.+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		values := strings.Split(res[0][1], ",")
		if !slices.Contains(values, value) {
			return errors.New("value must be one of " + res[0][1])
		}

		return nil
	}

	panic(fmt.Sprintf("unknown tag: %v", tag))
}

func validateInt(value int, tag string) error {
	var res [][]string

	res = regexp.MustCompile(`^min:(\d+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		val, err := strconv.Atoi(res[0][1])
		if err == nil {
			if value < val {
				return errors.New("value must be greater than " + strconv.Itoa(val))
			}

			return nil
		}
	}

	res = regexp.MustCompile(`^max:(\d+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		val, err := strconv.Atoi(res[0][1])
		if err == nil {
			if value > val {
				return errors.New("value must be less than " + strconv.Itoa(val))
			}

			return nil
		}
	}

	res = regexp.MustCompile(`^in:([\d,]+)$`).FindAllStringSubmatch(tag, -1)
	if res != nil {
		values := strings.Split(res[0][1], ",")

		for _, v := range values {
			val, err := strconv.Atoi(v)
			if err != nil {
				panic(fmt.Sprintf("unknown tag: %s", tag))
			}

			if value == val {
				return nil
			}
		}

		return errors.New("value must be one of " + res[0][1])
	}

	panic(fmt.Sprintf("unknown tag: %s", tag))
}
