// Package calculation implements entry totals and deductions for custom form entries.
//
// # Calculation methods: NET vs GROSS
//
// The form's calculation method (NET or GROSS) determines how deductions are applied:
//
//   - NET: commission-based (e.g. independent contractors). RunEntryCalculation produces
//     commission, gstOnCommission, totalPaymentReceived; persisted in tbl_entry_net_details.
//
//   - GROSS: service/facility fee based on a percentage of net fee. RunEntryCalculation produces
//     serviceFeeBase, gstOnServiceFee, totalServiceFee, subtotalAfterDeductions, remittedAmount,
//     reductionBreakdown (expense GST paid by clinic), reimbursementBreakdown (expense paid by owner),
//     additionalReductionBreakdown (reduction-section fields), and optional outwork charges.
//     These are persisted in tbl_entry_gross_details, tbl_entry_gross_reduction,
//     tbl_entry_gross_reimbursement, tbl_entry_gross_additional_reduction,
//     tbl_entry_gross_reductions_summary, and tbl_entry_gross_outwork.
//
// # When gross method runs
//
// Gross calculations run when:
//   - Form calculation method is GROSS and form type has income (INCOME or BOTH), and
//   - Net method is not forced (e.g. by commission in deductions).
//
// Deductions can override: serviceFacilityFeePercent, serviceFeeOverride, entryPaymentResponsibility.
// Outwork is applied when form has outwork enabled and a rate; expense GST is then consolidated
// into a single outwork charge.
//
// # Persistence (normalized path)
//
// Create/update entry runs RunEntryCalculation, then ConvertJSONBToNormalized maps the result
// to NormalizedEntry; the adapter writes to tbl_entry_header, tbl_entry_field_value,
// tbl_entry_field_calculation, tbl_entry_summary, gross/net tables, and tbl_entry_deductions.
// Recalculation from DB: POST /custom-form/entries/:entryId/recalculate loads stored values
// and deductions, runs RunEntryCalculation again, and UpdateEntry rewrites normalized child tables.
package calculation
