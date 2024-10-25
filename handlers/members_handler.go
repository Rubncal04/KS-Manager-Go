package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/middleware"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type MemberRequest struct {
	Name                 string             `json:"name"`
	LastName             string             `json:"last_name"`
	IdentificationNumber string             `json:"identification_number"`
	Address              string             `json:"address"`
	Email                string             `json:"email"`
	Birthday             string             `json:"birthday"`
	BaptizedBy           string             `json:"baptized_by"`
	BaptizedOn           string             `json:"baptized_on"`
	HolySpiritOn         string             `json:"holy_spirit_on"`
	Position             string             `json:"position"`
	NumChildren          int                `json:"num_children"`
	ChildrenNames        models.StringArray `json:"children_names"`
	PartnerName          string             `json:"partner_name"`
	Degree               string             `json:"degree"`
	Profession           string             `json:"profession"`
}

type MemberReponse struct {
	ID                   int                `json:"id"`
	ChurchId             int                `json:"church_id"`
	Name                 string             `json:"name"`
	LastName             string             `json:"last_name"`
	IdentificationNumber string             `json:"identification_number"`
	Address              string             `json:"address"`
	Email                string             `json:"email"`
	Birthday             string             `json:"birthday"`
	BaptizedBy           string             `json:"baptized_by"`
	BaptizedOn           string             `json:"baptized_on"`
	HolySpiritOn         string             `json:"holy_spirit_on"`
	Position             string             `json:"position"`
	NumChildren          int                `json:"num_children"`
	ChildrenNames        models.StringArray `json:"children_names"`
	PartnerName          string             `json:"partner_name"`
	Degree               string             `json:"degree"`
	Profession           string             `json:"profession"`
}

type SimpleReponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func IndexMembers(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := middleware.Authentication(c)
		id := fmt.Sprintf("%d", user.ChurchId)
		result, err := repo.FindAllMembers(id)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error fetching members",
				Code:    500,
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		var members []MemberReponse

		for _, row := range result {
			mem := MemberReponse{
				ID:                   int(row.ID),
				ChurchId:             row.ChurchId,
				Name:                 row.Name,
				LastName:             row.LastName,
				IdentificationNumber: row.IdentificationNumber,
				Address:              row.Address,
				Email:                row.Email,
				Birthday:             row.Birthday,
				BaptizedBy:           row.BaptizedBy,
				BaptizedOn:           row.BaptizedOn,
				HolySpiritOn:         row.HolySpiritOn,
				Position:             row.Position,
				NumChildren:          row.NumChildren,
				ChildrenNames:        row.ChildrenNames,
				PartnerName:          row.PartnerName,
				Degree:               row.Degree,
				Profession:           row.Profession,
			}

			members = append(members, mem)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return c.JSON(http.StatusOK, members)
	}
}

func FindOneMember(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		result, err := repo.FindOneMember(id)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "There was an error finding member",
				Code:    500,
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := MemberReponse{
			ID:                   int(result.ID),
			Name:                 result.Name,
			ChurchId:             result.ChurchId,
			LastName:             result.LastName,
			IdentificationNumber: result.IdentificationNumber,
			Address:              result.Address,
			Email:                result.Email,
			Birthday:             result.Birthday,
			Position:             result.Position,
			BaptizedOn:           result.BaptizedOn,
			BaptizedBy:           result.BaptizedBy,
			HolySpiritOn:         result.HolySpiritOn,
			Degree:               result.Degree,
			Profession:           result.Profession,
			PartnerName:          result.PartnerName,
			NumChildren:          result.NumChildren,
			ChildrenNames:        result.ChildrenNames,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func CreateMember(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		member := new(MemberRequest)
		user := middleware.Authentication(c)
		churchId := user.ChurchId

		if err := c.Bind(member); err != nil {
			errorRes := SimpleReponse{
				Message: "There're some required fields.",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		newMember := models.Member{
			Name:                 member.Name,
			LastName:             member.LastName,
			IdentificationNumber: member.IdentificationNumber,
			Address:              member.Address,
			ChurchId:             churchId,
			Email:                member.Email,
			Birthday:             member.Birthday,
			Position:             member.Position,
			BaptizedOn:           member.BaptizedOn,
			BaptizedBy:           member.BaptizedBy,
			HolySpiritOn:         member.HolySpiritOn,
			Degree:               member.Degree,
			Profession:           member.Profession,
			PartnerName:          member.PartnerName,
			NumChildren:          member.NumChildren,
			ChildrenNames:        models.StringArray(member.ChildrenNames),
		}

		result, err := repo.CreateMember(&newMember)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error create member",
				Code:    500,
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := MemberReponse{
			ID:                   int(result.ID),
			Name:                 result.Name,
			ChurchId:             result.ChurchId,
			LastName:             result.LastName,
			IdentificationNumber: result.IdentificationNumber,
			Address:              result.Address,
			Email:                result.Email,
			Birthday:             result.Birthday,
			Position:             result.Position,
			BaptizedOn:           result.BaptizedOn,
			BaptizedBy:           result.BaptizedBy,
			HolySpiritOn:         result.HolySpiritOn,
			Degree:               result.Degree,
			Profession:           result.Profession,
			PartnerName:          result.PartnerName,
			NumChildren:          result.NumChildren,
			ChildrenNames:        result.ChildrenNames,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func UpdateMember(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		member := new(MemberRequest)
		id := c.Param("id")
		if err := c.Bind(member); err != nil {
			errorRes := SimpleReponse{
				Message: "There're required fields.",
				Code:    400,
			}

			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		fields := models.Member{
			Name:                 member.Name,
			LastName:             member.LastName,
			IdentificationNumber: member.IdentificationNumber,
			Address:              member.Address,
			Email:                member.Email,
			Birthday:             member.Birthday,
			Position:             member.Position,
			BaptizedOn:           member.BaptizedOn,
			BaptizedBy:           member.BaptizedBy,
			HolySpiritOn:         member.HolySpiritOn,
			Degree:               member.Degree,
			Profession:           member.Profession,
			PartnerName:          member.PartnerName,
			NumChildren:          member.NumChildren,
			ChildrenNames:        models.StringArray(member.ChildrenNames),
		}

		result, err := repo.UpdateMember(id, &fields)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error update member",
				Code:    500,
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := MemberReponse{
			ID:                   int(result.ID),
			Name:                 result.Name,
			ChurchId:             result.ChurchId,
			LastName:             result.LastName,
			IdentificationNumber: result.IdentificationNumber,
			Address:              result.Address,
			Email:                result.Email,
			Birthday:             result.Birthday,
			Position:             result.Position,
			BaptizedOn:           result.BaptizedOn,
			BaptizedBy:           result.BaptizedBy,
			HolySpiritOn:         result.HolySpiritOn,
			Degree:               result.Degree,
			Profession:           result.Profession,
			PartnerName:          result.PartnerName,
			NumChildren:          result.NumChildren,
			ChildrenNames:        result.ChildrenNames,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func DeleteMember(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		result, err := repo.DeleteMember(id)

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error deleting member",
				Code:    500,
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := SimpleReponse{
			Message: fmt.Sprintf("Member with ID %s was deleted successfully", result),
			Code:    202,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}
