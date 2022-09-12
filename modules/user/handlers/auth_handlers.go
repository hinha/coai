package handlers

import "github.com/gofiber/fiber/v2"

type authHandler struct{}

// UserSignIn method to login user.
// @Description login user.
// @Summary login user
// @Success 200 {string} status "ok"
// @Router /v1/auth/login [post]
func (h *authHandler) SignIn(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"msg": "ok"})
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}
