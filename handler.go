package main

import (
	"net/http"

	"github.com/chinathaip/assesment/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateMainHandler() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/expenses", AddNewExpense)
	return e
}

func AddNewExpense(c echo.Context) error {

	expense := db.Expense{}
	if err := c.Bind(&expense); err != nil {
		return err
	}
	db.InsertExpense(&expense)
	return c.JSON(http.StatusCreated, expense)
}

// db.Expense{
// 		Title:  "hi",
// 		Amount: 1.4,
// 		Note:   "some note",
// 		Tags:   []string{"tag1", "tag2"},
// 	}
