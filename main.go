package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"text/template"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type User struct {
	gorm.Model
	FirstName   string
	LastName    string
	MiddleName  string
	PhoneNumber string
	Email       string
	Username    string
	Password    string
}

func main() {
	// Инициализация базы данных SQLite
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})

	// Инициализация роутера
	r := mux.NewRouter()

	// CRUD Роуты
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", GetUser).Methods("GET")
	r.HandleFunc("/user", CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")
	r.HandleFunc("/welcome", Welcome).Methods("POST")

	// Страница поиска
	r.HandleFunc("/search", SearchUsers).Methods("GET")

	// Страница авторизации
	r.HandleFunc("/login", Login).Methods("POST")

	// Запуск веб-сервера
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Вспомогательная функция для отправки электронных писем через Gmail
func sendWelcomeEmail(email string) error {
	from := "your-email@gmail.com" // Ваш адрес электронной почты на Gmail
	password := "your-password"    // Ваш пароль

	msg := "Subject: Добро пожаловать\n\nДобро пожаловать в наше приложение!"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

// Обработчики CRUD операций
func GetUsers(w http.ResponseWriter, r *http.Request) {
}
func GetUser(w http.ResponseWriter, r *http.Request) {
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
}

// Обработчик страницы поиска
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	// Получение параметров поиска из запроса
	query := r.URL.Query().Get("q")

	// Поиск в базе данных
	var users []User
	db.Where("FirstName LIKE ? OR LastName LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users)

	// Отображение результатов
	renderTemplate(w, "search.html", users)
}

// Обработчик страницы авторизации
func Login(w http.ResponseWriter, r *http.Request) {
	// Обработка данных формы авторизации
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Проверка логина и пароля (просто для примера, не используйте в реальном приложении)
	if username == "admin" && password == "admin" {
		// В случае успешной авторизации отправляем приветственное письмо на почту
		err := sendWelcomeEmail("admin@example.com")
		if err != nil {
			http.Error(w, "Failed to send welcome email", http.StatusInternalServerError)
			return
		}

		// Редирект на страницу поиска
		http.Redirect(w, r, "/search?q=admin", http.StatusSeeOther)
	} else {
		// В случае неудачной авторизации
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
	}
}

// Вспомогательная функция для отображения HTML-шаблона
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl, err := template.New("").ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

// Обработчик страницы приветствия
func Welcome(w http.ResponseWriter, r *http.Request) {
	// Получение адреса электронной почты из формы
	r.ParseForm()
	email := r.Form.Get("email")

	// Отправка приветственного письма
	err := sendWelcomeEmail(email)
	if err != nil {
		http.Error(w, "Failed to send welcome email", http.StatusInternalServerError)
		return
	}

	// Вывод сообщения об успешной отправке письма
	fmt.Fprintln(w, "Welcome email sent successfully!")
}
