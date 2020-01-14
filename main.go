package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// const (
// 	host = "localhost"
// 	port = 5432
// 	user = "postgres"
// 	password = "andres1993"

// )
type User struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

type PostGres struct {
	DB *sql.DB
}

const (
	dbhost     = "localhost"
	dbport     = "5432"
	dbuser     = "postgres"
	dbname     = "postgres"
	dbpassword = ""
)

func (p *PostGres) GetUserMethod(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	users := []User{}
	// n := 5
	// for index := 0; index < n; index++ {
	// 	fmt.Printf("JulioVargs0%v", index)
	// 	users = append(users, User{FirstName: "Julio", LastName: "Vargas", Email: "juliovargs0%v/n@place.com"})
	// }
	m := make(map[string][]User)
	rows, err := p.DB.Query("select * from users")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber); err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
		users = append(users, user)
	}
	m["users"] = users
	json.NewEncoder(response).Encode(m)

}
func (p *PostGres) AddUserMethod(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	postData := getDataFromRequest(request)
	sqlStatement := `
	INSERT INTO users(FirstName, LastName, Email, PhoneNumber)
	VALUES ($1, $2, $3, $4)`
	add, err := p.DB.Exec(sqlStatement, postData.FirstName, postData.LastName, postData.Email, postData.PhoneNumber)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
	m := make(map[string]string)
	m["status"] = "Systems are go"
	json.NewEncoder(response).Encode(m)
}
func (p *PostGres) DeleteUserMethod(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	n := 5
	for index := 0; index < n; index++ {
		fmt.Println("JulioVargs" + "%v")
	}
	response.WriteHeader(http.StatusOK)
}
func (p *PostGres) UpdateUserMethod(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	updateData := getDataFromRequest(request)
	fmt.Println(updateData)
	updateQuery := `UPDATE users SET PhoneNumber = $1 WHERE PhoneNumber = $2;`
	val, err := p.DB.Exec(updateQuery, "66666666666666", updateData.PhoneNumber)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val.RowsAffected())
	response.WriteHeader(http.StatusOK)
}
func getDataFromRequest(req *http.Request) User {
	user := User{}
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}
	return user
}
func (p *PostGres) initDb() {
	// config := dbConfig()
	// psqlInfo := fmt.Sprintf("host=%s port=%s "+
	// 	"dbname=%s sslmode=disable",
	// 	config["host"], config["port"], config["name"])
	psqlInfo := "user=ivantrujillo dbname=postgres host=localhost port=5432 sslmode=disable"
	pgconnect, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	err = pgconnect.Ping()
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	p.DB = pgconnect
	fmt.Println("DB connected")
}
func dbConfig() map[string]string {
	conf := make(map[string]string)
	conf["host"] = dbhost
	conf["port"] = dbport
	// conf["user"] = dbuser
	conf["name"] = dbname
	// conf["password"] = dbpassword
	return conf
}
func main() {
	db := PostGres{}
	fmt.Println("This is the program")
	router := mux.NewRouter()
	db.initDb()
	router.HandleFunc("/get-user", db.GetUserMethod).Methods("GET")
	router.HandleFunc("/add-user", db.AddUserMethod).Methods("POST")
	router.HandleFunc("/delete-user", db.DeleteUserMethod).Methods("Update")
	router.HandleFunc("/edit-user", db.UpdateUserMethod).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
