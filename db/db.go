package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

type Expense struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

var DB *sql.DB

const createScript = `CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY,title TEXT,amount FLOAT,note TEXT,tags TEXT[]);`

func CreateTable() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("cnanot connect to db %v", err)
		return
	}

	_, erro := DB.Exec(createScript)
	if erro != nil {
		fmt.Printf("cannot create table %v", erro)
		return
	}
	fmt.Println("create table successfully")
}

func Disconnect() {
	DB.Close()
}

func InsertExpense(expense *Expense) (*Expense, error) {
	tags := expense.Tags
	var id int64
	err := DB.QueryRow("INSERT INTO expenses(title, amount, note, tags) VALUES ($1,$2,$3,$4) RETURNING id",
		expense.Title, expense.Amount, expense.Note, pq.Array(&tags)).Scan(&id)
	if err != nil {
		log.Printf("error inserting expense to table %v\n", err)
		return nil, err
	}
	expense.ID = id
	return expense, nil
}

func GetExpenseById(id int) (*Expense, error) {
	row := DB.QueryRow("SELECT * FROM expenses WHERE ID = $1", id)

	var eid int64
	var title string
	var amount float64
	var note string
	var tags []string
	err := row.Scan(&eid, &title, &amount, &note, pq.Array(&tags))
	if err != nil {
		log.Printf("error retriving expense by id %v\n", err)
		return nil, err
	}
	return &Expense{ID: eid, Title: title, Amount: amount, Note: note, Tags: tags}, nil
}

func UpdateExpenseById(id int, expense Expense) (*Expense, error) {
	tags := expense.Tags
	pq.Array(&tags)
	_, err := DB.Exec("UPDATE expenses SET id = $1, title = $2, amount = $3, note = $4, tags = $5 WHERE id = $6",
		expense.ID, expense.Title, expense.Amount, expense.Note, pq.Array(&tags), id)
	if err != nil {
		log.Printf("error updating expense %v\n", err)
		return nil, err
	}
	return &expense, nil
}

func GetAllExpenses() ([]Expense, error) {
	rows, err := DB.Query("SELECT * FROM expenses")
	if err != nil {
		log.Printf("error retrieving expenses %v\n", err)
		return nil, err
	}
	expenses := []Expense{}

	for rows.Next() {
		var eid int64
		var title string
		var amount float64
		var note string
		var tags []string
		erro := rows.Scan(&eid, &title, &amount, &note, pq.Array(&tags))
		if erro != nil {
			log.Printf("error scanning row %v\n", erro)
			return nil, erro
		}
		expenses = append(expenses, Expense{ID: eid, Title: title, Amount: amount, Note: note, Tags: tags})
	}
	return expenses, nil
}
