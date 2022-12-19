package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

const createScript = `CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY,title TEXT,amount FLOAT,note TEXT,tags TEXT[]);`

func CreateTable() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("cnanot connect to db %v", err)
		return
	}
	defer db.Close()

	_, erro := db.Exec(createScript)
	if erro != nil {
		fmt.Printf("cannot create table %v", erro)
		return
	}
	fmt.Println("create table successfully")
}
