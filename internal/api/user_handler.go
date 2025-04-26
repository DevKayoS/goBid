package api

import (
	"errors"
	"net/http"

	"github.com/DevKayoS/goBid/internal/services"
	"github.com/DevKayoS/goBid/internal/useCase/user"
	"github.com/DevKayoS/goBid/internal/utils"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[user.CreateUserRequest](r)

	if err != nil {
		utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)

	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			utils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
			return
		}
	}

	utils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[user.LoginUserRequest](r)

	if err != nil {
		utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "Invalid credentials",
			})
			return
		}
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticatedUserId", id)

	utils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged in successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticatedUserId")
	utils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged out successfully",
	})

}
