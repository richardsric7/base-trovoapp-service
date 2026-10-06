// These mirror tm-api's /provider-settings endpoints, which write the rows
// app-backend reads for outside providers: kyc_configs (Sumsub, Doja) and
// stablerail_configs (bank deposits and withdrawals). Keys are write-only:
// tm-api only ever returns whether one is set and a short hint.

export interface SecretStatus {
  set: boolean;
  hint: string;
}

export interface KycProviderSettings {
  service_provider: string;
  token: SecretStatus;
  secret_key: SecretStatus;
}

export interface StablerailSettings {
  enabled: boolean;
  base_url: string;
  fintech_id: string;
  api_key: SecretStatus;
}

export interface ProviderSettings {
  kyc: KycProviderSettings[];
  stablerail: StablerailSettings;
}

export interface SaveKycProviderRequest {
  provider: string;
  token?: string;
  secret_key?: string;
}

export interface SaveStablerailRequest {
  enabled: boolean;
  base_url: string;
  fintech_id: string;
  api_key?: string;
}
