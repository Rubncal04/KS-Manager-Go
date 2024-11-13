package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func Register(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		newUser := new(UserRequest)
		if err := c.Bind(newUser); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if newUser.Password != newUser.ConfirmPassword {
			errorRes := SimpleReponse{
				Message: "Passwords do not match",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		role, err := repo.FindRoleByName("assistant")

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error finding role for this user",
				Code:    500,
			}
			log.Println(err)

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		user := models.User{
			Name:     newUser.Name,
			Email:    newUser.Email,
			UserName: newUser.UserName,
			ChurchId: newUser.ChurchId,
			Password: newUser.Password,
			RoleId:   int(role.ID),
		}

		result, err := repo.CreateUser(&user)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error registering",
				Code:    500,
			}
			log.Println(err)

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		token, er := middleware.GenerateJWT(result, *repo)

		if er != nil {
			errorRes := SimpleReponse{
				Message: "Error generating token",
				Code:    500,
			}
			log.Println(er)

			c.Response().WriteHeader(http.StatusInternalServerError)
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
		acceptHeader := c.Request().Header.Get(echo.HeaderAccept)
		user := new(LoginRequest)
		if err := c.Bind(user); err != nil {
			errorRes := SimpleReponse{
				Message: "There're some required fields",
				Code:    400,
			}

			if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
				return c.Render(http.StatusBadRequest, "error.html", errorRes)
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		userReq := models.User{
			UserName: user.UserName,
		}

		loggedUser, err := repo.FindUserBy(&userReq, userReq.UserName, "user_name")

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Something went wrong login user",
				Code:    500,
			}
			log.Println(err)

			if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
				return c.Render(http.StatusInternalServerError, "error.html", errorRes)
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		confirm := models.CheckPasswordHash(user.Password, loggedUser.Password)

		if !confirm {
			errorRes := SimpleReponse{
				Message: "Your password is incorrect",
				Code:    401,
			}
			log.Println(err)

			if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
				return c.Render(http.StatusBadRequest, "error.html", errorRes)
			}
			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		token, er := middleware.GenerateJWT(loggedUser, *repo)

		if er != nil {
			errorRes := SimpleReponse{
				Message: "Error generating token",
				Code:    500,
			}
			if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
				return c.Render(http.StatusInternalServerError, "error.html", errorRes)
			}
			log.Println(er)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := UserReponse{
			Name:     loggedUser.Name,
			UserName: loggedUser.UserName,
			Token:    token,
		}

		if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
			log.Println("Accept header html")
			return c.Render(http.StatusOK, "user.html", res)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}
