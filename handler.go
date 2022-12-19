package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/chinathaip/assesment/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateMainHandler() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/expenses", AddNewExpense)
	e.GET("/expenses/:id", GetExpenseById)
	return e
}

func AddNewExpense(c echo.Context) error {
	expense := db.Expense{}
	if err := c.Bind(&expense); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	db.InsertExpense(&expense)
	return c.JSON(http.StatusCreated, expense)
}

func GetExpenseById(c echo.Context) error {
	query := c.Param("id")
	id, _ := strconv.Atoi(query)
	fmt.Println(id)
	expense, err := db.GetExpenseById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	return c.JSON(http.StatusOK, expense)
}

// db.Expense{
// 		Title:  "hi",
// 		Amount: 1.4,
// 		Note:   "some note",
// 		Tags:   []string{"tag1", "tag2"},
// 	}
