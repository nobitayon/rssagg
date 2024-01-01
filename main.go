package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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

func main() {

	godotenv.Load(".env")

	portString:=os.Getenv("PORT")
	if portString==""{
		log.Fatal("PORT is not found in the environment")
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

	router.Mount("/v1",v1Router)

	srv:=&http.Server{
		Handler: router,
		Addr:":"+portString,
	}

	log.Printf("Server starting on port %v",portString)
	err := srv.ListenAndServe()
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Println("Port:",portString)
}