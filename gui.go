package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//HTTPDataTable структура для вывода на экран
type HTTPDataTable struct {
	Description string
	Value       string
}
type Index struct {
	Table   map[string]*HTTPDataTable
	Devices map[string]string
	Device  string
}

var index Index
var defDevice = "du"
var defRegister = 0
var mu sync.Mutex

func parseDirTemplates(start string) (*template.Template, error) {
	templ := template.New("")
	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return templ, err
	}
	return templ, nil
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	templ, err := parseDirTemplates("templates") ///
	if err != nil {
		fmt.Println(err.Error())
		// fmt.Println(templ)
		return
	}
	query := r.URL.Query()
	dev, ok := query["id"]
	if ok {
		switch dev[0] {
		case "coils":
			defRegister = 0
		case "di":
			defRegister = 1
		case "ir":
			defRegister = 2
		case "hr":
			defRegister = 3
		default:
			defDevice = dev[0]
		}
	}
	setDevice(defDevice, defRegister)
	templ.ExecuteTemplate(w, "index", index)
}

//func editHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.Query())
// 	templ.ExecuteTemplate(w, "index", index)
// }

func setDevice(dev string, reg int) {
	driver, ok := Drivers[dev]
	if !ok {
		return
	}
	index.Table = make(map[string]*HTTPDataTable)
	index.Device = driver.Description
	switch reg {
	case 0:
		index.Device += " Coil Registers"
	case 1:
		index.Device += " Discrets Input "
	case 2:
		index.Device += " Input Registers"
	case 3:
		index.Device += " Holding Registers"
	}

	descrs := driver.GetDescription()
	values := driver.GetValues()
	types := driver.GetTypes()
	for name, descr := range descrs {
		if types[name] != reg {
			continue
		}
		res := new(HTTPDataTable)
		res.Value = values[name]
		res.Description = descr
		index.Table[name] = res
	}

}

func gui() {
	index.Devices = make(map[string]string)
	for name, dev := range Drivers {
		index.Devices[name] = dev.Description
	}

	fmt.Println("listering on port :8080")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8080", nil)
}
