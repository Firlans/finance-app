package budgets

import (
	"time"

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
	budgetsGroup := app.Group("/api/budgets", authMiddleware)
	
	budgetsGroup.Post("/", h.UpsertBudget)
	budgetsGroup.Delete("/:id", h.DeleteBudget)
	budgetsGroup.Get("/summary", h.GetBudgetSummaries)
}

// UpsertBudget godoc
// @Summary      Create or update a budget
// @Description  Creates a new budget or updates an existing one for the authenticated user
// @Tags         Budgets
// @Accept       json
// @Produce      json
// @Param        request body CreateBudgetRequest true "Budget Payload"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /budgets [post]
func (h *Handler) UpsertBudget(c *fiber.Ctx) error {
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

	var req CreateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.UpsertBudget(c.Context(), userID, &req)
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
		"message":    "Budget updated successfully",
		"request_id": c.Locals("request_id"),
	})
}

// DeleteBudget godoc
// @Summary      Delete a budget
// @Description  Deletes an existing budget for the authenticated user
// @Tags         Budgets
// @Produce      json
// @Param        id path string true "Budget ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /budgets/{id} [delete]
func (h *Handler) DeleteBudget(c *fiber.Ctx) error {
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
	budgetID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid Budget ID format",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.DeleteBudget(c.Context(), userID, budgetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"details":    err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Budget deleted successfully",
		"request_id": c.Locals("request_id"),
	})
}

// GetBudgetSummaries godoc
// @Summary      Get budget summaries
// @Description  Retrieves budget summaries for the authenticated user for a specific date
// @Tags         Budgets
// @Accept       json
// @Produce      json
// @Param        date query string false "Date in YYYY-MM-DD format. Defaults to today."
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /budgets/summary [get]
func (h *Handler) GetBudgetSummaries(c *fiber.Ctx) error {
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

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	summaries, err := h.useCase.GetBudgetSummaries(c.Context(), userID, dateStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"details":    err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Get budget summaries",
		"data":       summaries,
		"request_id": c.Locals("request_id"),
	})
}
