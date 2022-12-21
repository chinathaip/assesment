package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chinathaip/assesment/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// higher order func (check if allow --> then call next)
func middlewareCheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization != "November 10, 2009" {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func NewApplication(handler *handler.Handler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middlewareCheckAuth)

	e.POST("/expenses", handler.HandleAddNewExpense)
	e.GET("/expenses", handler.HandleGetAllExpenses)
	e.GET("/expenses/:id", handler.HandleGetExpenseById)
	e.PUT("/expenses/:id", handler.HandleUpdateExpenseById)

	return e
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("error connecting to DB")
		return
	}
	fmt.Println("connect to DB successfully")

	service := handler.NewService(db)
	handler := handler.New(*service)
	handler.Service.CreateTable()
	e := NewApplication(handler)

	//graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	srv := http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: e,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	fmt.Println("Server started")
	<-signals
	defer db.Close()
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("error shutting down %v\n", err)
		return
	}
	fmt.Println("Shutdown success")
}
