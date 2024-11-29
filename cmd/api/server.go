package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SETUP SERVER CONFIGS

type application struct {
	server *http.Server
	db     *sqlx.DB
}

func (a *application) Start() {
	log.Printf("server is running on port: %s", a.server.Addr)
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Printf("Listen: %s", err)
			} else {
				log.Printf("could not start server: %v", err)
			}
		}
	}()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting Down Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Forced To ShutDown, %v", err)
	}
	log.Println("Server Exiting!")
}

func NewApplication(handler *gin.Engine, db *sqlx.DB) *application {
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	return &application{
		server: &http.Server{
			Addr:    PORT,
			Handler: handler,
		},
		db: db,
	}
}
