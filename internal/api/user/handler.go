package user

import (
	"strconv"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	userRepo interfaces.UserInterface
	logger   *zap.Logger
}

func NewUserHandler(userRepo interfaces.UserInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

// ListUsers godoc
// @Summary List all users (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userRepo.List(limit, offset)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		return errors.DatabaseError("Failed to list users", err)
	}

	return c.JSON(users)
}

// GetUser godoc
// @Summary Get user by ID (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.BadRequest("Invalid user ID", err)
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int64("user_id", id), zap.Error(err))
		return errors.NotFound("User not found", err)
	}

	return c.JSON(user)
}

// UpdateUser godoc
// @Summary Update user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.BadRequest("Invalid user ID", err)
	}

	// Check if user exists
	_, err = h.userRepo.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int64("user_id", id), zap.Error(err))
		return errors.NotFound("User not found", err)
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return errors.BadRequest("Invalid request body", err)
	}

	user.ID = id

	if err := h.userRepo.Update(&user); err != nil {
		h.logger.Error("Failed to update user", zap.Int64("user_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to update user", err)
	}

	h.logger.Info("User updated successfully", zap.Int64("user_id", id))
	return c.JSON(user)
}

// DeleteUser godoc
// @Summary Delete user (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.BadRequest("Invalid user ID", err)
	}

	// Check if user exists
	_, err = h.userRepo.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int64("user_id", id), zap.Error(err))
		return errors.NotFound("User not found", err)
	}

	if err := h.userRepo.Delete(id); err != nil {
		h.logger.Error("Failed to delete user", zap.Int64("user_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to delete user", err)
	}

	h.logger.Info("User deleted successfully", zap.Int64("user_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}
