package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

var (
	logs      []LogEntry
	logsMutex sync.Mutex
)

func main() {
	http.HandleFunc("/log", logHandler)   // Обработчик для POST-запроса (для добавления логов)
	http.HandleFunc("/getlogs", getLogsHandler) // Новый обработчик для GET-запроса (для отображения логов)

	http.ListenAndServe(":5434", nil) // Запуск HTTP-сервера на порту 5434
}

// Обработчик для POST-запроса, чтобы добавить новые логи
func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var logEntry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	logEntry.Timestamp = time.Now()

	logsMutex.Lock()
	logs = append(logs, logEntry) // Добавление лога в список
	logsMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// Новый обработчик для GET-запроса, чтобы вывести логи
func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	logsMutex.Lock()
	defer logsMutex.Unlock()

	// Шаблон для отображения логов
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Log Service</title>
    <meta http-equiv="refresh" content="1"> <!-- Обновление страницы каждую секунду -->
</head>
<body>
    <h1>Logs</h1>
    <ul>
        {{range .}}
        <li>[{{.Timestamp.Format "2006-01-02 15:04:05"}}] {{.Message}}</li>
        {{end}}
    </ul>
</body>
</html>
`

	// Парсинг шаблона
	t, err := template.New("logs").Parse(tmpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка логов на страницу
	if err := t.Execute(w, logs); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}