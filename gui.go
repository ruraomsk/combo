package main

import (
	"encoding/json"
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
	// Id          string
	Description string
	Value       string
}

//Index для вывода
type Index struct {
	Table   map[string]*HTTPDataTable
	Devices map[string]string
	Device  string
}

var index Index

// var defDevice = "du"

// var defRegister = 0
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

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	var defRegister = 0
// 	mu.Lock()
// 	defer mu.Unlock()
// 	templ, err := parseDirTemplates("templates") ///
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		// fmt.Println(templ)
// 		return
// 	}
// 	query := r.URL.Query()
// 	dev, ok := query["id"]
// 	if ok {
// 		switch dev[0] {
// 		case "coils":
// 			defRegister = 0
// 		case "di":
// 			defRegister = 1
// 		case "ir":
// 			defRegister = 2
// 		case "hr":
// 			defRegister = 3
// 		default:
// 			defDevice = dev[0]
// 		}
// 	}
// 	setDevice(defDevice, defRegister)
// 	templ.ExecuteTemplate(w, "index", index)
// }

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
	// var i = 0
	for name, descr := range descrs {
		if types[name] != reg {
			continue
		}
		res := new(HTTPDataTable)
		// res.Id = name
		res.Value = values[name]
		res.Description = descr
		index.Table[name] = res

		// index.Table[i] = res
		// i++
	}

}

func responJSON(w http.ResponseWriter, r *http.Request) {
	//Register variable hasn't to GLOBAL
	//Because will be simular change page all users
	var defRegister = 0
	var defDevice = "du"

	mu.Lock()
	defer mu.Unlock()
	query := r.URL.Query()
	reg, ok := query["reg"]
	if ok {
		switch reg[0] {
		case "coils":
			defRegister = 0
		case "di":
			defRegister = 1
		case "ir":
			defRegister = 2
		case "hr":
			defRegister = 3
		default:
			defRegister = 0
		}
	}
	dev, ok := query["dev"]
	if ok {
		defDevice = dev[0]
	} else {
		defDevice = "baz1"
	}

	// if ok {
	// 	switch dev[0] {
	// 	case "baz1":
	// 		defDevice = "baz1"
	// 	case "baz2":
	// 		defDevice = "baz2"
	// 	case "ctrl":
	// 		defDevice = "ctrl"
	// 	case "dozs":
	// 		defDevice = "dozs"
	// 	case "du":
	// 		defDevice = "du"
	// 	default:
	// 		defDevice = "baz1"
	// 	}
	// }

	setDevice(defDevice, defRegister)
	b, err := json.Marshal(index)
	if err != nil {
		fmt.Print("Misstake is created JSON answer")
	}
	w.Write(b)
}

// func resp(w http.ResponseWriter, r *http.Request) {
// 	var defRegister = 0
// 	query := r.URL.Query()
// 	dev, ok := query["id"]
// 	if ok {
// 		switch dev[0] {
// 		case "coils":
// 			defRegister = 0
// 		case "di":
// 			defRegister = 1
// 		case "ir":
// 			defRegister = 2
// 		case "hr":
// 			defRegister = 3
// 		default:
// 			defDevice = dev[0]
// 		}
// 	}
// 	defRegister = 1
// 	http.ServeFile(w, r, "./index.html")

// }

func gui() {
	index.Devices = make(map[string]string)
	for name, dev := range Drivers {
		index.Devices[name] = dev.Description
	}

	fmt.Println("listering on port :8181")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	// http.HandleFunc("/", resp)
	// http.HandleFunc("/style.css", func(response http.ResponseWriter, request *http.Request) {
	// 	http.ServeFile(response, request, "/assets/css/style.css")

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, "./index.html")
	})

	http.HandleFunc("/data.json", responJSON)

	http.ListenAndServe(":8181", nil)
}
