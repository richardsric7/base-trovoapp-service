// types/interface.ts or redux/api/country/interface.ts

export interface Country {
  id: number;
  country_code: string;
  region_name: string;
  country_name: string;
  quote_currency_code: string;
  regulator_name: string;
  fiat_label: string;
  fiat_glyph: string;
  sec_tokenization_fee_percent: number;
  sec_trade_fee_percent: number;
  sec_tokenization_fee: number;
  sec_trade_fee: number;
  min_tokenization_fee: number;
  tokenization_application_fee: number;
  tokenization_application_fee_asset: string;
  legal_and_professional_fee_percent: number;
  rating_agency_fee_percent: number;
  vat_percent: number;
  sec_tokenization_fee_fixed: number;
  sec_trade_fee_fixed: number;
  legal_and_professional_fee_fixed: number;
  rating_agency_fee_fixed: number;
}

export type CountryPayload = Country[];

export interface ICountryPayload {
  action: "create" | "update";
  id?: number;

  country_code: string;
  country_name: string;
  fiat_glyph: string;
  fiat_label: string;
  quote_currency_code: string;
  region_name: string;
  regulator_name: string;

  legal_and_professional_fee_fixed: number;
  legal_and_professional_fee_percent: number;

  rating_agency_fee_fixed: number;
  rating_agency_fee_percent: number;

  min_tokenization_fee: number;

  sec_tokenization_fee: number;
  sec_tokenization_fee_fixed: number;
  sec_tokenization_fee_percent: number;

  sec_trade_fee: number;
  sec_trade_fee_fixed: number;
  sec_trade_fee_percent: number;

  tokenization_application_fee: number;
  tokenization_application_fee_asset: string;

  vat_percent: number;
}

export interface ICountrySaveResponse {
  message: string;
  timestamp: string;
  status: string;
}
