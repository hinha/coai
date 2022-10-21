package adapters

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	app "github.com/hinha/coai/core/users/application"
	"github.com/hinha/coai/core/users/application/command"
	"github.com/hinha/coai/core/users/application/query"
	"github.com/hinha/coai/core/users/domain"
	"time"
)

func NewUserHTTP(app app.Application) *UserHTTP {
	return &UserHTTP{app: app}
}

type UserHTTP struct {
	app app.Application
}

// UserAll method to get all users.
// @Description login authority.
// @Summary get all users
// @Success 200 {string} status "ok"
// @Router /v1/users/all [get]
func (h *UserHTTP) UserAll(c *fiber.Ctx) error {
	users, err := h.app.Queries.AllUsers.Handle(c.UserContext(), query.AllUsers{})
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(fiber.Map{"msg": "ok", "data": users})
}

// Register method to register users.
// @Description login authority.
// @Summary register users
// @Success 200 {string} status "ok"
// @Router /v1/users/register [get]
func (h *UserHTTP) Register(c *fiber.Ctx) error {
	err := h.app.Commands.Register.Handle(c.UserContext(), command.RegisterUser{User: domain.User{
		LastLogin:    time.Now(),
		UserGroupsID: 1,
	}})
	if err != nil {
		return c.JSON(fiber.Map{"error": err})
	}

	return c.JSON(fiber.Map{"msg": "ok"})
}
