export type PartnerType =
  | "asset_manager"
  | "asset_issuing_house"
  | "approved_asset_custodian"
  | "legal_and_professionals"
  | "rating_agency"
  | "trustees"
  | "legal_adviser"
  | "financial_adviser";

export interface AssetManager {
  id: number;
  asset_manager_name: string;
  asset_manager_address: string;
  asset_manager_country: string;
  fee_percent: number;
  fee_fixed: number;
}

export interface AssetIssuingHouse {
  id: number;
  asset_issuing_house_name: string;
  asset_issuing_house_address: string;
  asset_issuing_house_country: string;
  fee_percent: number;
  fee_fixed: number;
}

export interface AssetCustodian {
  id: number;
  asset_custodian_name: string;
  asset_custodian_address: string;
  asset_custodian_country: string;
  requirement_document: string;
  fee_percent: string | number;
  fee_fixed: string | number;
}

// --- Legal and Professional ---
export interface LegalAndProfessional {
  id: number;
  partner_name: string;
  partner_address: string;
  partner_country: string;
  fee_percent: number;
  fee_fixed: number;
}

export interface RatingAgency {
  id: number;
  agency_name: string;
  agency_address: string;
  agency_country: string;
  fee_percent: number;
  fee_fixed: number;
}

export interface Trustee {
  id: number;
  trustee_name: string;
  trustee_address: string;
  trustee_country: string;
  fee_percent: number;
  fee_fixed: number;
}

// Legal Adviser and Financial Adviser are distinct A5 stakeholder roles,
// separate from legal_and_professionals. They are NOT granted portal access.
export interface LegalAdviser {
  id: number;
  adviser_name: string;
  adviser_address: string;
  adviser_country: string;
  fee_percent: number;
  fee_fixed: number;
}

export interface FinancialAdviser {
  id: number;
  adviser_name: string;
  adviser_address: string;
  adviser_country: string;
  fee_percent: number;
  fee_fixed: number;
}

export interface PartnerListResponse {
  approved_asset_custodian: AssetCustodian[];
  asset_issuing_house: AssetIssuingHouse[];
  asset_manager: AssetManager[];
  legal_and_professionals: LegalAndProfessional[];
  rating_agency: RatingAgency[];
  trustees: Trustee[];
  legal_adviser: LegalAdviser[];
  financial_adviser: FinancialAdviser[];
}

export type StakeholderData =
  | AssetManager
  | AssetIssuingHouse
  | AssetCustodian
  | LegalAndProfessional
  | RatingAgency
  | Trustee
  | LegalAdviser
  | FinancialAdviser;
