package goals

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler, rateLimiter fiber.Handler) {
	goalsGroup := app.Group("/api/goals", authMiddleware)

	goalsGroup.Get("/", h.GetGoals)
	goalsGroup.Post("/", h.CreateGoal)
	goalsGroup.Put("/:id", h.UpdateGoal)
	goalsGroup.Delete("/:id", h.DeleteGoal)
}

// GetGoals godoc
// @Summary      Get all goals
// @Description  Retrieves all savings goals for the authenticated user
// @Tags         Goals
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /goals [get]
func (h *Handler) GetGoals(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Invalid User ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	goals, err := h.useCase.GetGoals(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"details":    err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Get goals",
		"data":       goals,
		"request_id": c.Locals("request_id"),
	})
}

// CreateGoal godoc
// @Summary      Create a goal
// @Description  Creates a new savings goal for the authenticated user
// @Tags         Goals
// @Accept       json
// @Produce      json
// @Param        request body CreateGoalRequest true "Goal Payload"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /goals [post]
func (h *Handler) CreateGoal(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Invalid User ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	var req CreateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.CreateGoal(c.Context(), userID, &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"details":    err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"details":    err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Goal created successfully",
		"request_id": c.Locals("request_id"),
	})
}

// UpdateGoal godoc
// @Summary      Update a goal
// @Description  Updates an existing savings goal for the authenticated user
// @Tags         Goals
// @Accept       json
// @Produce      json
// @Param        id path string true "Goal ID (UUID)"
// @Param        request body UpdateGoalRequest true "Goal Payload"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /goals/{id} [put]
func (h *Handler) UpdateGoal(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Invalid User ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	idStr := c.Params("id")
	goalID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid Goal ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	var req UpdateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.UpdateGoal(c.Context(), userID, goalID, &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"details":    err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"details":    err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Goal updated successfully",
		"request_id": c.Locals("request_id"),
	})
}

// DeleteGoal godoc
// @Summary      Delete a goal
// @Description  Deletes an existing savings goal for the authenticated user
// @Tags         Goals
// @Produce      json
// @Param        id path string true "Goal ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /goals/{id} [delete]
func (h *Handler) DeleteGoal(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Invalid User ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	idStr := c.Params("id")
	goalID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid Goal ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.DeleteGoal(c.Context(), userID, goalID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"details":    err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Goal deleted successfully",
		"request_id": c.Locals("request_id"),
	})
}
