package payments

import (
	"strconv"

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
	paymentsGroup := app.Group("api/payments", authMiddleware)

	paymentsGroup.Get("/", h.getPayments)
	paymentsGroup.Get("/:id", h.getPaymentByID)
	paymentsGroup.Post("/", h.createPayment)
	paymentsGroup.Put("/:id", h.updatePayment)
	paymentsGroup.Delete("/:id", h.deletePayment)
}

func (h *Handler) getPayments(c *fiber.Ctx) error {
	loanIDQuery := c.Query("loan_id")
	if loanIDQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "loan_id query parameter is required",
			"request_id": c.Locals("request_id"),
		})
	}

	loanID, err := strconv.Atoi(loanIDQuery)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid loan_id parameter",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetPaymentsByLoanID(c.Context(), loanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get payments",
		"data":    res,
	})
}

func (h *Handler) getPaymentByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid payment ID",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetPaymentByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	if res == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "Payment not found",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get payment by ID",
		"data":    res,
	})
}

func (h *Handler) createPayment(c *fiber.Ctx) error {
	var req CreatePaymentRequest
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

	id, err := h.useCase.Save(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create payment",
		"data":    CreatePaymentResponse{ID: id},
	})
}

func (h *Handler) updatePayment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid payment ID",
			"request_id": c.Locals("request_id"),
		})
	}

	var req UpdatePaymentRequest
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

	if err := h.useCase.Update(c.Context(), id, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update payment",
	})
}

func (h *Handler) deletePayment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid payment ID",
			"request_id": c.Locals("request_id"),
		})
	}

	if err := h.useCase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete payment",
	})
}
