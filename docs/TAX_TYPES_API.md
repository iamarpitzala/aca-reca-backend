# Tax Types API

Single endpoint for listing tax types (form field lookups).  
Table: **`tbl_tax_type`** (INCLUSIVE, EXCLUSIVE, MANUAL).

---

## Endpoint

| Method | Endpoint | Description |
|--------|----------|-------------|
| **GET** | `/api/v1/tax-types` | List all tax types |

**Authentication:** Required (Bearer token).

---

## GET /api/v1/tax-types

Returns all tax types. No path or query parameters. No request body.

### Request

- **Method:** `GET`
- **URL:** `/api/v1/tax-types`
- **Headers:** `Authorization: Bearer <token>`
- **Body:** None

### Response

**Success (200)**

Standard envelope:

| Field     | Type    | Description        |
|-----------|---------|--------------------|
| `success` | boolean | `true`             |
| `message` | string  | Success message    |
| `data`    | array   | List of tax types  |

Each item in `data` (TaxTypeResponse):

| Field        | Type    | Description                    |
|--------------|---------|--------------------------------|
| `id`         | number  | Tax type ID (1, 2, 3)          |
| `name`       | string  | Display name                   |
| `type`       | string  | `INCLUSIVE` \| `EXCLUSIVE` \| `MANUAL` |
| `description`| string  | Optional description           |
| `createdAt`  | string  | ISO 8601 timestamp             |
| `updatedAt`  | string  | ISO 8601 timestamp             |
| `deletedAt`  | string? | Null if not deleted           |

**Example response (200):**

```json
{
  "success": true,
  "message": "tax types retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Inclusive",
      "type": "INCLUSIVE",
      "description": "Inclusive tax type",
      "createdAt": "2026-02-18T07:15:54Z",
      "updatedAt": "2026-02-18T07:15:54Z",
      "deletedAt": null
    },
    {
      "id": 2,
      "name": "Exclusive",
      "type": "EXCLUSIVE",
      "description": "Exclusive tax type",
      "createdAt": "2026-02-18T07:15:54Z",
      "updatedAt": "2026-02-18T07:15:54Z",
      "deletedAt": null
    },
    {
      "id": 3,
      "name": "Manual",
      "type": "MANUAL",
      "description": "Manual tax type",
      "createdAt": "2026-02-18T07:15:54Z",
      "updatedAt": "2026-02-18T07:15:54Z",
      "deletedAt": null
    }
  ]
}
```

**Error (401):** Unauthorized – missing or invalid token.

**Error (500):** Internal server error; `success: false`, `message` and optional `error` set.

---

## Summary

| Method | Endpoint            | Request body | Response body      |
|--------|---------------------|--------------|--------------------|
| GET    | `/api/v1/tax-types` | None         | `{ success, message, data: TaxTypeResponse[] }` |
