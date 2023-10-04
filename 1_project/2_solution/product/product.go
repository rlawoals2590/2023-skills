package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// Database Connect
var readDB *sql.DB
var writeDB *sql.DB

// User Insert and Select
type User struct {
	ID   int
	Name string
}

// DB config
type DBConfig struct {
	Read_Host  string
	Write_Host string
	Port       int
	User       string
	Password   string
	Name       string
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func GetDBEnv() *DBConfig {
	return &DBConfig{
		Read_Host:  os.Getenv("DB_READ_HOST"),
		Write_Host: os.Getenv("DB_WRITE_HOST"),
		Port:       getEnvAsInt("DB_PORT", 3306),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}
}

// Main and User Handler
func main() {
	var err error

	config := GetDBEnv()
	dbInfo := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"

	readDB, err = sql.Open("mysql", fmt.Sprintf(dbInfo, config.User, config.Password, config.Read_Host, config.Port, config.Name))
	if err != nil {
		log.Fatalf("Error initializing Read database: %v", err)
	}

	writeDB, err = sql.Open("mysql", fmt.Sprintf(dbInfo, config.User, config.Password, config.Write_Host, config.Port, config.Name))
	if err != nil {
		log.Fatalf("Error initializing Write database: %v", err)
	}

	err = writeDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	http.HandleFunc("/v1/users", userGet)
	http.HandleFunc("/v1/registry", userRegistry)
	http.HandleFunc("/healthz", health)
	http.ListenAndServe(":8080", nil)
}

func CreateUser(name string) error {
	_, err := writeDB.Exec("INSERT INTO users(name) VALUES(?)", name)
	return err
}

func GetAllUsers() ([]User, error) {
	rows, err := readDB.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func userRegistry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	// Access form values
	name := r.FormValue("userName")

	if name == "" {
		http.Error(w, "'name' is required", http.StatusBadRequest)
		return
	}

	err = CreateUser(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	// Display an alert and redirect to '/'
	w.Write([]byte(`<script>alert("User successfully created!"); window.location.href="/";</script>`))
}

func userGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	json.NewEncoder(w).Encode(users)
}

// Health check
func health(w http.ResponseWriter, req *http.Request) {
	response := map[string]string{"status": "ok"}

	// 응답을 JSON으로 마샬링
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	// JSON 응답을 보냅니다.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}