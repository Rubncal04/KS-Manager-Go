package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/middleware"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type OfferingRequest struct {
	Name       string    `json:"name"`
	Date       time.Time `json:"date"`
	CategoryId int       `json:"category_id"`
	Value      int       `json:"value"`
}

type OfferingResponse struct {
	ID         int       `json:"id"`
	ChurchId   int       `json:"church_id"`
	WorshipId  int       `json:"worship_id"`
	CategoryId int       `json:"category_id"`
	Name       string    `json:"name"`
	Date       time.Time `json:"date"`
	Value      int       `json:"value"`
}

func CreateOffering(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		offer := new(OfferingRequest)
		wId := c.Param("worship_id")
		worshipId, _ := strconv.Atoi(wId)
		user := middleware.Authentication(c)

		if err := c.Bind(offer); err != nil {
			log.Println("Error:", err)
			errorRes := SimpleReponse{
				Message: "There're some required fields",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		offering := models.Offering{
			Name:             offer.Name,
			ChurchId:         user.ChurchId,
			CategoryId:       offer.CategoryId,
			WorshipServiceId: worshipId,
			Value:            offer.Value,
			Date:             offer.Date,
		}

		newOffering, err := repo.CreateOffering(&offering)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error creating Offering",
				Code:    500,
			}
			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := OfferingResponse{
			ID:         int(newOffering.ID),
			ChurchId:   newOffering.ChurchId,
			WorshipId:  newOffering.WorshipServiceId,
			CategoryId: newOffering.CategoryId,
			Name:       newOffering.Name,
			Value:      newOffering.Value,
			Date:       newOffering.Date,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func FindAllOffering(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		worshipId := c.Param("worship_id")
		user := middleware.Authentication(c)
		churchId := fmt.Sprintf("%d", user.ChurchId)

		offerings, err := repo.FindAllOffering(worshipId, churchId)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error fetching Offerings",
				Code:    500,
			}
			log.Println(err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		var res []OfferingResponse

		for _, o := range offerings {
			offer := OfferingResponse{
				ID:         int(o.ID),
				ChurchId:   o.ChurchId,
				WorshipId:  o.WorshipServiceId,
				CategoryId: o.CategoryId,
				Name:       o.Name,
				Value:      o.Value,
				Date:       o.Date,
			}

			res = append(res, offer)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}
