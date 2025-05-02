package product

import (
	"context"
	"time"

	"github.com/DevKayoS/goBid/internal/validator"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

func (req CreateProductRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	// TODO - terminar a validcao
	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field cannot be blank")
	eval.CheckField(validator.NotBlank(req.Description), "description", "this field cannot be blank")

	eval.CheckField(
		validator.MinChars(req.Description, 10) &&
			validator.MaxChars(req.Description, 255),
		"description", "this field must have a length between 10 and 255",
	)

	eval.CheckField()

	return eval
}
