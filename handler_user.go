package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nobitayon/rssagg/internal/database"
	"github.com/nobitayon/rssagg/internal/auth"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Name string `json:"name"`
	}
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:%s",err))
		return
	}

	user,err:=apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID:uuid.New(),
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
		Name: params.Name,
	})
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create user:%s",err))
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request){
	apiKey, err:=auth.GetAPIKey(r.Header)
	if err!=nil{
		respondWithError(w,403,fmt.Sprintf("Auth error:%s",err))
		return 
	}
	user,err:=apiCfg.DB.GetUserByAPIKey(r.Context(),apiKey)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't get user:%v",err))
		return
	}
	respondWithJSON(w,200,databaseUserToUser(user))
}