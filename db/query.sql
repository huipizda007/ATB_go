-- name: GetFruit :one
SELECT * FROM fruits
WHERE id = $1 LIMIT 1;

-- name: ListFruits :many
SELECT * FROM fruits
ORDER BY name;

-- name: CreateFruit :one
INSERT INTO fruits (
    name, brand, price_per_kg, stock_kg
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateFruitStock :one
UPDATE fruits
SET stock_kg = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateFruitPrice :one
UPDATE fruits
SET price_per_kg = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteFruit :exec
DELETE FROM fruits
WHERE id = $1;

-- name: GetCheapFruits :many
SELECT * FROM fruits
WHERE price_per_kg < $1
ORDER BY price_per_kg ASC;