package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "api_transactions_created_total",
		Help: "The total number of created transactions",
	})
)

func main() {
	prometheus.MustRegister(opsProcessed)

	// Setup real DB
	psqlInfo := buildConnString()
	sqlDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Real DB implementation
	real := &realDB{db: sqlDB}
	if err := real.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to Postgres successfully!")

	real.initializeDBTable()

	// Assign global var db to the real DB
	db = real

	http.HandleFunc("/api/transaction/", createTransactionHandler)
	http.Handle("/metrics", promhttp.Handler())

	port := getEnv("PORT", "8080")
	log.Printf("API Service is listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Only POST is allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var t Transaction
	if err := decodeJSON(r, &t); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if t.TransactionID == "" {
		uuidVal, _ := uuid.NewUUID()
		t.TransactionID = uuidVal.String()
	}
	if t.Timestamp == "" {
		t.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	stmt := `INSERT INTO transactions (transaction_id, amount, timestamp) VALUES ($1, $2, $3)`
	err := db.Exec(stmt, t.TransactionID, t.Amount, t.Timestamp)
	if err != nil {
		log.Printf("Insert failed: %v", err)
		http.Error(w, `{"error":"Insert failed"}`, http.StatusInternalServerError)
		return
	}

	opsProcessed.Inc()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Transaction created successfully"}`))
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type Transaction struct {
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount,string"`
	Timestamp     string  `json:"timestamp"`
}

func buildConnString() string {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "transactions_db")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
