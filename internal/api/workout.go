package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"log"

	"github.com/Rizzwaan/workoutVerse/internal/store"
	"github.com/Rizzwaan/workoutVerse/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: workStore, logger: logger}
}

func (wh *WorkoutHandler) HandleWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout ID"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)

	if err != nil {
		wh.logger.Printf("Error fetching workout with ID %d: %v", workoutID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to fetch workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"workout": workout})

}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error decoding workout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "failed to decode workout"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error creating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelop{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {

	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout ID"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error fetching workout with ID %d: %v", workoutID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to fetch workout"})
		return
	}

	if existingWorkout == nil {
		wh.logger.Printf("Workout with ID %d not found", workoutID)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelop{"error": "workout not found"})
		return
	}

	// at this point we can assume we are able to find an existing workout
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("Error decoding update workout request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "failed to decode workout"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("Error updating workout with ID %d: %v", workoutID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to update the workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"workout": existingWorkout})

}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout ID"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		wh.logger.Printf("Workout with ID %d not found", workoutID)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelop{"error": "workout not found"})
		return
	}

	if err != nil {
		wh.logger.Printf("Error deleting workout with ID %d: %v", workoutID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "error deleting workout"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelop{
		"message": "workout deleted successfully",
	})

}
