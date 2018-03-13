package proc

import (
	"fmt"
	"strings"
)

type anyCode struct {
	funcName  string
	parametrs []string
	rezult    string
}

func (c anyCode) ToString() string {
	s := fmt.Sprintf("%s = %s (", c.rezult, c.funcName)
	p := " "
	for _, name := range c.parametrs {
		s += fmt.Sprintf("%s %s", p, name)
		p = ","
	}
	s += ")"
	return s
}
func allCodeToString() string {
	s := "Все формулы модели \n"
	for i, line := range lineCodes {
		s += fmt.Sprintf("%4d ", i) + line.ToString() + "\n"
	}
	return s
}
func oneCode(namefunc string, parametrs []string, rezult string) {
	code := new(anyCode)
	code.funcName = namefunc
	code.parametrs = parametrs
	code.rezult = rezult
	lineCodes = append(lineCodes, code)
	// Создаем если нужно данные
	for _, name := range parametrs {
		makeData(name)
	}
	makeData(rezult)
}
func makeData(name string) {
	c := strings.Split(name, separator)
	if len(c) == 2 {
		if _, ok := variables[name]; ok {
			return
		}
		s := fmt.Sprintf("Переменая %s с устройтва %s ", c[1], c[0])
		oneData(name, s, " ")
	}

}
func (c anyCode) exec() (res string, err error) {
	res = "0"
	param := make([]string, len(c.parametrs))
	for i, par := range c.parametrs {
		param[i], err = getValue(par)
		if err != nil {
			return
		}
	}
	ex, ok := funcs[c.funcName]
	if !ok {
		return "0", fmt.Errorf("Отстутсвут определение функции %s", c.funcName)
	}
	res, err = ex.exfun(param)
	if err != nil {
		return "0", err
	}
	err = setValue(c.rezult, res)
	return
}
