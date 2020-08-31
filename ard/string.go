package ard

import (
	"fmt"
	"strconv"
	"time"
)

func ValueToString(data Value) string {
	switch data_ := data.(type) {
	case bool:
		return strconv.FormatBool(data_)
	case int64:
		return strconv.FormatInt(data_, 10)
	case int32:
		return strconv.FormatInt(int64(data_), 10)
	case int8:
		return strconv.FormatInt(int64(data_), 10)
	case int:
		return strconv.FormatInt(int64(data_), 10)
	case uint64:
		return strconv.FormatUint(data_, 10)
	case uint32:
		return strconv.FormatUint(uint64(data_), 10)
	case uint8:
		return strconv.FormatUint(uint64(data_), 10)
	case uint:
		return strconv.FormatUint(uint64(data_), 10)
	case float64:
		return strconv.FormatFloat(data_, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(data_), 'g', -1, 32)
	case time.Time:
		return data_.String()
	default:
		return fmt.Sprintf("%s", data_)
	}
}
