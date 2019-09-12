package driver

import (
	"bytes"
	"fmt"
	"rura/combo/cmb"
	"rura/combo/dt"
	"time"
)

//DeviceInOut интерфейс всех внешних драйверов ввода вывода
type DeviceInOut interface {
	start()
	get() ([]bool, []bool, []uint16, []uint16)
	writeVariable(*Register, string) error
	worked() bool
	lock()
	unlock()
}

// Driver драйвер одного устройства
type Driver struct {
	name        string
	Description string
	drivertype  int //0 -slave 1- master 2 - dubbus 3 - Дозиметр 4 - rtu
	lenCoil     int
	lenDI       int
	lenIR       int
	lenHR       int
	registers   map[string]*Register
	connect     bool
	work        bool
	chanel      int
	tr          DeviceInOut
	tr2         DeviceInOut
	Step        int
	Restart     int

	IP    string
	Port  int
	IP2   string
	Port2 int
	RTUP  RTUParam
}

// Init подготавливает драйвер к работе
func Init(name string, DT dt.DataTable, dev cmb.Device) (*Driver, error) {
	driver := new(Driver)
	driver.name = name
	driver.Description = dev.Description
	driver.Step = dev.Step
	driver.Restart = dev.Restart
	driver.IP = dev.IP
	driver.IP2 = dev.IP2
	driver.Port = dev.Port
	driver.Port2 = dev.Port2
	driver.RTUP.baud = dev.Baud
	driver.RTUP.tty = dev.TTY
	driver.RTUP.databits = dev.DataBits
	driver.RTUP.parity = dev.Parity
	driver.RTUP.stopbits = dev.StopBits
	// fmt.Println("driver " + driver.name)
	driver.registers = make(map[string]*Register, DT.Len())
	//fmt.Println("")
	// Загружаем регистры
	var err error
	//fmt.Println(DT.Name)
	for i := 0; i < DT.Len(); i++ {
		reg := new(Register)
		rec, _ := DT.ReadRecod(i)
		reg.name, err = rec.GetString("name")
		if err != nil {
			cmb.Logger.Println("name not found")
		}
		if reg.description, err = rec.GetString("description"); err != nil {
			cmb.Logger.Println("description not found")
		}
		if reg.regtype, err = rec.GetInt("type"); err != nil {
			cmb.Logger.Println("type not found")
		}
		if reg.format, err = rec.GetInt("format"); err != nil {
			cmb.Logger.Println("format not found")
		}
		if reg.address, err = rec.GetInt("address"); err != nil {
			cmb.Logger.Println("addres not found")
		}
		if reg.size, err = rec.GetInt("size"); err != nil {
			cmb.Logger.Println("size not found")
		}
		if reg.unitID, err = rec.GetInt("unitId"); err != nil {
			cmb.Logger.Println("unitId not found")
		}
		reg.count(driver)
		driver.registers[reg.name] = reg
	}
	//fmt.Println("count",driver.lenCoil, driver.lenDI, driver.lenIR, driver.lenHR)
	switch dev.DevType {
	case "rtu":
		driver.drivertype = 4
		d, err := rtu(driver)
		if err != nil {
			cmb.Logger.Println(err.Error())
			return driver, err
		}
		driver.tr = d
	case "mono":
		driver.drivertype = 1
		con := fmt.Sprintf("%s:%d", dev.IP, dev.Port)
		d, err := master(driver, con)
		if err != nil {
			cmb.Logger.Println(err.Error())
			return driver, err
		}
		driver.tr = d
	case "dub":
		driver.drivertype = 2
		con := fmt.Sprintf("%s:%d", dev.IP, dev.Port)
		d, err := master(driver, con)
		if err != nil {
			cmb.Logger.Println(err.Error())
			return driver, err
		}
		driver.tr = d
		con = fmt.Sprintf("%s:%d", dev.IP2, dev.Port2)
		d, err = master(driver, con)
		if err != nil {
			cmb.Logger.Println(err.Error())
			return driver, err
		}
		driver.tr2 = d
	case "slave":
		driver.drivertype = 0
		con := fmt.Sprintf(":%d", dev.Port)
		s, err := slave(driver, con)
		if err != nil {
			cmb.Logger.Println(err.Error())
			return driver, err
		}
		driver.tr = s
	case "doza":
		driver.drivertype = 3
		con := fmt.Sprintf(":%d", dev.Port)
		s, err := doza(driver, con)
		if err != nil {
			cmb.Logger.Println(err.Error())
			return driver, err
		}
		driver.tr = s
	}
	driver.work = true
	return driver, nil
}

func (d *Driver) loop() {

	step := 60 * time.Second
	if d.Restart != 0 {
		step = time.Duration(d.Restart) * time.Second
	}

	for {
		//start := time.Now()
		if !d.work {
			return
		}
		switch d.drivertype {
		case 2:
			if d.tr.worked() {
				d.chanel = 0
			} else {
				d.tr.lock()
				con := fmt.Sprintf("%s:%d", d.IP, d.Port)
				dv, err := master(d, con)
				d.tr.unlock()
				if err != nil {
					cmb.Logger.Println(err.Error())
				} else {
					d.tr = dv
					dv.start()
				}
			}
			if d.tr2.worked() {
				d.chanel = 1
			} else {
				d.tr2.lock()
				con := fmt.Sprintf("%s:%d", d.IP2, d.Port2)
				dv, err := master(d, con)
				d.tr2.unlock()
				if err != nil {
					cmb.Logger.Println(err.Error())
				} else {
					d.tr2 = dv
					dv.start()
				}
			}
		case 1:
			if !d.tr.worked() {
				d.tr.lock()
				con := fmt.Sprintf("%s:%d", d.IP, d.Port)
				dv, err := master(d, con)
				d.tr.unlock()
				if err != nil {
					cmb.Logger.Println(err.Error())
				} else {
					d.tr = dv
					dv.start()
				}
			}
		case 4:
			if !d.tr.worked() {
				d.tr.lock()
				dv, err := rtu(d)
				d.tr.unlock()
				if err != nil {
					cmb.Logger.Println(err.Error())
				} else {
					d.tr = dv
					dv.start()
				}
			}

		}
		//stop := time.Now()
		//elapsed := stop.Sub(start)
		// fmt.Println("driver "+d.name)
		time.Sleep(step)

	}
}

// Run запускает драйвер
func (d *Driver) Run() {
	d.tr.start()
	if d.drivertype == 2 {
		d.tr2.start()
	}
	go d.loop()
	//fmt.Println("Запустили Драйвер " + d.name)
}

// Stop останавливает драйвер
func (d *Driver) Stop() {
	d.work = false
}

// Status возвращает статус драйвера
func (d *Driver) Status() bool {
	return d.work
}

//ToString выводит описание драйвера
func (d *Driver) ToString() string {
	s := bytes.NewBufferString("driver:" + d.name + fmt.Sprintf(" rgs=%d %d", len(d.registers), d.drivertype) + "\n")
	ss := fmt.Sprintf("%d %d %d %d \n", d.lenCoil, d.lenDI, d.lenIR, d.lenHR)
	s.WriteString(ss)
	for _, reg := range d.registers {
		s.WriteString(reg.ToString() + " \n")
	}
	return s.String()
}
