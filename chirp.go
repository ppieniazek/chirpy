package main

import (
	"encoding/json/v2"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirpParam struct {
		Body string `json:"body"`
	}
	err := json.UnmarshalRead(r.Body, &chirpParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if chirpParam.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid payload ('body' missing)", nil)
		return
	}
	if len(chirpParam.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long (140 characters limit)", nil)
		return
	}

	var chirpVal struct {
		CleanedBody string `json:"cleaned_body"`
	}
	bodyWords := strings.Split(chirpParam.Body, " ")
	for i, word := range bodyWords {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			bodyWords[i] = "****"
		}
	}
	chirpVal.CleanedBody = strings.Join(bodyWords, " ")

	respondWithJSON(w, http.StatusOK, chirpVal)
}
