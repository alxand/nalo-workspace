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

// setupTestApp initializes an in-memory Fiber app with the DailyTaskRepository and routes registered.
func setupTestApp(t *testing.T) (*fiber.App, repository.DailyTaskRepository) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "failed to open in-memory database")

	err = db.AutoMigrate(&models.DailyTask{})
	assert.NoError(t, err, "failed to auto-migrate DailyTask model")

	repo := repository.NewDailyTaskRepository(db)
	app := fiber.New()
	RegisterDailyTaskRoutes(app, *repo)
	// RegisterDailyTaskRoutes(app, *repo)

	return app, *repo
}

func TestCreateAndGetLog(t *testing.T) {
	app, _ := setupTestApp(t)

	t.Run("Create Daily Log", func(t *testing.T) {
		logJSON := `{"day":"Monday","date":"2025-05-24"}`
		req := httptest.NewRequest("POST", "/tasks/", strings.NewReader(logJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Get Daily tasks by Date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/2025-05-24", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var tasks []models.DailyTask
		err = json.NewDecoder(resp.Body).Decode(&tasks)
		assert.NoError(t, err)
		assert.NotEmpty(t, tasks)
		assert.Equal(t, "Monday", tasks[0].Day)
	})
}
