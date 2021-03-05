package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	conf "servers/config"

	"github.com/gorilla/mux"
)

// Employee -
type Employee struct {
	ID      string `json:"Id"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
	Salary  string `json:"Salary"`
}

// Employees -
var Employees []Employee

var stopHTTPServerChan chan bool

func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Welcome to the API Home page!</h1><br><a href='/exit'>Exit</a>")
	if err != nil {
		panic(err)
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Bye from Channels</h1>")
	if err != nil {
		panic(err)
	}

	stopHTTPServerChan <- true
}

func getAllEmps(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Employees)
}

func getEmp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	for _, employee := range Employees {
		if employee.ID == key {
			json.NewEncoder(w).Encode(employee)
		}
	}
}

func createEmp(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var employee Employee
	json.Unmarshal(reqBody, &employee)
	Employees = append(Employees, employee)
	json.NewEncoder(w).Encode(employee)

	fmt.Fprintln(w, "Created")
}

func updateEmp(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var employee Employee
	json.Unmarshal(reqBody, &employee)

	vars := mux.Vars(r)
	id := vars["id"]

	yid := checkID(Employees, id)
	if !yid {
		fmt.Fprintln(w, "id is not finded")
	} else {
		for index, emp := range Employees {
			if employee.ID == Employees[index].ID {
				fmt.Fprintln(w, "This already exists, use another")
			} else if emp.ID == id {
				Employees[index] = employee
				json.NewEncoder(w).Encode(Employees)

				fmt.Fprintln(w, "Updated")
			}
		}
	}
}

func deleteEmp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	yid := checkID(Employees, id)
	if !yid {
		fmt.Fprintln(w, "id is not finded")
	}

	for index, emp := range Employees {
		if emp.ID == id {
			Employees = append(Employees[:index], Employees[index+1:]...)
			json.NewEncoder(w).Encode(Employees)

			fmt.Fprintln(w, "Deleted")
		}
	}
}

func checkID(emp []Employee, empid string) bool {
	var idemp bool
	for _, id := range emp {
		if id.ID == empid {
			idemp = true
		}
	}
	return idemp
}

func startServer(cf conf.HTTPConfig) {
	stopHTTPServerChan = make(chan bool)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/exit", exitHandler)

	r.HandleFunc("/employees", getAllEmps)

	/* curl -X POST -H "Content-Type: application/json" \
	   -d '{"Id":"1","Name":"Papa", "Address":"Poland Cuve Street 88", "Salary":"18000"}' http://localhost:8000/employee */
	r.HandleFunc("/employee", createEmp).Methods("POST")

	/* curl -X PUT http://localhost:8000/employee/1 -H "Content-Type: application/json" -H "Accept: application/json" \
	   -d '{"Id": "6","Name": "Otec","Address":"CocoChabma","Salary": "18000"}' */
	r.HandleFunc("/employee/{id}", updateEmp).Methods("PUT")

	//curl -X DELETE http://localhost:8000/employee/1
	r.HandleFunc("/employee/{id}", deleteEmp).Methods("DELETE")

	//curl http://localhost:8000/employee/1
	r.HandleFunc("/employee/{id}", getEmp)

	fmt.Println("Server started at", cf.Host, cf.Port)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", cf.Host, cf.Port),
		// Таймауты для создаваемых серверов
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		// При закрытии всегда возвращает ошибку
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Непредвиденная ошибка или порт занят
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Ждем пока не получен сигнал
	<-stopHTTPServerChan
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // сбой/тайм-аут корректное завершение работы сервера
	}
	fmt.Println("Server closed - Channels")
}

func main() {
	cnf, err := conf.NewConfig("../config/config.yaml")
	if err != nil {
		log.Fatalf("Can't read the config: %s", err)
	}

	Employees = []Employee{
		Employee{ID: "1", Name: "Jhon Smith", Address: "New Jersy USA", Salary: "20000"},
		Employee{ID: "2", Name: "William", Address: "Wellington Newziland", Salary: "12000"},
		Employee{ID: "3", Name: "Adam", Address: "London England", Salary: "15000"},
	}

	startServer(cnf.HTTP)
}
