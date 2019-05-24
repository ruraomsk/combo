package proc

import (
	"strings"
)

func loadInputData() {
	getters := make([]string, 0)
	for name := range variables {
		c := strings.Split(name, separator)
		if len(c) == 2 {
			getters = append(getters, name)
		}
	}
	// передаем запрос получаем ответ
	rep := getNeedValues(getters)
	for name, value := range rep {
		variables[name].value = value
	}
}

func getNeedValues(values []string) map[string]string {
	ret := make(map[string]string)
	for namedev, device := range drivers {
		req := make(map[string]interface{}, 0)
		for _, name := range values {
			c := strings.Split(name, ":")
			if strings.Compare(c[0], namedev) == 0 {
				req[c[1]] = nil
			}
		}
		if len(req) == 0 {
			continue
		}
		res := device.GetNamedValues(req)
		for name, value := range res {
			ret[namedev+":"+name] = value
		}
	}
	return ret
}

func storeOutputDate() {
	setters := make(map[string]string)
	for name, value := range variables {
		if !value.ischange {
			continue
		}
		value.ischange = false
		c := strings.Split(name, separator)
		if len(c) == 2 {
			setters[name] = value.value
		}
	}
	setNeedValues(setters)
}
func setNeedValues(values map[string]string) {
	for namedev, device := range drivers {
		req := make(map[string]string, 0)
		for name, value := range values {
			c := strings.Split(name, ":")
			if strings.Compare(c[0], namedev) == 0 {
				req[c[1]] = value
			}
		}
		if len(req) == 0 {
			continue
		}
		device.SetNamedValues(req)
	}
}
