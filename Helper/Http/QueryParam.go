package helperhttp

import (
	"errors"
	"html"
	"strconv"
	"strings"
)

type URLQueryParam map[string][]string

var (
	ParamNotFound       = errors.New("URL param not found")
	ParamMultipleValues = errors.New("URL multiple values not supported")
)

func (obj URLQueryParam) GetString(key string) (string, error) {
	values, found := obj[key]
	if found == false {
		return "", ParamNotFound
	} else if len(values) == 0 {
		return "", nil
	} else if len(values) > 1 {
		return "", ParamMultipleValues
	}

	return html.UnescapeString(values[0]), nil
}

func (obj URLQueryParam) GetBool(key string) (bool, error) {
	value, err := obj.GetString(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(strings.ToLower(value))
}

func (obj URLQueryParam) GetInt(key string) (int, error) {
	value, err := obj.GetString(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}
