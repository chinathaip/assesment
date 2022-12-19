package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
			body, _ := json.Marshal(test.Input)

			resp, err := http.Post(
				fmt.Sprintf("%s/expenses", srv.URL),
				"application/json",
				bytes.NewReader(body))

			assert.Equal(t, test.Expect, resp.StatusCode)
			assert.NoError(t, err)
		})
	}

}
