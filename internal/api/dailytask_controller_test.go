package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	repository "github.com/alxand/nalo-workspace/internal/repository/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp() (*fiber.App, repository.DailyTaskRepository) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(
		&models.DailyTask{},
		&models.Deliverable{},
		&models.Activity{},
		&models.ProductFocus{},
		&models.NextStep{},
		&models.Challenge{},
		&models.Note{},
		&models.Comment{},
	)

	repo := repository.NewDailyTaskRepository(db)
	app := fiber.New()
	RegisterDailyTaskRoutes(app, *repo)
	return app, *repo
}

func TestCreateDailyTask(t *testing.T) {
	app, _ := setupTestApp()

	payload := map[string]interface{}{
		"day":                "Monday",
		"date":               "2025-05-26T00:00:00Z",
		"start_time":         "2025-05-26T08:00:00Z",
		"end_time":           "2025-05-26T17:00:00Z",
		"status":             "completed",
		"score":              90,
		"productivity_score": 85,
		"deliverables": []map[string]interface{}{
			{"title": "Report Draft", "description": "Initial draft of the quarterly report"},
		},
		"activities": []map[string]interface{}{
			{"name": "Team Meeting", "duration_minutes": 60},
		},
		"product_focus": []map[string]interface{}{
			{"area": "UX Design", "focus_level": "High"},
		},
		"next_steps": []map[string]interface{}{
			{"description": "Revise the report based on feedback"},
		},
		"challenges": []map[string]interface{}{
			{"description": "Slow internet in the afternoon"},
		},
		"notes": []map[string]interface{}{
			{"content": "Need to follow up with client X"},
		},
		"comments": []map[string]interface{}{
			{"author": "John Doe", "message": "Great progress!"},
		},
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
