package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lambda/db"
	"lambda/routes"
)

func main() {
	db.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/api/bookOrder", routes.BookOrderHandler)
	http.HandleFunc("/api/rejectOrder", routes.RejectOrderHandler)
	http.HandleFunc("/api/orderQueue", routes.InsertOrderQueueHandler)
	http.HandleFunc("/api/orderProcess", routes.DeleteOrderProcessHandler)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-c
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Disconnect MongoDB client
	if err := db.Client.Disconnect(ctx); err != nil {
		log.Fatal("Error disconnecting MongoDB client:", err)
	}

	log.Println("Server exiting")
}
