# Stakeholder-Organization Linkage Implementation Guide

## Overview
Organizations in the Tokenization Module can now be linked to existing stakeholders in the Trovo database. This provides a clear reference between organizations and their corresponding stakeholder entities.

---

## Database Changes

### Organizations Table
Two new columns have been added:
```sql
ALTER TABLE organizations
ADD COLUMN stakeholder_id BIGINT,
ADD COLUMN stakeholder_type VARCHAR(50);
```

- **`stakeholder_id`**: References the ID of the stakeholder in the Trovo database
- **`stakeholder_type`**: Specifies which stakeholder table to reference

---

## Valid Stakeholder Types

| Stakeholder Type Value | References Table | Example Use Case |
|------------------------|------------------|------------------|
| `asset_manager` | `asset_managers` | Asset management companies |
| `asset_issuing_house` | `asset_issuing_houses` | Issuing houses |
| `approved_asset_custodian` | `approved_asset_custodians` | Custodian services |
| `legal_and_professionals` | `legal_and_professionals_partners` | Legal/professional firms |
| `rating_agency` | `rating_agencies` | Rating agencies |
| `trustees` | `trustees` | Trustee services |

---

## API Usage

### 1. Creating an Organization with Stakeholder Linkage

**Endpoint:** `POST /api/v1/organizations/invite/root-user`

**Request Body:**
```json
{
  "organization_name": "RMB Nigeria Asset Management Limited",
  "organization_type": "ASSET_MANAGER",
  "admin_first_name": "John",
  "admin_last_name": "Doe",
  "admin_email": "john.doe@rmb.com",
  "stakeholder_id": 1,
  "stakeholder_type": "asset_manager"
}
```

**Success Response (200):**
```json
{
  "message": "You have successfully sent an invitation to john.doe@rmb.com to join Trovo Manager as ASSET_MANAGER",
  "invite_id": "uuid-here",
  "organization_id": "org-uuid-here"
}
```

---

### 2. Creating an Organization WITHOUT Stakeholder Linkage

**Request Body:**
```json
{
  "organization_name": "New Organization",
  "organization_type": "ASSET_MANAGER",
  "admin_first_name": "Jane",
  "admin_last_name": "Smith",
  "admin_email": "jane.smith@neworg.com"
}
```
Simply omit both `stakeholder_id` and `stakeholder_type`.

---

### 3. Listing Organizations (with stakeholder info)

**Endpoint:** `GET /api/v1/organizations`

**Response:**
```json
{
  "organizations": [
    {
      "id": "org-uuid",
      "name": "RMB Nigeria Asset Management Limited",
      "email": "john.doe@rmb.com",
      "type": "ASSET_MANAGER",
      "status": "ACTIVE",
      "stakeholder_id": 1,
      "stakeholder_type": "asset_manager",
      "created_at": "2025-10-01T10:00:00Z",
      "updated_at": "2025-10-01T10:00:00Z",
      "team_member_count": 5
    }
  ],
  "total": 1,
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

---

## Validation Rules

### Rule 1: Both Fields Required Together
- If you provide `stakeholder_id`, you **must** provide `stakeholder_type`
- If you provide `stakeholder_type`, you **must** provide `stakeholder_id`

**❌ Invalid:**
```json
{
  "stakeholder_id": 1
  // Missing stakeholder_type
}
```

**Error Response:**
```json
{
  "error": "stakeholder_type is required when stakeholder_id is provided"
}
```

### Rule 2: Valid Stakeholder Type
The `stakeholder_type` must be one of the valid values listed above.

**❌ Invalid:**
```json
{
  "stakeholder_id": 1,
  "stakeholder_type": "invalid_type"
}
```

**Error Response:**
```json
{
  "error": "invalid stakeholder_type 'invalid_type'. Valid values are: asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees"
}
```

### Rule 3: Stakeholder Must Exist
The stakeholder must already exist in the Trovo database (created via `/partners/save` endpoint).

**❌ Invalid:**
```json
{
  "stakeholder_id": 999,
  "stakeholder_type": "asset_manager"
}
```

**Error Response:**
```json
{
  "error": "stakeholder validation failed: asset_manager with ID 999 does not exist. Please create the stakeholder first using /partners/save endpoint"
}
```

---

## Workflow for Frontend Developers

### Recommended Flow:

1. **First, check if stakeholder exists:**
   ```
   GET /api/v1/partners/list?type=asset_manager
   ```
   This returns all existing stakeholders of that type.

2. **If stakeholder doesn't exist, create it:**
   ```
   POST /api/v1/partners/save?type=asset_manager&action=create
   Body: {
     "asset_manager_name": "New Manager",
     "asset_manager_address": "Address here",
     "asset_manager_country": "NG",
     "fee_percent": 1.0,
     "fee_fixed": 100
   }
   ```

3. **Then create organization with stakeholder linkage:**
   ```
   POST /api/v1/organizations/invite/root-user
   Body: {
     "organization_name": "New Manager",
     "organization_type": "ASSET_MANAGER",
     "stakeholder_id": 1,  // ID from step 1 or 2
     "stakeholder_type": "asset_manager",
     ...
   }
   ```

---

## Getting Stakeholder List for Dropdown

To populate a dropdown for stakeholder selection:

**Get all stakeholders of a specific type:**
```
GET /api/v1/partners/list?type=asset_manager
```

**Response:**
```json
[
  {
    "id": 1,
    "asset_manager_name": "RMB Nigeria Asset Management Limited",
    "asset_manager_address": "2nd Floor, Wings Office Complex...",
    "asset_manager_country": "NG",
    "fee_percent": 1.0,
    "fee_fixed": 2.0
  }
]
```

**Get all stakeholders (all types):**
```
GET /api/v1/partners/list
```

**Response:**
```json
{
  "asset_manager": [...],
  "asset_issuing_house": [...],
  "approved_asset_custodian": [...],
  "legal_and_professionals": [...],
  "rating_agency": [...],
  "trustees": [...]
}
```

---

## Example Use Cases

### Use Case 1: Creating Asset Manager Organization
```json
POST /api/v1/organizations/invite/root-user
{
  "organization_name": "RMB Nigeria Asset Management Limited",
  "organization_type": "ASSET_MANAGER",
  "admin_first_name": "Michael",
  "admin_last_name": "Johnson",
  "admin_email": "mjohnson@rmb.com",
  "stakeholder_id": 1,
  "stakeholder_type": "asset_manager"
}
```

### Use Case 2: Creating Custodian Organization
```json
POST /api/v1/organizations/invite/root-user
{
  "organization_name": "Afrinvest Custody Services",
  "organization_type": "ASSET_CUSTODIAN",
  "admin_first_name": "Sarah",
  "admin_last_name": "Williams",
  "admin_email": "swilliams@afrinvest.com",
  "stakeholder_id": 1,
  "stakeholder_type": "approved_asset_custodian"
}
```

### Use Case 3: Creating Organization Without Stakeholder Link
```json
POST /api/v1/organizations/invite/root-user
{
  "organization_name": "New Startup Organization",
  "organization_type": "ASSET_MANAGER",
  "admin_first_name": "Alice",
  "admin_last_name": "Brown",
  "admin_email": "alice@startup.com"
  // No stakeholder_id or stakeholder_type provided
}
```

---

## Important Notes for Frontend Developers

1. **Optional Fields**: Both `stakeholder_id` and `stakeholder_type` are optional. You can create organizations without linking them to stakeholders.

2. **Atomic Requirement**: If you're implementing stakeholder linkage in the UI, ensure both fields are provided together or both are omitted.

3. **Type Matching**: The stakeholder types are **lowercase with underscores** (e.g., `asset_manager`), not the organization types which are **UPPERCASE with underscores** (e.g., `ASSET_MANAGER`).

4. **Validation Feedback**: The API provides clear error messages if validation fails, so display these to users.

5. **Pre-existence Check**: Always verify the stakeholder exists before linking (use `/partners/list` endpoint).

---

## Testing Examples

### Test 1: Valid Stakeholder Linkage
```bash
curl -X POST "http://localhost:8080/api/v1/organizations/invite/root-user" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Test Org",
    "organization_type": "ASSET_MANAGER",
    "admin_first_name": "Test",
    "admin_last_name": "User",
    "admin_email": "test@example.com",
    "stakeholder_id": 1,
    "stakeholder_type": "asset_manager"
  }'
```

### Test 2: Invalid Stakeholder Type
```bash
curl -X POST "http://localhost:8080/api/v1/organizations/invite/root-user" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Test Org",
    "organization_type": "ASSET_MANAGER",
    "admin_first_name": "Test",
    "admin_last_name": "User",
    "admin_email": "test@example.com",
    "stakeholder_id": 1,
    "stakeholder_type": "wrong_type"
  }'
```
**Expected Error:** 400 Bad Request with message about invalid stakeholder_type

### Test 3: Missing Stakeholder ID
```bash
curl -X POST "http://localhost:8080/api/v1/organizations/invite/root-user" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Test Org",
    "organization_type": "ASSET_MANAGER",
    "admin_first_name": "Test",
    "admin_last_name": "User",
    "admin_email": "test@example.com",
    "stakeholder_type": "asset_manager"
  }'
```
**Expected Error:** 400 Bad Request with message about missing stakeholder_id

---

## Implementation Summary

### What Changed:
1. Organizations table now has `stakeholder_id` and `stakeholder_type` columns
2. Organization creation endpoint accepts optional stakeholder linkage
3. Organization list endpoint returns stakeholder linkage data
4. Full validation ensures data integrity across databases

### What Stayed the Same:
1. Existing `/partners/*` endpoints unchanged
2. Stakeholders still managed in Trovo database
3. Organizations still managed in Tokenization database
4. No breaking changes to existing functionality

---

## Questions?
Contact the backend team for any questions or clarifications.

