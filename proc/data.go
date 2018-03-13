package proc

import (
	"fmt"
	"strings"
)

type anyValue struct {
	name        string
	description string
	ischange    bool
	value       string
}

func (v anyValue) ToString() string {
	return fmt.Sprintf("[ %s=%s %t// %s ]\n", v.name, v.value, v.ischange, v.description)
}

func setValue(name string, value string) error {
	if _, ok := variables[name]; !ok {
		return fmt.Errorf("Нет переменной %s", name)
	}
	if strings.Compare(variables[name].value, value) != 0 {
		variables[name].value = value
		variables[name].ischange = true
	}
	return nil
}
func getValue(name string) (string, error) {
	if _, ok := variables[name]; !ok {

		return "", fmt.Errorf("Нет переменной %s", name)
	}
	return variables[name].value, nil
}
func oneData(name string, description string, value string) {
	oval := new(anyValue)
	oval.name = name
	oval.value = value
	oval.description = description
	variables[name] = oval
}
func allDataToString() string {
	s := "Все переменные модели \n"
	for _, value := range variables {
		s += value.ToString()
	}
	return s
}
func cleanBlank(value string) (res string) {
	res = value
	for {
		j := strings.LastIndex(res, " ")
		if j < 0 {
			return
		}
		res = res[0:j]
	}
}
