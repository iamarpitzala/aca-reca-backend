package service

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
// 	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
// 	"github.com/jmoiron/sqlx"
// )

// type FinancialCalculationService struct {
// 	db *sqlx.DB
// }

// func NewFinancialCalculationService(db *sqlx.DB) *FinancialCalculationService {
// 	return &FinancialCalculationService{
// 		db: db,
// 	}
// }

// // CalculateFinancial performs calculation based on form configuration
// func (fcs *FinancialCalculationService) CalculateFinancial(ctx context.Context, formID uuid.UUID, input domain.CalculationInput, userID *uuid.UUID) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	// Get form
// 	form, err := repository.GetFinancialFormByID(ctx, fcs.db, formID)
// 	if err != nil {
// 		return nil, nil, errors.New("financial form not found")
// 	}

// 	if !form.IsActive {
// 		return nil, nil, errors.New("financial form is not active")
// 	}

// 	// Validate input
// 	if err := ValidateCalculationInput(form.Configuration, input, form.CalculationMethod); err != nil {
// 		return nil, nil, err
// 	}

// 	var gst *domain.GST
// 	if form.GSTID != nil {
// 		gst, err = repository.GetGSTByID(ctx, fcs.db, *form.GSTID)
// 		form.GST = gst
// 		if err != nil {
// 			return nil, nil, errors.New("GST not found")
// 		}
// 	}

// 	var result *domain.CalculationResult
// 	var basMapping *domain.BASMapping

// 	// Perform calculation based on method
// 	if form.CalculationMethod == "net" {
// 		result, basMapping, err = fcs.calculateNetMethod(form.Configuration, input)
// 	} else if form.CalculationMethod == "gross" {
// 		result, basMapping, err = fcs.calculateGrossMethod(form, input)
// 	} else {
// 		return nil, nil, errors.New("invalid calculation method")
// 	}

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Optionally save calculation history
// 	if userID != nil {
// 		inputMap, _ := structToMap(input)
// 		resultMap, _ := structToMap(result)
// 		basMap, _ := structToMap(basMapping)

// 		calculation := &domain.FinancialCalculation{
// 			ID:              uuid.New(),
// 			FinancialFormID: formID,
// 			InputData:       inputMap,
// 			CalculatedData:  resultMap,
// 			BASMapping:      basMap,
// 			CreatedAt:       time.Now(),
// 			CreatedBy:       userID,
// 		}

// 		_ = repository.CreateFinancialCalculation(ctx, fcs.db, calculation)
// 	}

// 	return result, basMapping, nil
// }

// // calculateNetMethod performs Net Method calculations
// func (fcs *FinancialCalculationService) calculateNetMethod(config map[string]interface{}, input domain.CalculationInput) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	// Parse config
// 	configJSON, _ := json.Marshal(config)
// 	var netConfig domain.NetMethodConfig
// 	if err := json.Unmarshal(configJSON, &netConfig); err != nil {
// 		return nil, nil, errors.New("invalid net method configuration")
// 	}

// 	// Set defaults
// 	if netConfig.CommissionPercent == 0 {
// 		netConfig.CommissionPercent = 40.0
// 	}
// 	if netConfig.GSTOnCommissionPercent == 0 {
// 		netConfig.GSTOnCommissionPercent = 10.0
// 	}
// 	if netConfig.SuperHoldingEnabled && netConfig.SuperPercent == 0 {
// 		netConfig.SuperPercent = 12.0
// 	}

// 	result := &domain.CalculationResult{}
// 	basMapping := &domain.BASMapping{}

// 	// A = Gross Patient Fee
// 	A := input.GrossPatientFee

// 	// B = Lab Fee (if enabled)
// 	B := 0.0
// 	if netConfig.LabFeeEnabled {
// 		B = input.LabFee
// 	}

// 	// C = Net Patient Fee
// 	C := A - B

// 	// D = Commission for Dentist
// 	D := C * (netConfig.CommissionPercent / 100.0)

// 	if netConfig.SuperHoldingEnabled {
// 		// A2. With Super Holding
// 		// F = Commission Component = D ÷ 1.12
// 		F := D / (1 + netConfig.SuperPercent/100.0)

// 		// E = Super Component = F × 12%
// 		E := F * (netConfig.SuperPercent / 100.0)

// 		// G = Total for Reconciliation = E + F
// 		G := E + F

// 		// H = GST on Commission = F × GST%
// 		H := F * (netConfig.GSTOnCommissionPercent / 100.0)

// 		// I = Total Payment to Dentist = F + H
// 		I := F + H

// 		result.CommissionComponent = F
// 		result.SuperComponent = E
// 		result.TotalForReconciliation = G
// 		result.GSTOnCommission = H
// 		result.TotalCommission = I

// 		// BAS Mapping
// 		basMapping.Field1A = H
// 		basMapping.FieldG1 = I
// 	} else {
// 		// A1. Without Super Holding
// 		// E = GST on Commission = D × GST%
// 		E := D * (netConfig.GSTOnCommissionPercent / 100.0)

// 		// F = Total Commission = D + E
// 		F := D + E

// 		result.CommissionForDentist = D
// 		result.GSTOnCommission = E
// 		result.TotalCommission = F

// 		// BAS Mapping
// 		basMapping.Field1A = E
// 		basMapping.FieldG1 = F
// 	}

// 	result.NetPatientFee = C

// 	return result, basMapping, nil
// }

// // calculateGrossMethod performs Gross Method calculations
// func (fcs *FinancialCalculationService) calculateGrossMethod(form *domain.FinancialForm, input domain.CalculationInput) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	// Parse config
// 	configJSON, _ := json.Marshal(form.Configuration)
// 	var grossConfig domain.GrossMethodConfig
// 	if err := json.Unmarshal(configJSON, &grossConfig); err != nil {
// 		return nil, nil, errors.New("invalid gross method configuration")
// 	}

// 	if grossConfig.OutworkRatePercent == nil {
// 		grossConfig.OutworkRatePercent = 40
// 	}

// 	// A = Gross Patient Fee
// 	A := input.GrossPatientFee

// 	// Determine which variant to use
// 	// if grossConfig.OutworkChargeRateEnabled {
// 	// 	return fcs.calculateGrossMethodB5(grossConfig, input, A)
// 	// } else if grossConfig.GSTOnPatientFee && grossConfig.LabFeePaidBy == "dentist" {
// 	// 	return fcs.calculateGrossMethodB4(grossConfig, input, A)
// 	// } else if grossConfig.MerchantFeeEnabled {
// 	// 	return fcs.calculateGrossMethodB3(grossConfig, input, A)
// 	// } else
// 	if grossConfig.GSTOnLabFee {
// 		return fcs.calculateGrossMethodB2(grossConfig, input, A)
// 	} else {
// 		return fcs.calculateGrossMethodB1(grossConfig, input, A)
// 	}
// }

// // B1. Standard (Lab Fee Paid by Clinic)
// func (fcs *FinancialCalculationService) calculateGrossMethodB1(config domain.GrossMethodConfig, input domain.CalculationInput, A float64) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	result := &domain.CalculationResult{}
// 	basMapping := &domain.BASMapping{}

// 	// B = Lab Fee (if enabled)
// 	B := 0.0
// 	if config.LabFeeEnabled {
// 		B = input.LabFee
// 	}

// 	// C = Net Patient Fee
// 	C := A - B

// 	// D = Service & Facility Fee
// 	D := C * (config.ServiceFeePercent / 100.0)

// 	// E = GST on Service Fee
// 	E := D * (config.GSTOnServiceFeePercent / 100.0)

// 	// F = Total Service Fee
// 	F := D + E

// 	// G = Amount Remitted to Dentist
// 	G := C - F

// 	result.ServiceFacilityFee = D
// 	result.GSTOnServiceFee = E
// 	result.TotalServiceFee = F
// 	result.AmountRemittedToDentist = G

// 	// BAS Mapping
// 	basMapping.Field1B = E
// 	basMapping.FieldG1 = A
// 	basMapping.FieldG11 = B + F

// 	return result, basMapping, nil
// }

// // B2. With GST on Lab Fee
// func (fcs *FinancialCalculationService) calculateGrossMethodB2(config domain.GrossMethodConfig, input domain.CalculationInput, A float64) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	result := &domain.CalculationResult{}
// 	basMapping := &domain.BASMapping{}

// 	// B = Lab Fee
// 	B := input.LabFee

// 	// G = GST on Lab Fee
// 	G := input.GSTOnLabFee

// 	// C = Net Patient Fee
// 	C := A - B

// 	// D = Service & Facility Fee
// 	D := C * (config.ServiceFeePercent / 100.0)

// 	// E = GST on Service Fee
// 	E := D * (config.GSTOnServiceFeePercent / 100.0)

// 	// F = Total Service Fee
// 	F := D + E

// 	// I = Amount Remitted to Dentist
// 	I := C - F - G

// 	result.ServiceFacilityFee = D
// 	result.GSTOnServiceFee = E
// 	result.TotalServiceFee = F
// 	result.AmountRemittedToDentist = I

// 	// BAS Mapping
// 	basMapping.Field1B = E + G
// 	basMapping.FieldG1 = A
// 	basMapping.FieldG11 = B + F + G

// 	return result, basMapping, nil
// }

// // B3. With Merchant Fee / Bank Fee
// func (fcs *FinancialCalculationService) calculateGrossMethodB3(config domain.GrossMethodConfig, input domain.CalculationInput, A float64) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	result := &domain.CalculationResult{}
// 	basMapping := &domain.BASMapping{}

// 	// B = Lab Fee (optional)
// 	B := 0.0
// 	if config.LabFeeEnabled {
// 		B = input.LabFee
// 	}

// 	// G = Merchant Fee (Incl GST)
// 	G := input.MerchantFeeInclGST

// 	// K = Bank Fee
// 	K := input.BankFee

// 	// C = Net Patient Fee
// 	C := A - B

// 	// D = Service & Facility Fee
// 	D := C * (config.ServiceFeePercent / 100.0)

// 	// E = GST on Service Fee
// 	E := D * (config.GSTOnServiceFeePercent / 100.0)

// 	// F = Total Service Fee
// 	F := D + E

// 	// Calculate Merchant Fee GST Component
// 	I := 0.0
// 	J := G
// 	if config.GSTOnMerchantFee {
// 		// I = G × (GST% × 100) / (100 + GST% × 100)
// 		gstPercent := config.GSTOnServiceFeePercent
// 		I = G * (gstPercent * 100.0) / (100.0 + gstPercent*100.0)
// 		J = G - I
// 	}

// 	// H = Amount Remitted to Dentist
// 	H := C - F - G - K

// 	result.ServiceFacilityFee = D
// 	result.GSTOnServiceFee = E
// 	result.TotalServiceFee = F
// 	result.MerchantFeeGSTComponent = I
// 	result.NetMerchantFee = J
// 	result.AmountRemittedToDentist = H

// 	// BAS Mapping
// 	basMapping.Field1B = E + I
// 	basMapping.FieldG1 = A
// 	basMapping.FieldG11 = F + G + K

// 	return result, basMapping, nil
// }

// // B4. GST on Patient Fee + Lab Fee Paid by Dentist
// func (fcs *FinancialCalculationService) calculateGrossMethodB4(config domain.GrossMethodConfig, input domain.CalculationInput, A float64) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	result := &domain.CalculationResult{}
// 	basMapping := &domain.BASMapping{}

// 	// B = GST on Patient Fee
// 	B := input.GSTOnPatientFee

// 	// I = Lab Fee Paid by Dentist
// 	I := input.LabFeePaidByDentist

// 	// C = Patient Fee Excl GST
// 	C := A - B

// 	// D = Lab Fee Paid by Dentist
// 	D := I

// 	// E = Net Patient Fee
// 	E := C - D

// 	// F = Service & Facility Fee
// 	F := E * (config.ServiceFeePercent / 100.0)

// 	// G = GST on Service Fee
// 	G := F * (config.GSTOnServiceFeePercent / 100.0)

// 	// H = Total Service Fee
// 	H := F + G

// 	// J = Amount Remitted to Dentist
// 	J := E - H + I + B

// 	result.ServiceFacilityFee = F
// 	result.GSTOnServiceFee = G
// 	result.TotalServiceFee = H
// 	result.AmountRemittedToDentist = J

// 	// BAS Mapping
// 	basMapping.Field1A = B
// 	basMapping.Field1B = G
// 	basMapping.FieldG1 = A
// 	basMapping.FieldG11 = H + I

// 	return result, basMapping, nil
// }

// // B5. Outwork Charge Rate
// func (fcs *FinancialCalculationService) calculateGrossMethodB5(config domain.GrossMethodConfig, input domain.CalculationInput, A float64) (*domain.CalculationResult, *domain.BASMapping, error) {
// 	result := &domain.CalculationResult{}
// 	basMapping := &domain.BASMapping{}

// 	// w = Lab Fee
// 	w := input.OutworkLabFee

// 	// x = GST on Lab Fee
// 	x := input.OutworkGSTOnLabFee

// 	// y = GST on Merchant Fee
// 	y := input.OutworkGSTOnMerchantFee

// 	// B = Total Outwork Charge
// 	B := w + x + y

// 	// C = Net Patient Fee
// 	C := A - B

// 	// D = Service & Facility Fee
// 	D := C * (config.ServiceFeePercent / 100.0)

// 	// E = Lab & Other Cost Charge
// 	E := B * (config.OutworkRatePercent / 100.0)

// 	// F = Total Service Fee + Other Charges
// 	F := D + E

// 	// G = GST on Service Fee
// 	G := F * (config.GSTOnServiceFeePercent / 100.0)

// 	// H = Total Charges incl GST
// 	H := F + G

// 	// I = Amount Remitted to Dentist
// 	I := A - H

// 	result.ServiceFacilityFee = D
// 	result.LabOtherCostCharge = E
// 	result.TotalServiceFeeOtherCharges = F
// 	result.GSTOnServiceFee = G
// 	result.TotalChargesInclGST = H
// 	result.AmountRemittedToDentist = I
// 	result.TotalOutworkCharge = B

// 	// BAS Mapping
// 	basMapping.Field1B = G
// 	basMapping.FieldG1 = A
// 	basMapping.FieldG11 = H

// 	return result, basMapping, nil
// }

// // GetCalculationHistory retrieves calculation history for a form
// func (fcs *FinancialCalculationService) GetCalculationHistory(ctx context.Context, formID uuid.UUID) ([]domain.FinancialCalculation, error) {
// 	return repository.GetCalculationsByFormID(ctx, fcs.db, formID)
// }
