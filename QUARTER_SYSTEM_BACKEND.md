# Backend Quarter System Implementation

## Overview

The backend now supports **system-driven quarter calculation** based on clinic financial settings. Quarters are calculated dynamically, never stored manually.

## New Endpoints

### Calculate Quarters for Clinic
```
GET /api/v1/clinic/:clinicId/quarters?yearsBack=1&yearsForward=1
```

**Query Parameters:**
- `yearsBack` (optional, default: 1) - Number of years to look back
- `yearsForward` (optional, default: 1) - Number of years to look forward

**Response:**
```json
{
  "data": [
    {
      "id": "FY2025-07-Q1",
      "quarterIndex": 0,
      "label": "Q1 (1 Jul – 30 Sep 2025)",
      "startDate": "2025-07-01T00:00:00Z",
      "endDate": "2025-09-30T00:00:00Z",
      "status": "open",
      "financialYearStart": "2025-07-01T00:00:00Z",
      "financialYearEnd": "2026-06-30T00:00:00Z"
    },
    ...
  ],
  "message": "quarters calculated successfully"
}
```

### Get Quarter for Date
```
GET /api/v1/clinic/:clinicId/quarter/date?date=2025-02-13
```

**Query Parameters:**
- `date` (required) - Date in RFC3339 format (e.g., `2025-02-13T00:00:00Z`) or date format (e.g., `2025-02-13`)

**Response:**
```json
{
  "data": {
    "id": "FY2025-07-Q3",
    "quarterIndex": 2,
    "label": "Q3 (1 Jan – 31 Mar 2026)",
    "startDate": "2026-01-01T00:00:00Z",
    "endDate": "2026-03-31T00:00:00Z",
    "status": "open",
    "financialYearStart": "2025-07-01T00:00:00Z",
    "financialYearEnd": "2026-06-30T00:00:00Z"
  },
  "message": "quarter found successfully"
}
```

## Implementation Details

### Quarter Calculation Logic

Quarters are calculated deterministically:

```go
quarter_start = financial_year_start + (quarter_index * 3 months)
quarter_end   = quarter_start + 3 months - 1 day
```

**Quarter Index:**
- 0 → Q1
- 1 → Q2
- 2 → Q3
- 3 → Q4

### Financial Year Start

- **JULY** (Australian default): Financial year runs July 1 - June 30
- **JANUARY** (Calendar year): Financial year runs January 1 - December 31

### Quarter Status

- **open**: Editable, current or future quarter
- **locked**: Finalised, locked by lockDate
- **draft**: Contains lockDate but not fully locked

Status is determined by comparing quarter dates to the clinic's `lockDate` in financial settings.

### Quarter ID Format

Deterministic ID format: `FY{year}-{month}-Q{quarter}`

Example: `FY2025-07-Q1` (Financial year starting July 2025, Quarter 1)

### Quarter Label Format

Always formatted as: `Q{n} ({start_date} – {end_date} {year})`

Example: `Q3 (1 Jan – 31 Mar 2026)`

## Code Structure

### Utils Package
- `backend/internal/utils/quarter.go` - Core quarter calculation functions

### Use Case Layer
- `backend/internal/application/usecase/quarter.go` - Business logic for quarter operations
  - `CalculateQuartersForClinic()` - Calculate quarters for a clinic
  - `GetQuarterForDate()` - Find quarter containing a date

### HTTP Layer
- `backend/internal/http/quarter.go` - HTTP handlers
  - `CalculateForClinic()` - Handler for calculating quarters
  - `GetQuarterForDate()` - Handler for finding quarter by date

### Routes
- `backend/route/quarter/quarter.go` - Route registration
  - `/api/v1/clinic/:clinicId/quarters` - Calculate quarters
  - `/api/v1/clinic/:clinicId/quarter/date` - Get quarter for date

## Legacy Code Removal

All legacy quarter CRUD endpoints have been removed. The system now exclusively uses calculated quarters based on clinic financial settings.

**Removed endpoints:**
- `POST /api/v1/quarter` - Removed
- `GET /api/v1/quarter/:id` - Removed
- `PUT /api/v1/quarter/:id` - Removed
- `DELETE /api/v1/quarter/:id` - Removed
- `GET /api/v1/quarter` - Removed

**Note:** All quarter operations now use the clinic-based endpoints that calculate quarters from financial settings.

## Integration with Financial Settings

The quarter system automatically uses:
- `financial_year_start` from `tbl_clinic_financial_settings` (JULY or JANUARY)
- `lock_date` from `tbl_clinic_financial_settings` (determines quarter status)

If financial settings don't exist for a clinic, defaults are used:
- Financial Year Start: JULY (July-June)
- Lock Date: nil (all quarters open)

## Error Handling

- **Invalid clinic ID**: Returns 400 Bad Request
- **Financial settings not found**: Uses defaults (JULY, no lock date)
- **Invalid date format**: Returns 400 Bad Request
- **Quarter not found for date**: Returns 404 Not Found

## Example Usage

### Calculate quarters for a clinic
```bash
curl -X GET "http://localhost:8080/api/v1/clinic/{clinicId}/quarters?yearsBack=1&yearsForward=1" \
  -H "Authorization: Bearer {token}"
```

### Find quarter for a specific date
```bash
curl -X GET "http://localhost:8080/api/v1/clinic/{clinicId}/quarter/date?date=2025-02-13" \
  -H "Authorization: Bearer {token}"
```

## Migration Notes

1. **Frontend Integration**: Frontend should use the new calculated quarter endpoints
2. **Database**: The `tbl_quarter` table is no longer used by the application. It can be dropped in a future migration if desired
3. **Legacy Code**: All legacy quarter CRUD operations have been removed from the codebase
4. **QuarterRepository**: The `QuarterRepository` interface and implementation remain in the codebase but are no longer used. They can be removed in a future cleanup if needed
