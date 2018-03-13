package proc

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"ruraomsk/combo/cmb"
	"ruraomsk/combo/driver"
)

// Основной цикл выполнения  процедуры управления
// Состоит из трех фаз
// 1- прием данных из подсистемы ввода
// 2- проведения расчетов как правило это функции
// 3- передача всех измененных значений в подситему вводв-вывода
// Относительно имен переменных оно состоит зи двух компонент
// локация:собственно_имя
// где локация - это имя устройства если отсутствует то это внутренняя переменная
const separator = ":"

var step time.Duration
var drivers map[string]*driver.Driver

var variables map[string]*anyValue
var lineCodes []*anyCode

func initDatas() {
	variables = make(map[string]*anyValue)
}

func initCodes() {
	lineCodes = make([]*anyCode, 0)
}
func Procedure(Cmb *cmb.Combo, drvs map[string]*driver.Driver) {
	step = time.Duration(Cmb.Extend.Step) * time.Millisecond
	drivers = drvs
	initCodes()
	initDatas()
	for _, load := range Cmb.Extend.Loads {
		file, err := os.Open(Cmb.Server.Path + load.File)
		if err != nil {
			break
		}
		f := bufio.NewReader(file)
		// fmt.Println(load)
		for {
			line, err := f.ReadString(byte(10))
			// fmt.Println(line)
			toking(line)
			if err != nil {
				break
			}
		}
		file.Close()
	}
	go mainCycle()
}
func mainCycle() {
	start()
	for {
		loadInputData()
		for _, line := range lineCodes {
			_, err := line.exec()
			if err != nil {
				fmt.Println(err)
			}
		}
		// fmt.Println(allDataToString())
		storeOutputDate()
		time.Sleep(step)
	}
}
