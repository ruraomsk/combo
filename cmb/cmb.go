package cmb

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"rura/combo/dt"
	"rura/teprol/logger"
)

//Logger главная переменная логов
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
	Name          string `xml:"name,attr" json:"name"`
	Description   string `xml:"description,attr" json:"description"`
	Path          string `xml:"path,attr" json:"path"`
	Timezone      string `xml:"timezone,attr" json:"timezone"`
	Project       string `xml:"project,attr" json:"project"`
	MasterStep    int    `xml:"masterstep,attr" json:"masterstep"`
	MasterRestart int    `xml:"masterrestart,attr" json:"masterrestart"`
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
	DT          dt.DataTable
}

//Project описание одного проекта системы
type Project struct {
	Subs []Sub `xml:"subs" json:"subs"`
	// General General `xml:"general" json:"general"`
}

//General описание заголока проекта
type General struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"desription"`
}

//Sub описание одной подсистемы
type Sub struct {
	Name        string `xml:"name,attr" json:"name"`
	Path        string `xml:"path,attr" json:"path"`
	File        string `xml:"file,attr" json:"file"`
	Description string `xml:"description,attr" json:"description"`
	Main        string `xml:"main,attr" json:"main"`
	Second      string `xml:"second,attr" json:"second"`
}

//Subsystem полное описание одной подсистемы
type Subsystem struct {
	Modbuses []Modbus `xml:"modbus" json:"modbus"`
}

//Modbus Описание модбасов в подсистеме
type Modbus struct {
	Name        string `xml:"name,attr" json:"name"`
	Description string `xml:"description,attr" json:"description"`
	Type        string `xml:"type,attr" json:"type"`
	Port        int    `xml:"port,attr" json:"port"`
	XML         string `xml:"xml,attr" json:"xml"`
}

//LoadServer производит загрузку настроечных XML
func LoadServer(namefile string) (*Combo, error) {
	sl := new(Combo)
	buf, err := ioutil.ReadFile(namefile)
	if err != nil {
		logger.Error.Println(err.Error())
		return sl, err
	}
	err = xml.Unmarshal(buf, &sl)
	if err != nil {
		logger.Error.Println(err.Error())
		return sl, err
	}
	for i, dev := range sl.Devices.DeviceList {
		dev.Load = sl.Server.Path + dev.Load
		sl.Devices.DeviceList[i] = dev
	}
	if len(sl.Server.Project) != 0 {
		pr := new(Project)
		namefile := sl.Server.Project + "/main.xml"
		buf, err := ioutil.ReadFile(namefile)
		if err != nil {
			logger.Error.Println(err.Error())
			return sl, err
		}
		err = xml.Unmarshal(buf, &pr)
		if err != nil {
			logger.Error.Println(err.Error())
			return sl, err
		}
		for _, sub := range pr.Subs {
			newpath := sl.Server.Project + "/" + sub.Path
			newfilexml := newpath + "/" + sub.File + ".xml"
			ssub := new(Subsystem)
			buf, err := ioutil.ReadFile(newfilexml)
			if err != nil {
				logger.Error.Println(err.Error())
				return sl, err
			}
			err = xml.Unmarshal(buf, &ssub)
			if err != nil {
				logger.Error.Println(err.Error())
				return sl, err
			}
			for _, mod := range ssub.Modbuses {
				if mod.Type != "slave" {
					continue
				}
				dev := new(Device)
				dev.Name = mod.Name
				dev.Description = sub.Name + ":" + mod.Description
				dev.DevType = "dub"
				dev.IP = sub.Main
				dev.IP2 = sub.Second
				dev.Port = mod.Port
				dev.Port2 = mod.Port
				dev.Step = sl.Server.MasterStep
				dev.Restart = sl.Server.MasterRestart
				dev.Load = newpath + "/" + mod.XML + ".xml"
				sl.Devices.DeviceList = append(sl.Devices.DeviceList, *dev)

			}
		}

	}
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
