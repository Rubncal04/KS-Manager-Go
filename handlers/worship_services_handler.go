package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/middleware"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type WorshipService struct {
	Name string         `json:"name"`
	Day  models.WeekDay `json:"day"`
}

type WorshipServiceReponse struct {
	ID   int            `json:"id"`
	Name string         `json:"name"`
	Day  models.WeekDay `json:"day"`
}

func CreateWorshipService(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		worship := new(WorshipService)
		user := middleware.Authentication(c)

		if err := c.Bind(worship); err != nil {
			errorRes := SimpleReponse{
				Message: "There're some required fields",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		service := models.WorshipService{
			Name:     worship.Name,
			Day:      worship.Day,
			ChurchId: user.ChurchId,
		}

		if !service.IsValidWeekday() {
			errorRes := SimpleReponse{
				Message: "Invalid weekday",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		result, err := repo.CreateWorshipService(&service)

		res := WorshipServiceReponse{
			ID:   int(result.ID),
			Name: result.Name,
			Day:  result.Day,
		}

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error creating worship service",
				Code:    500,
			}
			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func IndexWorshipService(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := middleware.Authentication(c)
		id := fmt.Sprintf("%d", user.ChurchId)

		worships, err := repo.FindAllWorship(id)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error fetching worship services",
				Code:    500,
			}

			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		var foundWorship []WorshipServiceReponse

		for _, w := range worships {
			ws := WorshipServiceReponse{
				ID:   int(w.ID),
				Name: w.Name,
				Day:  w.Day,
			}

			foundWorship = append(foundWorship, ws)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, foundWorship)
	}
}

func UpdateWorshipService(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		worship := new(WorshipService)
		id := c.Param("id")

		if err := c.Bind(worship); err != nil {
			errorRes := SimpleReponse{
				Message: "There're some required fields",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		service := models.WorshipService{
			Name: worship.Name,
			Day:  worship.Day,
		}

		if !service.IsValidWeekday() {
			errorRes := SimpleReponse{
				Message: "Invalid weekday",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		updatedService, err := repo.UpdateWorship(id, &service)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error updating worship service",
				Code:    500,
			}

			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := WorshipServiceReponse{
			ID:   int(updatedService.ID),
			Name: updatedService.Name,
			Day:  updatedService.Day,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}

func DeleteWorshipService(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		deleted, err := repo.DeleteWorship(id)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error deleting worship service",
				Code:    500,
			}

			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := SimpleReponse{
			Message: fmt.Sprintf("Service with ID %s was deleted successfully", deleted),
			Code:    200,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}

func FindOneWorshipService(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		worshihp, err := repo.FindWorshipByID(id)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error finding worship service",
				Code:    500,
			}

			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := WorshipServiceReponse{
			ID:   int(worshihp.ID),
			Name: worshihp.Name,
			Day:  worshihp.Day,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}
