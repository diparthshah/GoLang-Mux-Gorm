package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name string
	Age  string
}

func getAllStudents(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("testgo.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	var students []Student
	db.Find(&students)
	json.NewEncoder(w).Encode(students)
}

func addNewStudent(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("testgo.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	vars := mux.Vars(r)
	rName := vars["name"]
	rAge := vars["age"]
	db.Create(&Student{Name: rName, Age: rAge})
	fmt.Fprintf(w, "New Record Added")
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("testgo.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	vars := mux.Vars(r)
	rName := vars["name"]
	rAge := vars["age"]
	var student Student
	db.Where("name = ?", rName).Find(&student)
	db.Model(&student).Update("age", rAge)
	fmt.Fprintf(w, "Updated Record")
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("testgo.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	vars := mux.Vars(r)
	rName := vars["name"]
	var student Student
	db.Where("name = ?", rName).Find(&student)
	db.Delete(&student)
	fmt.Fprintf(w, "Deleted Record")
}

func setupDB() {
	db, err := gorm.Open(sqlite.Open("testgo.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	db.AutoMigrate(&Student{})
}

func spinOff() {
	router := mux.NewRouter()
	router.HandleFunc("/students", getAllStudents).Methods("GET")
	router.HandleFunc("/student/{name}/{age}", addNewStudent).Methods("POST")
	router.HandleFunc("/student/{name}/{age}", updateStudent).Methods("PUT")
	router.HandleFunc("/student/{name}", deleteStudent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8083", router))
}

func main() {
	fmt.Println("Starting skillbox server")
	setupDB()
	spinOff()
}
