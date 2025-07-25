// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: order_item.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrderItem = `-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, item_id, quantity, unit_price, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, order_id, item_id, quantity, unit_price, notes, created_at, updated_at
`

type CreateOrderItemParams struct {
	OrderID   int32          `json:"order_id"`
	ItemID    int32          `json:"item_id"`
	Quantity  int32          `json:"quantity"`
	UnitPrice pgtype.Numeric `json:"unit_price"`
	Notes     pgtype.Text    `json:"notes"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (*OrderItem, error) {
	row := q.db.QueryRow(ctx, createOrderItem,
		arg.OrderID,
		arg.ItemID,
		arg.Quantity,
		arg.UnitPrice,
		arg.Notes,
	)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ItemID,
		&i.Quantity,
		&i.UnitPrice,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const deleteOrderItem = `-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1
`

func (q *Queries) DeleteOrderItem(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteOrderItem, id)
	return err
}

const deleteOrderItemsByOrderID = `-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items
WHERE order_id = $1
`

func (q *Queries) DeleteOrderItemsByOrderID(ctx context.Context, orderID int32) error {
	_, err := q.db.Exec(ctx, deleteOrderItemsByOrderID, orderID)
	return err
}

const getOrderItemByID = `-- name: GetOrderItemByID :one
SELECT id, order_id, item_id, quantity, unit_price, notes, created_at, updated_at FROM order_items
WHERE id = $1
`

func (q *Queries) GetOrderItemByID(ctx context.Context, id int32) (*OrderItem, error) {
	row := q.db.QueryRow(ctx, getOrderItemByID, id)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ItemID,
		&i.Quantity,
		&i.UnitPrice,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getOrderItemByOrderAndItem = `-- name: GetOrderItemByOrderAndItem :one
SELECT id, order_id, item_id, quantity, unit_price, notes, created_at, updated_at FROM order_items
WHERE order_id = $1 AND item_id = $2
LIMIT 1
`

type GetOrderItemByOrderAndItemParams struct {
	OrderID int32 `json:"order_id"`
	ItemID  int32 `json:"item_id"`
}

func (q *Queries) GetOrderItemByOrderAndItem(ctx context.Context, arg GetOrderItemByOrderAndItemParams) (*OrderItem, error) {
	row := q.db.QueryRow(ctx, getOrderItemByOrderAndItem, arg.OrderID, arg.ItemID)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ItemID,
		&i.Quantity,
		&i.UnitPrice,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getOrderItemsByOrderID = `-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, item_id, quantity, unit_price, notes, created_at, updated_at FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetOrderItemsByOrderID(ctx context.Context, orderID int32) ([]*OrderItem, error) {
	rows, err := q.db.Query(ctx, getOrderItemsByOrderID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*OrderItem{}
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ItemID,
			&i.Quantity,
			&i.UnitPrice,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrderItem = `-- name: UpdateOrderItem :one
UPDATE order_items
SET 
    order_id = $2,
    item_id = $3,
    quantity = $4,
    unit_price = $5,
    notes = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, order_id, item_id, quantity, unit_price, notes, created_at, updated_at
`

type UpdateOrderItemParams struct {
	ID        int32          `json:"id"`
	OrderID   int32          `json:"order_id"`
	ItemID    int32          `json:"item_id"`
	Quantity  int32          `json:"quantity"`
	UnitPrice pgtype.Numeric `json:"unit_price"`
	Notes     pgtype.Text    `json:"notes"`
}

func (q *Queries) UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (*OrderItem, error) {
	row := q.db.QueryRow(ctx, updateOrderItem,
		arg.ID,
		arg.OrderID,
		arg.ItemID,
		arg.Quantity,
		arg.UnitPrice,
		arg.Notes,
	)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ItemID,
		&i.Quantity,
		&i.UnitPrice,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
