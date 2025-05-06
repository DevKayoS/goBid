package api

import (
	"net/http"

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

	id, err := api.ProductService.CreateProduct(
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

	utils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"msg":        "product created with successs",
		"product_id": id,
	})
	return

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
