package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h Handler) HandleAddNewExpense(c echo.Context) error {
	expense := Expense{}
	if err := c.Bind(&expense); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	InsertExpense(h.DB, &expense)
	return c.JSON(http.StatusCreated, expense)
}

func (h Handler) HandleGetExpenseById(c echo.Context) error {
	query := c.Param("id")
	id, _ := strconv.Atoi(query)
	expense, err := GetExpenseById(h.DB, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	return c.JSON(http.StatusOK, expense)
}

func (h Handler) HandleUpdateExpenseById(c echo.Context) error {
	query := c.Param("id")
	id, _ := strconv.Atoi(query)
	expense := Expense{}
	if err := c.Bind(&expense); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	response, err := UpdateExpenseById(h.DB, id, expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (h Handler) HandleGetAllExpenses(c echo.Context) error {
	expenses, err := GetAllExpenses(h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "something is wrong on our end, try again later")
	}
	return c.JSON(http.StatusOK, expenses)
}
