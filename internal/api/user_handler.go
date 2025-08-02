package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/Rizzwaan/workoutVerse/internal/store"
	"github.com/Rizzwaan/workoutVerse/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

// Constructor for UserHandler
func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

// server side validation for user registration
func (h *UserHandler) validateRegisterUserRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "Invalid request body"})
		return
	}

	err = h.validateRegisterUserRequest(&req)
	if err != nil {
		h.logger.Printf("Validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// how do we deal with password
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "Internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelop{"message": "User created successfully", "user_id": user.ID})

}
