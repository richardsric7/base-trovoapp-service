export type TokenizationData = {
    assetCustodians:               AssetCustodian[];
    assetIssuingHouses:            AssetIssuingHouse[];
    assetManagers:                 AssetManager[];
    assetProceedCycle:             AssetProceedCycle[];
    assetProtectionOptions:        AssetProceedCycle[];
    assetSectors:                  AssetSector[];
    assetSubSectors:               AssetSubSector[];
    assetTypes:                    AssetType[];
    countries:                     Country[];
    countryConfigs:                CountryConfig[];
    feePaymentMethods:             FeePaymentMethod[];
    publicListingAllowedCountries: AssetProceedCycle[];
    tokenizationCurrencies:        TokenizationCurrency[];
    tokenizationDocumentTypes:     TokenizationDocumentType[];
    tokenizationFees:              TokenizationFee[];
    tokenizationStatuses:          TokenizationStatus[];
}

export type AssetCustodian = {
    id:                    number;
    assetCustodianName:    string;
    assetCustodianAddress: string;
    assetCustodianCountry: string;
    requirementDocument:   string;
    FeePercent:            number;
    FeeFixed:              number;
}

export type AssetIssuingHouse = {
    id:                       number;
    assetIssuingHouseName:    string;
    assetIssuingHouseAddress: string;
    assetIssuingHouseCountry: string;
    FeePercent:               number;
    FeeFixed:                 number;
}

export type AssetManager = {
    id:                  number;
    assetManagerName:    string;
    assetManagerAddress: string;
    assetManagerCountry: string;
    FeePercent:          number;
    FeeFixed:            number;
}

export type AssetProceedCycle = {
    id: string;
}

export type AssetSector = {
    sector:              string;
    requirementDocument: string;
}

export type AssetSubSector = {
    subSector:     string;
    assetSectorId: string;
}

export type AssetType = {
    id:               number;
    assetSubSectorId: string;
    assetType:        string;
}

export type Country = {
    countryCode:       string;
    regionName:        string;
    countryName:       string;
    quoteCurrencyCode: string;
    fiatLabel:         string;
    fiatGlyph:         string;
}

export type CountryConfig = {
    countryCode:                              string;
    SECTokenizationFeePercent:                number;
    SECTokenizationFeeFixed:                  number;
    SECTradeFeePercent:                       number;
    SECTradeFeeFixed:                         number;
    regulatorName:                            string;
    quoteCurrencyCode:                        string;
    fiatLabel:                                string;
    fiatGlyph:                                string;
    minTokenizationFee:                       number;
    minTROVBalanceForTokenizationApplication: number;
    tokenizationApplicationFee:               number;
    tokenizationApplicationFeeAsset:          string;
    vatPercent:                               number;
    fiatActivationAmount:                     number;
    trovTokenActivationPercent:               number;
}

export type FeePaymentMethod = {
    id:               string;
    feeDescription:   string;
    extraDescription: string;
}

export type TokenizationCurrency = {
    assetCode:   string;
    contractAddress: string;
    label:       string;
}

export type TokenizationDocumentType = {
    documentType:            string;
    documentTypeDescription: string;
    documentCategory:        DocumentCategory;
    isPublic:                number;
}

export enum DocumentCategory {
    AssetImages = "Asset Images",
    AssetVerificationDocuments = "Asset Verification Documents",
    DocumentFile = "Document File",
}

export type TokenizationFee = {
    id:                 number;
    feeFiatPercentage:  number;
    feeFiatCap:         number;
    feeAssetPercentage: number;
    feeDescription:     string;
    countryCode:        string;
}

export type TokenizationStatus = {
    id:          number;
    description: string;
}
