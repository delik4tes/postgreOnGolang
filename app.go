package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

//FIXME работа с app.go:
// -1.Написать структуры для каждой группы пользователей, в каждой структуре будут свои команды, которые потом можно будет повторно использовать
// -2.Переименовать в таблице branch колонку address
// Учитель:
//  -1.Договоры на которые он назначен (Сравнивать id учителя контракта и таблицы)
//  -2.Информацию о своих клиентах (Нужно сравнивать id учителя контракта и таблицы, получать id клиента и выводить все данные связанные с ним)
// Клиент:
//  -1.Просмотр доступных языков и их преподавателей
// Администратор:
//  -1.Просмотр и редактирование всех договоров
//  -2.Установка зарплаты для учителя
//  -3.Ставить статус договора
// Владелец:
//  -1.Полный доступ (объединить все предыдущие возможности + добавить недостающие)

//TODO: Добавить функцию в PostgreSQL сколько всего учиться учеников в каждом филиале

const password = " password=postgres"
const user = "user=postgres"
const dbname = " dbname=tr_zhelagin_12"
const sslmode = " sslmode=disable"

const connectParam = user + password + dbname + sslmode

type Branch struct {
	Id                                        uint16
	Address, Name, Surname, Patronymic, Login string
	Salary                                    float32
}

type Client struct {
	Id, Branch                              uint16
	Name, Surname, Patronymic, Phone, Login string
}

type Contract struct {
	Id, Client, Teacher uint16
	Quantity            uint32
	Language, Status    string
	Price               float32
}

type Teacher struct {
	Id, Branch, Experience                     uint16
	Name, Surname, Patronymic, Language, Login string
	Salary                                     float32
}

type Login struct {
	Email, Password, Login string
	Status                 rune
}

var tableBranches []Branch
var tableClients []Client
var tableContracts []Contract
var tableTeachers []Teacher
var tableLogins []Login

type Parameter struct {
	RegAlert string
	LogAlert string
	Alrt     bool

	Authorization bool
}

//-------------------

var parameters Parameter = Parameter{"Вы успешно зарегистрировались", "Вы успешно вошли", false, false}

//fmt.Println(request.URL.Query())

// -------------------

func mainPage(w http.ResponseWriter, request *http.Request) {

	main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header_new.html")
	err = main.ExecuteTemplate(w, "main", nil)
	if err != nil {
		panic(err)
	}

}

func aboutPage(writer http.ResponseWriter, request *http.Request) {
	about, err := template.ParseFiles("templates/about.html", "templates/footer.html", "templates/header_new.html")
	err = about.ExecuteTemplate(writer, "about", nil)

	if err != nil {
		panic(err)
	}
}
func loginPage(writer http.ResponseWriter, request *http.Request) {

	login, err := template.ParseFiles("templates/login.html", "templates/footer.html")
	err = login.ExecuteTemplate(writer, "login", nil)

	if err != nil {
		panic(err)
	}
}

func registrationPage(writer http.ResponseWriter, request *http.Request) {

	db, err := sql.Open("postgres", connectParam)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	result, err := db.Query("SELECT address FROM branch")
	var addresses []string
	for result.Next() {
		var address string
		err := result.Scan(&address)
		if err != nil {
			return
		}
		addresses = append(addresses, address)
	}

	registration, err := template.ParseFiles("templates/registration.html", "templates/footer.html")
	err = registration.ExecuteTemplate(writer, "registration", addresses)

	if err != nil {
		panic(err)
	}
}

func saveRegistrationForm(writer http.ResponseWriter, request *http.Request) {

	db, err := sql.Open("postgres", connectParam)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	result := request.URL.Query()

	if result["position"][0] == "student" {
		_, err = db.Exec("INSERT INTO logins (email,password,status) VALUES ($1,$2,$3)", result["mail"][0], result["password"][0], "S")
		if err != nil {
			panic(err)
		}
		id := db.QueryRow("SELECT id FROM branch WHERE address = $1", result["address"][0])
		var idBranch string
		err = id.Scan(&idBranch)
		if err != nil {
			panic(err)
		}
		mail := db.QueryRow("SELECT login FROM logins WHERE email = $1", result["mail"][0])
		var studentLogin string
		err = mail.Scan(&studentLogin)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("INSERT INTO clients (name,surname,patronymic,branch,phone, login) VALUES ($1,$2,$3,$4,$5,$6)", result["name"][0], result["surname"][0], result["patronymic"][0], idBranch, result["phone"][0], studentLogin)
	} else if result["position"][0] == "teacher" {
		_, err = db.Exec("INSERT INTO ")
	} else if result["position"][0] == "admin" {
		_, err = db.Exec("INSERT INTO ")
	}

	http.Redirect(writer, request, "/success/", http.StatusSeeOther)
}

func contractPage(writer http.ResponseWriter, request *http.Request) {
	contract, err := template.ParseFiles("templates/contract.html", "templates/footer.html", "templates/header_new.html")
	err = contract.ExecuteTemplate(writer, "contract", nil)

	if err != nil {
		panic(err)
	}
}

func successPage(writer http.ResponseWriter, request *http.Request) {

	success, err := template.ParseFiles("templates/success.html", "templates/footer.html", "templates/header.html")
	err = success.ExecuteTemplate(writer, "success", parameters)

	if err != nil {
		panic(err)
	}
}

func teacherCabinet(writer http.ResponseWriter, request *http.Request) {

	db, err := sql.Open("postgres", connectParam)
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	result, err := db.Query("SELECT * FROM teachers")

	for result.Next() {
		var tmp Teacher
		err = result.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		tableTeachers = append(tableTeachers, tmp)
	}

	teacher, err := template.ParseFiles("templates/teacher.html")
	err = teacher.Execute(writer, tableTeachers)

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
	router.HandleFunc("/saveRegistrationForm/", saveRegistrationForm)
	router.HandleFunc("/contract/", contractPage)
	router.HandleFunc("/success/", successPage)

	router.HandleFunc("/teacher/", teacherCabinet)

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
