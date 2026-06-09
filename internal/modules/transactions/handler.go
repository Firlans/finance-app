package transactions

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase  UseCase
	validate *validator.Validate
}

func NewHandler(useCase UseCase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: validate,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler, rateLimiter fiber.Handler) {
	transactionsGroup := app.Group("api/transactions", authMiddleware)

	transactionsGroup.Get("/", h.getTransactions)
	transactionsGroup.Get("/:id", h.getTransactionByID)
	transactionsGroup.Post("/", h.createTransaction)
	transactionsGroup.Put("/:id", h.updateTransaction)
	transactionsGroup.Delete("/:id", h.deleteTransaction)
}

func (h *Handler) getTransactions(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetTransactions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get all transactions",
		"data":    res,
	})
}

func (h *Handler) getTransactionByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid transaction ID",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetTransactionByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	if res == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "Transaction not found",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get transaction by ID",
		"data":    res,
	})
}

func (h *Handler) createTransaction(c *fiber.Ctx) error {
	var req CreateTransactionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"details":    validationErrors.Error(),
				"request_id": c.Locals("request_id"),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	now := time.Now().UTC()
	categoryID := req.CategoryID
	transaction := &Transaction{
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		Description:     req.Description,
		AccountID:       req.AccountID,
		CategoryID:      &categoryID,
		UserID:          userID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := h.useCase.Save(c.Context(), transaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create transaction",
	})
}

func (h *Handler) updateTransaction(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid transaction ID",
			"request_id": c.Locals("request_id"),
		})
	}

	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Transaction ID is required",
			"request_id": c.Locals("request_id"),
		})
	}

	var req UpdateTransactionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"message":    validationErrors.Error(),
				"request_id": c.Locals("request_id"),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}
	transaction := &Transaction{
		ID:              id,
		Amount:          *req.Amount,
		TransactionType: *req.TransactionType,
		Description:     *req.Description,
		CategoryID:      req.CategoryID,
		AccountID:       *req.AccountID,
	}

	err = h.useCase.UpdateTransaction(c.Context(), transaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update transaction",
	})
}

func (h *Handler) deleteTransaction(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid transaction ID",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.DeleteTransaction(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete transaction",
	})
}
