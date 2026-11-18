package main

import (
	"encoding/json/v2"
	"net/http"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userParam struct {
		Email string `json:"email"`
	}
	err := json.UnmarshalRead(r.Body, &userParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if userParam.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid payload ('email' missing)", nil)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), userParam.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
