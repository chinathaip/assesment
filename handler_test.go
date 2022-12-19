package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/chinathaip/assesment/db"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func MockHandler() *echo.Echo {
	e := echo.New()

	e.POST("/expenses", func(c echo.Context) error {
		expense := db.Expense{}
		if err := c.Bind(&expense); err != nil {
			return c.JSON(http.StatusBadRequest, "bad request")
		}
		return c.JSON(http.StatusCreated, expense)
	})

	e.GET("/expenses/:id", func(c echo.Context) error {
		query := c.Param("id")
		id, _ := strconv.Atoi(query)

		fmt.Printf("QUERY IS : %v", id)
		if id == 0 {
			return c.JSON(http.StatusBadRequest, "bad request")
		}
		return c.JSON(http.StatusOK, "OK")
	})

	return e
}

// Story 1: As a user, I want to add a new expense So that I can track my expenses
func TestAddNewExpense(t *testing.T) {
	handler := MockHandler()
	srv := httptest.NewServer(handler)
	defer srv.Close()
	tests := []struct {
		TestName string
		Input    interface{}
		Expect   int
	}{
		{
			TestName: "add valid expense should return status created",
			Input: db.Expense{
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
func TestGetExpenseById(t *testing.T) {
	handler := MockHandler()
	srv := httptest.NewServer(handler)
	defer srv.Close()

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
			fmt.Printf("URL IS %v\n", url)
			//assert
			assert.Equal(t, test.Expect, resp.StatusCode)
			assert.NoError(t, err)
		})
	}
}
