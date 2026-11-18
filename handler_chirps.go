package main

import (
	"encoding/json/v2"
	"net/http"
	"strings"

	"github.com/ppieniazek/chirpy/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirpParams database.CreateChirpParams
	err := json.UnmarshalRead(r.Body, &chirpParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if chirpParams.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid payload ('body' or 'user_id' missing)", nil)
		return
	}
	if len(chirpParams.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long (140 characters limit)", nil)
		return
	}

	bodyWords := strings.Split(chirpParams.Body, " ")
	for i, word := range bodyWords {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			bodyWords[i] = "****"
		}
	}
	chirpParams.Body = strings.Join(bodyWords, " ")

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp in the database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
