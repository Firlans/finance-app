package transaction

import "github.com/gofiber/fiber/v2"

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// CreateTransaction godoc
// @Summary      Create Transaction
// @Description  Record a new expense for a specific account
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTransactionRequest true "Transaction Payload"
// @Success      201  {object}  Transaction
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /transactions [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	transaction, err := h.useCase.CreateTransaction(c.Context(), userID, &req)
	if err != nil {
		if err == ErrForbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": transaction})
}

// UpdateTransaction godoc
// @Summary      Update Transaction
// @Description  Update existing expense record
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Transaction ID"
// @Param        request body UpdateTransactionRequest true "Update Payload"
// @Success      200  {object}  Transaction
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /transactions/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Transaction ID required"})
	}
	var req UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	transaction, err := h.useCase.UpdateTransaction(c.Context(), userID, transactionID, &req)
	if err != nil {
		if err == ErrTransactionNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		if err == ErrForbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": transaction})
}

// ListTransactions godoc
// @Summary      List Transactions
// @Description  Get list of expenses for a specific account
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        account_id query string true "Account ID (UUID)"
// @Success      200  {object}  []Transaction
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /transactions [get]
func (h *Handler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	req := &ListTransactionRequest{
		AccountID: c.Query("account_id"),
	}
	transactions, err := h.useCase.GetTransactions(c.Context(), userID, req)
	if err != nil {
		if err == ErrForbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if transactions == nil {
		transactions = []*Transaction{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": transactions})
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	api := app.Group("/api/transactions", authMiddleware)

	api.Post("/", h.Create)
	api.Put("/:id", h.Update)
	api.Get("/", h.List)
}
