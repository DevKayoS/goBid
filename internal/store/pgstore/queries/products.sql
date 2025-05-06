-- name: CreateProduct :one
INSERT INTO products (
    seller_id,
    product_name,
    description,
    baseprice,
    auction_end
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: ListAvailableProducts :many
SELECT * FROM products WHERE auction_end >= NOW() and is_sold = false;