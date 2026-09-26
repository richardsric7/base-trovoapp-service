// These mirror tm-api's /assets/curated endpoints
// (usermetrics/db/curated_assets.go), which write directly to the same
// curated_assets table app-backend's P2P module (and every other
// wallet feature) reads from.

export interface ICuratedAsset {
  id: number;
  createdAt: string;
  updatedAt: string;
  assetCode: string;
  assetName: string;
  contractAddress: string;
  description: string;
  imageUrl: string | null;
  website: string;
  assetConditions: string;
  assetLimit: number;
  assetRedemptionInstructions: string;
  contactEmail: string;
  priority: number;
  assetClassId: number;
  organization: string;
  withdrawable: number;
  generateDepositAddress: number;
  decimalPlaces: number;
  realAssetImageUrl: string | null;
  inactive: number;
  closedGroup: string | null;
  // Gates whether this asset can be used to create a P2P offer or is
  // found in marketplace search - see app-backend's
  // CuratedAsset.P2PEnabled doc for the full contract.
  p2pEnabled: boolean;
}

export interface CuratedAssetListQueryParams {
  page?: number;
  pageSize?: number;
  assetCode?: string;
  assetClassId?: number;
  p2pEnabled?: "true" | "false";
}

export interface CuratedAssetListResponse {
  message: string;
  data: {
    data: ICuratedAsset[];
    total: number;
    page: number;
    pageSize: number;
  };
  status: string;
  timestamp: string;
}

export interface CuratedAssetResponse {
  message: string;
  data: ICuratedAsset;
  status: string;
  timestamp: string;
}

export interface CuratedAssetRequest {
  action: "create" | "update";
  id?: number;
  assetCode: string;
  assetName?: string;
  contractAddress?: string;
  description?: string;
  imageUrl?: string;
  website?: string;
  assetConditions?: string;
  assetLimit?: number;
  assetRedemptionInstructions?: string;
  contactEmail?: string;
  priority?: number;
  assetClassId?: number;
  organization?: string;
  withdrawable?: boolean;
  generateDepositAddress?: boolean;
  decimalPlaces?: number;
  realAssetImageUrl?: string;
  inactive?: boolean;
  closedGroup?: string;
  p2pEnabled?: boolean;
}

export interface IAssetClass {
  id: number;
  assetClass: string;
}

export interface AssetClassListResponse {
  message: string;
  data: IAssetClass[];
  status: string;
  timestamp: string;
}
