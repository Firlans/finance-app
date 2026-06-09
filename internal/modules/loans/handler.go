package loans

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
	loansGroup := app.Group("api/loans", authMiddleware)

	loansGroup.Get("/", h.getLoans)
	loansGroup.Get("/:id", h.getLoanByID)
	loansGroup.Post("/", h.createLoan)
	loansGroup.Put("/:id", h.updateLoan)
	loansGroup.Delete("/:id", h.deleteLoan)
}

func (h *Handler) getLoans(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetLoans(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get all loans",
		"data":    res,
	})
}

func (h *Handler) getLoanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid loan ID",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetLoanByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	if res == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "Loan not found",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get loan by ID",
		"data":    res,
	})
}

func (h *Handler) createLoan(c *fiber.Ctx) error {
	var req CreateLoanRequest
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

	req.UserID = userID

	id, err := h.useCase.Save(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create loan",
		"data":    CreateLoanResponse{ID: id},
	})
}

func (h *Handler) updateLoan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid loan ID",
			"request_id": c.Locals("request_id"),
		})
	}

	var req UpdateLoanRequest
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
		"message": "Update loan",
	})
}

func (h *Handler) deleteLoan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid loan ID",
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
		"message": "Delete loan",
	})
}
