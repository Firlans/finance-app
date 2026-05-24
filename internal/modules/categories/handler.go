package categories

import "github.com/gofiber/fiber/v2"

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler, rateLimiter fiber.Handler) {
	categoriesGroup := app.Group("api/categories", authMiddleware)

	categoriesGroup.Get("/", h.getCategories)
	categoriesGroup.Get("/:id", h.getCategoryByID)
	categoriesGroup.Post("/", h.createCategory)
	categoriesGroup.Put("/:id", h.updateCategory)
	categoriesGroup.Delete("/:id", h.deleteCategory)
}

func (h *Handler) getCategories(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetCategories(c.Context(), &userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get all categories",
		"data":    res,
	})
}
func (h *Handler) getCategoryByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid category ID",
			"request_id": c.Locals("request_id"),
		})
	}

	res, err := h.useCase.GetCategoryByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	if res == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "Category not found",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get category by ID",
		"data":    res,
	})
}
func (h *Handler) createCategory(c *fiber.Ctx) error {
	var req CreateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
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

	category := &Category{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
	}

	err := h.useCase.Save(c.Context(), category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create category",
		"data":    category.ID,
	})
}
func (h *Handler) updateCategory(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid category ID",
			"request_id": c.Locals("request_id"),
		})
	}

	var req UpdateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	category := &Category{
		ID:          id,
		Name:        *req.Name,
		Description: req.Description,
		UserID:      userID,
	}
	err = h.useCase.UpdateCategory(c.Context(), category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update category",
	})
}
func (h *Handler) deleteCategory(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "User ID not found in JWT token",
			"request_id": c.Locals("request_id"),
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid category ID",
			"request_id": c.Locals("request_id"),
		})
	}

	err = h.useCase.DeleteCategory(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      err.Error(),
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delete category",
	})
}
