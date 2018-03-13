package dt

import (
	"bytes"
	"errors"
	"fmt"
)

// DataTable описание всей таблицы
// Упрощено описание
type DataTable struct {
	Name       string
	fields     map[string]*fieldFormat
	dataStore  []*Data
	MaxRecords int
}

type fieldFormat struct {
	name         string
	description  string
	format       int8 // 0-bool 1-int 2-float 4-int64 4-string 5-date
	defaultValue string
}

// Data собственно дата
type Data struct {
	values map[string]*Value
}

//Value хранит одно значение
type Value struct {
	format int8
	value  string
}

const (
	booltype int8 = iota
	integertype
	floattype
	longtype
	stringtype
	datetime
)

// AddField добавление описания поля
func (dt *DataTable) AddField(name string, description string, format string, defaultValue string) {
	ff := new(fieldFormat)
	ff.name = name
	ff.description = description
	switch format {
	case "B":
		ff.format = booltype
	case "I":
		ff.format = integertype
	case "F":
		ff.format = floattype
	case "L":
		ff.format = longtype
	case "D":
		ff.format = datetime
	case "S":
		ff.format = stringtype
	}
	ff.defaultValue = defaultValue
	dt.fields[name] = ff
}

// NewDT создает новую таблицу данных
func NewDT(name string) *DataTable {
	dt := new(DataTable)
	dt.Name = name
	dt.dataStore = make([]*Data, 0)
	dt.fields = make(map[string]*fieldFormat)
	return dt
}

// Len возвращает кол-во записей в таблице
func (dt *DataTable) Len() int {
	return len(dt.dataStore)
}

// ReadRecod возвращает прочитанную запись из таблице
func (dt *DataTable) ReadRecod(nomer int) (*Data, error) {
	d := new(Data)
	if nomer >= 0 && nomer < dt.Len() {
		d = dt.dataStore[nomer]
		return d, nil
	}
	println("Номер записи вне таблицы")
	return d, errors.New("Номер записи вне таблицы")
}

// NewRecord создает  запись
func (dt *DataTable) NewRecord() *Data {
	rec := new(Data)
	rec.values = make(map[string]*Value)
	for _, df := range dt.fields {
		rec.MakeField(df)
	}
	return rec
}

//AddRecord записывает запись в таблицу
func (dt *DataTable) AddRecord(d *Data) {
	if dt.dataStore == nil {
		dt.dataStore = make([]*Data, 0)
	}
	dt.dataStore = append(dt.dataStore, d)
}
func (dt *DataTable) ToString() string {

	s := bytes.NewBufferString(fmt.Sprintf("DataTable %s %d \n", dt.Name, len(dt.dataStore)))
	for _, field := range dt.fields {
		s.WriteString(field.name + ":" + field.description + fmt.Sprintf(" %d", field.format) + "\n")
	}
	for i := 0; i < dt.Len(); i++ {
		data, _ := dt.ReadRecod(i)
		for _, field := range dt.fields {
			ss, _ := data.GetValue(field.name)
			s.WriteString(ss + " ")
		}
		s.WriteString("\n")
	}
	return s.String()
}
