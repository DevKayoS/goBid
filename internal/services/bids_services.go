package services

import (
	"context"
	"errors"

	"github.com/DevKayoS/goBid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidServices struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsServices(pool *pgxpool.Pool) BidServices {
	return BidServices{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrBidIsTooLow = errors.New("the bid value is too low")

func (bs *BidServices) PlaceBid(ctx context.Context, product_Id, bidder_Id uuid.UUID, amount float64) (pgstore.Bid, error) {
	//Amount > previus_amount
	//Amount > baseprice
	product, err := bs.queries.GetProductById(ctx, product_Id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	highestBid, err := bs.queries.GetHighestBidByProductId(ctx, product.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	if product.Baseprice >= amount || highestBid.BidAmount >= amount {
		return pgstore.Bid{}, ErrBidIsTooLow
	}

	highestBid, err = bs.queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: product_Id,
		BidderID:  bidder_Id,
		BidAmount: amount,
	})

	if err != nil {
		return pgstore.Bid{}, err
	}

	return highestBid, nil
}
