package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type (
	ChurchRequest struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	ChurchResponse struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Address string `json:"address"`
	}
)

func Index(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := repo.FindAllChurches()

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching churches")
		}

		var res []ChurchResponse

		for _, row := range result {
			chur := ChurchResponse{
				ID:      row.ID,
				Name:    row.Name,
				Address: row.Address,
			}

			res = append(res, chur)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}

func Create(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		church := new(ChurchRequest)
		if err := c.Bind(church); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		newChurch := models.Church{
			Name:    church.Name,
			Address: church.Address,
		}

		result, err := repo.CreateChurch(&newChurch)

		res := ChurchResponse{
			ID:      result.ID,
			Name:    result.Name,
			Address: result.Address,
		}

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching churches")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func FindOne(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		result, err := repo.FindOneChurch(id)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching church with id: "+id)
		}

		res := ChurchResponse{
			ID:      result.ID,
			Name:    result.Name,
			Address: result.Address,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}
