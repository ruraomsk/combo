package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"ruraomsk/combo/cmb"
	"ruraomsk/combo/driver"
	"ruraomsk/combo/dt"
	"ruraomsk/combo/logger"
	"ruraomsk/combo/proc"
	"ruraomsk/combo/router"
)

// Cmb описание всей системы
var Cmb *cmb.Combo

//Drivers все драйвера системы
var Drivers map[string]*driver.Driver

// DBRouters все роутеры баз данных
var DBRouters map[string]*router.DBRouter

const fsql = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

func paramSQL() string {
	return fmt.Sprintf(fsql, Cmb.DataBase.IP, Cmb.DataBase.Port, Cmb.DataBase.User, Cmb.DataBase.Password, Cmb.DataBase.Base)
}

func main() {

	var err error
	floger, err := logger.LogOpen("data/")
	if err != nil {
		fmt.Printf("error opening file logger: %v \n", err)
	}
	defer floger.Close()

	cmb.Logger = log.New(floger, "combo:", log.LstdFlags)
	Cmb, err = cmb.LoadServer("data/slaves.xml")
	if err != nil {
		cmb.Logger.Fatalf("%v\n", err)
		return
	}

	// fmt.Println(Cmb)

	clearDB := Cmb.DataBase.MakeDB
	fmt.Println(Cmb.Server.Name + " : " + Cmb.Server.Description)
	Drivers = make(map[string]*driver.Driver)
	DBRouters = make(map[string]*router.DBRouter)
	logger.LoadLogger(Cmb.Loggers, paramSQL(), Cmb.Loggers.Name, clearDB)
	for _, device := range Cmb.Devices.DeviceList {
		device.DT, err = dt.LoadTableXML(Cmb.Server.Path + device.Load)
		if err != nil {
			cmb.Logger.Fatalf("Не могли загрузить таблицу " + device.Load)
			return
		}
		//cmb.Logger.Println(device.DT.ToString())
		dr, err := driver.Init(device.Name, device.DT, device)
		if err != nil {
			cmb.Logger.Fatalf("Не могли загрузить Драйвер " + device.Name)
			return
		}
		Drivers[device.Name] = dr
		//cmb.Logger.Println("Load driver:", dr.ToString())
		rdb, err := router.Init(device.Name, dr, paramSQL(), clearDB)
		if err != nil {
			cmb.Logger.Printf("Не могли загрузить Роутер " + device.Name)
		} else {
			DBRouters[device.Name] = rdb
		}
	}
	for name, dr := range Drivers {
		dr.Run()
		cmb.Logger.Println("Запустили Драйвер " + name)
	}
	for name, rdb := range DBRouters {
		rdb.Run()
		cmb.Logger.Println("Запустили Роутер " + name)

	}

	go gui()
	
	go proc.Procedure(Cmb, Drivers)

	for {
		time.Sleep(10 * time.Hour)
		now := time.Now()
		old := now.AddDate(0, 0, -Cmb.DataBase.Period)
		for _, rdb := range DBRouters {
			rdb.KillOld(old)
		}
		//logger.KillOldLog(old)
	}
}
