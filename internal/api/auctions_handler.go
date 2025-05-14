package api

import (
	"errors"
	"net/http"

	"github.com/DevKayoS/goBid/internal/services"
	"github.com/DevKayoS/goBid/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductId := chi.URLParam(r, "product_Id")

	productId, err := uuid.Parse(rawProductId)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"msg": "invalid product id - must be a valid uuid",
		})
		return
	}

	_, err = api.ProductService.GetProductById(r.Context(), productId)

	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			utils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"msg": "no product with given id",
			})
			return
		}
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"msg": "unexpected error, try again later",
		})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)

	if !ok {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"msg": "unexpected error, try again later",
		})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productId]
	api.AuctionLobby.Unlock()

	if !ok {
		utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"msg": "the auction has ended",
		})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"msg": "could not upgrade connection to a websocket protocol",
		})
		return
	}

	client := services.NewClient(room, conn, userId)

	room.Register <- client
	go client.ReadEventLoop()
	go client.WriteEventLoop()
}
