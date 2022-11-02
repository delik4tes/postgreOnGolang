package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

//FIXME работа с app.go:
// Написать логику в header через {{ }} на демонстрацию кнопки определенной, а не два файла
// Нормально настроить все css файлы, чтобы не было конфликтов между классами и айди
// Написать структуры для каждой группы пользователей, в каждой структуре будут свои команды, которые потом можно будет повторно использовать

//TODO: Добавить функцию в PostgreSQL сколько всего учиться учеников в каждом фелиале

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

type Parameter struct {
	Reg_alert string
	Log_alert string
	Alrt      bool

	Authorization bool
}

type Test struct {
	Massive [3]string
}

var reg bool = false

var parametrs Parameter = Parameter{"Вы успешно зарегистрировались", "Вы успешно вошли", false, false}

var test Test = Test{[3]string{"One", "Two", "Three"}}

func mainPage(w http.ResponseWriter, r *http.Request) {

	if !reg {
		main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header_new.html")
		err = main.ExecuteTemplate(w, "main", parametrs)
		if err != nil {
			panic(err)
		}
	} else {
		main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header.html")
		err = main.ExecuteTemplate(w, "main", parametrs)
		if err != nil {
			panic(err)
		}
	}

	//Если сделана полностью страница html header+main+footer
	// test, err := template.ParseFiles("templates/hello_page.html")
	// err = test.Execute(w, nil)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	about, err := template.ParseFiles("templates/about.html", "templates/footer.html", "templates/header_new.html")
	err = about.ExecuteTemplate(w, "about", nil)

	//fmt.Println(r.URL.Query())

	if err != nil {
		panic(err)
	}
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	login, err := template.ParseFiles("templates/login.html", "templates/footer.html")
	err = login.ExecuteTemplate(w, "login", nil)

	if err != nil {
		panic(err)
	}
}

func registrationPage(w http.ResponseWriter, r *http.Request) {
	registration, err := template.ParseFiles("templates/registration.html", "templates/footer.html")
	err = registration.ExecuteTemplate(w, "registration", test)

	reg, parametrs.Alrt = true, true

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

func successPage(w http.ResponseWriter, r *http.Request) {

	success, err := template.ParseFiles("templates/success.html", "templates/footer.html", "templates/header.html")
	err = success.ExecuteTemplate(w, "success", parametrs)

	fmt.Println(r.URL.Query())

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
	router.HandleFunc("/success/", successPage)

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
