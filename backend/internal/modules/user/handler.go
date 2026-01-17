package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(useCase UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new user account
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register Payload"
// @Success      201  {object}  RegisterResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.useCase.Register(c.Context(), &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": err.Error(),
			})
		}
		if err == ErrEmailTaken || err == ErrUsernameTaken {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": user})
}

// Login godoc
// @Summary      Login User
// @Description  Login and get JWT Token
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Payload"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request Body",
		})
	}

	resp, err := h.useCase.Login(c.Context(), &req)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
	})
}

// Logout godoc
// @Summary      Logout User
// @Description  Invalidate current session (client-side token deletion recommended)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]string "message: Logout successful"
// @Failure      401  {object}  map[string]string "error: Unauthorized"
// @Router       /users/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	// JWT is stateless, so logout is handled client-side
	// Client should delete token from storage (localStorage/sessionStorage)

	// Validate that user is authenticated
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Optional: Log logout event for audit trail
	// h.log.Info("User logged out", "user_id", userID)

	// Response success
	// Client MUST delete token after receiving this response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful. Please delete your access token.",
		"user_id": userID, // Optional: confirm which user logged out
	})
}

// GetMe godoc
// @Summary      Get Current User Profile
// @Description  Get currently logged in user profile based on JWT Token
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UserResponse
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/current [get]
func (h *Handler) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.useCase.GetMe(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}

// RegisterRoutes registers all user-related routes
func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	api := app.Group("/api/users")

	// Public routes
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)

	// Protected routes - require authentication
	api.Post("/logout", authMiddleware, h.Logout)
	api.Get("/current", authMiddleware, h.GetMe)
}
