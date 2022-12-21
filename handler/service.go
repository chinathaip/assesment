package handler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s Service) CreateTable() {
	createScript := `CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY,title TEXT,amount FLOAT,note TEXT,tags TEXT[]);`
	_, erro := s.DB.Exec(createScript)
	if erro != nil {
		fmt.Printf("cannot create table %v", erro)
		return
	}
	fmt.Println("create table successfully")
}

func (s Service) InsertExpense(expense *Expense) (*Expense, error) {
	tags := expense.Tags
	var id int64
	err := s.DB.QueryRow("INSERT INTO expenses(title, amount, note, tags) VALUES ($1,$2,$3,$4) RETURNING id",
		expense.Title, expense.Amount, expense.Note, pq.Array(&tags)).Scan(&id)
	if err != nil {
		log.Printf("error inserting expense to table %v\n", err)
		return nil, err
	}
	expense.ID = id
	return expense, nil
}

func (s Service) GetExpenseById(id int) (*Expense, error) {
	row := s.DB.QueryRow("SELECT * FROM expenses WHERE ID = $1", id)
	expense, err := ScanRow(row)
	return expense, err
}

func (s Service) UpdateExpenseById(id int, expense Expense) (*Expense, error) {
	tags := expense.Tags
	pq.Array(&tags)
	_, err := s.DB.Exec("UPDATE expenses SET id = $1, title = $2, amount = $3, note = $4, tags = $5 WHERE id = $6",
		expense.ID, expense.Title, expense.Amount, expense.Note, pq.Array(&tags), id)
	if err != nil {
		log.Printf("error updating expense %v\n", err)
		return nil, err
	}
	return &expense, nil
}

func (s Service) GetAllExpenses() ([]Expense, error) {
	rows, err := s.DB.Query("SELECT * FROM expenses")
	if err != nil {
		log.Printf("error retrieving expenses %v\n", err)
		return nil, err
	}
	expenses := []Expense{}

	for rows.Next() {
		expense, err := ScanRow(rows)
		if err != nil {
			log.Printf("error occurred while scanning %v\n", err)
		}
		expenses = append(expenses, *expense)
	}
	return expenses, nil
}

type Scanner interface {
	Scan(dest ...any) error
}

func ScanRow(sc Scanner) (*Expense, error) {
	var eid int64
	var title string
	var amount float64
	var note string
	var tags []string
	erro := sc.Scan(&eid, &title, &amount, &note, pq.Array(&tags))
	if erro != nil {
		log.Printf("error scanning row %v\n", erro)
		return nil, erro
	}
	return &Expense{ID: eid, Title: title, Amount: amount, Note: note, Tags: tags}, nil
}
