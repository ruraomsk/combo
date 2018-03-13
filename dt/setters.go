package dt

import (
	"fmt"
	"time"
)

func (data *Data) MakeField(ff *fieldFormat) {
	val := new(Value)
	val.format = ff.format
	switch ff.format {
	case integertype:
		val.setInt("0")
		if len(ff.defaultValue) > 0 {
			val.setInt(ff.defaultValue)
		}
	case booltype:
		val.setBool("false")
		if len(ff.defaultValue) > 0 {
			val.setBool(ff.defaultValue)
		}
	case floattype:
		val.setFloat("0.0")
		if len(ff.defaultValue) > 0 {
			val.setFloat(ff.defaultValue)
		}
	case longtype:
		val.setLong("0")
		if len(ff.defaultValue) > 0 {
			val.setLong(ff.defaultValue)
		}
	case datetime:
		tt := time.Now()
		s := tt.Format(time.RFC3339Nano)
		if len(ff.defaultValue) > 0 {
			val.setDate(ff.defaultValue)
		} else {
			val.setDate(s)
		}
	case stringtype:
		val.setString("")
		if len(ff.defaultValue) > 0 {
			val.setString(ff.defaultValue)
		}
	default:
		return
	}
	val.format = ff.format
	data.values[ff.name] = val
}
func (val *Value) setValue(v string) {
	switch val.format {
	case integertype:
		val.setInt(v)
	case booltype:
		val.setBool(v)
	case floattype:
		val.setFloat(v)
	case longtype:
		val.setLong(v)
	case datetime:
		val.setDate(v)
	case stringtype:
		val.setString(v)
	default:
		return
	}
}

func (val *Value) setInt(v interface{}) {
	if val.format != integertype {
		return
	}
	if v == nil {
		return
	}
	switch vv := v.(type) {
	case int:
		val.value = fmt.Sprintf("%d", vv)
	case string:
		val.value = vv
	default:
		val.value = "0"
	}
}
func (val *Value) setLong(v interface{}) {
	if val.format != longtype {
		return
	}
	if v == nil {
		return
	}
	switch vv := v.(type) {
	case int:
		val.value = fmt.Sprintf("%d", vv)
	case int64:
		val.value = fmt.Sprintf("%d", vv)
	case string:
		val.value = vv
	default:
		val.value = "0"
	}
}

func (val *Value) setFloat(v interface{}) {
	if val.format != floattype {
		return
	}
	if v == nil {
		return
	}
	switch vv := v.(type) {
	case float32:
		val.value = fmt.Sprintf("%g", vv)
	case float64:
		val.value = fmt.Sprintf("%g", vv)
	case string:
		val.value = vv
	default:
		val.value = "0.0"
	}
}

func (val *Value) setBool(v interface{}) {
	if val.format != booltype {
		return
	}
	if v == nil {
		return
	}
	switch vv := v.(type) {
	case bool:
		val.value = "false"
		if vv {
			val.value = "true"
		}
	case string:
		val.value = vv
	default:
		val.value = "false"
	}
}
func (val *Value) setString(v interface{}) {
	if val.format != stringtype {
		return
	}
	if v == nil {
		return
	}
	switch vv := v.(type) {
	case string:
		val.value = vv
	default:
		val.value = ""
	}
}

func (val *Value) setDate(v interface{}) {
	if val.format != datetime {
		return
	}
	if v == nil {
		return
	}
	switch vv := v.(type) {
	case string:
		val.value = vv
	default:
		//		val.value = time.Format(time.RFC3339Nano, vv)
	}
}
