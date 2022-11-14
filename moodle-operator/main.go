package main

import (
	"context"
	"log"
	"moodle/handler"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"flag"
	"fmt"
	k8scli "moodle/cli"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)

func main() {
	// cli
	if len(os.Args) > 1 {
		k8scli := k8scli.NewK8sCli()
		k8scli.Run()
		return
	}

	// api
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.NewClient(
		options.Client().ApplyURI(handler.MONGOURI))
	if err != nil {
		log.Fatalf("Error creating mongo client: %+v", err)
	}
	defer client.Disconnect(ctx)
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %+v", err)
	}
	myHandler := handler.New(client.Database(handler.DBNAME))
	r := mux.NewRouter()
	r.HandleFunc("/health", myHandler.HealthCheck).Methods("GET")
	r.HandleFunc("/api/moodles", myHandler.ListMoodle).
	    Queries(
        "userid", "{userid}",
    ).Methods("GET")
	r.HandleFunc("/api/moodles", myHandler.CreateMoodle).Methods("POST")
	r.HandleFunc("/api/moodles", myHandler.DeleteMoodle).Methods("DELETE")
	r.HandleFunc("/api/moodles", myHandler.ScaleMoodle).Methods("PUT")
	r.HandleFunc("/api/users", myHandler.UserAdd).Methods("POST")
	r.HandleFunc("/api/users", myHandler.UserDelete).Methods("DELETE")
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:		  listenAddr,
		Handler:	  tracing(nextRequestID)(logging(logger)(r)),
		ErrorLog:	  logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}


func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

