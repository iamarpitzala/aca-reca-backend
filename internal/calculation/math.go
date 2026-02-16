package calculation

import (
	"math"
	"strconv"
)

func round2(n float64) float64 { return math.Round(n*100) / 100 }

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
