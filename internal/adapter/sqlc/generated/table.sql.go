// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: table.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkTableHasOrders = `-- name: CheckTableHasOrders :one
SELECT EXISTS(
    SELECT 1 FROM orders 
    WHERE table_id = $1
) as has_orders
`

func (q *Queries) CheckTableHasOrders(ctx context.Context, tableID int32) (bool, error) {
	row := q.db.QueryRow(ctx, checkTableHasOrders, tableID)
	var has_orders bool
	err := row.Scan(&has_orders)
	return has_orders, err
}

const createTable = `-- name: CreateTable :one
INSERT INTO tables (table_number, qr_code, seating, is_active)
VALUES ($1, $2, $3, $4)
RETURNING id, table_number, qr_code, seating, is_active, created_at, updated_at
`

type CreateTableParams struct {
	TableNumber int32       `json:"table_number"`
	QrCode      string      `json:"qr_code"`
	Seating     int32       `json:"seating"`
	IsActive    pgtype.Bool `json:"is_active"`
}

func (q *Queries) CreateTable(ctx context.Context, arg CreateTableParams) (*Table, error) {
	row := q.db.QueryRow(ctx, createTable,
		arg.TableNumber,
		arg.QrCode,
		arg.Seating,
		arg.IsActive,
	)
	var i Table
	err := row.Scan(
		&i.ID,
		&i.TableNumber,
		&i.QrCode,
		&i.Seating,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const deleteTable = `-- name: DeleteTable :exec
DELETE FROM tables
WHERE id = $1
`

func (q *Queries) DeleteTable(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteTable, id)
	return err
}

const getTableByID = `-- name: GetTableByID :one
SELECT id, table_number, qr_code, seating, is_active, created_at, updated_at FROM tables
WHERE id = $1
`

func (q *Queries) GetTableByID(ctx context.Context, id int32) (*Table, error) {
	row := q.db.QueryRow(ctx, getTableByID, id)
	var i Table
	err := row.Scan(
		&i.ID,
		&i.TableNumber,
		&i.QrCode,
		&i.Seating,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getTableByNumber = `-- name: GetTableByNumber :one
SELECT id, table_number, qr_code, seating, is_active, created_at, updated_at FROM tables
WHERE table_number = $1
`

func (q *Queries) GetTableByNumber(ctx context.Context, tableNumber int32) (*Table, error) {
	row := q.db.QueryRow(ctx, getTableByNumber, tableNumber)
	var i Table
	err := row.Scan(
		&i.ID,
		&i.TableNumber,
		&i.QrCode,
		&i.Seating,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getTableByQRCode = `-- name: GetTableByQRCode :one
SELECT id, table_number, qr_code, seating, is_active, created_at, updated_at FROM tables
WHERE qr_code = $1
`

func (q *Queries) GetTableByQRCode(ctx context.Context, qrCode string) (*Table, error) {
	row := q.db.QueryRow(ctx, getTableByQRCode, qrCode)
	var i Table
	err := row.Scan(
		&i.ID,
		&i.TableNumber,
		&i.QrCode,
		&i.Seating,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listTables = `-- name: ListTables :many
SELECT id, table_number, qr_code, seating, is_active, created_at, updated_at FROM tables
ORDER BY table_number ASC
`

func (q *Queries) ListTables(ctx context.Context) ([]*Table, error) {
	rows, err := q.db.Query(ctx, listTables)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Table{}
	for rows.Next() {
		var i Table
		if err := rows.Scan(
			&i.ID,
			&i.TableNumber,
			&i.QrCode,
			&i.Seating,
			&i.IsActive,
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

const updateTable = `-- name: UpdateTable :one
UPDATE tables
SET 
    table_number = $2,
    qr_code = $3,
    seating = $4,
    is_active = $5
WHERE id = $1
RETURNING id, table_number, qr_code, seating, is_active, created_at, updated_at
`

type UpdateTableParams struct {
	ID          int32       `json:"id"`
	TableNumber int32       `json:"table_number"`
	QrCode      string      `json:"qr_code"`
	Seating     int32       `json:"seating"`
	IsActive    pgtype.Bool `json:"is_active"`
}

func (q *Queries) UpdateTable(ctx context.Context, arg UpdateTableParams) (*Table, error) {
	row := q.db.QueryRow(ctx, updateTable,
		arg.ID,
		arg.TableNumber,
		arg.QrCode,
		arg.Seating,
		arg.IsActive,
	)
	var i Table
	err := row.Scan(
		&i.ID,
		&i.TableNumber,
		&i.QrCode,
		&i.Seating,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
