package services

import (
	"errors"
	"numtostr/gotodo/app/dal"
	"numtostr/gotodo/app/types"
	"numtostr/gotodo/utils"
	"numtostr/gotodo/utils/jwt"
	"numtostr/gotodo/utils/password"

	"github.com/gofiber/fiber"
	"gorm.io/gorm"
)

// Login service logs in a user
func Login(ctx *fiber.Ctx) {
	b := new(types.LoginDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		ctx.Next(err)
		return
	}

	u := &types.UserResponse{}

	err := dal.FindUserByEmail(u, b.Email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Next(fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password"))
		return
	}

	if err := password.Verify(u.Password, b.Password); err != nil {
		ctx.Next(fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password"))
		return
	}

	t := jwt.Generate(&jwt.TokenPayload{
		ID: u.ID,
	})

	ctx.JSON(&types.AuthResponse{
		User: u,
		Auth: &types.AccessResponse{
			Token: t,
		},
	})
}

// Signup service creates a user
func Signup(ctx *fiber.Ctx) {
	b := new(types.SignupDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		ctx.Next(err)
		return
	}

	err := dal.FindUserByEmail(&struct{ ID string }{}, b.Email).Error

	// If email already exists, return
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Next(fiber.NewError(fiber.StatusConflict, "Email already exists"))
		return
	}

	user := &dal.User{
		Name:     b.Name,
		Password: password.Generate(b.Password),
		Email:    b.Email,
	}

	// Create a user, if error return
	if err := dal.CreateUser(user); err.Error != nil {
		ctx.Next(fiber.NewError(fiber.StatusConflict, err.Error.Error()))
		return
	}

	// generate access token
	t := jwt.Generate(&jwt.TokenPayload{
		ID: user.ID,
	})

	ctx.JSON(&types.AuthResponse{
		User: &types.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Auth: &types.AccessResponse{
			Token: t,
		},
	})
}