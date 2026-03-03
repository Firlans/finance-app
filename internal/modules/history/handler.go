package history

import "github.com/gofiber/fiber/v2"

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// CreateHistory godoc
// @Summary      Create Transaction History
// @Description  Record a new expense for a specific budget
// @Tags         History
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateHistoryRequest true "History Payload"
// @Success      201  {object}  History
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /histories [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CreateHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	history, err := h.useCase.CreateHistory(c.Context(), userID, &req)
	if err != nil {
		if err == ErrForbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": history})
}

// UpdateHistory godoc
// @Summary      Update Transaction History
// @Description  Update existing expense record
// @Tags         History
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "History ID"
// @Param        request body UpdateHistoryRequest true "Update Payload"
// @Success      200  {object}  History
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /histories/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	historyID := c.Params("id")
	if historyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "History ID required"})
	}

	var req UpdateHistoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	history, err := h.useCase.UpdateHistory(c.Context(), userID, historyID, &req)
	if err != nil {
		if err == ErrHistoryNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		if err == ErrForbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": history})
}

// ListHistory godoc
// @Summary      List Transaction Histories
// @Description  Get list of expenses for a specific budget
// @Tags         History
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id query string true "Account ID (UUID)"
// @Success      200  {object}  []History
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /histories [get]
func (h *Handler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	req := &ListHistoryRequest{
		AccountID: c.Query("account_id"),
	}

	histories, err := h.useCase.GetHistories(c.Context(), userID, req)
	if err != nil {
		if err == ErrForbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if histories == nil {
		histories = []*History{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": histories})
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	api := app.Group("/api/histories", authMiddleware)

	api.Post("/", h.Create)
	api.Put("/:id", h.Update)
	api.Get("/", h.List)
}
