//go:build unit
// +build unit

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockService struct{}

func (m MockService) CreateTable() {}
func (m MockService) InsertExpense(expense *Expense) (*Expense, error) {
	expense.ID = 1
	return expense, nil
}
func (m MockService) GetExpenseById(id int) (*Expense, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return &Expense{}, nil
}
func (m MockService) UpdateExpenseById(id int, expense Expense) (*Expense, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	return &Expense{}, nil
}
func (m MockService) GetAllExpenses() ([]Expense, error) {
	return []Expense{{}}, nil
}

func setup() *httptest.Server {
	handler := MockHandler()
	srv := httptest.NewServer(handler)
	return srv
}

func teardown(srv *httptest.Server) {
	srv.Close()
}

func MockHandler() *echo.Echo {
	e := echo.New()
	handler := New(MockService{})

	e.POST("/expenses", handler.HandleAddNewExpense)

	e.GET("/expenses", handler.HandleGetAllExpenses)

	e.GET("/expenses/:id", handler.HandleGetExpenseById)

	e.PUT("/expenses/:id", handler.HandleUpdateExpenseById)

	return e
}

// Story 1: As a user, I want to add a new expense So that I can track my expenses
func TestHandleAddNewExpense(t *testing.T) {
	srv := setup()
	defer teardown(srv)

	tests := []struct {
		TestName string
		Input    interface{}
		Expect   int
	}{
		{
			TestName: "add valid expense should return status created",
			Input: Expense{
				Title:  "hi",
				Amount: 1.4,
				Note:   "some note",
				Tags:   []string{"tag1", "tag2"},
			},
			Expect: http.StatusCreated,
		},
		{
			TestName: "add invalid expense should return status created",
			Input:    "hello test test",
			Expect:   http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			//arrange
			url := fmt.Sprint(srv.URL + "/expenses")
			body, _ := json.Marshal(test.Input)
			//act
			resp, err := http.Post(
				url,
				"application/json",
				bytes.NewReader(body))
			//assert
			assert.Equal(t, test.Expect, resp.StatusCode)
			assert.NoError(t, err)
		})
	}

}

// Story 2: As a user, I want to see my expense by using expense ID So that I can check my expense information
func TestHandleGetExpenseById(t *testing.T) {
	srv := setup()
	defer teardown(srv)

	tests := []struct {
		TestName string
		Input    string
		Expect   int
	}{
		{
			TestName: "valid id should return status OK",
			Input:    "1",
			Expect:   http.StatusOK,
		},
		{
			TestName: "invalid id should return status bad request",
			Input:    "yoyo",
			Expect:   http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			//arrange
			url := fmt.Sprint(srv.URL + "/expenses/" + test.Input)
			//act
			resp, err := http.Get(url)
			//assert
			assert.Equal(t, test.Expect, resp.StatusCode)
			assert.NoError(t, err)
		})
	}
}

// Story 3: As a user, I want to update my expense So that I can correct my expense information
func TestHandleUpdateExpenseById(t *testing.T) {
	srv := setup()
	defer teardown(srv)

	tests := []struct {
		TestName string
		ID       string
		Body     interface{}
		Expect   int
	}{
		{
			TestName: "valid id and body should get status OK",
			ID:       "1",
			Body: Expense{
				ID:     1,
				Title:  "hi",
				Amount: 1.4,
				Note:   "some note",
				Tags:   []string{"tag1", "tag2"},
			},
			Expect: http.StatusOK,
		},
		{
			TestName: "valid id but invalid body should get status bad request",
			ID:       "1",
			Body:     "ayoyoyo",
			Expect:   http.StatusBadRequest,
		},
		{
			TestName: "invalid id but valid body should get status bad request",
			ID:       "hello test tset test",
			Body: Expense{
				ID:     1,
				Title:  "hi",
				Amount: 1.4,
				Note:   "some note",
				Tags:   []string{"tag1", "tag2"},
			},
			Expect: http.StatusBadRequest,
		},
		{
			TestName: "invalid id and invalid body should get status bad request",
			ID:       "hello test tset test",
			Body:     "test again ",
			Expect:   http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			url := fmt.Sprint(srv.URL + "/expenses/" + test.ID)
			body, _ := json.Marshal(test.Body)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, erro := client.Do(req)

			assert.Equal(t, test.Expect, resp.StatusCode)
			assert.NoError(t, err)
			assert.NoError(t, erro)
		})
	}
}

// Story 4: As a user, I want to see all my expenses So that I can check my expense information
func TestHandleGetAllExpense(t *testing.T) {
	t.Run("Normal server return status OK", func(t *testing.T) {
		srv := setup()
		defer teardown(srv)

		url := fmt.Sprint(srv.URL + "/expenses")
		resp, err := http.Get(url)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, err)
	})
}
