package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"pento/code-challenge/application/handlers"
	"pento/code-challenge/domain/tracker/services"
	"pento/code-challenge/repositories/postgresql"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	pgsqlAddr = ""
	pgsqlPort = 0
	kafkaAddr = ""
	kafkaPort = 0
)

//SetupAPI ...
func SetupAPI() {

	getEnvironmentVariables()

	connString := fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=postgres sslmode=disable", pgsqlAddr, pgsqlPort)

	pool, err := sql.Open("pgx", connString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	store := postgresql.NewTrackerStore(pool)
	service := services.NewTrackerService(store)
	handler := handlers.NewTrackerHandler(service)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/tracker/{id}", handler.GetTracker).Methods("GET")
	router.HandleFunc("/api/v1/tracker", handler.ListTrackers).Methods("GET")
	router.HandleFunc("/api/v1/tracker", handler.CreateTracker).Methods("POST")
	router.HandleFunc("/api/v1/tracker/{id}", handler.UpdateTracker).Methods("PUT")
	router.HandleFunc("/api/v1/tracker/{id}", handler.DeleteTracker).Methods("DELETE")

	headersOk := gHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := gHandlers.AllowedOrigins([]string{"*"})
	methodsOk := gHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Println("starting tracker API")
	log.Fatal(http.ListenAndServe(":8080", gHandlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func getEnvironmentVariables() {
	env := os.Getenv("ENVIRONMENT")

	if env == "docker" {
		pgsqlAddr = "psql"
		pgsqlPort = 5432
	} else {
		pgsqlAddr = "localhost"
		pgsqlPort = 5434
	}
}
