package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var db *sql.DB

type URLschema struct {
	Id        int64
	LongURL   string
	CreatedAt time.Time
}

type ShortenRequest struct {
	LongURL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

type LongResponse struct {
	id      int64
	LongURL string
}

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(n int64) string {

	//in cse the n is 0 we need to speical case this.
	//otherwise it will loop infinite below
	if n == 0 {
		return string(base62Chars[0])
	}
	var result []byte

	//first get the remainder append it to the result
	//n = remainder is the stopper
	for n > 0 {
		remainder := n % 62
		result = append(result, base62Chars[remainder])
		n = n / 62
	}

	//reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func decodeBase62(ns string) int64 {

	var result int64

	for _, n := range ns {
		index := strings.IndexByte(base62Chars, byte(n))
		result = result*62 + int64(index)
	}

	return result

}

func initDB() {

	godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var errOpen error
	db, errOpen = sql.Open("pgx", connStr)
	if errOpen != nil {
		log.Fatal("failed to open database:", errOpen)
	}

	if errPing := db.Ping(); errPing != nil {
		log.Fatal("failed to connect to database:", errPing)
	}

	log.Println("connected to database successfully")
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	var id int64

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := db.QueryRow(
		"Insert Into urls (original_url) VALUES ($1) RETURNING id",
		req.LongURL,
	).Scan(&id); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortURL: encodeBase62(id),
		LongURL:  req.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {

	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	id := decodeBase62(shortCode)

	var longURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE id = $1", id).Scan(&longURL)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusNotFound)
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next(w, r)
	}
}

func main() {

	initDB()

	http.HandleFunc("/shorten", enableCORS(handleShorten))
	http.HandleFunc("/", enableCORS(handleRedirect))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
