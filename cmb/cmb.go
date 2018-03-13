package cmb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"ruraomsk/combo/dt"
)

var Logger *log.Logger

// Combo главная структура программы
type Combo struct {
	Server   Server   `xml:"server" jsons:"server"`
	Devices  Devices  `xml:"devices" jsons:"devices"`
	DataBase DataBase `xml:"database" jsons:"database"`
	Loggers  Loggers  `xml:"loggers" jsons:"loggers"`
	Extend   Extend   `xml:"extends" jsons:"extends"`
}

// Devices описывает устройства которые есть у нас
type Devices struct {
	DeviceList []Device `xml:"device" jsons:"device"`
	//	TableFormat
}

// Loggers описывает устройства приема логов
type Loggers struct {
	LoggerList []Log  `xml:"logger" jsons:"logger"`
	Name       string `xml:"name,attr" jsons:"name"`
}

//Extend описывает дополнительные файлы
type Extend struct {
	Step  int    `xml:"step,attr" json:"step"`
	Loads []Load `xml:"load" jsons:"load"`
}

//Load описывает загружаемые файлы
type Load struct {
	File string `xml:"file,attr" json:"file"`
}

// Log описывает устройства которые есть у нас
type Log struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"description"`
	Port        int    `xml:"port,attr" json:"port"`
}

// Server опмсывает параметры самого сервера
type Server struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"description"`
	Path        string `xml:"path,attr" json:"path"`
	Timezone    string `xml:"timezone,attr" json:"timezone"`
}

// DataBase описывает соединение с базой данных
type DataBase struct {
	IP       string `xml:"ip,attr" json:"ip"`
	Port     int    `xml:"port,attr" json:"port"`
	Base     string `xml:"base,attr" json:"base"`
	User     string `xml:"user,attr" json:"user"`
	Password string `xml:"password,attr" json:"password"`
	Period   int    `xml:"period,attr" json:"period"`
	MakeDB   bool   `xml:"make,attr" json:"make"`
}

// Device описывает устройства которые есть у нас
type Device struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"description"`
	IP          string `xml:"ip,attr" json:"ip"`
	Port        int    `xml:"port,attr" json:"port"`
	IP2         string `xml:"ip2,attr,omitempty" json:"ip2"`
	Port2       int    `xml:"port2,attr,omitempty" json:"port2"`
	Load        string `xml:"load,attr" json:"load"`
	Value       string `xml:",chardata" json:"value"`
	Step        int    `xml:"step,attr" json:"step"`
	Restart     int    `xml:"restart,attr" json:"restart"`
	DevType     string `xml:"type,attr" json:"type"`
	TTY         string `xml:"tty,attr"  json:"tty"`
	Baud        int    `xml:"baud,attr" json:"baud"`
	DataBits    int    `xml:"databits,attr" json:"databits"`
	Parity      string `xml:"parity,attr" json:"parity"`
	StopBits    int    `xml:"stopbits,attr" jspn:"stopbits"`
	DT          *dt.DataTable
}

//LoadServer производит загрузку настроечных XML
func LoadServer(namefile string) (*Combo, error) {
	sl := new(Combo)
	buf := bytes.NewBuffer(nil)
	file, err := os.Open(namefile)
	if err != nil {
		return sl, err
	}
	io.Copy(buf, file)
	file.Close()
	xml.Unmarshal(buf.Bytes(), &sl)
	return sl, nil
}

//ToString строка описания настроек
func (c *Combo) ToString() string {
	return fmt.Sprint(c)
}

// SwapBuffer переворачивает порядок байтов для Модбас
func SwapBuffer(buffer []byte) []byte {
	for i := 0; i < len(buffer)-1; i += 2 {
		buffer[i], buffer[i+1] = buffer[i+1], buffer[i]
	}
	return buffer
}

//SwapUint16 переворачивает байты внутри регистров
func SwapUint16(buffer []uint16) []uint16 {
	for i, val := range buffer {
		// val := buffer[i]
		buffer[i] = ((val & 0xff) << 8) | ((val >> 8) & 0xff)
	}
	return buffer
}
