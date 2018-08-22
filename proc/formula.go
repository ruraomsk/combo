package proc

import (
	"fmt"
	"strings"
)

func outFloat(in float64, blink bool) (buf []uint16) {
	var sres string
	buf = make([]uint16, 4)
	if in < 0.0 || in > 99999.0 {
		sres = "99999"
	} else {
		if in < 1000.0 {
			sres = fmt.Sprintf("%5.2f", in)
		} else {
			if in < 10000.0 {
				sres = fmt.Sprintf("%5.1f", in)
			} else {
				sres = fmt.Sprintf("%5.f", in)
			}
		}
	}
	rbyte := convertToDisplay(sres)
	buf = upack(rbyte, blink)
	return
}
func convertToDisplay(sres string) (rez []byte) {
	lenb := len(sres)
	if strings.Contains(sres, ".") {
		lenb++
	}
	rez = make([]byte, lenb)
	j := 0
	for i := 0; i < len(sres); i++ {
		b := byte(0xf)
		switch sres[i] {
		case '0':
			b = byte(0)
		case '1':
			b = byte(1)
		case '2':
			b = byte(2)
		case '3':
			b = byte(3)
		case '4':
			b = byte(4)
		case '5':
			b = byte(5)
		case '6':
			b = byte(6)
		case '7':
			b = byte(7)
		case '8':
			b = byte(8)
		case '9':
			b = byte(9)
		case ':':
			b = byte(0x2f)
		case ' ':
			b = byte(0x3f)
		case ',':
			b = byte(0xff)
		case '.':
			b = byte(0xff)
		}
		if b != 0xff {
			rez[j] = b
			j++
		} else {
			rez[j-1] = rez[j-1] | 0x40
		}
	}
	return
}
func upack(rbyte []byte, blink bool) (rez []uint16) {
	lenb := int(len(rbyte) / 2)
	if len(rbyte)%2 != 0 {
		lenb++
	}
	rez = make([]uint16, lenb+1)
	if blink {
		rez[0] = 0
	} else {
		rez[0] = 1
	}
	j := 1
	for i := 0; i < len(rbyte); i++ {
		if i%2 == 0 {
			rez[j] = rez[j] | (uint16(rbyte[i]) << 8)
		} else {
			rez[j] = rez[j] | uint16(rbyte[i])
			j++
		}
	}
	return
}
