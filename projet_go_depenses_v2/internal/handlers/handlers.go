package handlers

import (
	"database/sql"
	"depenses-personnelles/internal/models"
	"depenses-personnelles/internal/repository"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	expenses   *repository.ExpenseRepository
	categories *repository.CategoryRepository
	tmpl       *template.Template
}
type PageData struct {
	Title                               string
	Expenses                            []models.Expense
	Categories                          []models.Category
	Error                               string
	Amount, Description, Date, From, To string
	CategoryID, FilterCat               int
}

func New(e *repository.ExpenseRepository, c *repository.CategoryRepository) (*App, error) {
	t, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}
	return &App{e, c, t}, nil
}
func (a *App) render(w http.ResponseWriter, n string, d PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.tmpl.ExecuteTemplate(w, n, d); err != nil {
		http.Error(w, "Erreur interne du serveur", 500)
	}
}
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	a.render(w, "home.html", PageData{Title: "Accueil"})
}
func (a *App) About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("Учёт личных трат — application web de gestion des dépenses personnelles."))
}
func (a *App) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Méthode non autorisée : utilisez GET", 405)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("pong"))
}
func (a *App) Expenses(w http.ResponseWriter, r *http.Request) {
	cats, err := a.categories.GetAll()
	if err != nil {
		http.Error(w, "Erreur catégories", 500)
		return
	}
	cat, _ := strconv.Atoi(r.URL.Query().Get("category"))
	from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")
	items, err := a.expenses.Filter(cat, from, to)
	if err != nil {
		http.Error(w, "Erreur dépenses", 500)
		return
	}
	a.render(w, "expenses.html", PageData{Title: "Mes dépenses", Expenses: items, Categories: cats, FilterCat: cat, From: from, To: to})
}
func (a *App) NewExpenseForm(w http.ResponseWriter, r *http.Request) {
	cats, err := a.categories.GetAll()
	if err != nil {
		http.Error(w, "Erreur catégories", 500)
		return
	}
	a.render(w, "new_expense.html", PageData{Title: "Ajouter une dépense", Categories: cats})
}
func (a *App) CreateExpense(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", 400)
		return
	}
	amt, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64)
	cat, ce := strconv.Atoi(r.FormValue("category_id"))
	desc, date := strings.TrimSpace(r.FormValue("description")), strings.TrimSpace(r.FormValue("date"))
	if err == nil && ce == nil {
		err = a.expenses.Create(models.Expense{Amount: amt, Description: desc, Date: date, CategoryID: cat})
	} else {
		err = sql.ErrNoRows
	}
	if err != nil {
		cats, _ := a.categories.GetAll()
		msg := "Montant ou catégorie invalide."
		if err != sql.ErrNoRows {
			msg = err.Error()
		}
		a.render(w, "new_expense.html", PageData{Title: "Ajouter une dépense", Categories: cats, Error: msg, Amount: r.FormValue("amount"), Description: desc, Date: date, CategoryID: cat})
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
func (a *App) EditExpenseForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	e, err := a.expenses.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	cats, err := a.categories.GetAll()
	if err != nil {
		http.Error(w, "Erreur catégories", 500)
		return
	}
	a.render(w, "edit_expense.html", PageData{Title: "Modifier une dépense", Categories: cats, Amount: strconv.FormatFloat(e.Amount, 'f', 2, 64), Description: e.Description, Date: e.Date, CategoryID: e.CategoryID})
}
func (a *App) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if err = r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", 400)
		return
	}
	amt, pe := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64)
	cat, ce := strconv.Atoi(r.FormValue("category_id"))
	desc, date := strings.TrimSpace(r.FormValue("description")), strings.TrimSpace(r.FormValue("date"))
	if pe != nil || ce != nil {
		a.editError(w, amt, desc, date, cat, "Montant ou catégorie invalide.")
		return
	}
	err = a.expenses.Update(models.Expense{ID: id, Amount: amt, Description: desc, Date: date, CategoryID: cat})
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		a.editError(w, amt, desc, date, cat, err.Error())
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
func (a *App) editError(w http.ResponseWriter, amt float64, desc, date string, cat int, msg string) {
	cats, _ := a.categories.GetAll()
	a.render(w, "edit_expense.html", PageData{Title: "Modifier une dépense", Categories: cats, Amount: strconv.FormatFloat(amt, 'f', -1, 64), Description: desc, Date: date, CategoryID: cat, Error: msg})
}
func (a *App) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if err = a.expenses.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Suppression impossible", 500)
		}
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
