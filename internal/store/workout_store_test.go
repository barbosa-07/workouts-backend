package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	err = Migrate(db, "../../migrations/")

	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	_, err = db.Exec("TRUNCATE workouts, workout_entries CASCADE") // Clear the workouts table before each test
	if err != nil {
		t.Fatalf("Failed to clear workouts table: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "Create valid workout",
			workout: &Workout{
				Title:           "Morning Run",
				Description:     "5km run in the park",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       Float64Ptr(130.55),
						Notes:        "Felt good",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{name: "Create workout with empty title", workout: &Workout{
			Title:           "",
			Description:     "5km run in the park",
			DurationMinutes: 30,
			CaloriesBurned:  300,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench Press",
					Sets:         3,
					Reps:         IntPtr(10),
					Weight:       Float64Ptr(130.55),
					Notes:        "Felt good",
					OrderIndex:   1,
				},
			},
		}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrievedWorkout, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrievedWorkout.ID)
			assert.Equal(t, createdWorkout.Title, retrievedWorkout.Title)
			assert.Equal(t, createdWorkout.Description, retrievedWorkout.Description)
			assert.Equal(t, len(tt.workout.Entries), len(retrievedWorkout.Entries))

			for i, entry := range tt.workout.Entries {
				assert.Equal(t, entry.ExerciseName, retrievedWorkout.Entries[i].ExerciseName)
				assert.Equal(t, entry.Sets, retrievedWorkout.Entries[i].Sets)
				assert.Equal(t, entry.Reps, retrievedWorkout.Entries[i].Reps)
				assert.Equal(t, entry.Weight, retrievedWorkout.Entries[i].Weight)
				assert.Equal(t, entry.Notes, retrievedWorkout.Entries[i].Notes)
				assert.Equal(t, entry.OrderIndex, retrievedWorkout.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func Float64Ptr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
