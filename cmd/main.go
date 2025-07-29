package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MicaelaJofre/gocourse_enrollment/internal/enrollment"
	"github.com/MicaelaJofre/gocourse_enrollment/pkg/bootstrap"
	"github.com/MicaelaJofre/gocourse_enrollment/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: no se pudo cargar .env (%v)", err)
	}

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)
	l := bootstrap.InitLogger()


	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := bootstrap.DBConecction()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	pagLimitDef := os.Getenv("PAGINATION_DEFAULT_PAGE")
	if pagLimitDef == "" {
		l.Fatal("PAGINATION_DEFAULT_PAGE environment variable not set.")
	}

	ctx := context.Background()

	enrollmentRepo := enrollment.NewRepository(db, l)
	enrollmentSrv := enrollment.NewService(l,enrollmentRepo)
	h :=handler.NewEnrollmentHTTpServer(ctx, enrollment.MakeEndpoint(enrollmentSrv, enrollment.Config{LimPageDef: pagLimitDef}))


	server := &http.Server{
		Addr:         address,
		Handler:      accessControl(h), 
		ReadTimeout:  5 * time.Second, 
		WriteTimeout: 5 * time.Second, 
		IdleTimeout: 120 * time.Second, 
	}

	errCh := make(chan error)
	go func(){
		l.Println("Listen in", address)
		errCh <- server.ListenAndServe()
	}()
	
	err = <-errCh
	if err != nil {
		l.Fatal(err)
	}
	
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type, X-Requested-With")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}