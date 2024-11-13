package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateCategory(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		category := new(CategoryRequest)

		if err := c.Bind(category); err != nil {
			errorRes := SimpleReponse{
				Message: "There're some required fields",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		newCategory := models.Category{
			Name: category.Name,
		}

		result, err := repo.CreateCategory(&newCategory)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error creating category",
				Code:    500,
			}
			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := CategoryResponse{
			ID:   int(result.ID),
			Name: result.Name,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func FindAllCategory(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		categories, err := repo.FindAllCategory()

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error getting categories",
				Code:    500,
			}
			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		var res []CategoryResponse

		for _, c := range categories {
			cat := CategoryResponse{
				ID:   int(c.ID),
				Name: c.Name,
			}

			res = append(res, cat)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}
