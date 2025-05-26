package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	repository "github.com/alxand/nalo-workspace/internal/repository/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp() (*fiber.App, *repository.GormLogRepository) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.DailyTask{})
	repo := repository.NewGormLogRepository(db)
	app := fiber.New()
	RegisterLogRoutes(app, repo)
	return app, repo
}

func TestCreateAndGetLog(t *testing.T) {
	app, _ := setupTestApp()

	// Create a log
	logJSON := `{"day":"Monday","date":"2025-05-24"}`
	req := httptest.NewRequest("POST", "/logs", strings.NewReader(logJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)

	// Get logs by date
	req2 := httptest.NewRequest("GET", "/logs/2025-05-24", nil)
	resp2, _ := app.Test(req2)
	assert.Equal(t, 200, resp2.StatusCode)

	var logs []models.DailyTask
	json.NewDecoder(resp2.Body).Decode(&logs)
	assert.NotEmpty(t, logs)
	assert.Equal(t, "Monday", logs[0].Day)
}
