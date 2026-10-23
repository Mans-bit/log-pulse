package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type SystemStatus struct {
	Timestamp     string `json:"timestamp"`
	OS            string `json:"os"`
	Architecture  string `json:"architecture"`
	NumCPU        int    `json:"num_cpu"`
	NumGoroutines int    `json:"active_goroutines"`
	AllocMemory   MB     `json:"allocated_memory_mb"`
	Status        string `json:"status"`
}

type MB float64

const htmlPage = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Log-Pulse Панель</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #0f172a; color: #f8fafc; margin: 0; padding: 40px; display: flex; justify-content: center; }
        .card { background-color: #1e293b; border-radius: 16px; padding: 30px; width: 450px; box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3); border: 1px solid #334155; }
        h1 { margin-top: 0; font-size: 24px; color: #38bdf8; border-bottom: 2px solid #334155; padding-bottom: 15px; display: flex; justify-content: space-between; align-items: center; }
        .status-badge { font-size: 14px; background-color: #10b981; color: white; padding: 4px 12px; border-radius: 9999px; animation: pulse 2s infinite; }
        .metric { display: flex; justify-content: space-between; margin: 18px 0; font-size: 16px; }
        .label { color: #94a3b8; }
        .value { font-weight: 600; color: #e2e8f0; }
        @keyframes pulse { 0% { opacity: 0.6; } 50% { opacity: 1; } 100% { opacity: 0.6; } }
    </style>
</head>
<body>
    <div class="card">
        <h1>Log-Pulse <span class="status-badge">{{.Status}}</span></h1>
        <div class="metric"><span class="label">Операционная система:</span><span class="value">{{.OS}} ({{.Architecture}})</span></div>
        <div class="metric"><span class="label">Доступно ядер CPU:</span><span class="value">{{.NumCPU}}</span></div>
        <div class="metric"><span class="label">Активные горутины:</span><span class="value">{{.NumGoroutines}}</span></div>
        <div class="metric"><span class="label">Выделено памяти:</span><span class="value">{{printf "%.2f" .AllocMemory}} MB</span></div>
        <div class="metric"><span class="label">Время обновления:</span><span class="value">{{.Timestamp}}</span></div>
    </div>
</body>
</html>
`

func getStatus() SystemStatus {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return SystemStatus{
		Timestamp:     time.Now().Format("15:04:05"),
		OS:            runtime.GOOS,
		Architecture:  runtime.GOARCH,
		NumCPU:        runtime.NumCPU(),
		NumGoroutines: runtime.NumGoroutine(),
		AllocMemory:   MB(m.Alloc) / 1024 / 1024,
		Status:        "ONLINE",
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getStatus())
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("ui").Parse(htmlPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, getStatus())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", uiHandler)
	http.HandleFunc("/api/v1/status", apiHandler)

	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
