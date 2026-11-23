package main

import (
	"database/sql"
	"encoding/json/v2"
	"errors"
	"net/http"
	"strings"
	"uuid"

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

func (cfg *apiConfig) listChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get the chirps from the database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse uuid from url", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp doesn't exist", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp from the database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
