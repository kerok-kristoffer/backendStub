-- name: CreateIngredientTag :one
INSERT INTO ingredient_tags (
                             name, user_id
) VALUES ($1, $2) RETURNING *;

-- name: GetIngredientTag :one
SELECT * FROM ingredient_tags
WHERE id = $1 LIMIT 1;

-- name: GetIngredientTagByName :one
SELECT * FROM ingredient_tags
WHERE name ILIKE $1 AND user_id = $2 LIMIT 1;

-- name: ListIngredientTags :many
SELECT * FROM ingredient_tags
WHERE user_id = $1;

-- name: ListIngredientTagsByIngredientId :many
SELECT i."name" as "ingredient_name",
       i.id as "ingredient_id",
       it."name" as "ingredient_tag",
       it.id as "ingredient_tags_id",
       itm.id as "ingredient_tag_maps_id"
FROM "public".ingredient_tag_maps itm
         JOIN ingredients i on i.id = itm.ingredient_id
         JOIN ingredient_tags it on it.id = itm.ingredient_tag_id
WHERE i.id = $1
ORDER BY it.name;

-- name: UpdateIngredientTag :one
UPDATE ingredient_tags
SET (name) =
        ($2)
WHERE id = $1
RETURNING *;

-- name: DeleteIngredientTag :exec
DELETE FROM ingredient_tags
WHERE id = $1;