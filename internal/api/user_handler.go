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
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)

	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrPassword) {
			_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
			return
		}
	}

	_ = utils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT IMPLEMENT")
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("TODO - NOT IMPLEMENT")
}
