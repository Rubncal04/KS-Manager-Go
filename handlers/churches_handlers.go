package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type (
	ChurchRequest struct {
		Name      string `json:"name"`
		Address   string `json:"address"`
		CountryId int    `json:"country_id"`
		StateId   int    `json:"state_id"`
		CityId    int    `json:"city_id"`
	}

	ChurchResponse struct {
		ID        uint   `json:"id"`
		CountryId int    `json:"country_id"`
		StateId   int    `json:"state_id"`
		CityId    int    `json:"city_id"`
		Name      string `json:"name"`
		Address   string `json:"address"`
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
				ID:        row.ID,
				Name:      row.Name,
				Address:   row.Address,
				CityId:    row.CityId,
				StateId:   row.StateId,
				CountryId: row.CountryId,
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
			Name:      church.Name,
			Address:   church.Address,
			CityId:    church.CityId,
			StateId:   church.StateId,
			CountryId: church.CountryId,
		}

		result, err := repo.CreateChurch(&newChurch)

		res := ChurchResponse{
			ID:        result.ID,
			Name:      result.Name,
			Address:   result.Address,
			CityId:    result.CityId,
			StateId:   result.StateId,
			CountryId: result.CountryId,
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
		acceptHeader := c.Request().Header.Get(echo.HeaderAccept)
		id := c.Param("id")
		result, err := repo.FindOneChurch(id)
		sess, _ := session.Get("church", c)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Something went wrong",
				Code:    400,
			}

			if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
				return c.Render(http.StatusBadRequest, "error.html", errorRes)
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := ChurchResponse{
			ID:        result.ID,
			Name:      result.Name,
			Address:   result.Address,
			CityId:    result.CityId,
			StateId:   result.StateId,
			CountryId: result.CountryId,
		}

		if acceptHeader == "" || strings.Contains(acceptHeader, echo.MIMETextHTML) {
			churchData, err := json.Marshal(res)
			if err != nil {
				log.Println("Error serializing user data:", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to serialize church data"})
			}
			sess.Values["church"] = string(churchData)
			sess.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusSeeOther, "/church")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}

func ShowChurch() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("church", c)

		userData, ok := sess.Values["church"].(string)
		if !ok || userData == "" {
			return c.Redirect(http.StatusSeeOther, "/user")
		}

		var userRes UserReponse
		if err := json.Unmarshal([]byte(userData), &userRes); err != nil {
			log.Println("Error deserializing user data:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to deserialize user data"})
		}

		return c.Render(http.StatusOK, "base.html", userRes)
	}
}

func Update(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		church := new(ChurchRequest)
		if err := c.Bind(church); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id := c.Param("id")
		fields := models.Church{
			Name:      church.Name,
			Address:   church.Address,
			CityId:    church.CityId,
			StateId:   church.StateId,
			CountryId: church.CountryId,
		}

		result, err := repo.UpdateChurch(id, &fields)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching church with id: "+id)
		}

		res := ChurchResponse{
			ID:        result.ID,
			Name:      result.Name,
			Address:   result.Address,
			CityId:    result.CityId,
			StateId:   result.StateId,
			CountryId: result.CountryId,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, res)
	}
}
