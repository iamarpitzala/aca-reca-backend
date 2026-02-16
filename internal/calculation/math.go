package calculation

import (
	"math"
	"strconv"
)

// round2 rounds to 2 decimal places. Uses math.Round (half to even).
func round2(n float64) float64 {
	if n == 0 {
		return 0
	}
	return math.Round(n*100) / 100
}

// round2HalfUp rounds to 2 decimals with 0.5 rounding up (standard for currency).
// Ensures e.g. 100.00 * 40% = 40.00 not 39.97 when float drift would give 39.996.
func round2HalfUp(n float64) float64 {
	if n == 0 {
		return 0
	}
	// Scale to cents, round 0.5 up, scale back
	cents := n * 100
	if cents < 0 {
		return math.Ceil(cents-0.5) / 100
	}
	return math.Floor(cents+0.5) / 100
}

// parseFloat parses a numeric value from JSON (number or string).
func parseFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	}
	return 0
}

func calcGST(amount, rate float64, gstType string, manualGst *float64) (base, gst, total float64) {
	if gstType == "manual" {
		m := 0.0
		if manualGst != nil {
			m = *manualGst
		}
		return round2(amount), round2(m), round2(amount+m)
	}
	if rate == 0 {
		return amount, 0, amount
	}
	rateDec := rate / 100
	if gstType == "inclusive" {
		gst = amount - (amount / (1 + rateDec))
		base = amount - gst
		return round2(base), round2(gst), amount
	}
	gst = amount * rateDec
	total = amount + gst
	return amount, round2(gst), round2(total)
}

func ptrFloat64(f float64) *float64 { return &f }
