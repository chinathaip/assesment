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
)

func main() {
	db.CreateTable()

	handler := CreateMainHandler()

	//graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	srv := http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: handler,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	fmt.Println("Server started")
	<-signals
	defer db.DB.Close()
	err := srv.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("error shutting down %v\n", err)
		return
	}
	fmt.Println("Shutdown success")
}
