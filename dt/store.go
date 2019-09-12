package dt

//TableToString функция
func (dt *DataTable) TableToString() string {
	var s string
	for _, f := range dt.fields {
		s += f.name + "\t"
	}
	s += "\n"
	for _, rec := range dt.dataStore {
		for _, ff := range dt.fields {
			s += rec.values[ff.name] + "\t"
		}
		s += "\n"
	}
	s += "==============================================="
	return s
}
