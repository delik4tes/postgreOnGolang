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

var login bool = false
var registration bool = false
var test string = "1"

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
	main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header_new.html")
	err = main.ExecuteTemplate(w, "main", nil)

	//Если сделана полностью страница html header+main+footer
	// test, err := template.ParseFiles("templates/hello_page.html")
	// err = test.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	about, err := template.ParseFiles("templates/about.html", "templates/footer.html", "templates/header_new.html")
	err = about.ExecuteTemplate(w, "about", nil)

	if err != nil {
		panic(err)
	}
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	login, err := template.ParseFiles("templates/login.html", "templates/footer.html")
	err = login.ExecuteTemplate(w, "login", nil)

	test = r.Form.Get("password")

	if err != nil {
		panic(err)
	}
}

func registrationPage(w http.ResponseWriter, r *http.Request) {
	registration, err := template.ParseFiles("templates/registration.html", "templates/footer.html")
	err = registration.ExecuteTemplate(w, "registration", nil)

	if err != nil {
		panic(err)
	}
}

func contractPage(w http.ResponseWriter, r *http.Request) {
	contract, err := template.ParseFiles("templates/contract.html", "templates/footer.html", "templates/header_new.html")
	err = contract.ExecuteTemplate(w, "contract", nil)

	if err != nil {
		panic(err)
	}
}

func handlerRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", mainPage).Methods("GET")
	router.HandleFunc("/about/", aboutPage).Methods("GET")
	router.HandleFunc("/login/", loginPage)
	router.HandleFunc("/registration/", registrationPage)
	router.HandleFunc("/contract/", contractPage)

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
