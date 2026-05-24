package user

import (
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra/middleware"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	user, err := h.useCase.Register(c.Context(), &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"details":    err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}
		if err == ErrEmailTaken || err == ErrUsernameTaken {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":      err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":       user,
		"request_id": c.Locals("request_id"),
	})
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
			"error":      "Invalid request Body",
			"request_id": c.Locals("request_id"),
		})
	}

	resp, err := h.useCase.Login(c.Context(), &req)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       resp,
		"request_id": c.Locals("request_id"),
	})
}

// Logout godoc
// @Summary      Logout User
// @Description  Invalidate current session by blacklisting token
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]string "message: Logout successful"
// @Failure      401  {object}  map[string]string "error: Unauthorized"
// @Router       /users/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	// Validate that user is authenticated
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"request_id": c.Locals("request_id"),
		})
	}

	// Get token from context (set by auth middleware)
	token, ok := c.Locals("access_token").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Token not found",
			"request_id": c.Locals("request_id"),
		})
	}

	// Get token expiry
	tokenExp, ok := c.Locals("token_exp").(time.Time)
	if !ok {
		// Default to 24 hours if not found
		tokenExp = time.Now().Add(24 * time.Hour)
	}

	// Add token to blacklist
	middleware.GetBlacklist().Add(token, tokenExp)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Logout successful. Token has been invalidated.",
		"user_id":    userID,
		"request_id": c.Locals("request_id"),
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"request_id": c.Locals("request_id"),
		})
	}

	resp, err := h.useCase.GetMe(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       resp,
		"request_id": c.Locals("request_id"),
	})
}

func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"request_id": c.Locals("request_id"),
		})
	}

	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	err := h.useCase.ResetPassword(c.Context(), userID, &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"details":    err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}
		if err == ErrInvalidCurrentPassword {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Password reset successful",
		"request_id": c.Locals("request_id"),
	})
}

func (h *Handler) ForgetPassword(c *fiber.Ctx) error {
	var req ForgetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	err := h.useCase.ForgetPassword(c.Context(), &req)
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
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Password reset link sent to your email",
		"request_id": c.Locals("request_id"),
	})
}

func (h *Handler) ResetPasswordWithToken(c *fiber.Ctx) error {
	var req ResetPasswordWithTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid request body",
			"request_id": c.Locals("request_id"),
		})
	}

	err := h.useCase.ResetPasswordWithToken(c.Context(), &req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      "Validation failed",
				"details":    err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}
		if err == ErrInvalidPasswordResetToken {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":      err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":      err.Error(),
				"request_id": c.Locals("request_id"),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Locals("request_id"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Password reset successful",
		"request_id": c.Locals("request_id"),
	})
}

// RefreshToken godoc
// @Summary      Refresh Access Token
// @Description  Issue a new access token using the current valid token. Old token is invalidated.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  LoginResponse
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/refresh-token [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Unauthorized",
			"request_id": c.Locals("request_id"),
		})
	}

	oldToken, ok := c.Locals("access_token").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":      "Token not found",
			"request_id": c.Locals("request_id"),
		})
	}

	tokenExp, ok := c.Locals("token_exp").(time.Time)
	if !ok {
		tokenExp = time.Now().Add(24 * time.Hour)
	}

	resp, err := h.useCase.RefreshToken(c.Context(), userID, oldToken, tokenExp)
	if err != nil {
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "Unauthorized",
				"request_id": c.Locals("request_id"),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Internal Server Error",
			"request_id": c.Locals("request_id"),
		})
	}

	// Blacklist the old token
	middleware.GetBlacklist().Add(oldToken, tokenExp)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       resp,
		"request_id": c.Locals("request_id"),
	})
}

// RegisterRoutes registers all user-related routes
func (h *Handler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler, authRateLimiter fiber.Handler) {
	api := app.Group("/api/users")

	// Public routes with stricter rate limiting
	api.Post("/register", authRateLimiter, h.Register)
	api.Post("/login", authRateLimiter, h.Login)
	api.Post("/forget-password", authRateLimiter, h.ForgetPassword)
	api.Post("/reset-password", authRateLimiter, h.ResetPasswordWithToken)

	// Protected routes - require authentication
	api.Post("/logout", authMiddleware, h.Logout)
	api.Get("/current", authMiddleware, h.GetMe)
	api.Post("/change-password", authMiddleware, h.ChangePassword)
	api.Post("/refresh-token", authMiddleware, h.RefreshToken)
}
