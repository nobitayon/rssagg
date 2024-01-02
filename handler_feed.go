package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nobitayon/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:%s",err))
		return
	}

	feed,err:=apiCfg.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID:uuid.New(),
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
		Name: params.Name,
		Url:params.URL,
		UserID: user.ID,
	})
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create feed:%s",err))
		return
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}