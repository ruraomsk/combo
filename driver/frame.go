package driver

import (
	"encoding/binary"

	"github.com/tbrandon/mbserver"
)

func registerAddressAndNumber(frame mbserver.Framer) (register int, numRegs int, endRegister int) {
	data := frame.GetData()
	register = int(binary.BigEndian.Uint16(data[0:2]))
	numRegs = int(binary.BigEndian.Uint16(data[2:4]))
	endRegister = register + numRegs
	return register, numRegs, endRegister
}

func registerAddressAndValue(frame mbserver.Framer) (int, uint16) {
	data := frame.GetData()
	register := int(binary.BigEndian.Uint16(data[0:2]))
	value := binary.BigEndian.Uint16(data[2:4])
	return register, value
}
func bitAtPosition(value uint8, pos uint) uint8 {
	return (value >> pos) & 0x01
}

// Функции расширения запросов
//
func writeDiscretRegister(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, value := registerAddressAndValue(frame)
	if value != 0 {
		value = 1
	}
	s.DiscreteInputs[register] = byte(value)
	return frame.GetData()[0:4], &mbserver.Success
}
func writeDiscretRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, numRegs, endRegister := registerAddressAndNumber(frame)
	valueBytes := frame.GetData()[5:]
	if endRegister > 65535 {
		return []byte{}, &mbserver.IllegalDataAddress
	}
	bitCount := 0
	for i, value := range valueBytes {
		for bitPos := uint(0); bitPos < 8; bitPos++ {
			s.DiscreteInputs[register+(i*8)+int(bitPos)] = bitAtPosition(value, bitPos)
			bitCount++
			if bitCount >= numRegs {
				break
			}
		}
		if bitCount >= numRegs {
			break
		}
	}

	return frame.GetData()[0:4], &mbserver.Success
}
func writeSingleRegister(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, value := registerAddressAndValue(frame)
	s.InputRegisters[register] = value
	return frame.GetData()[0:4], &mbserver.Success

}
func writeRegisters(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
	register, numRegs, _ := registerAddressAndNumber(frame)
	valueBytes := frame.GetData()[5:]
	var exception *mbserver.Exception
	var data []byte

	if len(valueBytes)/2 != numRegs {
		exception = &mbserver.IllegalDataAddress
	}

	// Copy data to memroy
	values := mbserver.BytesToUint16(valueBytes)
	valuesUpdated := copy(s.InputRegisters[register:], values)
	if valuesUpdated == numRegs {
		exception = &mbserver.Success
		data = frame.GetData()[0:4]
	} else {
		exception = &mbserver.IllegalDataAddress
	}

	return data, exception

}
