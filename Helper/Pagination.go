package helper

import (
	"time"
)

type Pagination struct {
	Offset int
	Limit  int
	Period struct {
		From time.Time
		To   time.Time
	}
}

func (obj Pagination) NoLimit() bool {
	if obj.Offset != 0 || obj.Limit != 0 {
		return false
	}
	if obj.Period.From.IsZero() && obj.Period.To.IsZero() {
		return false
	}
	return true
}
