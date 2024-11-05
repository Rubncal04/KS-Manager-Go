package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/models"
	"github.com/labstack/echo/v4"
)

type RoleRequest struct {
	Name        string             `json:"name"`
	Permissions models.Permissions `json:"permissions"`
}

type RoleResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func CreateRole(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := new(RoleRequest)

		if err := c.Bind(role); err != nil {
			errorRes := SimpleReponse{
				Message: "Some fields are required",
				Code:    500,
			}
			log.Println(err)

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		newRole := models.Role{
			Name:        role.Name,
			Permissions: role.Permissions,
		}

		result, err := repo.CreateRole(&newRole)
		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error creating role",
				Code:    500,
			}
			log.Println(err)

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := RoleResponse{
			ID:   int64(result.ID),
			Name: result.Name,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func GetAllRoles(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		roles, err := repo.FindAllRoles()

		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error finding role",
				Code:    500,
			}
			log.Println(err)

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		var res []RoleResponse

		for _, r := range roles {
			role := RoleResponse{
				ID:   int64(r.ID),
				Name: r.Name,
			}
			res = append(res, role)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func DeleteRole(repo *db.PostgresRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		deletedRole, err := repo.DeleteRole(id)
		if err != nil {
			errorRes := SimpleReponse{
				Message: "Error deleting role",
				Code:    500,
			}
			log.Println(err)

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(errorRes)
		}

		res := SimpleReponse{
			Message: fmt.Sprintf("Role with ID %s was deleted successfully", deletedRole),
			Code:    201,
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}
