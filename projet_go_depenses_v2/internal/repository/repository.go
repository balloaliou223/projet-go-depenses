package repository

import (
	"database/sql"
	"depenses-personnelles/internal/models"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type ExpenseRepository struct{ db *sql.DB }

func NewExpenseRepository(db *sql.DB) *ExpenseRepository { return &ExpenseRepository{db: db} }
func Migrate(db *sql.DB) error {
	b, err := os.ReadFile("internal/repository/migrations/001_create_tables.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b))
	return err
}
func SeedCategories(db *sql.DB) error {
	for _, n := range []string{"Alimentation", "Transport", "Logement", "Loisirs", "Santé", "Autre"} {
		if _, err := db.Exec("INSERT OR IGNORE INTO categories(name) VALUES (?)", n); err != nil {
			return err
		}
	}
	return nil
}
func validate(e models.Expense) error {
	if e.Amount <= 0 {
		return errors.New("la somme doit être supérieure à zéro")
	}
	if strings.TrimSpace(e.Description) == "" {
		return errors.New("la description est obligatoire")
	}
	if _, err := time.Parse("2006-01-02", e.Date); err != nil {
		return errors.New("la date doit être au format AAAA-MM-JJ")
	}
	if e.CategoryID <= 0 {
		return errors.New("veuillez choisir une catégorie")
	}
	return nil
}
func (r *ExpenseRepository) Create(e models.Expense) error {
	if err := validate(e); err != nil {
		return err
	}
	_, err := r.db.Exec("INSERT INTO expenses(amount,description,date,category_id) VALUES(?,?,?,?)", e.Amount, e.Description, e.Date, e.CategoryID)
	return err
}

const selectExpenses = "SELECT e.id,e.amount,e.description,e.date,e.category_id,c.name FROM expenses e JOIN categories c ON c.id=e.category_id"

func scanExpenses(rows *sql.Rows) ([]models.Expense, error) {
	defer rows.Close()
	var out []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query(selectExpenses + " ORDER BY e.date DESC,e.id DESC")
	if err != nil {
		return nil, err
	}
	return scanExpenses(rows)
}
func (r *ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(selectExpenses+" WHERE e.id=?", id).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName)
	return e, err
}
func (r *ExpenseRepository) Update(e models.Expense) error {
	if err := validate(e); err != nil {
		return err
	}
	res, err := r.db.Exec("UPDATE expenses SET amount=?,description=?,date=?,category_id=? WHERE id=?", e.Amount, e.Description, e.Date, e.CategoryID, e.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}
func (r *ExpenseRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM expenses WHERE id=?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return err
}
func (r *ExpenseRepository) Filter(cat int, from, to string) ([]models.Expense, error) {
	q := selectExpenses
	var where []string
	var args []any
	if cat > 0 {
		where = append(where, "e.category_id=?")
		args = append(args, cat)
	}
	if from != "" {
		where = append(where, "e.date>=?")
		args = append(args, from)
	}
	if to != "" {
		where = append(where, "e.date<=?")
		args = append(args, to)
	}
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	q += " ORDER BY e.date DESC,e.id DESC"
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	return scanExpenses(rows)
}

type CategoryRepository struct{ db *sql.DB }

func NewCategoryRepository(db *sql.DB) *CategoryRepository { return &CategoryRepository{db: db} }
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id,name FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

var _ = fmt.Sprintf
