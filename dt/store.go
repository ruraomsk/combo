package dt

import (
	"encoding/xml"
	"os"

	"github.com/clbanning/mxj"
)

//TableToString функция
func (dt *DataTable) TableToString() string {
	var s string
	names := make([]string, len(dt.fields))
	i := 0
	for _, f := range dt.fields {
		s += f.name + "\t"
		names[i] = f.name
		i++
	}
	s += "\n"
	for _, rec := range dt.dataStore {
		for _, ff := range names {
			s += rec.values[ff] + "\t"
		}
		s += "\n"
	}
	s += "==============================================="
	return s
}

//TableToXML выгружает таблицу в XML
func (dt *DataTable) TableToXML() error {
	dt.table.RecordListXML = make([]RecordXML, 0)
	for _, rec := range dt.dataStore {
		var rlx RecordXML
		rlx.ValueListXML = make([]ValueXML, 0)
		for name, value := range rec.values {
			var rxml ValueXML
			rxml.Name = name
			rxml.Value = value
			rlx.ValueListXML = append(rlx.ValueListXML, rxml)
		}
		dt.table.RecordListXML = append(dt.table.RecordListXML, rlx)
	}
	result, err := xml.Marshal(dt.table)
	if err != nil {
		return err
	}
	result, err = mxj.BeautifyXml(result, "", "\t")
	if err != nil {
		return err
	}
	file, err := os.Create(dt.Name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(result)
	return err
}
