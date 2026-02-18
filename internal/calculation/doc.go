// Package calculation implements entry totals for custom form entries.
//
// RunEntryCalculation computes field totals (base, GST, total) from form definition and raw values.
// It calculates totals for all fields, handles GST inclusive/exclusive, and builds BAS mapping.
//
// # Persistence (normalized path)
//
// Create/update entry runs RunEntryCalculation, then ConvertJSONBToNormalized maps the result
// to NormalizedEntry; the adapter persists to normalized entry tables (metadata, field values,
// field calculations, summary, and deductions).
// Recalculation from DB: POST /custom-form/entries/:entryId/recalculate loads stored values
// and deductions, runs RunEntryCalculation again, and UpdateEntry rewrites normalized child tables.
package calculation
