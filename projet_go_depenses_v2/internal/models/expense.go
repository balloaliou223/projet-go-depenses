package models

type Category struct {
	ID   int
	Name string
}
type Expense struct {
	ID           int
	Amount       float64
	Description  string
	Date         string
	CategoryID   int
	CategoryName string
}
