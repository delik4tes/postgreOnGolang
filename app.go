package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//FIXME работа с app.go:
// Учитель:
//  -1.Договоры на которые он назначен (Сравнивать id учителя контракта и таблицы)
//  -2.Информацию о своих клиентах (Нужно сравнивать id учителя контракта и таблицы, получать id клиента и выводить все данные связанные с ним)
// Администратор:
//  -1.Просмотр и редактирование всех договоров
//  -2.Установка зарплаты для учителя
//  -3.Ставить статус договора
// Владелец:
//  -1.Полный доступ (объединить все предыдущие возможности + добавить недостающие)
//  -2.Установка зарплаты для администратора

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
	Id, Client, Teacher    uint16
	Quantity               uint32
	Language, Status, Date string
	Price                  float32
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

type Parameter struct {
	Message                                      string
	Out, Login, Registration                     bool
	OutMask, LoginMask, RegistrationMask         bool
	Authorization, checkStudent, successContract bool
}

type User struct {
	Status, Login string
}

type Language struct {
	English []Teacher
	Germany []Teacher
	French  []Teacher
	Spanish []Teacher
	China   []Teacher
	Japan   []Teacher
	Hindi   []Teacher
	Hebrew  []Teacher
	Kazakh  []Teacher
	Chuvash []Teacher
	Turkish []Teacher
	Arabic  []Teacher
}

type ContractsAndClients struct {
	Contracts []Contract
	Clients   []Client
}

var tableBranches []Branch
var tableClients []Client
var tableContracts []Contract
var tableTeachers []Teacher
var tableLogins []Login

var parameters = Parameter{
	"",
	true,
	false,
	false,
	false,
	false,
	false,
	true,
	true,
	false}

var currentUser User
var languages Language
var contractsAndClients ContractsAndClients

func mainPage(w http.ResponseWriter, request *http.Request) {

	if parameters.Login || parameters.Registration {
		main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header.html")
		err = main.ExecuteTemplate(w, "main", nil)
		if err != nil {
			panic(err)
		}

	} else if parameters.Out {
		main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header_new.html")
		err = main.ExecuteTemplate(w, "main", nil)
		if err != nil {
			panic(err)
		}
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

func checkLoginForm(writer http.ResponseWriter, request *http.Request) {

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

	tmp := db.QueryRow("SELECT login FROM logins WHERE email = $1", result["login"][0])
	var check string
	err = tmp.Scan(&check)
	if err != nil {
		panic(err)
	}

	if check != "" {
		currentUser.Login = check
		parameters.Login = true
		parameters.LoginMask = true
		parameters.Registration = false
		parameters.RegistrationMask = false
		parameters.Out = false
		parameters.OutMask = false
		parameters.Authorization = true

	} else {
		parameters.Login = false
		parameters.LoginMask = true
		parameters.Registration = false
		parameters.Out = true
		parameters.OutMask = false
	}

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
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
		_, err = db.Exec("INSERT INTO logins (email,password,status) VALUES ($1,$2,$3)",
			result["mail"][0], result["password"][0], "S")
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
		_, err = db.Exec("INSERT INTO clients (name,surname,patronymic,branch,phone, login) VALUES ($1,$2,$3,$4,$5,$6)",
			result["name"][0], result["surname"][0], result["patronymic"][0], idBranch, result["phone"][0], studentLogin)
		if err != nil {
			panic(err)
		}

		currentUser.Login = studentLogin

	} else if result["position"][0] == "teacher" {

		tmp := db.QueryRow("SELECT login FROM logins WHERE email = $1", result["mail"][0])
		var check string
		err = tmp.Scan(&check)

		_, err = db.Exec("INSERT INTO logins (email,password,status) VALUES ($1,$2,$3)",
			result["mail"][0], result["password"][0], "T")
		if err != nil {
			fmt.Println(result)
			panic(err)
		}
		id := db.QueryRow("SELECT id FROM branch WHERE address = $1", result["address"][0])
		var idBranch string
		err = id.Scan(&idBranch)
		mail := db.QueryRow("SELECT login FROM logins WHERE email = $1", result["mail"][0])
		var teacherLogin string
		err = mail.Scan(&teacherLogin)
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("INSERT INTO teachers (name,surname,patronymic,language,experience,login,branch) VALUES ($1,$2,$3,$4,$5,$6,$7)",
			result["name"][0], result["surname"][0], result["patronymic"][0], result["language"][0], result["experience"][0], teacherLogin, idBranch)
		if err != nil {
			panic(err)
		}

		currentUser.Login = teacherLogin

	} else if result["position"][0] == "admin" {
		_, err = db.Exec("INSERT INTO logins (email,password,status) VALUES ($1,$2,$3)",
			result["mail"][0], result["password"][0], "A")
		if err != nil {
			panic(err)
		}
		mail := db.QueryRow("SELECT login FROM logins WHERE email = $1", result["mail"][0])
		var adminLogin string
		err = mail.Scan(&adminLogin)
		_, err = db.Exec("INSERT INTO branch (name,surname,patronymic,address,login) VALUES ($1,$2,$3,$4,$5)",
			result["name"][0], result["surname"][0], result["patronymic"][0], result["insert-address"][0], adminLogin)
		if err != nil {
			panic(err)
		}

		currentUser.Login = adminLogin
	}

	parameters.Authorization = true
	parameters.Login = true
	parameters.LoginMask = false
	parameters.Registration = true
	parameters.RegistrationMask = true
	parameters.Out = false
	parameters.OutMask = false

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func checkOut(writer http.ResponseWriter, request *http.Request) {

	parameters.Login = false
	parameters.LoginMask = false
	parameters.Registration = false
	parameters.RegistrationMask = false
	parameters.Out = true
	parameters.OutMask = true

	currentUser = User{"", ""}
	parameters.Authorization = false

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func contractPage(writer http.ResponseWriter, request *http.Request) {

	if !parameters.Authorization {
		http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
	}

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

	status := db.QueryRow("SELECT status FROM logins WHERE login = $1", currentUser.Login)
	err = status.Scan(&currentUser.Status)
	if err != nil {
		panic(err)
	}

	if currentUser.Status != "S" {
		parameters.checkStudent = false
		http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
	}

	english, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Английский язык")
	for english.Next() {
		var tmp Teacher
		_ = english.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.English = append(languages.English, tmp)
	}

	germany, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Немецкий язык")
	for germany.Next() {
		var tmp Teacher
		_ = germany.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Germany = append(languages.Germany, tmp)
	}

	french, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Французский язык")
	for french.Next() {
		var tmp Teacher
		_ = french.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.French = append(languages.French, tmp)
	}

	spanish, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Испанский язык")
	for spanish.Next() {
		var tmp Teacher
		_ = spanish.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Spanish = append(languages.Spanish, tmp)
	}

	china, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Китайский язык")
	for china.Next() {
		var tmp Teacher
		_ = china.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.China = append(languages.China, tmp)
	}

	japan, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Японский язык")
	for japan.Next() {
		var tmp Teacher
		_ = japan.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Japan = append(languages.Japan, tmp)
	}

	hindi, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Хинди")
	for hindi.Next() {
		var tmp Teacher
		_ = hindi.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Hindi = append(languages.Hindi, tmp)
	}

	hebrew, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Иврит")
	for hebrew.Next() {
		var tmp Teacher
		_ = hebrew.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Hebrew = append(languages.Hebrew, tmp)
	}

	kazakh, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Казахский язык")
	for kazakh.Next() {
		var tmp Teacher
		_ = kazakh.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Kazakh = append(languages.Kazakh, tmp)
	}

	chuvash, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Чувашский язык")
	for chuvash.Next() {
		var tmp Teacher
		_ = chuvash.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Chuvash = append(languages.Chuvash, tmp)
	}

	turkish, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Турецкий язык")
	for turkish.Next() {
		var tmp Teacher
		_ = turkish.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Turkish = append(languages.Turkish, tmp)
	}

	arabic, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Арабский язык")
	for arabic.Next() {
		var tmp Teacher
		_ = arabic.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
		languages.Arabic = append(languages.Arabic, tmp)
	}

	parameters.Authorization = true
	parameters.checkStudent = true
	parameters.successContract = false

	contract, err := template.ParseFiles("templates/contract.html", "templates/footer.html", "templates/header.html")
	err = contract.ExecuteTemplate(writer, "contract", languages)

	if err != nil {
		panic(err)
	}
}

func saveContract(writer http.ResponseWriter, request *http.Request) {

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

	//idClient := db.QueryRow()
	//idTeacher := db.QueryRow()

	contractPrice := strings.Trim(result["contract-dynamic_contract"][0], "₽")
	fmt.Println(contractPrice)

	_, err = db.Exec("INSERT INTO contract (client,teacher,language,quantity,dynamic_contract) VALUES ($1,$2,$3,$4,$5)",
		result["contract-language"][0], result["contract-month"][0], contractPrice)

	parameters.successContract = true

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func alertPage(writer http.ResponseWriter, request *http.Request) {

	if parameters.LoginMask {
		if parameters.Login {
			parameters.Message = "Успешная авторизация в систему"
		} else {
			parameters.Message = "Вход в систему не выполнен"
		}
	}

	if parameters.RegistrationMask {
		if parameters.Registration {
			parameters.Message = "Успешная регистрация в систему"
		} else {
			parameters.Message = "Регистрация не выполнена"
		}
	}

	if parameters.OutMask {
		if parameters.Out {
			parameters.Message = "Успешный выход из личного кабинета"
		}
	}

	if parameters.successContract {
		parameters.Message = "Успешная запись на курс"
	}
	if !parameters.checkStudent {
		parameters.Message = "Записаться на курс можно только учеником"
	}

	if !parameters.Authorization {
		parameters.Message = "Сначала нужно авторизоваться в системе"
	}

	if parameters.Out {
		alert, err := template.ParseFiles("templates/alert.html", "templates/footer.html", "templates/header_new.html")
		err = alert.ExecuteTemplate(writer, "alert", parameters)

		if err != nil {
			panic(err)
		}
	} else {
		alert, err := template.ParseFiles("templates/alert.html", "templates/footer.html", "templates/header.html")
		err = alert.ExecuteTemplate(writer, "alert", parameters)

		if err != nil {
			panic(err)
		}
	}

}

func checkStatus(writer http.ResponseWriter, request *http.Request) {

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

	status := db.QueryRow("SELECT status FROM logins WHERE login = $1", currentUser.Login)
	err = status.Scan(&currentUser.Status)
	if err != nil {
		panic(err)
	}

	if currentUser.Status == "admin" {
		http.Redirect(writer, request, "/director/", http.StatusSeeOther)
	} else if currentUser.Status == "A" {
		http.Redirect(writer, request, "/admin/", http.StatusSeeOther)
	} else if currentUser.Status == "S" {
		http.Redirect(writer, request, "/student/", http.StatusSeeOther)
	} else if currentUser.Status == "T" {
		http.Redirect(writer, request, "/teacher/", http.StatusSeeOther)
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

	id := db.QueryRow("SELECT id FROM teachers WHERE login = $1", currentUser.Login)
	var idTeacher string
	_ = id.Scan(&idTeacher)

	res, _ := db.Query("SELECT * FROM contracts WHERE teacher = $1", idTeacher)
	for res.Next() {
		var contract Contract
		_ = res.Scan(&contract.Id, &contract.Client, &contract.Teacher,
			&contract.Language, &contract.Quantity, &contract.Price,
			&contract.Date, &contract.Status)
		tableContracts = append(tableContracts, contract)
	}

	for _, contract := range tableContracts {
		res, _ := db.Query("SELECT * FROM clients WHERE id = $1", contract.Client)
		for res.Next() {
			var client Client
			_ = res.Scan(&client.Id, &client.Name, &client.Surname, &client.Patronymic,
				&client.Branch, &client.Phone, &client.Login)
			tableClients = append(tableClients, client)
		}
	}

	contractsAndClients.Clients = tableClients
	contractsAndClients.Contracts = tableContracts

	teacher, err := template.ParseFiles("templates/teacher.html")
	err = teacher.Execute(writer, contractsAndClients)

	if err != nil {
		panic(err)
	}
}

// Полный доступ ко всему просмотру
func directorCabinet(writer http.ResponseWriter, request *http.Request) {

}

func adminCabinet(writer http.ResponseWriter, request *http.Request) {

}

func studentCabinet(writer http.ResponseWriter, request *http.Request) {

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

	id := db.QueryRow("SELECT id FROM clients WHERE login = $1", currentUser.Login)
	var idStudent string
	_ = id.Scan(&idStudent)

	res, err := db.Query("SELECT * FROM contracts WHERE client = $1", idStudent)
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var contract Contract
		_ = res.Scan(&contract.Id, &contract.Client, &contract.Teacher,
			&contract.Language, &contract.Quantity, &contract.Price,
			&contract.Date, &contract.Status)
		tableContracts = append(tableContracts, contract)
	}

	student, err := template.ParseFiles("templates/student.html")
	err = student.Execute(writer, tableContracts)

	if err != nil {
		panic(err)
	}
}

func handlerRequest() {
	router := mux.NewRouter()
	currentUser = User{"", ""}
	router.HandleFunc("/", mainPage).Methods("GET")
	router.HandleFunc("/about/", aboutPage).Methods("GET")
	router.HandleFunc("/login/", loginPage)
	router.HandleFunc("/checkLoginForm/", checkLoginForm)
	router.HandleFunc("/registration/", registrationPage)
	router.HandleFunc("/saveRegistrationForm/", saveRegistrationForm)
	router.HandleFunc("/contract/", contractPage)
	router.HandleFunc("/checkOut/", checkOut)
	router.HandleFunc("/saveContract/", saveContract)

	router.HandleFunc("/checkStatus/", checkStatus)
	router.HandleFunc("/admin/", adminCabinet)
	router.HandleFunc("/director/", directorCabinet)
	router.HandleFunc("/student/", studentCabinet)
	router.HandleFunc("/teacher/", teacherCabinet)

	router.HandleFunc("/alert/", alertPage)

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
