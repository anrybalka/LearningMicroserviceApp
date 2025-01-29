package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	userServiceURL    = "http://000-user-service:5432"
	logServiceURL     = "http://001-log-service:5434/log"
	robotUsernameBase = "robot_user_"
	robotPassword     = "robot_password"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Функция для отправки логов в сервис
func sendLog(message string) {
	// Добавляем префикс "002-robot-tester" к каждому сообщению
	message = "002-robot-tester: " + message

	logEntry := LogEntry{
		Timestamp: time.Now(),
		Message:   message,
	}

	logData, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Println("Error marshalling log:", err)
		return
	}

	_, err = http.Post(logServiceURL, "application/json", bytes.NewBuffer(logData))
	if err != nil {
		fmt.Println("Error sending log:", err)
	}
}

// Функция для регистрации пользователя
func registerUser(username string) {
	// Создание нового пользователя
	user := map[string]string{
		"username": username,
		"password": robotPassword,
	}

	// Отправка POST-запроса для регистрации
	loginJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshalling registration data:", err)
		return
	}

	resp, err := http.Post(userServiceURL+"/register", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		fmt.Println("Error during registration:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		sendLog(fmt.Sprintf("User registered: %s", username)) // Логирование успешной регистрации
	} else {
		sendLog(fmt.Sprintf("Registration failed for user: %s", username)) // Логирование неудачной регистрации
	}
}

// Функция для авторизации пользователя
func loginUser(username string) string {
	// Создание запроса для авторизации
	loginData := map[string]string{
		"username": username,
		"password": robotPassword,
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println("Error marshalling login data:", err)
		return ""
	}

	resp, err := http.Post(userServiceURL+"/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		fmt.Println("Error during login:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var loginResponse LoginResponse
		err := json.NewDecoder(resp.Body).Decode(&loginResponse)
		if err != nil {
			fmt.Println("Error decoding login response:", err)
			return ""
		}

		sendLog(fmt.Sprintf("User logged in: %s, Token: %s", username, loginResponse.Token)) // Логирование успешного входа
		return loginResponse.Token
	} else {
		sendLog(fmt.Sprintf("Login failed for user: %s", username)) // Логирование неудачной попытки входа
		return ""
	}
}

// Функция для проверки токена
func checkToken(token string) {
	req, err := http.NewRequest("GET", userServiceURL+"/check-token", nil)
	if err != nil {
		fmt.Println("Error creating check token request:", err)
		return
	}

	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error during token check:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		sendLog(fmt.Sprintf("Token is valid: %s", token)) // Логирование успешной проверки токена
	} else {
		sendLog(fmt.Sprintf("Token is invalid: %s", token)) // Логирование невалидного токена
	}
}

func main() {
	// Каждую минуту выполняем регистрацию, авторизацию и проверку токена
	for {
		// Генерируем случайное имя пользователя для регистрации
		username := robotUsernameBase + strconv.Itoa(rand.Intn(10000))

		// Регистрация пользователя
		registerUser(username)

		// Попытка авторизации
		token := loginUser(username)

		if token != "" {
			// Проверка токена
			checkToken(token)
		}

		// Ожидаем 5 секунд перед следующей итерацией
		time.Sleep(5 * time.Second)
	}
}
