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
// Владелец:
//  -1.Полный доступ (объединить все предыдущие возможности + добавить недостающие)
//  -2.Установка зарплаты для администратора
//TODO: добавить кнопку для администратора и директора, что именно надо показывать им

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
	Id                                              uint16
	Name, Surname, Patronymic, Phone, Login, Branch string
}

type Contract struct {
	Id, Client, Teacher    uint16
	Quantity               uint32
	Language, Status, Date string
	Price                  float32
}

type Teacher struct {
	Id, Experience                                     uint16
	Name, Surname, Patronymic, Language, Login, Branch string
	Salary                                             float32
}

type Login struct {
	Email, Password, Login, Status string
}

type Parameter struct {
	Message                                                  string
	Out, Login, Registration                                 bool
	OutMask, LoginMask, RegistrationMask                     bool
	Authorization, checkStudent, successContract             bool
	AuthorizationMask, checkStudentMask, successContractMask bool
	TeacherCabinet, StudentCabinet                           bool
	TeacherCabinetMask, StudentCabinetMask                   bool
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
	Id                               uint16
	Quantity, Price                  uint32
	Language, Status, Date, Branch   string
	Name, Surname, Patronymic, Phone string
}

type TeacherInfo struct {
	Main                Teacher
	ContractsAndClients []ContractsAndClients
}

type ContractsAndTeachers struct {
	Id, Client                                     uint16
	Quantity, Price                                uint32
	Language, Status, Date, Branch                 string
	NameTeacher, SurnameTeacher, PatronymicTeacher string
}

type ClientInfo struct {
	Main                 Client
	ContractsAndTeachers []ContractsAndTeachers
}

type AdminInfo struct {
	TableClients   []Client
	TableContracts []Contract
	TableTeachers  []Teacher
}

type DirectorInfo struct {
	TableBranches  []Branch
	TableClients   []Client
	TableContracts []Contract
	TableTeachers  []Teacher
	TableLogins    []Login
}

var parameters = Parameter{
	"",
	true,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false,
	false}

var currentUser User
var languages Language
var teacherInfo TeacherInfo
var clientInfo ClientInfo
var adminInfo AdminInfo
var directorInfo DirectorInfo

func mainPage(w http.ResponseWriter, _ *http.Request) {

	if !parameters.Out {
		main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header.html")
		err = main.ExecuteTemplate(w, "main", nil)
		if err != nil {
			panic(err)
		}

	} else {
		main, err := template.ParseFiles("templates/main.html", "templates/footer.html", "templates/header_new.html")
		err = main.ExecuteTemplate(w, "main", nil)
		if err != nil {
			panic(err)
		}
	}
}

func aboutPage(writer http.ResponseWriter, _ *http.Request) {

	if !parameters.Out {
		about, err := template.ParseFiles("templates/about.html", "templates/footer.html", "templates/header.html")
		err = about.ExecuteTemplate(writer, "about", nil)
		if err != nil {
			panic(err)
		}
	} else {
		about, err := template.ParseFiles("templates/about.html", "templates/footer.html", "templates/header_new.html")
		err = about.ExecuteTemplate(writer, "about", nil)
		if err != nil {
			panic(err)
		}
	}
}

func loginPage(writer http.ResponseWriter, _ *http.Request) {

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

	tmp := db.QueryRow("SELECT EXISTS(SELECT login FROM logins WHERE email = $1)", result["login"][0])
	var exist bool
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		tmp = db.QueryRow("SELECT login FROM logins WHERE email = $1", result["login"][0])
		err = tmp.Scan(&currentUser.Login)
		if err != nil {
			panic(err)
		}
		parameters.Login = true
		parameters.LoginMask = true
		parameters.Registration = false
		parameters.RegistrationMask = false
		parameters.Out = false
		parameters.OutMask = false
		parameters.checkStudent = false
		parameters.checkStudentMask = false
		parameters.successContract = false
		parameters.successContractMask = false
		parameters.Authorization = true
		parameters.AuthorizationMask = false
		parameters.TeacherCabinet = false
		parameters.TeacherCabinetMask = false
		parameters.StudentCabinet = false
		parameters.StudentCabinetMask = false
	} else {
		parameters.Login = false
		parameters.LoginMask = true
		parameters.Registration = false
		parameters.RegistrationMask = false
		parameters.Out = true
		parameters.OutMask = false
		parameters.checkStudent = false
		parameters.checkStudentMask = false
		parameters.successContract = false
		parameters.successContractMask = false
		parameters.Authorization = false
		parameters.AuthorizationMask = false
		parameters.TeacherCabinet = false
		parameters.TeacherCabinetMask = false
		parameters.StudentCabinet = false
		parameters.StudentCabinetMask = false
	}

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func registrationPage(writer http.ResponseWriter, _ *http.Request) {

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
	fmt.Println(result)

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

	parameters.Login = true
	parameters.LoginMask = false
	parameters.Registration = true
	parameters.RegistrationMask = true
	parameters.Out = false
	parameters.OutMask = false
	parameters.checkStudent = false
	parameters.checkStudentMask = false
	parameters.successContract = false
	parameters.successContractMask = false
	parameters.Authorization = true
	parameters.AuthorizationMask = false

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func checkOut(writer http.ResponseWriter, request *http.Request) {

	parameters.Login = false
	parameters.LoginMask = false
	parameters.Registration = false
	parameters.RegistrationMask = false
	parameters.Out = true
	parameters.OutMask = true
	parameters.checkStudent = false
	parameters.checkStudentMask = false
	parameters.successContract = false
	parameters.successContractMask = false
	parameters.Authorization = false
	parameters.AuthorizationMask = false
	parameters.TeacherCabinet = false
	parameters.TeacherCabinetMask = false
	parameters.StudentCabinet = false
	parameters.StudentCabinetMask = false

	currentUser = User{"", ""}
	languages = Language{}

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func contractPage(writer http.ResponseWriter, request *http.Request) {

	if parameters.Authorization {
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
			parameters.checkStudentMask = true
			http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
		}

		tmp := db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Английский язык")
		if err != nil {
			panic(err)
		}
		var exist bool
		_ = tmp.Scan(&exist)
		fmt.Println(exist)
		if exist {
			english, err := db.Query("SELECT * FROM teachers WHERE language = $1", "Английский язык")
			for english.Next() {
				var tmp Teacher
				err = english.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
				if err != nil {
					panic(err)
				}
				languages.English = append(languages.English, tmp)

			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Немецкий язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			germany, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Немецкий язык")
			if germany != nil {
				for germany.Next() {
					var tmp Teacher
					err = germany.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Germany = append(languages.Germany, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Французский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			french, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Французский язык")
			if french != nil {
				for french.Next() {
					var tmp Teacher
					err = french.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.French = append(languages.French, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Испанский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			spanish, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Испанский язык")
			if spanish != nil {
				for spanish.Next() {
					var tmp Teacher
					err = spanish.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Spanish = append(languages.Spanish, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Китайский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			china, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Китайский язык")
			if china != nil {
				for china.Next() {
					var tmp Teacher
					err = china.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.China = append(languages.China, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Японский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			japan, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Японский язык")
			if japan != nil {
				for japan.Next() {
					var tmp Teacher
					err = japan.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Japan = append(languages.Japan, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Хинди")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			hindi, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Хинди")
			if hindi != nil {
				for hindi.Next() {
					var tmp Teacher
					err = hindi.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Hindi = append(languages.Hindi, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Иврит")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			hebrew, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Иврит")
			if hebrew != nil {
				for hebrew.Next() {
					var tmp Teacher
					err = hebrew.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Hebrew = append(languages.Hebrew, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Казахский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			kazakh, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Казахский язык")
			if kazakh != nil {
				for kazakh.Next() {
					var tmp Teacher
					err = kazakh.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Kazakh = append(languages.Kazakh, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Чувашский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			chuvash, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Чувашский язык")
			if chuvash != nil {
				for chuvash.Next() {
					var tmp Teacher
					err = chuvash.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Chuvash = append(languages.Chuvash, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Турецкий язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			turkish, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Турецкий язык")
			if turkish != nil {
				for turkish.Next() {
					var tmp Teacher
					err = turkish.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Turkish = append(languages.Turkish, tmp)
				}
			}
		}

		tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers WHERE language = $1 LIMIT 1)", "Арабский язык")
		if err != nil {
			panic(err)
		}
		_ = tmp.Scan(&exist)

		if exist {
			arabic, err := db.Query("SELECT * FROM teachers WHERE language=$1", "Арабский язык")
			if arabic != nil {
				for arabic.Next() {
					var tmp Teacher
					err = arabic.Scan(&tmp.Id, &tmp.Branch, &tmp.Name, &tmp.Surname, &tmp.Patronymic, &tmp.Language, &tmp.Salary, &tmp.Experience, &tmp.Login)
					if err != nil {
						panic(err)
					}
					languages.Arabic = append(languages.Arabic, tmp)
				}
			}
		}

		contract, err := template.ParseFiles("templates/contract.html", "templates/footer.html", "templates/header.html")
		err = contract.ExecuteTemplate(writer, "contract", languages)

		if err != nil {
			panic(err)
		}
	}

	parameters.AuthorizationMask = true
	parameters.Authorization = false
	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)

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

	id := db.QueryRow("SELECT id FROM clients WHERE login = $1", currentUser.Login)
	var idClient string
	err = id.Scan(&idClient)
	if err != nil {
		panic(err)
	}

	var loginTeacher string
	switch result["contract-lang"][0] {
	case "Английский язык":
		loginTeacher = result["contract-english"][0]
	case "Немецкий язык":
		loginTeacher = result["contract-germany"][0]
	case "Французский язык":
		loginTeacher = result["contract-french"][0]
	case "Испанский язык":
		loginTeacher = result["contract-spanish"][0]
	case "Китайский язык":
		loginTeacher = result["contract-china"][0]
	case "Японский язык":
		loginTeacher = result["contract-japan"][0]
	case "Хинди":
		loginTeacher = result["contract-hindi"][0]
	case "Иврит":
		loginTeacher = result["contract-hebrew"][0]
	case "Казахский язык":
		loginTeacher = result["contract-kazakh"][0]
	case "Чувашский язык":
		loginTeacher = result["contract-chuvash"][0]
	case "Турецкий язык":
		loginTeacher = result["contract-turkish"][0]
	case "Арабский язык":
		loginTeacher = result["contract-arabic"][0]

	}

	id = db.QueryRow("SELECT id FROM teachers WHERE login = $1", loginTeacher)

	var idTeacher string
	err = id.Scan(&idTeacher)
	if err != nil {
		panic(err)
	}
	contractPrice := strings.Trim(result["contract-price"][0], "₽")

	_, err = db.Exec("INSERT INTO contract (client,teacher,language,quantity,price) VALUES ($1,$2,$3,$4,$5)",
		idClient, idTeacher, result["contract-lang"][0], result["contract-month"][0], contractPrice)
	if err != nil {
		panic(err)
	}

	parameters.successContract = true
	parameters.successContractMask = true

	http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
}

func alertPage(writer http.ResponseWriter, _ *http.Request) {

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

	if parameters.successContractMask {
		if parameters.successContract {
			parameters.Message = "Успешная запись на курс"
		} else {
			parameters.Message = "Не удалось записаться на курс"
		}
	}

	if parameters.checkStudentMask {
		if !parameters.checkStudent {
			parameters.Message = "Записаться на курс можно только учеником"
		}
	}

	if parameters.AuthorizationMask {
		if !parameters.Authorization {
			parameters.Message = "Сначала нужно авторизоваться в системе"
		} else {
			parameters.Message = "Успешная авторизация в системе"
		}
	}

	if parameters.TeacherCabinetMask {
		if !parameters.TeacherCabinet {
			parameters.Message = "У учителя нет учеников"
		}
	}

	if parameters.StudentCabinetMask {
		if !parameters.StudentCabinet {
			parameters.Message = "У ученика нет записей на курсы"
		}
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

	teacherInfo = TeacherInfo{}

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
	err = id.Scan(&idTeacher)
	if err != nil {
		panic(err)
	}

	tmp := db.QueryRow("SELECT EXISTS(SELECT * FROM contract WHERE teacher = $1 LIMIT 1)", idTeacher)
	var exist bool
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	var idBranch uint16
	if exist {
		res, err := db.Query("SELECT contract.id,contract.language,contract.quantity,contract.price,contract.date,contract.status, clients.name, clients.surname, clients.patronymic, clients.branch, clients.phone FROM contract LEFT JOIN clients ON contract.client = clients.id WHERE contract.teacher = $1 ORDER BY contract.id", idTeacher)
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var tmp ContractsAndClients
			err = res.Scan(&tmp.Id, &tmp.Language, &tmp.Quantity, &tmp.Price, &tmp.Date, &tmp.Status, &tmp.Name,
				&tmp.Surname, &tmp.Patronymic, &idBranch, &tmp.Phone)
			if err != nil {
				panic(err)
			}
			address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
			err = address.Scan(&tmp.Branch)
			if err != nil {
				panic(err)
			}

			teacherInfo.ContractsAndClients = append(teacherInfo.ContractsAndClients, tmp)
		}

		teacher := db.QueryRow("SELECT * FROM teachers WHERE id = $1", idTeacher)
		err = teacher.Scan(&teacherInfo.Main.Id, &idBranch, &teacherInfo.Main.Name, &teacherInfo.Main.Surname, &teacherInfo.Main.Patronymic,
			&teacherInfo.Main.Language, &teacherInfo.Main.Salary, &teacherInfo.Main.Experience, &teacherInfo.Main.Login)
		if err != nil {
			panic(err)
		}

		address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
		err = address.Scan(&teacherInfo.Main.Branch)
		if err != nil {
			panic(err)
		}

		parameters.TeacherCabinet = true
		parameters.TeacherCabinetMask = true

	} else {
		parameters.TeacherCabinet = false
		parameters.TeacherCabinetMask = true
		http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
	}

	teacher, err := template.ParseFiles("templates/teacher.html")
	err = teacher.Execute(writer, teacherInfo)

	if err != nil {
		panic(err)
	}
}

func directorCabinet(writer http.ResponseWriter, _ *http.Request) {

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

	directorInfo = DirectorInfo{}

	tmp := db.QueryRow("SELECT EXISTS(SELECT * FROM contract LIMIT 1)")
	var exist bool
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM contract ORDER BY id")
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var contract Contract
			err = res.Scan(&contract.Id, &contract.Client, &contract.Teacher, &contract.Language, &contract.Quantity,
				&contract.Price, &contract.Date, &contract.Status)
			if err != nil {
				panic(err)
			}
			contract.Date = contract.Date[:10]
			directorInfo.TableContracts = append(directorInfo.TableContracts, contract)
		}
	}

	tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM clients LIMIT 1)")
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM clients ORDER BY id")
		if err != nil {
			panic(err)
		}
		var idBranch int
		for res.Next() {
			var client Client
			err = res.Scan(&client.Id, &client.Name, &client.Surname, &client.Patronymic, &idBranch,
				&client.Phone, &client.Login)
			if err != nil {
				panic(err)
			}

			address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
			err = address.Scan(&client.Branch)
			if err != nil {
				panic(err)
			}

			directorInfo.TableClients = append(directorInfo.TableClients, client)
		}
	}

	tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers LIMIT 1)")
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM teachers ORDER BY id")
		if err != nil {
			panic(err)
		}
		var idBranch int
		for res.Next() {
			var teacher Teacher
			err = res.Scan(&teacher.Id, &idBranch, &teacher.Name, &teacher.Surname, &teacher.Patronymic,
				&teacher.Language, &teacher.Salary, &teacher.Experience, &teacher.Login)
			if err != nil {
				panic(err)
			}

			address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
			err = address.Scan(&teacher.Branch)
			if err != nil {
				panic(err)
			}

			directorInfo.TableTeachers = append(directorInfo.TableTeachers, teacher)
		}
	}

	tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM branch LIMIT 1)")
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM branch ORDER BY id")
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var branch Branch
			err = res.Scan(&branch.Id, &branch.Address, &branch.Name, &branch.Surname, &branch.Patronymic,
				&branch.Salary, &branch.Login)
			if err != nil {
				panic(err)
			}

			directorInfo.TableBranches = append(directorInfo.TableBranches, branch)
		}
	}

	tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM logins LIMIT 1)")
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM logins")
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var login Login
			err = res.Scan(&login.Email, &login.Password, &login.Login, &login.Status)
			if err != nil {
				panic(err)
			}

			directorInfo.TableLogins = append(directorInfo.TableLogins, login)
		}
	}

	director, err := template.ParseFiles("templates/director.html")
	err = director.Execute(writer, directorInfo)

}

func editDirector(writer http.ResponseWriter, request *http.Request) {

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

	fmt.Println(request.URL.Query())

	http.Redirect(writer, request, "/director/", http.StatusSeeOther)
}

func adminCabinet(writer http.ResponseWriter, _ *http.Request) {

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

	adminInfo = AdminInfo{}

	tmp := db.QueryRow("SELECT EXISTS(SELECT * FROM contract LIMIT 1)")
	var exist bool
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM contract ORDER BY id")
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var contract Contract
			err = res.Scan(&contract.Id, &contract.Client, &contract.Teacher, &contract.Language, &contract.Quantity,
				&contract.Price, &contract.Date, &contract.Status)
			if err != nil {
				panic(err)
			}
			adminInfo.TableContracts = append(adminInfo.TableContracts, contract)
		}
	}

	tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM clients LIMIT 1)")
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM clients ORDER BY id")
		if err != nil {
			panic(err)
		}
		var idBranch int
		for res.Next() {
			var client Client
			err = res.Scan(&client.Id, &client.Name, &client.Surname, &client.Patronymic, &idBranch,
				&client.Phone, &client.Login)
			if err != nil {
				panic(err)
			}

			address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
			err = address.Scan(&client.Branch)
			if err != nil {
				panic(err)
			}

			adminInfo.TableClients = append(adminInfo.TableClients, client)
		}
	}

	tmp = db.QueryRow("SELECT EXISTS(SELECT * FROM teachers LIMIT 1)")
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {
		res, err := db.Query("SELECT * FROM teachers ORDER BY id")
		if err != nil {
			panic(err)
		}
		var idBranch int
		for res.Next() {
			var teacher Teacher
			err = res.Scan(&teacher.Id, &idBranch, &teacher.Name, &teacher.Surname, &teacher.Patronymic,
				&teacher.Language, &teacher.Salary, &teacher.Experience, &teacher.Login)
			if err != nil {
				panic(err)
			}

			address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
			err = address.Scan(&teacher.Branch)
			if err != nil {
				panic(err)
			}

			adminInfo.TableTeachers = append(adminInfo.TableTeachers, teacher)
		}
	}

	admin, err := template.ParseFiles("templates/admin.html")
	err = admin.Execute(writer, adminInfo)
}

func editAdmin(writer http.ResponseWriter, request *http.Request) {

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
	for key, value := range result {
		words := strings.Split(key, " ")
		if words[0] == "status" {
			tmp := strings.Split(value[0], "_")
			_, err := db.Exec("UPDATE contract SET status = $1 WHERE id = $2",
				tmp[0], tmp[1])
			if err != nil {
				panic(err)
			}
		}
		if words[0] == "salary" {
			_, err := db.Exec("UPDATE teachers SET salary = $1 WHERE id = $2",
				value[0], words[1])
			if err != nil {
				panic(err)
			}
		}
	}

	http.Redirect(writer, request, "/admin/", http.StatusSeeOther)
}

func studentCabinet(writer http.ResponseWriter, request *http.Request) {

	clientInfo = ClientInfo{}

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

	id := db.QueryRow("SELECT id,branch FROM clients WHERE login = $1", currentUser.Login)
	var idStudent uint16
	var idBranch uint16
	err = id.Scan(&idStudent, &idBranch)
	if err != nil {
		panic(err)
	}

	tmp := db.QueryRow("SELECT EXISTS(SELECT * FROM contract WHERE client = $1 LIMIT 1)", idStudent)
	var exist bool
	err = tmp.Scan(&exist)
	if err != nil {
		panic(err)
	}

	if exist {

		res, err := db.Query("SELECT * FROM contract WHERE client = $1 ORDER BY id", idStudent)
		if err != nil {
			panic(err)
		}

		for res.Next() {
			var idTeacher uint16
			var tmp ContractsAndTeachers
			err = res.Scan(&tmp.Id, &tmp.Client, &idTeacher,
				&tmp.Language, &tmp.Quantity, &tmp.Price,
				&tmp.Date, &tmp.Status)
			if err != nil {
				panic(err)
			}
			teacher := db.QueryRow("SELECT name,surname,patronymic FROM teachers WHERE id = $1",
				idTeacher)
			err = teacher.Scan(&tmp.NameTeacher, &tmp.SurnameTeacher, &tmp.PatronymicTeacher)
			if err != nil {
				panic(err)
			}

			address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
			err = address.Scan(&tmp.Branch)
			if err != nil {
				panic(err)
			}

			clientInfo.ContractsAndTeachers = append(clientInfo.ContractsAndTeachers, tmp)
		}

		client := db.QueryRow("SELECT * FROM clients WHERE id = $1", idStudent)
		err = client.Scan(&clientInfo.Main.Id, &clientInfo.Main.Name, &clientInfo.Main.Surname,
			&clientInfo.Main.Patronymic, &clientInfo.Main.Branch, &clientInfo.Main.Phone, &clientInfo.Main.Login)
		if err != nil {
			panic(err)
		}

		address := db.QueryRow("SELECT address FROM branch WHERE id = $1", idBranch)
		err = address.Scan(&clientInfo.Main.Branch)
		if err != nil {
			panic(err)
		}

		parameters.StudentCabinet = true
		parameters.StudentCabinetMask = true
	} else {
		parameters.StudentCabinet = false
		parameters.StudentCabinetMask = true
		http.Redirect(writer, request, "/alert/", http.StatusSeeOther)
	}

	student, err := template.ParseFiles("templates/student.html")
	err = student.Execute(writer, clientInfo)

	if err != nil {
		panic(err)
	}
}

func handlerRequest() {
	router := mux.NewRouter()
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

	router.HandleFunc("/editAdmin/", editAdmin)
	router.HandleFunc("/editDirector/", editDirector)
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
