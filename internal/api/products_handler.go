package api

import (
	"context"
	"net/http"

	"github.com/DevKayoS/goBid/internal/services"
	"github.com/DevKayoS/goBid/internal/useCase/product"
	"github.com/DevKayoS/goBid/internal/utils"
	"github.com/google/uuid"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[product.CreateProductRequest](r)

	if err != nil {
		utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)

	if !ok {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error, try again later",
		})
		return
	}

	productId, err := api.ProductService.CreateProduct(
		r.Context(),
		userId,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)

	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product auction try again later",
		})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)
	auctionRoom := services.NewAuctionRoom(ctx, api.BidService, productId)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[productId] = auctionRoom
	api.AuctionLobby.Unlock()

	utils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"msg":        "Auction has started with sucess",
		"product_id": productId,
	})
}

func (api *Api) handleListAvailableProduct(w http.ResponseWriter, r *http.Request) {
	products, err := api.ProductService.ListAvailableProducts(r.Context())

	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error, try again later",
		})
		return
	}

	utils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"msg":  "generated product list with successfuly",
		"data": products,
	})

}
