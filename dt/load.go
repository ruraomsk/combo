package dt

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
)

// TableXML описывает таблицу настройки регистров Modbus
type TableXML struct {
	Format FormatXML `xml:"format" json:"format"`
	//	Recods Records `xml:"records" json:"records"`
	RecordListXML []RecordXML `xml:"records>record" json:"record"`
}

// FormatXML сбор форматов
type FormatXML struct {
	MaxRecors int       `xml:"maxRecords,attr"`
	FieldsXML FieldsXML `xml:"fields" json:"fields"`
}

// FieldsXML описывает массив полей регистов
type FieldsXML struct {
	FieldListXML []FieldXML `xml:"field" json:"field"`
	Ret          int        `xml:"ret,attr"`
}

// FieldXML описывает отдельное поле
type FieldXML struct {
	Name           string `xml:"name,attr" json:"name"`
	Description    string `xml:"description,attr" json:"description"`
	Type           string `xml:"type,attr" json:"type"`
	Validators     string `xml:"validators,omitempty"`
	SelectionValue string `xml:"selectionValue,omitempty"`
	DefaultValue   string `xml:"defaultValue,omitempty" json:"defaultValue"`
}

// RecordXML описывает отдельное поле
type RecordXML struct {
	ValueListXML []ValueXML `xml:"value" jsons:"value"`
}

// ValueXML наполнение записи
type ValueXML struct {
	Name  string `xml:"name,attr" json:"name"`
	Value string `xml:",chardata" json:"value"`
}

// DataTable описание таблицы
var dt *DataTable
var table *TableXML

func loadFile(mainPath string, xmlf bool) error {
	dt = NewDT(mainPath)
	buffer := bytes.NewBuffer(nil)
	file, err := os.Open(mainPath)
	if err != nil {
		return err
	}
	io.Copy(buffer, file)
	file.Close()
	table = new(TableXML)
	if xmlf {
		xml.Unmarshal(buffer.Bytes(), table)
	}
	if !xmlf {
		json.Unmarshal(buffer.Bytes(), table)
	}
	// fmt.Println(table)
	makeFormatTable()
	loadDataValues()
	//fmt.Printf(device.DataTable.ToString())
	return nil

}

//LoadTableXML загружает таблицу из XML
func LoadTableXML(mainPath string) (*DataTable, error) {
	err := loadFile(mainPath, true)
	//fmt.Println(dt.TableToString())
	return dt, err
}

//LoadTableJSON загружает таблицу из JSON
func LoadTableJSON(mainPath string) (*DataTable, error) {
	err := loadFile(mainPath, false)
	return dt, err
}
func makeFormatTable() {
	for _, field := range table.Format.FieldsXML.FieldListXML {
		dt.AddField(field.Name, field.Description, field.Type, field.DefaultValue)
	}
}
func loadDataValues() {
	for _, rec := range table.RecordListXML {
		data := dt.NewRecord()
		for _, val := range rec.ValueListXML {
			sv, ok := data.values[val.Name]
			if ok {
				sv.setValue(val.Value)
				data.values[val.Name] = sv
			}
		}
		dt.AddRecord(data)
	}
}
