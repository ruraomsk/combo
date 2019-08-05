package driver

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"rura/combo/cmb"
	"rura/combo/modbus"
)

//Doza структура для мастера TCP modbus
type Doza struct {
	master      *modbus.TCPClientHandler
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
	con         string
}

func doza(d *Driver, con string) (*Doza, error) {
	m := new(Doza)
	m.Step = d.Step
	m.name = d.name
	m.con = con
	m.connected = false
	m.coils = make([]bool, d.lenCoil)
	m.di = make([]bool, d.lenDI)
	m.ir = make([]uint16, d.lenIR)
	m.hr = make([]uint16, d.lenHR)
	m.master = modbus.NewTCPClientHandler(m.con)
	m.master.Timeout = time.Second
	m.master.SlaveId = 1
	m.connected = false
	// m.master.Logger = cmb.Logger
	return m, nil
}
func (m *Doza) start() {
	m.work = true
	go m.run()

}
func (m *Doza) stop() {
	m.master.Close()
	m.work = false
	return
}
func (m *Doza) worked() bool {
	return m.work
}

func (m *Doza) get() (coils []bool, di []bool, ir []uint16, hr []uint16) {
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

func (m *Doza) run() {
	// listen on a port
	fmt.Println("Listering on port " + m.master.Address)
	ln, err := net.Listen("tcp", m.master.Address)
	if err != nil {
		cmb.Logger.Printf("При попытке слушать %s %s", m.con, err.Error())
		return
	}

	for {
		// accept a connection
		conn, err := ln.Accept()
		if err != nil {
			cmb.Logger.Printf("Doza  %s %s", m.con, err.Error())
			continue
		}
		go workData(m, conn)
	}

}
func workData(m *Doza, conn net.Conn) {
	var message string
	defer conn.Close()
	buff := make([]byte, 1024)
	for true {
		n, err := conn.Read(buff)
		if err != nil {
			if strings.Compare(err.Error(), "EOF") != 0 {
				cmb.Logger.Printf("Ошибка приема %s %s", m.name, err.Error())
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}
		message = string(buff[:n])
		// handle the connection
		// err = gob.NewDecoder(conn).Decode(&message)
		//bb := make([]byte, 56)
		// conn.Read(bb)
		// fmt.Println("=", message)
		ss := strings.Split(message, " ")
		if len(ss) != 6 {
			cmb.Logger.Printf("Неверное сообщение %s %s", m.name, message)
			return
		}
		for i := 0; i < len(ss); i++ {
			s := ss[i]
			for {
				j := strings.LastIndex(s, " ")
				if j < 0 {
					break
				}
				s = s[0:j]
			}
			ss[i] = s
		}
		if strings.Compare(ss[0], "[") != 0 {
			cmb.Logger.Printf("Неверное начало %s %s", m.name, message)
			return
		}
		if !strings.Contains(ss[5], "]") {
			cmb.Logger.Printf("Неверное завершение %s %s", m.name, message)
			return
		}
		if len(ss[1]) == 0 {
			return
		}
		v1, err := strconv.ParseFloat(ss[1], 32)
		if err != nil {
			cmb.Logger.Printf("Ошибка значения первого канала %s %s %s", m.name, message, err.Error())
			return
		}
		if len(ss[3]) == 0 {
			return
		}
		v2, err := strconv.ParseFloat(ss[3], 32)
		if err != nil {
			cmb.Logger.Printf("Ошибка значения второго канала %s %s %s", m.name, message, err.Error())
			return
		}
		b1, err := strconv.ParseBool(ss[2])
		if err != nil {
			cmb.Logger.Printf("Ошибка состояния первого канала %s %s %s", m.name, message, err.Error())
			return
		}
		b2, err := strconv.ParseBool(ss[4])
		if err != nil {
			cmb.Logger.Printf("Ошибка состояния второго канала %s %s %s", m.name, message, err.Error())
			return
		}
		// fmt.Printf("v1=%f v2=%f", v1, v2)
		m.mu.Lock()
		m.di[0] = b1
		m.di[1] = b2
		buf := toBuffer(float32(v1))
		for i := 0; i < 2; i++ {
			m.ir[i] = buf[i]
		}
		buf = toBuffer(float32(v2))
		for i := 0; i < 2; i++ {
			m.ir[2+i] = buf[i]
		}
		m.mu.Unlock()
	}
}
func toBuffer(val float32) (buffer []uint16) {
	b := math.Float32bits(val)
	// println(b)
	buffer = make([]uint16, 2)
	buffer[0] = uint16((b >> 16) & 0xffff)
	buffer[1] = uint16(b & 0xffff)
	return
}
func (m *Doza) writeVariable(reg *Register, value string) (err error) {
	err = fmt.Errorf("Для этого устройства %s запрещены операции записи", m.name)
	return
}
func (m *Doza) lock() {
	m.mu.Lock()
}
func (m *Doza) unlock() {
	m.mu.Unlock()
}
