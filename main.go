package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"database/sql"
	"github.com/nobitayon/rssagg/internal/database"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code>499{
		log.Println("Responding with 5xx errors:",msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error:msg,
	})
}

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")

	portString:=os.Getenv("PORT")
	if portString==""{
		log.Fatal("PORT is not found in the environment")
	}

	dbURL:=os.Getenv("DB_URL")
	if dbURL==""{
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err:= sql.Open("postgres",dbURL)
	if err!=nil{
		log.Fatal("Can't connect to database")
	}

	apiCfg:=apiConfig{
		DB: database.New(conn),
	}

	router:=chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:		[]string{"https://*","http:://*"},
		AllowedMethods:		[]string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders:		[]string{"*"},
		ExposedHeaders:		[]string{"Link"},
		AllowCredentials:	false,
		MaxAge:				300,
	}))

	v1Router:=chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handlerErr)
	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	router.Mount("/v1",v1Router)

	srv:=&http.Server{
		Handler: router,
		Addr:":"+portString,
	}

	log.Printf("Server starting on port %v",portString)
	err = srv.ListenAndServe()
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Println("Port:",portString)
}