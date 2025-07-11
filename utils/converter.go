// utils/converter.go (Updated with pgtype.Numeric handling)
package utils

import (
	"math"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToBool(pgBool pgtype.Bool) bool {
	if pgBool.Valid {
		return pgBool.Bool
	}
	return false
}

func ConvertToPGBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}

func ConvertToText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

func ConvertToPGNumericFromFloat(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Convert float64 to a big.Float for precise representation
	bf := new(big.Float).SetFloat64(f)
	// Convert big.Float to string with full precision
	str := bf.Text('f', -1) // 'f' format for decimal representation
	_ = n.Scan(str)
	return n
}

func ConvertToPGNumericFromString(s string) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(s) // เช่น "123.45"
	return n
}

func ConvertToPGNumericFromBigFloat(b *big.Float) pgtype.Numeric {
	var n pgtype.Numeric
	if b == nil {
		n.Valid = false
		return n
	}
	str := b.Text('f', -1) // convert to full-precision string
	_ = n.Scan(str)
	return n
}

// FromPgNumericToFloat converts pgtype.Numeric to float64
func FromPgNumericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	// สร้าง big.Float จาก big.Int
	val := new(big.Float).SetInt(n.Int)

	// 10^Exp เป็น big.Float
	scale := big.NewFloat(math.Pow10(int(n.Exp)))

	// คูณเข้าด้วยกัน
	val.Mul(val, scale)

	// แปลงเป็น float64
	f64, _ := val.Float64()
	return f64
}

// FromInterfaceToFloat converts interface{} (from SQLC) to float64
// This handles the case where SQLC returns interface{} for SUM/COALESCE results
func FromInterfaceToFloat(i interface{}) float64 {
	if i == nil {
		return 0
	}

	switch v := i.(type) {
	case pgtype.Numeric:
		return FromPgNumericToFloat(v)
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	default:
		// Try to convert pgtype.Numeric if it's wrapped in interface{}
		if numeric, ok := i.(pgtype.Numeric); ok {
			return FromPgNumericToFloat(numeric)
		}
		return 0
	}
}

func FromPgTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// ConvertToPGTimestamp converts *time.Time to pgtype.Timestamp
func ConvertToPGTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{
			Valid: false,
		}
	}
	return pgtype.Timestamp{
		Time:  *t,
		Valid: true,
	}
}

// ConvertFromPGTimestamp converts pgtype.Timestamp to *time.Time
func ConvertFromPGTimestamp(t pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
