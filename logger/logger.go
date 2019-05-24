package logger

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"rura/combo/cmb"

	"github.com/lib/pq"
)

type logmessage struct {
	prefix  string
	message string
}

var namelog string

const maxLenPipe = 1000

var work bool
var db *sql.DB
var chlog chan logmessage

func logdb(name string, condb string, init bool) {
	var err error
	work = false
	db, err := sql.Open("postgres", condb)
	if err != nil {
		cmb.Logger.Printf("logger %s DB %v\n", name, err)
		return
	}
	if err = db.Ping(); err != nil {
		cmb.Logger.Printf("logger %s DB %v\n", name, err)
		return
	}
	if init {
		ss := "drop table if exists " + name + ";"
		cmb.Logger.Println(ss)
		_, err := db.Query(ss)
		if err != nil {
			cmb.Logger.Printf("Error drop database %s", err.Error())
			return
		}
		//primary key not null
		s := "create table if not exists " + name + " ( tm timestamp primary key not null, prefix text,message text);"
		row, err := db.Query(s)
		if err != nil {
			cmb.Logger.Printf("Error create database %s", err.Error())
			return
		}
		row.Close()
		ss = "comment on table " + name + " is 'Логи всех сообщений об ошибках';"
		row, err = db.Query(ss)
		if err != nil {
			cmb.Logger.Printf("Error comment table %s database %s ", name, err.Error())
			return
		}
		row.Close()
		cmb.Logger.Printf("table %s database create ", name)
		work = true
	}
	for {
		message := <-chlog
		t := time.Now()
		s := "insert into " + name + " ( tm,prefix,message) values ('" + string(pq.FormatTimestamp(t)) + "',"
		s += "'" + message.prefix + "','" + message.message + "');"
		row, err := db.Query(s)
		if err != nil {
			cmb.Logger.Printf("Error table %s database %s ", name, err.Error())
			continue
		}
		row.Close()
		time.Sleep(1 * time.Millisecond)
	}
}

func worker(con net.Conn, name string) {
	defer con.Close()
	buffer := make([]byte, 2048)
	for {
		n, err := con.Read(buffer)
		if err != nil {
			if strings.Compare(err.Error(), "EOF") != 0 {
				cmb.Logger.Printf("Ошибка приема %s %s", name, err.Error())
				return
			}
		}
		message := string(buffer[:n])
		if len(message) == 0 {
			return
		}
		mes := logmessage{name, message}
		chlog <- mes
	}
}

func listenlog(port int, name string) {
	// listen on a port
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		cmb.Logger.Printf("При попытке слушать порт %d %s", port, err.Error())
		return
	}
	for {
		// accept a connection
		conn, err := ln.Accept()
		if err != nil {
			cmb.Logger.Printf("logger port %d %s", port, err.Error())
			fmt.Println(err)
			continue
		}
		// handle the connection
		go worker(conn, name)
	}
}

//LoadLogger загружает и запускает подсистему логирования
func LoadLogger(logs cmb.Loggers, dbcon string, name string, init bool) {
	chlog = make(chan logmessage, maxLenPipe)
	namelog = name
	go logdb(name, dbcon, init)

	for _, log := range logs.LoggerList {
		go listenlog(log.Port, log.Name)
	}
	cmb.Logger.Printf("logger %s запущен\n", name)
}

//KillOldLog функция
func KillOldLog(old time.Time) {
	if !work {
		return
	}
	req := "delete from " + namelog + " where tm<'" + string(pq.FormatTimestamp(old)) + "';"
	// println(d.name, req)
	rows, err := db.Query(req)
	if err != nil {
		cmb.Logger.Printf("Error delete old dates in database %s %s\n", namelog, err.Error())
		// fmt.Println(s.String())
		return
	}
	rows.Close()
}
