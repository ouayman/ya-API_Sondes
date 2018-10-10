package helper

import (
	"strconv"
	"time"
)

type TimestampMs struct {
	time.Time
}

func TimeFromTimestampMs(timestampMs int64) time.Time {
	return time.Unix(timestampMs/1000, (timestampMs%1000)*1e6)
}

func TimestampMsFromTime(value time.Time) int64 {
	timeSec := value.UTC().Unix()
	return timeSec * 1000
}

// TODO Case null à gérer !!!!

func (obj *TimestampMs) MarshalJSON() ([]byte, error) {
	stamp := strconv.FormatInt(TimestampMsFromTime(obj.Time), 10)
	return []byte(stamp), nil
}

func (obj *TimestampMs) UnmarshalJSON(value []byte) error {
	timestampMs, err := strconv.ParseInt(string(value), 10, 64)
	if err == nil {
		obj.Time = TimeFromTimestampMs(timestampMs)
	}
	return err
}
