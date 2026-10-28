package main

import (
	"database/sql"
	"depenses-personnelles/internal/handlers"
	"depenses-personnelles/internal/repository"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	if err = repository.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err = repository.SeedCategories(db); err != nil {
		log.Fatal(err)
	}
	app, err := handlers.New(repository.NewExpenseRepository(db), repository.NewCategoryRepository(db))
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.Home)
	mux.HandleFunc("GET /about", app.About)
	mux.HandleFunc("/ping", app.Ping)
	mux.HandleFunc("GET /expenses", app.Expenses)
	mux.HandleFunc("GET /expenses/new", app.NewExpenseForm)
	mux.HandleFunc("POST /expenses", app.CreateExpense)
	mux.HandleFunc("GET /expenses/{id}/edit", app.EditExpenseForm)
	mux.HandleFunc("POST /expenses/{id}/edit", app.UpdateExpense)
	mux.HandleFunc("POST /expenses/{id}/delete", app.DeleteExpense)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	srv := &http.Server{Addr: ":8080", Handler: logRequests(mux), ReadHeaderTimeout: 5 * time.Second}
	log.Println("Serveur démarré : http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
