package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nobitayon/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:%s",err))
		return
	}

	feedFollow,err:=apiCfg.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{
		ID:uuid.New(),
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC(),
		UserID: user.ID,
		FeedID:params.FeedID,
	})
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't create feed follow:%s",err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){

	feedFollows,err:=apiCfg.DB.GetFeedFollows(r.Context(),user.ID)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't get feed follow:%s",err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollowIDStr:=chi.URLParam(r,"feedFollowID")
	feedFollowID,err:=uuid.Parse(feedFollowIDStr)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't parse feed follow id:%s",err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollows(r.Context(),database.DeleteFeedFollowsParams{
		ID:feedFollowID,
		UserID:user.ID,
	})
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldn't delete feed follow:%s",err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}