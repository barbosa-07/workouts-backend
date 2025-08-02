package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Rizzwaan/workoutVerse/internal/store"
	"github.com/Rizzwaan/workoutVerse/internal/tokens"
	"github.com/Rizzwaan/workoutVerse/internal/utils"
)

type TokenHandler struct {
	tokenStore  store.TokenStore
	userHandler store.UserStore
	logger      *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Constructor for TokenHandler
func NewTokenHandler(tokenStore store.TokenStore, userHandler store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore:  tokenStore,
		userHandler: userHandler,
		logger:      logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "Invalid request body"})
		return
	}

	user, err := h.userHandler.GetUserByUsername(req.Username)
	if err != nil {
		h.logger.Printf("Error fetching user: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "Invalid username or password"})
		return
	}

	passwordMatches, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("Error comparing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "Internal server error"})
		return
	}

	if !passwordMatches {
		h.logger.Printf("Invalid password for user: %s", req.Username)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "Invalid username or password"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("Error creating token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "Could not create token"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"token": token})
}
