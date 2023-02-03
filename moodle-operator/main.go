package main

import (
	"context"
	"log"

	// mysqladapter "moodle/adapter/mariadb"
	k8scli "moodle/cli/k8s"
	mariacli "moodle/cli/maria"

	// k8sworker "moodle/worker/k8s"
	lbworker "moodle/worker/loadbalancer"
	trackingworker "moodle/worker/tracking"
	sendmailworker "moodle/worker/sendmail"

	mariaworker "moodle/worker/maria"

	"moodle/config"
	"moodle/handler"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slices"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)


var mariaCommands = []string{
	"list-instance",
	"create-instance",
	"delete-instance",
	"create-instance-from-backup",
}


var k8sCommands = []string{
	"list-namespaces",
	"list-pods",
	"list-storages-classes",
	"set-default-storage-class",
	"get-env-vars-values",
	"create-namespace",
	"delete-namespace",
	"get-service",
	"apply-pvc",
	"apply-service",
	"apply-statefulset",
	"delete-pvc",
	"delete-service",
	"delete-statefulset",
	"list-pvc",
	"get-podlogs",
	"list-statefulsets",
	"list-services",
	"patch-statefulset",
}


func main() {

	// cli
	if len(os.Args) > 1 {
		// check string in slice
		if slices.Contains(k8sCommands, os.Args[1]) {
			cli := k8scli.NewK8sCli()
			cli.Run()
			return
		} else if slices.Contains(mariaCommands, os.Args[1]) {
			mariacli := mariacli.NewMariaCli()
			mariacli.Run()
			return
		}
	}

	// api
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")
	ctx, cancel := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		logger.Fatal(err)
	}
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	ctx_, cancel_ := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel_()
	err = client.Ping(ctx_, nil)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("Connected to MongoDB!")

	myHandler := handler.New(client.Database(config.DBNAME))
	r := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	r.Use(MoodleMiddleware)
	r.HandleFunc("/health", myHandler.HealthCheck).Methods("GET", "OPTIONS")
	r.HandleFunc("/moodles", myHandler.ListMoodle). // search if have search params
	    Queries(
		"page", "{page}",
		"limit", "{limit}",
		"email", "{email}",
    ).Methods("GET", "OPTIONS")
	r.HandleFunc("/moodles", myHandler.CreateMoodle).Methods("POST", "OPTIONS")
	r.HandleFunc("/moodles", myHandler.DeleteMoodle).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/moodles", myHandler.GetMoodle).Queries(
		"id", "{id}",
	).Methods("GET", "OPTIONS")
	r.HandleFunc("/courses", myHandler.ListCourse).
	    Queries(
		"page", "{page}",
		"limit", "{limit}",
    ).Methods("GET", "OPTIONS")
	r.HandleFunc("/packages", myHandler.ListPackage).
	    Queries(
		"page", "{page}",
		"limit", "{limit}",
    ).Methods("GET", "OPTIONS")
	r.HandleFunc("/packages", myHandler.ChangeMoodlePackages).Methods("PUT", "OPTIONS")
	r.HandleFunc("/pre-installed-course", myHandler.ChangeMoodlePreInstalledCourse).Methods("PUT", "OPTIONS")
	r.HandleFunc("/autoscale", myHandler.ChangeMoodleAutoScale).Methods("PUT", "OPTIONS")
	r.HandleFunc("/document-storage-extra", myHandler.ChangeMoodleDocumentsStorageExtra).Methods("PUT", "OPTIONS")
	r.HandleFunc("/users", myHandler.UserAdd).Methods("POST", "OPTIONS")
	r.HandleFunc("/users", myHandler.UserDelete).Methods("DELETE", "OPTIONS")
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

	// lb worker
	go func() {
		for {
			time.Sleep(time.Duration(config.INTERVAL) * time.Second)
			lbWorker := lbworker.NewLBWorker(client.Database(config.DBNAME))
			// fmt.Println("lb worker is running...")
			lbWorker.GetLBInstances()
		}
	}()

	// sendmail worker
	go func() {
		for {
			time.Sleep(time.Duration(config.INTERVAL) * time.Second)
			sendmailWorker := sendmailworker.NewSendmailWorker(client.Database(config.DBNAME))
			// fmt.Println("sendmail worker is running...")
			sendmailWorker.SendMailWorkerPool()
		}
	}()

	// tracking worker
	go func() {
		for {
			time.Sleep(time.Duration(config.INTERVAL) * time.Second)
			trackingWorker := trackingworker.NewTrackingWorker(client.Database(config.DBNAME))
			// fmt.Println("tracking worker is running...")
			trackingWorker.UpdateLbStatus()
		}
	}()
 
	// maria worker
	go func() {
		for {
			time.Sleep(time.Duration(config.INTERVAL) * time.Second)
			mariaWorker := mariaworker.NewMariaWorker(client.Database(config.DBNAME))
			// fmt.Println("maria worker is running...")
			mariaWorker.Restore()
		}
	}()
 
// 	// k8s worker
// 	go func() {
// 		for {
// 			time.Sleep(time.Duration(config.INTERVAL) * time.Second)
// 			k8sWorker := k8sworker.NewK8sWorker(client.Database(config.DBNAME))
// 			fmt.Println("k8s worker is running...")
// 			k8sWorker.ApplyStatefulSetWorker(myHandler.GetClientset())
// 		}
// 	}()


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

func Auth(tenantName, token string) bool {
	if tenantName != config.USER {
		return false
	}

	if token != config.MOODLE_TOKEN {
		return false
	}
	return true
}


func MoodleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ignore health check
		if r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}

		// handle cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		// set response header
		w.Header().Set("Content-Type", "application/json")

		// handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		tenentName := r.Header.Get("X-Tenant-Name")
		token := r.Header.Get("X-Auth-Token")
		if !Auth(tenentName, token) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
