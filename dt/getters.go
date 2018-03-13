package dt

import (
	"strconv"
	"time"
)

func (v *Value) getInt() int {
	if v.format == integertype {
		res, _ := strconv.Atoi(v.value)
		return res
	}
	return 0
}
func (v *Value) getLong() int64 {
	if v.format == longtype {
		res, _ := strconv.Atoi(v.value)
		return int64(res)
	}
	return 0
}

func (v *Value) getFloat() float64 {
	if v.format == floattype {
		res, _ := strconv.ParseFloat(v.value, 64)
		return float64(res)
	}
	return 0.0
}
func (v *Value) getBool() bool {
	if v.format == booltype {
		res, _ := strconv.ParseBool(v.value)
		return res
	}
	return false
}
func (v *Value) getString() string {
	if v.format == stringtype {

		return v.value
	}
	return ""
}
func (v *Value) getDate() time.Time {
	if v.format == datetime {
		res, _ := time.Parse(time.RFC3339Nano, v.value)
		return res
	}
	return time.Now()

}
