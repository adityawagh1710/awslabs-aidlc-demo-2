package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"github.com/todo-app/auth-service/internal/service"
)

var validate = validator.New()

type AuthHandler struct{ svc service.AuthService }

func NewAuthHandler(svc service.AuthService) *AuthHandler { return &AuthHandler{svc} }

type registerRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	MFACode  string `json:"mfa_code"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type mfaVerifyRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	pair, err := h.svc.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			return fiber.NewError(fiber.StatusConflict, "email already registered")
		}
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(pair)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	pair, err := h.svc.Login(c.Context(), req.Email, req.Password, req.MFACode)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials, service.ErrInvalidMFA:
			return fiber.ErrUnauthorized
		case service.ErrAccountLocked:
			return fiber.ErrTooManyRequests
		case service.ErrMFARequired:
			return fiber.NewError(fiber.StatusUnprocessableEntity, "mfa_required")
		}
		return err
	}
	return c.JSON(pair)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	pair, err := h.svc.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	return c.JSON(pair)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req logoutRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	userID := c.Locals("userID").(string)
	accessToken := c.Get("Authorization")[len("Bearer "):]
	if err := h.svc.Logout(c.Context(), userID, accessToken, req.RefreshToken); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) MFAEnroll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	secret, qrURL, err := h.svc.EnrollMFA(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"secret": secret, "qr_url": qrURL})
}

func (h *AuthHandler) MFAVerify(c *fiber.Ctx) error {
	var req mfaVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	if err := validate.Struct(req); err != nil {
		return fiber.ErrUnprocessableEntity
	}
	userID := c.Locals("userID").(string)
	if err := h.svc.VerifyMFA(c.Context(), userID, req.Code); err != nil {
		return fiber.ErrUnauthorized
	}
	return c.SendStatus(fiber.StatusOK)
}
