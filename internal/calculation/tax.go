package calculation

type TaxResult struct {
	GSTAmount   float64
	GrossAmount float64
	NetAmount   float64
}

type Tax interface {
	Inclusive(amount float64) TaxResult
	Exclusive(amount float64) TaxResult
	Manual(amount float64, manualGst *float64) TaxResult
}

type GSTCalculator struct {
	Rate float64
}

func NewGSTCalculator(rate float64) Tax {
	return &GSTCalculator{
		Rate: rate,
	}
}

// Exclusive implements [Tax].
func (g *GSTCalculator) Exclusive(amount float64) TaxResult {
	gst := amount * g.Rate / (1 + g.Rate)
	net := amount - gst

	return TaxResult{
		GSTAmount:   gst,
		NetAmount:   net,
		GrossAmount: amount,
	}
}

// Inclusive implements [Tax].
func (g *GSTCalculator) Inclusive(amount float64) TaxResult {
	gst := amount * g.Rate
	gross := amount + gst

	return TaxResult{
		GSTAmount:   gst,
		NetAmount:   amount,
		GrossAmount: gross,
	}
}

// Manual implements [Tax].
func (g *GSTCalculator) Manual(amount float64, manualGst *float64) TaxResult {
	manualGstAmount := 0.0
	if manualGst != nil {
		manualGstAmount = *manualGst
	}
	return TaxResult{
		GSTAmount:   manualGstAmount,
		NetAmount:   amount - manualGstAmount,
		GrossAmount: amount + manualGstAmount,
	}
}
