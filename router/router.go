package router

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"

	"rura/combo/cmb"
	"rura/combo/driver"
)

// DBRouter роутер одной таблицы
type DBRouter struct {
	name        string
	description string
	driver      *driver.Driver
	db          *sql.DB
	work        bool
	step        int
	mu          sync.Mutex
}

// Init инициализирует роутер базы данных для данного устройства
// если необходимо создает таблицу в базе для записи
func Init(name string, drv *driver.Driver, paramSQL string, init bool) (*DBRouter, error) {
	r := new(DBRouter)
	r.name = name
	r.description = drv.Description
	r.driver = drv
	r.step = drv.Step
	var err error
	if r.db, err = sql.Open("postgres", paramSQL); err != nil {
		cmb.Logger.Printf("router %s DB %v\n", r.name, err)
		return r, err
	}
	if err = r.db.Ping(); err != nil {
		cmb.Logger.Printf("router %s DB %v\n", r.name, err)
		return r, err
	}
	if init {
		//make db table
		names := r.driver.GetNames()
		ss := "drop table if exists " + r.name + ";"
		cmb.Logger.Println(ss)
		row, err := r.db.Query(ss)
		if err != nil {
			cmb.Logger.Printf("Error drop database %s", err.Error())
			return r, err
		}
		row.Close()
		r.db.Close()
		if r.db, err = sql.Open("postgres", paramSQL); err != nil {
			cmb.Logger.Printf("router %s DB %v\n", r.name, err)
			return r, err
		}
		s := bytes.NewBufferString("create table if not exists " + r.name + " ( tm timestamp primary key not null ")

		//fmt.Printf("router %s len %d", r.name, len(names))
		for name, types := range names {
			//add define column
			s.WriteString(",")
			s.WriteString(name)
			switch types {
			case 0:
				s.WriteString(" boolean")
			case 1:
				s.WriteString(" integer")
			case 2:
				s.WriteString(" real")
				// fmt.Printf("%s real \n", name)
			case 3:
				s.WriteString(" bigint")
			}
		}
		s.WriteString(");")
		//println(s.String())
		row, err = r.db.Query(s.String())
		if err != nil {
			cmb.Logger.Printf("Error create database %s", err.Error())
			return r, err
		}
		row.Close()
		ss = "comment on table " + r.name + " is '" + r.description + "';"
		row, err = r.db.Query(ss)
		if err != nil {
			cmb.Logger.Printf("Error comment table %s database %s ", r.name, err.Error())
			return r, err
		}
		row.Close()

		desc := r.driver.GetDescription()
		for name, description := range desc {
			sd := "comment on column " + r.name + "." + name + " is '" + description + "';"
			row, err := r.db.Query(sd)
			if err != nil {
				cmb.Logger.Printf("Error comment %s database %s ", name, err.Error())
				return r, err
			}
			row.Close()
		}

	}
	cmb.Logger.Println("open table " + r.name)
	return r, nil
}

// Run собственно запускает обмен с базой данных
func (d *DBRouter) Run() {
	d.work = true
	go d.loop()
}
func (d *DBRouter) loop() {
	// old := d.driver.GetValues()
	step := time.Duration(d.step) * time.Millisecond
	for {
		time.Sleep(step)
		defer d.db.Close()
		//start := time.Now()
		if !d.work {
			return
		}
		if !d.driver.Status() {
			continue
		}
		// d.mu.Lock()
		// defer d.mu.Unlock()
		// fmt.Printf("route %s \n", d.name)
		values := d.driver.GetValues()
		t := time.Now()

		s := bytes.NewBufferString("insert into " + d.name + " ( tm")
		sv := bytes.NewBufferString(" values ('" + string(pq.FormatTimestamp(t)) + "'")
		for name, value := range values {
			s.WriteString(",")
			s.WriteString(name)
			sv.WriteString(",")
			ss := strings.Split(value, " ")
			if len(ss) == 1 {
				if value == "NaN" {
					sv.WriteString("0")
				} else {
					if strings.ToUpper(value) == "4294967096" {
						sv.WriteString("0")
					} else {
						sv.WriteString(value)
					}
				}
			} else {
				if ss[0] == "NaN" {
					sv.WriteString("0")
				} else {
					if strings.ToUpper(ss[0]) == "4294967096" {
						sv.WriteString("0")
					} else {
						sv.WriteString(ss[0])
					}
				}
			}
		}
		s.WriteString(")")
		sv.WriteString(");")
		s.WriteString(sv.String())
		rows, err := d.db.Query(s.String())
		if err != nil {
			cmb.Logger.Printf("Error write database %s %s\n", d.name, err.Error())
			cmb.Logger.Printf(s.String(), "\n")
			// d.work = false
			// return
		} else {
			rows.Close()
		}
		// for name, value := range values {
		// 	ov := old[name]
		// 	if strings.Compare(ov, value) != 0 {
		// 		fmt.Println(d.name, name, value)
		// 	}
		// }
		// old = values
		if !d.work {
			return
		}
		// d.mu.Unlock()

		//stop := time.Now()
		//elapsed := stop.Sub(start)
		// fmt.Println("router " + d.name)

	}
}

// Stop останавливает обмен с базой и закрывает соединение
func (d *DBRouter) Stop() {
	d.work = false
	d.db.Close()
}

// Status возвращает статус роутера
func (d *DBRouter) Status() bool {
	return d.work
}

// ToString роутер описание
func (d *DBRouter) ToString() string {
	return fmt.Sprintln(d)
}

// KillOld  удаляет старые записи старше dur
// delete from table where tm<oldtm;
func (d *DBRouter) KillOld(old time.Time) {
	if !d.work {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	req := "delete from " + d.name + " where tm<'" + string(pq.FormatTimestamp(old)) + "';"
	// println(d.name, req)
	rows, err := d.db.Query(req)
	if err != nil {
		cmb.Logger.Printf("Error delete old dates in database %s %s\n", d.name, err.Error())
		// fmt.Println(s.String())
		d.work = false
		return
	}
	rows.Close()
}
