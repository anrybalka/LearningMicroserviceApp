package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	userStore  = make(map[string]string) // Хранилище пользователей (username:hashed_password)
	tokenStore = make(map[string]string) // Хранилище токенов (token:username)
	storeMutex = sync.RWMutex{}
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Структура для логирования
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// Генерация случайного токена
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Отправка лога на сервис с префиксом "000-user-service"
func sendLog(message string) {
	// Добавляем префикс "000-user-service" к каждому сообщению
	message = "000-user-service: " + message

	logEntry := LogEntry{
		Timestamp: time.Now(),
		Message:   message,
	}

	logData, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Println("Error marshalling log:", err)
		return
	}

	// Отправляем POST-запрос на лог-сервис
	_, err = http.Post("http://001-log-service:5434/log", "application/json", bytes.NewBuffer(logData))
	if err != nil {
		fmt.Println("Error sending log:", err)
	}
}

// Регистрация пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storeMutex.Lock()
	defer storeMutex.Unlock()

	if _, exists := userStore[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	userStore[user.Username] = user.Password                   // Хранение пароля в открытом виде (небезопасно)
	sendLog(fmt.Sprintf("User registered: %s", user.Username)) // Логирование регистрации

	w.WriteHeader(http.StatusCreated)
}

// Аутентификация пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storeMutex.RLock()
	password, exists := userStore[user.Username]
	storeMutex.RUnlock()

	if !exists || password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		sendLog(fmt.Sprintf("Failed login attempt for user: %s", user.Username)) // Логирование неудачной попытки
		return
	}

	token, err := generateToken()
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	storeMutex.Lock()
	tokenStore[token] = user.Username
	storeMutex.Unlock()

	sendLog(fmt.Sprintf("User logged in: %s", user.Username)) // Логирование успешной авторизации

	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Проверка токена
func checkTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	storeMutex.RLock()
	username, exists := tokenStore[tokenStr]
	storeMutex.RUnlock()

	if !exists {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		sendLog(fmt.Sprintf("Invalid token attempt: %s", tokenStr)) // Логирование невалидного токена
		return
	}

	sendLog(fmt.Sprintf("Token is valid for user: %s", username)) // Логирование успешной проверки токена

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Token is valid for user: %s", username)
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/check-token", checkTokenHandler)

	fmt.Println("Starting server on :5432")
	if err := http.ListenAndServe(":5432", nil); err != nil {
		panic(err)
	}
}
