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

// MAYBE NO NEED --> JUST get the row and read it just like Update and Get
func InsertExpense(expense *Expense) Expense {
	tags := expense.Tags
	lastInsertId := 0
	err := DB.QueryRow("INSERT INTO expenses(title, amount, note, tags) VALUES ($1,$2,$3,$4) RETURNING id", expense.Title, expense.Amount, expense.Note, pq.Array(&tags)).Scan(&lastInsertId)
	if err != nil {
		log.Fatalf("error inserting expense to table %v\n", err)
	}
	expense.ID = int64(lastInsertId)
	return *expense
}

// MISSING tags!!!
func GetExpenseById(id int) (*Expense, error) {
	row := DB.QueryRow("SELECT id, title, amount, note FROM expenses WHERE ID = $1", id)

	var eid int64
	var title string
	var amount float64
	var note string
	// var tags []string
	err := row.Scan(&eid, &title, &amount, &note)
	if err != nil {
		log.Printf("error retriving expense by id %v", err)
		return nil, err
	}
	return &Expense{ID: eid, Title: title, Amount: amount, Note: note}, nil
}

func UpdateExpenseById(id int, expense Expense) (*Expense, error) {
	tags := expense.Tags
	pq.Array(&tags)
	_, err := DB.Exec("UPDATE expenses SET id = $1, title = $2, amount = $3, note = $4, tags = $5 WHERE id = $6", expense.ID, expense.Title, expense.Amount, expense.Note, pq.Array(&tags), id)
	if err != nil {
		log.Printf("error updating expense %v", err)
		return nil, err
	}
	return &expense, nil
}
