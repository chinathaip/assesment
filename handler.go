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
	e.PUT("/expenses/:id", UpdateExpenseById)
	e.GET("/expenses", GetAllExpenses)

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

func UpdateExpenseById(c echo.Context) error {
	query := c.Param("id")
	id, _ := strconv.Atoi(query)
	expense := db.Expense{}
	err := c.Bind(&expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response, erro := db.UpdateExpenseById(id, expense)
	if erro != nil {
		return c.JSON(http.StatusBadRequest, erro)
	}
	return c.JSON(http.StatusOK, response)
}

func GetAllExpenses(c echo.Context) error {
	expenses, err := db.GetAllExpenses()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "something is wrong on our end, try again later")
	}
	return c.JSON(http.StatusOK, expenses)
}
