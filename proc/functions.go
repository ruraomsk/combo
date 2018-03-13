package proc

import (
	"math"
	"strconv"
)

type funcexec struct {
	namefun string
	exfun   func([]string) (string, error)
}

var funcs map[string]*funcexec

func start() {
	funcs = make(map[string]*funcexec)
	loadFunc("NOT", bNOT)
	loadFunc("AND", bAND)
	loadFunc("OR", bOR)
	loadFunc("ISUB", iISUB)
	loadFunc("ISUM", iISUM)
	loadFunc("FSUM", fFSUM)
	loadFunc("FSUB", fFSUB)
	loadFunc("FORIF", FORIF)
	loadFunc("FORMAX", FORMAX)
	loadFunc("FloatStr", FloatStr)
	loadFunc("LIVE", iLIVE)
	loadFunc("SET", SET)

}
func loadFunc(name string, fu func([]string) (string, error)) {
	f := new(funcexec)
	f.exfun = fu
	f.namefun = name
	funcs[name] = f
}

// логическое NOT
func bNOT(param []string) (res string, err error) {
	res = "false"
	b, err := strconv.ParseBool(param[0])
	if err != nil {
		return
	}
	res = strconv.FormatBool(!b)
	return
}

// SET простое присвоение входного параметра на выход
func SET(param []string) (res string, err error) {
	res = param[0]
	return
}

//
func bAND(param []string) (res string, err error) {
	res = "false"
	b := true
	bb := true
	for _, par := range param {
		b, err = strconv.ParseBool(par)
		if err != nil {
			return
		}
		bb = bb && b
	}
	res = strconv.FormatBool(bb)
	return
}

func bOR(param []string) (res string, err error) {
	res = "false"
	b := false
	bb := false
	for _, par := range param {
		b, err = strconv.ParseBool(par)
		if err != nil {
			return
		}
		bb = bb || b
	}
	res = strconv.FormatBool(bb)
	return
}

func iISUM(param []string) (res string, err error) {
	res = "0"
	var ib int64
	var ibb int64
	for _, par := range param {
		par = cleanBlank(par)
		ib, err = strconv.ParseInt(par, 10, 64)
		if err != nil {
			return
		}
		ibb += ib
	}
	res = strconv.FormatInt(ibb, 10)
	return
}
func iISUB(param []string) (res string, err error) {
	res = "0"
	err = nil
	var ib int64
	var ibb int64
	f := true
	for _, par := range param {
		par = cleanBlank(par)
		ib, err = strconv.ParseInt(par, 10, 64)
		if err != nil {
			return
		}
		if f {
			ibb += ib
			f = false
		} else {
			ibb -= ib
		}
	}
	res = strconv.FormatInt(ibb, 10)
	return
}
func fFSUM(param []string) (res string, err error) {
	res = "0.0"
	fb := 0.0
	fbb := 0.0
	for _, par := range param {
		par = cleanBlank(par)
		fb, err = strconv.ParseFloat(par, 64)
		if err != nil {
			return
		}
		fbb += fb
	}
	res = strconv.FormatFloat(fbb, 'f', -5, 64)
	return
}
func fFSUB(param []string) (res string, err error) {
	res = "0.0"
	fb := 0.0
	fbb := 0.0
	f := true
	for _, par := range param {
		par = cleanBlank(par)
		fb, err = strconv.ParseFloat(par, 64)
		if err != nil {
			return
		}
		if f {
			fbb += fb
			f = false
		} else {
			fbb -= fb
		}
	}
	res = strconv.FormatFloat(fbb, 'f', -5, 64)
	return
}
func FORIF(param []string) (res string, err error) {
	res = "0.0"
	count := 0
	sum := 0.0
	for i := 0; i < len(param); i += 3 {
		var b bool
		b, err = strconv.ParseBool(cleanBlank(param[i]))
		if err != nil {
			return
		}
		var v1 float64
		var v2 float64

		v1, err = strconv.ParseFloat(cleanBlank(param[i+1]), 64)
		if err != nil {
			return
		}
		v2, err = strconv.ParseFloat(cleanBlank(param[i+2]), 64)
		if err != nil {
			return
		}
		if b {
			sum += v1
		} else {
			sum += v2
		}
		count++
	}
	sum = sum / float64(count)
	res = strconv.FormatFloat(sum, 'f', -5, 64)
	return
}
func FORMAX(param []string) (res string, err error) {
	res = "0.0"
	fb := 0.0
	fbb := 0.0 - math.MaxFloat64
	for _, par := range param {
		par = cleanBlank(par)
		fb, err = strconv.ParseFloat(par, 64)
		if err != nil {
			return
		}
		if fb > fbb {
			fbb = fb
		}
	}
	res = strconv.FormatFloat(fbb, 'f', -5, 64)
	return
}
func FloatStr(param []string) (res string, err error) {
	res = "0 0 0 0"
	var b bool
	var v float64
	b, err = strconv.ParseBool(cleanBlank(param[1]))
	if err != nil {
		return
	}
	v, err = strconv.ParseFloat(cleanBlank(param[0]), 64)
	if err != nil {
		return
	}
	buffer := outFloat(v, b)
	if len(buffer) == 0 {
		return
	}
	res = ""
	for _, val := range buffer {
		res += strconv.FormatInt(int64(val), 10) + " "
	}
	// println("-", res, "-")
	return
}
func iLIVE(param []string) (res string, err error) {
	res = "0"
	var ib int64
	var ibb int64
	for _, par := range param {
		par = cleanBlank(par)
		ib, err = strconv.ParseInt(par, 10, 64)
		if err != nil {
			return
		}
		ibb += ib
	}
	ibb++
	if ibb > 32000 {
		ibb = 1
	}
	res = strconv.FormatInt(ibb, 10)
	return
}
