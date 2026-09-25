// Organisation/stakeholder types used for scoping a compliance template to
// the org types it applies to. Kept local to the compliance template editor —
// other screens (invite-org modal, org filters) already maintain their own
// label/value lists with differing casing conventions; this does not attempt
// to unify them.
export const ORG_TYPES = [
  { value: "approved_asset_custodian", label: "Approved Asset Custodian" },
  { value: "rating_agency", label: "Rating Agency" },
  { value: "asset_issuing_house", label: "Asset Issuing House" },
  { value: "asset_manager", label: "Asset Manager" },
  { value: "trustees", label: "Trustees" },
  { value: "legal_adviser", label: "Legal Adviser" },
  { value: "financial_adviser", label: "Financial Adviser" },
] as const;
