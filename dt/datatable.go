package dt

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

// DataTable описание всей таблицы
// Упрощено описание
type DataTable struct {
	Name       string
	fields     map[string]FieldFormat
	dataStore  []Data
	MaxRecords int
	table      TableXML
}

//FieldFormat формат поля
type FieldFormat struct {
	name         string
	description  string
	format       string
	defaultValue string
}

// Data собственно дата
type Data struct {
	values map[string]string
}

// AddField добавление описания поля
func (dt *DataTable) AddField(name string, description string, format string, defaultValue string) {
	var ff FieldFormat
	ff.name = name
	ff.description = description
	ff.format = format
	ff.defaultValue = defaultValue
	dt.fields[name] = ff
}

// NewDT создает новую таблицу данных
func NewDT(name string) DataTable {
	var dt DataTable
	dt.Name = name
	dt.dataStore = make([]Data, 0)
	dt.fields = make(map[string]FieldFormat)
	return dt
}

// Len возвращает кол-во записей в таблице
func (dt *DataTable) Len() int {
	return len(dt.dataStore)
}

// ReadRecod возвращает прочитанную запись из таблице
func (dt *DataTable) ReadRecod(nomer int) (Data, error) {
	var d Data
	if nomer >= 0 && nomer < dt.Len() {
		d = dt.dataStore[nomer]
		return d, nil
	}
	return d, errors.New("Номер записи вне таблицы")
}

// NewRecord создает  запись
func (dt *DataTable) NewRecord() Data {
	var rec Data
	rec.values = make(map[string]string)
	for _, df := range dt.fields {
		rec.MakeField(df)
	}
	return rec
}

//AddRecord записывает запись в таблицу
func (dt *DataTable) AddRecord(d Data) {
	if dt.dataStore == nil {
		dt.dataStore = make([]Data, 0)
	}
	dt.dataStore = append(dt.dataStore, d)
}

//ToString тестовое представление таблицы
func (dt *DataTable) ToString() string {

	s := bytes.NewBufferString(fmt.Sprintf("DataTable %s %d \n", dt.Name, len(dt.dataStore)))
	for _, field := range dt.fields {
		s.WriteString(field.name + ":" + field.description + " " + field.format + "\n")
	}
	for i := 0; i < dt.Len(); i++ {
		record, _ := dt.ReadRecod(i)
		for _, value := range record.values {
			s.WriteString(value + " ")
		}

		s.WriteString("\n")
	}
	return s.String()
}

//MakeField функция
func (data *Data) MakeField(ff FieldFormat) {
	data.values[ff.name] = ""
	switch ff.format {
	case "I":
		data.SetInt(ff.name, 0)
	case "B":
		data.SetBool(ff.name, false)
	case "F":
		data.SetFloat(ff.name, 0.0)
	case "L":
		data.SetLong(ff.name, 0)
	case "D":
		data.SetDate(ff.name, time.Now())
	case "S":
		data.SetString(ff.name, "")
	}
	if len(ff.defaultValue) > 0 {
		data.values[ff.name] = ff.defaultValue
	}
}
