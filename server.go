package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chinathaip/assesment/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetAllExpenses(c echo.Context) error {
	return c.JSON(http.StatusOK, db.Expense{
		ID:     1,
		Title:  "test title",
		Amount: 13.5,
		Note:   "test note",
		Tags:   []string{"tag1", "tag2"},
	})
}

func main() {
	db.CreateTable()

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/expenses", GetAllExpenses)

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
	err := srv.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("error shutting down %v\n", err)
		return
	}
	fmt.Println("Shutdown success")
}
