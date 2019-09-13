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
	Fields     map[string]FieldFormat
	DataStore  []Data
	MaxRecords int
	table      TableXML
}

//FieldFormat формат поля
type FieldFormat struct {
	Name         string
	Description  string
	Format       string
	DefaultValue string
}

// Data собственно дата
type Data struct {
	Values map[string]string
}

// AddField добавление описания поля
func (dt *DataTable) AddField(name string, description string, format string, defaultValue string) {
	var ff FieldFormat
	ff.Name = name
	ff.Description = description
	ff.Format = format
	ff.DefaultValue = defaultValue
	dt.Fields[name] = ff
}

// NewDT создает новую таблицу данных
func NewDT(name string) DataTable {
	var dt DataTable
	dt.Name = name
	dt.DataStore = make([]Data, 0)
	dt.Fields = make(map[string]FieldFormat)
	return dt
}

// Len возвращает кол-во записей в таблице
func (dt *DataTable) Len() int {
	return len(dt.DataStore)
}

// ReadRecod возвращает прочитанную запись из таблице
func (dt *DataTable) ReadRecod(nomer int) (Data, error) {
	var d Data
	if nomer >= 0 && nomer < dt.Len() {
		d = dt.DataStore[nomer]
		return d, nil
	}
	return d, errors.New("Номер записи вне таблицы")
}

// NewRecord создает  запись
func (dt *DataTable) NewRecord() Data {
	var rec Data
	rec.Values = make(map[string]string)
	for _, df := range dt.Fields {
		rec.MakeField(df)
	}
	return rec
}

//AddRecord записывает запись в таблицу
func (dt *DataTable) AddRecord(d Data) {
	if dt.DataStore == nil {
		dt.DataStore = make([]Data, 0)
	}
	dt.DataStore = append(dt.DataStore, d)
}

//ToString тестовое представление таблицы
func (dt *DataTable) ToString() string {

	s := bytes.NewBufferString(fmt.Sprintf("DataTable %s %d \n", dt.Name, len(dt.DataStore)))
	for _, field := range dt.Fields {
		s.WriteString(field.Name + ":" + field.Description + " " + field.Format + "\n")
	}
	for i := 0; i < dt.Len(); i++ {
		record, _ := dt.ReadRecod(i)
		for _, value := range record.Values {
			s.WriteString(value + " ")
		}

		s.WriteString("\n")
	}
	return s.String()
}

//MakeField функция
func (data *Data) MakeField(ff FieldFormat) {
	data.Values[ff.Name] = ""
	switch ff.Format {
	case "I":
		data.SetInt(ff.Name, 0)
	case "B":
		data.SetBool(ff.Name, false)
	case "F":
		data.SetFloat(ff.Name, 0.0)
	case "L":
		data.SetLong(ff.Name, 0)
	case "D":
		data.SetDate(ff.Name, time.Now())
	case "S":
		data.SetString(ff.Name, "")
	}
	if len(ff.DefaultValue) > 0 {
		data.Values[ff.Name] = ff.DefaultValue
	}
}
