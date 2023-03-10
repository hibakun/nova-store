package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hibakun/nova-store/config"
	"github.com/hibakun/nova-store/models/model"
	"github.com/hibakun/nova-store/utils"
	"time"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identify string `form:"identify" validate:"required"`
		Password string `form:"password" validate:"required,min=6"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Password string `json:"password"`
	}

	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "check your input",
			"data":    err,
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error",
			"error":   err.Error(),
		})
	}

	identify := input.Identify
	password := input.Password
	email, phoneNumber, con := new(model.User), new(model.User), *new(bool)
	var user UserData

	if utils.ValidEmail(identify) {
		email, con = utils.GetUserByEmail(identify)
		if !con {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "user not found",
			})
		}
		user = UserData{
			ID:       email.ID,
			Password: email.Password,
		}
	} else {
		phoneNumber, con = utils.GetUserByPhoneNumber(identify)
		if !con {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "user not found",
			})
		}
		user = UserData{
			ID:       phoneNumber.ID,
			Password: phoneNumber.Password,
		}
	}

	if !utils.CheckHashPassword(user.Password, password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid password",
		})
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"identity": identify,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	token, err := t.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to create token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "USER_SESSION",
		Value:    token,
		MaxAge:   3600 * 24 * 7,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		SameSite: "lax",
	})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "success login",
		"token":   token,
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "USER_SESSION",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		SameSite: "lax",
	})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "success logout",
	})
}

func Validate(c *fiber.Ctx) error {
	cookie := c.Cookies("USER_SESSION")

	claims, err := utils.DecodeJWT(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "unauthorized",
		})
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims["id"],
		"identity": claims["identity"],
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	token, errToken := t.SignedString([]byte(config.Config("SECRET")))
	if errToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to create token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "USER_SESSION",
		Value:    token,
		MaxAge:   3600 * 24 * 7,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		SameSite: "lax",
	})

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "success validation",
		"token":   token,
	})
}
