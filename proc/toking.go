package proc

import (
	"strings"
	"unicode"
)

var f = func(c rune) bool {
	if c == ':' {
		return false
	}
	if c == '.' {
		return false
	}
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

func toking(l string) {
	if strings.HasPrefix(l, "//") {
		return
	}
	//	l = strings.ToUpper(l)
	s := strings.FieldsFunc(l, f)
	if len(s) == 0 {
		return
	}
	if strings.Compare(s[0], "set") == 0 {
		// name:=value
		if len(s) < 4 {
			s = append(s, "Нет описания ")
		}
		comment := ""
		for i := 3; i < len(s); i++ {
			comment += " " + s[i]
		}
		oneData(s[1], comment, s[2])
		return
	}
	if strings.Compare(s[0], "let") == 0 {
		// name:=value
		pars := make([]string, 0)
		if len(s) == 3 {
			pars = append(pars, s[2])
			oneCode("SET", pars, s[1])
			return
		}
		for i := 3; i < len(s); i++ {
			pars = append(pars, s[i])
		}
		oneCode(s[2], pars, s[1])
		return
	}

}
