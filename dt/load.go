package dt

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strings"
)

// TableXML описывает таблицу настройки регистров Modbus
type TableXML struct {
	XMLName xml.Name  `xml:"table" json:"table"`
	Format  FormatXML `xml:"format" json:"format"`
	//	Recods Records `xml:"records" json:"records"`
	RecordListXML []RecordXML `xml:"records>record" json:"record"`
}

// FormatXML сбор форматов
type FormatXML struct {
	MaxRecors int       `xml:"maxRecords,attr" json:"maxRecords"`
	FieldsXML FieldsXML `xml:"fields" json:"fields"`
}

// FieldsXML описывает массив полей регистов
type FieldsXML struct {
	FieldListXML []FieldXML `xml:"field" json:"field"`
	Ret          int        `xml:"ret,attr" json:"ret"`
}

// FieldXML описывает отдельное поле
type FieldXML struct {
	Name           string `xml:"name,attr" json:"name"`
	Description    string `xml:"description,attr" json:"description"`
	Type           string `xml:"type,attr" json:"type"`
	Validators     string `xml:"validators,omitempty"`
	SelectionValue string `xml:"selectionValue,omitempty"`
	// DefaultValue   string `xml:"defaultValue,attr,chardata" json:"defaultValue"`
	DefaultValue DefValue `xml:"defaultValue,omitempty" json:"defaultValue"`
}

//DefValue значение по умолчанию
type DefValue struct {
	Value string `xml:",chardata" json:"value"`
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
var dt DataTable

func loadFile(path string, xmlf bool) error {
	dt = NewDT(path)
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if xmlf {
		err = xml.Unmarshal(buffer, &dt.table)
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal(buffer, &dt.table)
		if err != nil {
			return err
		}

	}
	makeFormatTable()
	loadDataValues()
	return nil
}

//LoadTableXML загружает таблицу из XML
func LoadTableXML(path string) (DataTable, error) {
	if !strings.Contains(path, ".xml") {
		path += ".xml"
	}
	err := loadFile(path, true)

	//fmt.Println(dt.TableToString())
	return dt, err
}

//LoadTableFromJSON загружает таблицу из файла с JSON
func LoadTableFromJSON(mainPath string) (DataTable, error) {
	err := loadFile(mainPath, false)
	return dt, err
}

//FromJSON загруска из JSON
func (dt *DataTable) FromJSON(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, &dt)
}
func makeFormatTable() {
	for _, field := range dt.table.Format.FieldsXML.FieldListXML {
		dt.AddField(field.Name, field.Description, field.Type, field.DefaultValue.Value)
	}
}
func loadDataValues() {
	for _, rec := range dt.table.RecordListXML {
		data := dt.NewRecord()
		for _, val := range rec.ValueListXML {
			_, ok := data.Values[val.Name]
			if ok {
				data.Values[val.Name] = val.Value
			}
		}
		dt.AddRecord(data)
	}
}
