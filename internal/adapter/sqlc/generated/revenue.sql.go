// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: revenue.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getDailyRevenue = `-- name: GetDailyRevenue :one

SELECT 
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at < $2
`

type GetDailyRevenueParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

// internal/adapter/sqlc/queries/revenue.sql
func (q *Queries) GetDailyRevenue(ctx context.Context, arg GetDailyRevenueParams) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, getDailyRevenue, arg.PaidAt, arg.PaidAt_2)
	var total_revenue pgtype.Numeric
	err := row.Scan(&total_revenue)
	return total_revenue, err
}

const getDailyRevenueRange = `-- name: GetDailyRevenueRange :many
SELECT 
    DATE(paid_at) as revenue_date,
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY DATE(paid_at)
ORDER BY revenue_date ASC
`

type GetDailyRevenueRangeParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

type GetDailyRevenueRangeRow struct {
	RevenueDate  pgtype.Date    `json:"revenue_date"`
	TotalRevenue pgtype.Numeric `json:"total_revenue"`
}

func (q *Queries) GetDailyRevenueRange(ctx context.Context, arg GetDailyRevenueRangeParams) ([]*GetDailyRevenueRangeRow, error) {
	rows, err := q.db.Query(ctx, getDailyRevenueRange, arg.PaidAt, arg.PaidAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetDailyRevenueRangeRow{}
	for rows.Next() {
		var i GetDailyRevenueRangeRow
		if err := rows.Scan(&i.RevenueDate, &i.TotalRevenue); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHourlyRevenue = `-- name: GetHourlyRevenue :many
SELECT 
    EXTRACT(HOUR FROM paid_at)::integer as hour,
    COALESCE(SUM(amount), 0)::decimal as total_revenue,
    COUNT(*)::bigint as transaction_count
FROM payments
WHERE DATE(paid_at) = DATE($1)
GROUP BY EXTRACT(HOUR FROM paid_at)
ORDER BY hour ASC
`

type GetHourlyRevenueRow struct {
	Hour             int32          `json:"hour"`
	TotalRevenue     pgtype.Numeric `json:"total_revenue"`
	TransactionCount int64          `json:"transaction_count"`
}

func (q *Queries) GetHourlyRevenue(ctx context.Context, date interface{}) ([]*GetHourlyRevenueRow, error) {
	rows, err := q.db.Query(ctx, getHourlyRevenue, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetHourlyRevenueRow{}
	for rows.Next() {
		var i GetHourlyRevenueRow
		if err := rows.Scan(&i.Hour, &i.TotalRevenue, &i.TransactionCount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonthlyRevenue = `-- name: GetMonthlyRevenue :one
SELECT 
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at < $2
`

type GetMonthlyRevenueParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

func (q *Queries) GetMonthlyRevenue(ctx context.Context, arg GetMonthlyRevenueParams) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, getMonthlyRevenue, arg.PaidAt, arg.PaidAt_2)
	var total_revenue pgtype.Numeric
	err := row.Scan(&total_revenue)
	return total_revenue, err
}

const getMonthlyRevenueRange = `-- name: GetMonthlyRevenueRange :many
SELECT 
    DATE_TRUNC('month', paid_at)::timestamp as revenue_month,
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY DATE_TRUNC('month', paid_at)
ORDER BY revenue_month ASC
`

type GetMonthlyRevenueRangeParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

type GetMonthlyRevenueRangeRow struct {
	RevenueMonth pgtype.Timestamp `json:"revenue_month"`
	TotalRevenue pgtype.Numeric   `json:"total_revenue"`
}

func (q *Queries) GetMonthlyRevenueRange(ctx context.Context, arg GetMonthlyRevenueRangeParams) ([]*GetMonthlyRevenueRangeRow, error) {
	rows, err := q.db.Query(ctx, getMonthlyRevenueRange, arg.PaidAt, arg.PaidAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetMonthlyRevenueRangeRow{}
	for rows.Next() {
		var i GetMonthlyRevenueRangeRow
		if err := rows.Scan(&i.RevenueMonth, &i.TotalRevenue); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRevenueByPaymentMethod = `-- name: GetRevenueByPaymentMethod :many
SELECT 
    method,
    COALESCE(SUM(amount), 0)::decimal as total_revenue,
    COUNT(*)::bigint as transaction_count
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
GROUP BY method
ORDER BY total_revenue DESC
`

type GetRevenueByPaymentMethodParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

type GetRevenueByPaymentMethodRow struct {
	Method           PaymentMethod  `json:"method"`
	TotalRevenue     pgtype.Numeric `json:"total_revenue"`
	TransactionCount int64          `json:"transaction_count"`
}

func (q *Queries) GetRevenueByPaymentMethod(ctx context.Context, arg GetRevenueByPaymentMethodParams) ([]*GetRevenueByPaymentMethodRow, error) {
	rows, err := q.db.Query(ctx, getRevenueByPaymentMethod, arg.PaidAt, arg.PaidAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRevenueByPaymentMethodRow{}
	for rows.Next() {
		var i GetRevenueByPaymentMethodRow
		if err := rows.Scan(&i.Method, &i.TotalRevenue, &i.TransactionCount); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRevenueStats = `-- name: GetRevenueStats :one
SELECT 
    COUNT(*)::bigint as total_transactions,
    COALESCE(SUM(amount), 0)::decimal as total_revenue,
    COALESCE(AVG(amount), 0)::decimal as average_transaction,
    COALESCE(MIN(amount), 0)::decimal as min_transaction,
    COALESCE(MAX(amount), 0)::decimal as max_transaction
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
`

type GetRevenueStatsParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

type GetRevenueStatsRow struct {
	TotalTransactions  int64          `json:"total_transactions"`
	TotalRevenue       pgtype.Numeric `json:"total_revenue"`
	AverageTransaction pgtype.Numeric `json:"average_transaction"`
	MinTransaction     pgtype.Numeric `json:"min_transaction"`
	MaxTransaction     pgtype.Numeric `json:"max_transaction"`
}

func (q *Queries) GetRevenueStats(ctx context.Context, arg GetRevenueStatsParams) (*GetRevenueStatsRow, error) {
	row := q.db.QueryRow(ctx, getRevenueStats, arg.PaidAt, arg.PaidAt_2)
	var i GetRevenueStatsRow
	err := row.Scan(
		&i.TotalTransactions,
		&i.TotalRevenue,
		&i.AverageTransaction,
		&i.MinTransaction,
		&i.MaxTransaction,
	)
	return &i, err
}

const getTopSellingItems = `-- name: GetTopSellingItems :many
SELECT 
    mi.name as item_name,
    SUM(oi.quantity)::bigint as total_quantity,
    SUM(oi.quantity * oi.unit_price)::decimal as total_revenue
FROM order_items oi
JOIN menu_items mi ON oi.item_id = mi.id
JOIN orders o ON oi.order_id = o.id
JOIN payments p ON o.id = p.order_id
WHERE p.paid_at >= $1 AND p.paid_at <= $2
GROUP BY mi.id, mi.name
ORDER BY total_revenue DESC
LIMIT $3
`

type GetTopSellingItemsParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
	Limit    int32            `json:"limit"`
}

type GetTopSellingItemsRow struct {
	ItemName      string         `json:"item_name"`
	TotalQuantity int64          `json:"total_quantity"`
	TotalRevenue  pgtype.Numeric `json:"total_revenue"`
}

func (q *Queries) GetTopSellingItems(ctx context.Context, arg GetTopSellingItemsParams) ([]*GetTopSellingItemsRow, error) {
	rows, err := q.db.Query(ctx, getTopSellingItems, arg.PaidAt, arg.PaidAt_2, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetTopSellingItemsRow{}
	for rows.Next() {
		var i GetTopSellingItemsRow
		if err := rows.Scan(&i.ItemName, &i.TotalQuantity, &i.TotalRevenue); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTotalRevenue = `-- name: GetTotalRevenue :one
SELECT 
    COALESCE(SUM(amount), 0)::decimal as total_revenue
FROM payments
WHERE paid_at >= $1 AND paid_at <= $2
`

type GetTotalRevenueParams struct {
	PaidAt   pgtype.Timestamp `json:"paid_at"`
	PaidAt_2 pgtype.Timestamp `json:"paid_at_2"`
}

func (q *Queries) GetTotalRevenue(ctx context.Context, arg GetTotalRevenueParams) (pgtype.Numeric, error) {
	row := q.db.QueryRow(ctx, getTotalRevenue, arg.PaidAt, arg.PaidAt_2)
	var total_revenue pgtype.Numeric
	err := row.Scan(&total_revenue)
	return total_revenue, err
}
