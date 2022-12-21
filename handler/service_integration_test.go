package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// dont forget to put //go:build integration at the top
func TestInsertExpense(t *testing.T) {
	e := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://testdb:1234@db/test-it-db?sslmode=disable")
		if err != nil {
			log.Fatalf("something went wrong: %v\n", err)
		}

		s := NewService(db)
		h := New(*s)
		e.POST("/expenses", h.HandleAddNewExpense)
		e.Start(":2565")
	}(e)

	//wait for the server to respond
	time.Sleep(2 * time.Second)

	//arrange
	url := "http://localhost:2565/expenses"
	expense := Expense{
		Title:  "hi",
		Amount: 1.4,
		Note:   "some note",
		Tags:   []string{"tag1", "tag2"},
	}

	body, _ := json.Marshal(expense)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	//act
	resp, erro := client.Do(req)
	assert.NoError(t, erro)

	got := &Expense{}
	byteBody, erroo := ioutil.ReadAll(resp.Body)
	assert.NoError(t, erroo)
	json.Unmarshal([]byte(byteBody), &got)

	//assert
	want := &Expense{
		ID:     1,
		Title:  "hi",
		Amount: 1.4,
		Note:   "some note",
		Tags:   []string{"tag1", "tag2"},
	}
	assert.Equal(t, got, want)

}
