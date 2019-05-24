package driver

import (
	"github.com/tbrandon/mbserver"
)

//Slave Структура
type Slave struct {
	server *mbserver.Server
	name   string
	ret    chan int
	con    string
	coils  []bool
	di     []bool
	ir     []uint16
	hr     []uint16
	work   bool
}

func slave(d *Driver, con string) (*Slave, error) {
	s := new(Slave)
	s.con = con
	s.name = d.name
	s.coils = make([]bool, d.lenCoil)
	s.di = make([]bool, d.lenDI)
	s.ir = make([]uint16, d.lenIR)
	s.hr = make([]uint16, d.lenHR)
	s.work = true
	s.server = mbserver.NewServer()

	s.server.RegisterFunctionHandler(65, writeDiscretRegister)
	s.server.RegisterFunctionHandler(66, writeSingleRegister)
	s.server.RegisterFunctionHandler(67, writeDiscretRegisters)
	s.server.RegisterFunctionHandler(68, writeRegisters)

	// err := s.server.ListenTCP(con)
	// if err != nil {
	// 	cmb.Logger.Println(err.Error())
	// 	return s, err
	// }
	return s, nil
}

func (s *Slave) start() {
	_ = s.server.ListenTCP(s.con)
	s.work = true
	// fmt.Println("Slave:" + s.con)

}
func (s *Slave) stop() {
	s.server.Close()
}
func (s *Slave) worked() bool {
	return s.work
}

func (s *Slave) readAllCoils() {
	for i := range s.coils {
		if s.server.Coils[i] > 0 {
			s.coils[i] = true
		} else {
			s.coils[i] = false
		}
	}
}

func (s *Slave) readAllDI() {
	for i := range s.di {
		if s.server.DiscreteInputs[i] > 0 {
			s.di[i] = true
		} else {
			s.di[i] = false
		}
	}
}

func (s *Slave) readAllIR() {
	for i := range s.ir {
		s.ir[i] = s.server.InputRegisters[i]
	}
	// s.ir = cmb.SwapUint16(s.ir)
}

func (s *Slave) readAllHR() {
	for i := range s.hr {
		s.hr[i] = s.server.HoldingRegisters[i]
	}
	// s.hr = cmb.SwapUint16(s.hr)
}

func (s *Slave) get() (coils []bool, di []bool, ir []uint16, hr []uint16) {
	s.readAllCoils()
	s.readAllDI()
	s.readAllIR()
	s.readAllHR()
	coils = make([]bool, len(s.coils))
	for i, val := range s.coils {
		coils[i] = val
	}
	di = make([]bool, len(s.di))
	for i, val := range s.di {
		di[i] = val
	}
	ir = make([]uint16, len(s.ir))
	for i, val := range s.ir {
		ir[i] = val
	}
	hr = make([]uint16, len(s.hr))
	for i, val := range s.hr {
		hr[i] = val
	}
	return
}
func (s *Slave) writeVariable(reg *Register, value string) (err error) {
	// println(reg.ToString(), value)
	buffer, err := reg.SetValue(value)
	if err != nil {
		return
	}
	switch reg.regtype {
	case 0:
		for i := 0; i < reg.size; i++ {
			s.server.Coils[reg.address+i] = byte(buffer[i])
		}
	case 1:
		for i := 0; i < reg.size; i++ {
			s.server.DiscreteInputs[reg.address+i] = byte(buffer[i])
		}
	case 2:
		for i := 0; i < len(buffer); i++ {
			s.server.InputRegisters[reg.address+i] = buffer[i]
		}
	case 3:
		for i := 0; i < len(buffer); i++ {
			s.server.HoldingRegisters[reg.address+i] = buffer[i]
		}
	}
	return
}
func (s *Slave) lock() {
}
func (s *Slave) unlock() {
}
