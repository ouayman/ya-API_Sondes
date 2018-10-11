package helperhttp

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	".."
)

func ExtractPagination(request *http.Request) (pagination helper.Pagination, err error) {
	urlParams := URLQueryParam(request.URL.Query())

	if pagination.Offset, err = urlParams.GetInt("offset"); err != nil && err != ParamNotFound {
		return
	}
	if pagination.Limit, err = urlParams.GetInt("limit"); err != nil && err != ParamNotFound {
		return
	}
	if pagination.Offset < 0 || pagination.Limit < 0 {
		return pagination, errors.New("limit and offset must be positive values")
	}

	// Method to extract time from key, with support for day or week or month from now, support for a timestamp in second or a time in RFC3339
	var extractTime = func(key string) (timeValue time.Time, err2 error) {
		if value, err2 := urlParams.GetString(key); err2 == nil && err2 != ParamNotFound {
			//value = strings.ToLower(value)
			if value == "day" {
				timeValue = time.Now().AddDate(0, 0, -1)
			} else if value == "week" {
				timeValue = time.Now().AddDate(0, 0, -7)
			} else if value == "month" {
				timeValue = time.Now().AddDate(0, -1, 0)
			} else if timestamp, err2 := strconv.ParseInt(value, 10, 64); err2 == nil {
				timeValue = time.Unix(timestamp, 0)
			} else {
				timeValue, err2 = time.Parse(time.RFC3339, value)
			}
		}
		return
	}

	if pagination.Period.From, err = extractTime("since"); err != nil {
		return
	}
	if pagination.Period.To, err = extractTime("to"); err != nil {
		return
	}

	if pagination.Period.From.IsZero() == false && pagination.Period.To.IsZero() == false && pagination.Period.From.After(pagination.Period.To) {
		return pagination, errors.New("Consistency error between start and end dates")
	}

	return pagination, nil
}

func AddUrlPagination(url string, pagination helper.Pagination) string {
	// Init separator for first url query parameter
	separator := "?"
	// Change separator if url already contains query parameter
	if strings.Contains(url, "?") {
		separator = "&"
	}

	if pagination.Offset != 0 {
		url += separator + fmt.Sprintf("offset=%v", pagination.Offset)
		separator = "&"
	}

	if pagination.Limit != 0 {
		url += separator + fmt.Sprintf("limit=%v", pagination.Limit)
		separator = "&"
	}

	if pagination.Period.From.IsZero() == false {
		url += separator + fmt.Sprintf("since=%v", pagination.Period.From.Format(time.RFC3339))
		separator = "&"
	}

	if pagination.Period.To.IsZero() == false {
		url += separator + fmt.Sprintf("to=%v", pagination.Period.To.Format(time.RFC3339))
		separator = "&"
	}

	return url
}
