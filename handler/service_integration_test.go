package handler

// dont forget to put //go:build integration at the top
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

func newBody(t *testing.T, expense Expense) []byte {
	body, err := json.Marshal(expense)
	assert.NoError(t, err)
	return body
}

func newRequest(t *testing.T, method, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)
	return req
}

func sendRequest(t *testing.T, req *http.Request) *http.Response {
	client := http.Client{}
	resp, erro := client.Do(req)
	assert.NoError(t, erro)
	return resp
}

func readResponse(t *testing.T, resp *http.Response) *Expense {
	result := &Expense{}
	byteBody, erroo := ioutil.ReadAll(resp.Body)
	assert.NoError(t, erroo)
	json.Unmarshal(byteBody, result)
	return result
}

func TestGetAllExpenses(t *testing.T) {
	// start server
	e := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://testdb:1234@db/test-it-db?sslmode=disable")
		if err != nil {
			log.Fatalf("something went wrong: %v\n", err)
		}
		s := NewService(db)
		h := New(*s)
		e.POST("/expenses", h.HandleAddNewExpense)
		e.GET("/expenses", h.HandleGetAllExpenses)
		e.GET("/expenses/:id", h.HandleGetExpenseById)
		e.PUT("/expenses/:id", h.HandleUpdateExpenseById)
		e.Start(":2565")
	}(e)
	// wait for the server to start
	time.Sleep(2 * time.Second)

	//arrange
	url := "http://localhost:2565/expenses"

	//act
	resp, err := http.Get(url)
	assert.NoError(t, err)

	//assert
	got := []*Expense{}
	// exp := &Expense{}
	byteBody, erroo := ioutil.ReadAll(resp.Body)
	assert.NoError(t, erroo)
	json.Unmarshal(byteBody, &got)
	want := &Expense{
		ID:     1,
		Title:  "hi",
		Amount: 1.4,
		Note:   "some note",
		Tags:   []string{"tag1", "tag2"},
	}
	assert.Equal(t, want, got[0])
}

func TestInsertExpense(t *testing.T) {
	//start server
	e := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://testdb:1234@db/test-it-db?sslmode=disable")
		if err != nil {
			log.Fatalf("something went wrong: %v\n", err)
		}
		s := NewService(db)
		h := New(*s)
		e.POST("/expenses", h.HandleAddNewExpense)
		e.GET("/expenses", h.HandleGetAllExpenses)
		e.GET("/expenses/:id", h.HandleGetExpenseById)
		e.PUT("/expenses/:id", h.HandleUpdateExpenseById)
		e.Start(":2565")
	}(e)
	//wait for the server to start
	time.Sleep(2 * time.Second)

	//arrange
	url := "http://localhost:2565/expenses"
	expense := Expense{
		Title:  "yo",
		Amount: 1.4,
		Note:   "some note",
		Tags:   []string{"tag1", "tag2"},
	}
	body := newBody(t, expense)
	req := newRequest(t, http.MethodPost, url, body)

	//act
	resp := sendRequest(t, req)

	//assert
	got := readResponse(t, resp)
	want := &Expense{
		ID:     2, //ID 1 has already been initialized in scrips.sql
		Title:  "yo",
		Amount: 1.4,
		Note:   "some note",
		Tags:   []string{"tag1", "tag2"},
	}
	assert.Equal(t, want, got)
}

func TestGetExpenseById(t *testing.T) {
	// start server
	e := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://testdb:1234@db/test-it-db?sslmode=disable")
		if err != nil {
			log.Fatalf("something went wrong: %v\n", err)
		}
		s := NewService(db)
		h := New(*s)
		e.POST("/expenses", h.HandleAddNewExpense)
		e.GET("/expenses", h.HandleGetAllExpenses)
		e.GET("/expenses/:id", h.HandleGetExpenseById)
		e.PUT("/expenses/:id", h.HandleUpdateExpenseById)
		e.Start(":2565")
	}(e)
	// wait for the server to start
	time.Sleep(2 * time.Second)

	//arrange
	url := "http://localhost:2565/expenses/1"

	//act
	resp, err := http.Get(url)
	assert.NoError(t, err)

	//assert
	got := readResponse(t, resp)

	want := &Expense{
		ID:     1,
		Title:  "hi",
		Amount: 1.4,
		Note:   "some note",
		Tags:   []string{"tag1", "tag2"},
	}
	assert.Equal(t, want, got)
}

func TestUpdateExpenseById(t *testing.T) {
	// start server
	e := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://testdb:1234@db/test-it-db?sslmode=disable")
		if err != nil {
			log.Fatalf("something went wrong: %v\n", err)
		}
		s := NewService(db)
		h := New(*s)
		e.POST("/expenses", h.HandleAddNewExpense)
		e.GET("/expenses", h.HandleGetAllExpenses)
		e.GET("/expenses/:id", h.HandleGetExpenseById)
		e.PUT("/expenses/:id", h.HandleUpdateExpenseById)
		e.Start(":2565")
	}(e)
	// wait for the server to start
	time.Sleep(2 * time.Second)

	//arrange
	url := "http://localhost:2565/expenses/1"
	expense := Expense{
		ID:     1,
		Title:  "Im changing this one",
		Amount: 2,
		Note:   "some note after changed",
		Tags:   []string{"tag11", "tag22"},
	}
	body := newBody(t, expense)
	req := newRequest(t, http.MethodPut, url, body)

	//act
	resp := sendRequest(t, req)

	//assert
	got := readResponse(t, resp)
	want := &Expense{
		ID:     1,
		Title:  "Im changing this one",
		Amount: 2,
		Note:   "some note after changed",
		Tags:   []string{"tag11", "tag22"},
	}
	assert.Equal(t, want, got)
}
