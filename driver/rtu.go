package driver

import (
	"rura/combo/cmb"
	"sync"
	"time"

	"rura/combo/modbus"
)

//RTU структура для мастера RTU modbus
type RTU struct {
	master      *modbus.RTUClientHandler
	client      modbus.Client
	description string
	name        string
	coils       []bool
	di          []bool
	ir          []uint16
	hr          []uint16
	connected   bool
	work        bool
	mu          sync.Mutex
	Step        int
	con         RTUParam
}

//RTUParam структура для передачи параметров для соединения
type RTUParam struct {
	tty      string
	baud     int
	databits int
	parity   string
	stopbits int
}

func rtu(d *Driver) (*RTU, error) {
	m := new(RTU)
	m.Step = d.Step
	m.name = d.name
	m.con = d.RTUP
	m.connected = false
	m.coils = make([]bool, d.lenCoil)
	m.di = make([]bool, d.lenDI)
	m.ir = make([]uint16, d.lenIR)
	m.hr = make([]uint16, d.lenHR)
	m.master = modbus.NewRTUClientHandler(m.con.tty)
	m.master.BaudRate = m.con.baud
	m.master.DataBits = m.con.databits
	m.master.Parity = m.con.parity
	m.master.StopBits = m.con.stopbits
	m.master.Timeout = 5 * time.Second
	m.master.SlaveId = 1
	m.connected = false
	// m.master.Logger = cmb.Logger
	return m, nil
}
func (m *RTU) worked() bool {
	return m.work
}

func (m *RTU) start() {
	err := m.master.Connect()
	if err != nil {
		cmb.Logger.Println(m.name + " " + err.Error())
		return
	}
	m.client = modbus.NewClient(m.master)
	m.work = true
	go m.run()
	//fmt.Println("Master :" + m.con)

}
func (m *RTU) stop() {
	m.master.Close()
	m.work = false
	return
}

func (m *RTU) readAllCoils() {
	var coils = uint16(len(m.coils))
	if coils == 0 {
		return
	}

	buff, err := m.client.ReadCoils(0, coils)
	if err != nil {
		cmb.Logger.Println(m.name + err.Error())
		m.stop()
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	//cmb.SwapBuffer(buff)
	for i := range m.coils {
		m.coils[i] = getbool(buff, uint16(i))
	}
}

func (m *RTU) readAllDI() {
	var di = uint16(len(m.di))
	if di == 0 {
		return
	}
	buff, err := m.client.ReadDiscreteInputs(0, di)
	if err != nil {
		cmb.Logger.Println(m.name + err.Error())
		m.stop()
		return
	}
	//	cmb.SwapBuffer(buff)

	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.di {
		m.di[i] = getbool(buff, uint16(i))
	}
}

func (m *RTU) readAllIR() {
	if len(m.ir) == 0 {
		return
	}
	ref := uint16(0)
	for count := len(m.ir); count > 0; count -= 125 {
		len := count
		if count > 125 {
			len = 125
		}
		buff, err := m.client.ReadInputRegisters(ref, uint16(len))
		if err != nil {
			cmb.Logger.Println(m.name + err.Error())
			m.stop()
			return
		}
		pos := ref
		left := 0
		cmb.SwapBuffer(buff)
		m.mu.Lock()
		defer m.mu.Unlock()
		for i := 0; i < len; i++ {
			m.ir[pos] = (uint16(buff[left+1]) << 8) | uint16(buff[left])
			pos++
			left += 2
		}
		ref += 125
	}
}

func (m *RTU) readAllHR() {
	if len(m.hr) == 0 {
		return
	}
	ref := uint16(0)
	for count := len(m.hr); count > 0; count -= 125 {
		len := count
		if count > 125 {
			len = 125
		}
		buff, err := m.client.ReadHoldingRegisters(ref, uint16(len))
		if err != nil {
			cmb.Logger.Println(m.name + err.Error())
			m.stop()
			return
		}
		pos := ref
		left := 0
		cmb.SwapBuffer(buff)
		m.mu.Lock()
		defer m.mu.Unlock()
		for i := 0; i < len; i++ {
			m.hr[pos] = (uint16(buff[left+1]) << 8) | uint16(buff[left])
			pos++
			left += 2
		}
		ref += 125
	}
}
func (m *RTU) get() (coils []bool, di []bool, ir []uint16, hr []uint16) {
	m.mu.Lock()
	defer m.mu.Unlock()
	coils = make([]bool, len(m.coils))
	for i, val := range m.coils {
		coils[i] = val
	}
	di = make([]bool, len(m.di))
	for i, val := range m.di {
		di[i] = val
	}
	ir = make([]uint16, len(m.ir))
	for i, val := range m.ir {
		ir[i] = val
	}
	hr = make([]uint16, len(m.hr))
	for i, val := range m.hr {
		hr[i] = val
	}
	return
}

func (m *RTU) run() {
	step := time.Duration(m.Step) * time.Millisecond
	for {

		//start := time.Now()
		if !m.work {
			return
		}

		// m.readAllCoils()
		// m.readAllDI()
		// m.readAllIR()
		// m.readAllHR()

		if !m.work {
			return
		}

		//stop := time.Now()
		//elapsed := stop.Sub(start)
		// fmt.Println("master " + m.name)
		time.Sleep(step)

	}
}
func (m *RTU) writeVariable(reg *Register, value string) (err error) {
	if !m.work {
		return
	}
	buffer, err := reg.SetValue(value)
	if err != nil {
		return
	}
	buf := make([]byte, len(buffer)*2)
	pos := 0
	// print(reg.name + " [")
	for i := 0; i < len(buffer); i++ {
		// print(buffer[i], "->")
		buf[pos+1] = byte(buffer[i] & 0xff)
		buf[pos+0] = byte((buffer[i] >> 8) & 0xff)
		// fmt.Print(buf[pos], buf[pos+1], "-")
		pos += 2
	}
	// println("]")
	m.master.SlaveId = byte(reg.unitID)
	// fmt.Println("Id=", m.master.SlaveId)
	// fmt.Println("len=", len(buffer))
	if len(buffer) == 0 {
		return
	}
	if len(buffer) == 1 {
		// buffer[0] = ((buffer[0] & 0xff) << 8) | ((buffer[0] >> 8) & 0xff)
		switch reg.regtype {
		case 0:
			// println(reg.name, buffer[0])
			_, err = m.client.WriteSingleCoil(uint16(reg.address&0xffff), buffer[0])
		case 1:
			_, err = m.client.WriteSingleDiscreteInput(uint16(reg.address&0xffff), buffer[0])
		case 2:
			_, err = m.client.WriteSingleInputRegister(uint16(reg.address&0xffff), buffer[0])
		case 3:
			_, err = m.client.WriteSingleRegister(uint16(reg.address&0xffff), buffer[0])
		}
		return
	}
	switch reg.regtype {
	case 0:
		_, err = m.client.WriteMultipleCoils(uint16(reg.address&0xffff), uint16(len(buffer)), buf)
	case 1:
		_, err = m.client.WriteMultipleDiscreteInput(uint16(reg.address&0xffff), uint16(len(buffer)), buf)
	case 2:
		_, err = m.client.WriteMultipleInputRegisters(uint16(reg.address&0xffff), uint16(len(buffer)), buf)
	case 3:
		// fmt.Println(reg.ToString())
		_, err = m.client.WriteMultipleRegisters(uint16(reg.address&0xffff), uint16(len(buffer)), buf)
	}
	return
}
func (m *RTU) lock() {
	m.mu.Lock()
}
func (m *RTU) unlock() {
	m.mu.Unlock()
}
