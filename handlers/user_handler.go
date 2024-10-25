package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/middleware"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type UserRequest struct {
	Name            string `json:"name"`
	UserName        string `json:"user_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	ChurchId        int    `json:"church_id"`
}

type UserReponse struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Token    string `json:"token"`
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func Register(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		newUser := new(UserRequest)
		if err := c.Bind(newUser); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if newUser.Password != newUser.ConfirmPassword {
			return echo.NewHTTPError(http.StatusBadRequest, "Mismatch passwords")
		}

		user := models.User{
			Name:     newUser.Name,
			Email:    newUser.Email,
			UserName: newUser.UserName,
			ChurchId: newUser.ChurchId,
			Password: newUser.Password,
		}

		result, err := repo.CreateUser(&user)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error registering",
				Code:    500,
			}
			log.Println(err)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		token, er := middleware.GenerateJWT(result)

		if er != nil {
			errorRes := SimpleReponse{
				Message: "Error generating token",
				Code:    500,
			}
			log.Println(er)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := UserReponse{
			Name:     result.Name,
			UserName: result.UserName,
			Token:    token,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func Login(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(LoginRequest)
		if err := c.Bind(user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		userReq := models.User{
			UserName: user.UserName,
		}

		loggedUser, err := repo.FindUser(&userReq)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Something went wrong login user",
				Code:    500,
			}
			log.Println(err)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		confirm := models.CheckPasswordHash(user.Password, loggedUser.Password)

		if confirm == false {
			errorRes := SimpleReponse{
				Message: "Your password is incorrect",
				Code:    401,
			}
			log.Println(err)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		token, er := middleware.GenerateJWT(loggedUser)

		if er != nil {
			errorRes := SimpleReponse{
				Message: "Error generating token",
				Code:    500,
			}
			log.Println(er)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := UserReponse{
			Name:     loggedUser.Name,
			UserName: loggedUser.UserName,
			Token:    token,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}
