package dailytask

import (
	"strconv"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/errors"
	"github.com/alxand/nalo-workspace/internal/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TaskHandler struct {
	Repo   interfaces.DailyTaskInterface
	Logger *zap.Logger
}

func NewTDailyTaskHandler(repo interfaces.DailyTaskInterface, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		Repo:   repo,
		Logger: logger,
	}
}

// CreateDailyTask godoc
// @Summary Create a new daily task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.DailyTask true "Daily task data"
// @Success 201 {object} models.DailyTask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dailytask [post]
func (h *TaskHandler) CreateDailyTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user == nil {
		return errors.Unauthorized("User not found in context", nil)
	}

	var task models.DailyTask

	if err := c.BodyParser(&task); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	task.UserID = user.ID

	if err := validation.ValidateDailyTask(&task); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	if err := h.Repo.Create(&task); err != nil {
		h.Logger.Error("Failed to create task", zap.Error(err))
		return errors.DatabaseError("Failed to create task", err)
	}

	h.Logger.Info("Task created successfully", zap.Int64("task_id", task.ID), zap.Int64("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTasksByDate godoc
// @Summary Get tasks by date for the authenticated user
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param date path string true "Date in YYYY-MM-DD format"
// @Success 200 {array} models.DailyTask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dailytask/{date} [get]
func (h *TaskHandler) GetTasksByDate(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user == nil {
		return errors.Unauthorized("User not found in context", nil)
	}

	date := c.Params("date")
	if date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Date parameter is required"})
	}

	tasks, err := h.Repo.GetByDateAndUser(date, user.ID)
	if err != nil {
		h.Logger.Error("Failed to get tasks by date", zap.String("date", date), zap.Int64("user_id", user.ID), zap.Error(err))
		return errors.DatabaseError("Failed to get tasks", err)
	}

	h.Logger.Info("Tasks retrieved successfully", zap.String("date", date), zap.Int64("user_id", user.ID), zap.Int("count", len(tasks)))
	return c.JSON(tasks)
}

// UpdateTask godoc
// @Summary Update a task (only if owned by the authenticated user)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param task body models.DailyTask true "Updated task data"
// @Success 200 {object} models.DailyTask
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dailytask/{id} [put]
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user == nil {
		return errors.Unauthorized("User not found in context", nil)
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	existingTask, err := h.Repo.GetByID(id)
	if err != nil {
		h.Logger.Error("Failed to get task", zap.Int64("task_id", id), zap.Error(err))
		return errors.NotFound("Task not found", err)
	}

	if existingTask.UserID != user.ID {
		h.Logger.Error("User trying to update task they don't own", zap.Int64("task_id", id), zap.Int64("user_id", user.ID), zap.Int64("task_user_id", existingTask.UserID))
		return errors.Forbidden("You can only update your own tasks", nil)
	}

	var task models.DailyTask
	if err := c.BodyParser(&task); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	task.ID = id
	task.UserID = user.ID

	if err := validation.ValidateDailyTask(&task); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	updatedTask, err := h.Repo.Update(&task)
	if err != nil {
		h.Logger.Error("Failed to update task", zap.Int64("task_id", id), zap.Int64("user_id", user.ID), zap.Error(err))
		return errors.DatabaseError("Failed to update task", err)
	}

	h.Logger.Info("Task updated successfully", zap.Int64("task_id", id), zap.Int64("user_id", user.ID))
	return c.JSON(updatedTask)
}

// DeleteTask godoc
// @Summary Delete a task (only if owned by the authenticated user)
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dailytask/{id} [delete]
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user == nil {
		return errors.Unauthorized("User not found in context", nil)
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	existingTask, err := h.Repo.GetByID(id)
	if err != nil {
		h.Logger.Error("Failed to get task", zap.Int64("task_id", id), zap.Error(err))
		return errors.NotFound("Task not found", err)
	}

	if existingTask.UserID != user.ID {
		h.Logger.Error("User trying to delete task they don't own", zap.Int64("task_id", id), zap.Int64("user_id", user.ID), zap.Int64("task_user_id", existingTask.UserID))
		return errors.Forbidden("You can only delete your own tasks", nil)
	}

	if err := h.Repo.Delete(id); err != nil {
		h.Logger.Error("Failed to delete task", zap.Int64("task_id", id), zap.Int64("user_id", user.ID), zap.Error(err))
		return errors.DatabaseError("Failed to delete task", err)
	}

	h.Logger.Info("Task deleted successfully", zap.Int64("task_id", id), zap.Int64("user_id", user.ID))
	return c.SendStatus(fiber.StatusNoContent)
}
