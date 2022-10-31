package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

const password = " password=postgres"
const user = "user=postgres"
const dbname = " dbname=tr_zhelagin_12"
const sslmode = " sslmode=disable"

const connectParam = user + password + dbname + sslmode

type Branch struct {
}

type Clients struct {
}

type Contract struct {
}

type Teachers struct {
}

type Logins struct {
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	//test, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header_new.html")
	//err = test.ExecuteTemplate(w, "main", nil)
	test, err := template.ParseFiles("templates/hello_page.html")
	err = test.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func handlerRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", mainPage).Methods("GET")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("start server")
	handlerRequest()
}
