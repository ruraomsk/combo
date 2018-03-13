package driver

import (
	"ruraomsk/combo/cmb"
)

//GetNames возвращает имена переменных с типом переменной
func (d *Driver) GetNames() map[string]int {
	r := make(map[string]int)
	for name, reg := range d.registers {
		r[name] = reg.Regtype()
	}
	return r
}

//GetDescription возвращает имена переменных с описанием
func (d *Driver) GetDescription() map[string]string {
	r := make(map[string]string)
	for name, reg := range d.registers {
		r[name] = reg.description
	}
	return r
}

//GetTypes возвращает имена переменных с типом регистра
func (d *Driver) GetTypes() map[string]int {
	r := make(map[string]int)
	for name, reg := range d.registers {
		r[name] = reg.regtype
	}
	return r
}

//GetValues возвращает все переменные в символьном виде
func (d *Driver) GetValues() map[string]string {
	r := make(map[string]string)
	coils := make([]bool, d.lenCoil)
	di := make([]bool, d.lenDI)
	ir := make([]uint16, d.lenIR)
	hr := make([]uint16, d.lenHR)
	if d.drivertype != 2 {
		coils, di, ir, hr = d.tr.get()
	} else {
		if d.chanel == 0 {
			coils, di, ir, hr = d.tr.get()
		} else {
			coils, di, ir, hr = d.tr2.get()
		}

	}
	//fmt.Printf("%d %d %d %d ", len(coils), len(di), len(ir), len(hr))
	for name, reg := range d.registers {
		val := ""
		switch reg.regtype {
		case 0:
			for i := 0; i < reg.size; i++ {
				v := "false "
				if reg.GetBool(coils, i) {
					v = "true "
				}
				val += v
			}
		case 1:
			for i := 0; i < reg.size; i++ {
				v := "false "
				if reg.GetBool(di, i) {
					v = "true "
				}
				val += v
			}
		case 2:
			val = reg.GetValue(ir)
		case 3:
			val = reg.GetValue(hr)
		}
		r[name] = val
	}
	return r
}

//GetNamedValues возвращает запрошенные переменные в символьном виде
func (d *Driver) GetNamedValues(names map[string]interface{}) map[string]string {
	r := make(map[string]string)
	coils := make([]bool, d.lenCoil)
	di := make([]bool, d.lenDI)
	ir := make([]uint16, d.lenIR)
	hr := make([]uint16, d.lenHR)
	if d.drivertype != 2 {
		coils, di, ir, hr = d.tr.get()
	} else {
		if d.chanel == 0 {
			coils, di, ir, hr = d.tr.get()
		} else {
			coils, di, ir, hr = d.tr2.get()
		}

	}
	//fmt.Printf("%d %d %d %d ", len(coils), len(di), len(ir), len(hr))
	for name, reg := range d.registers {
		if _, ok := names[name]; !ok {
			continue
		}
		val := ""
		switch reg.regtype {
		case 0:
			for i := 0; i < reg.size; i++ {
				v := "false"
				if reg.GetBool(coils, i) {
					v = "true"
				}
				val += v
			}
		case 1:
			for i := 0; i < reg.size; i++ {
				v := "false"
				if reg.GetBool(di, i) {
					v = "true"
				}
				val += v
			}
		case 2:
			val = reg.GetValue(ir)
		case 3:
			val = reg.GetValue(hr)
		}
		r[name] = val
	}
	return r
}

//SetNamedValues записывает на вывод запрошенные переменные в символьном виде
func (d *Driver) SetNamedValues(names map[string]string) {
	// need insert code
	for name, value := range names {
		// fmt.Printf("dev:%s name:%s value=%s \n", d.name, name, value)
		reg, _ := d.registers[name]

		err := d.tr.writeVariable(reg, value)
		if err != nil {
			cmb.Logger.Println(err)
		}
		if d.drivertype == 2 {
			err := d.tr2.writeVariable(reg, value)
			if err != nil {
				cmb.Logger.Println(err)
			}

		}
	}
}
