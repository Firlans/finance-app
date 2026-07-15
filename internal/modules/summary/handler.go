package summary

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler, rateLimiter fiber.Handler) {
	summaryGroup := app.Group("api/summary", authMiddleware)
	summaryGroup.Get("/", h.GetSummary)
}

// GetSummary godoc
// @Summary      Get transactions summary
// @Description  Get transactions summary by module (debit/credit) and date range
// @Tags         Summary
// @Accept       json
// @Produce      json
// @Param        module query string true "Module (debit/credit)" Enums(debit, credit)
// @Param        from query string true "Start date (YYYY-MM-DD)"
// @Param        to query string true "End date (YYYY-MM-DD)"
// @Success      200  {object}  TransactionCategoryBalanceListResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /summary [get]
func (h *Handler) GetSummary(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	module := c.Query("module")
	if module != string(ModuleDebit) && module != string(ModuleCredit) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid module. Must be 'debit' or 'credit'",
			"request_id": c.Locals("request_id"),
		})
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Query params 'from' and 'to' are required",
			"request_id": c.Locals("request_id"),
		})
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid 'from'. Expected format YYYY-MM-DD",
			"request_id": c.Locals("request_id"),
		})
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid 'to'. Expected format YYYY-MM-DD",
			"request_id": c.Locals("request_id"),
		})
	}

	fromUTC := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toUTC := time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, time.UTC)

	balances, err := h.useCase.GetSummaryByModule(c.Context(), userID, Module(module), fromUTC, toUTC)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(TransactionCategoryBalanceListResponse{
		Message:   "Get transactions summary",
		Data:      balances,
		RequestID: c.Locals("request_id").(string),
	})
}
