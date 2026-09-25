const stakeholderDashboardRoutes: Record<string, string> = {
  asset_manager: "/organisation/assetmanager",
  approved_asset_custodian: "/organisation/assetcustodian",
  asset_custodian: "/organisation/assetcustodian",
  trustees: "/organisation/trustee",
  trustee: "/organisation/trustee",
  rating_agency: "/organisation/ratingagency",
  asset_issuing_house: "/organisation/issuinghouse",
  issuing_house: "/organisation/issuinghouse",
};

const stakeholderTokenizedAssetRoutes: Record<string, string> = {
  asset_manager: "/organisation/assetmanager/tokenizedasset",
};

const normalizeStakeholderType = (stakeholderType?: string | null) =>
  stakeholderType?.trim().toLowerCase() ?? "";

export const getOrganizationDashboardRoute = (
  stakeholderType?: string | null,
): string => {
  if (!stakeholderType) return "/organisation";

  return (
    stakeholderDashboardRoutes[normalizeStakeholderType(stakeholderType)] ??
    "/organisation"
  );
};

export const getOrganizationTokenizedAssetRoute = (
  stakeholderType?: string | null,
): string =>
  stakeholderTokenizedAssetRoutes[normalizeStakeholderType(stakeholderType)] ??
  "/organisation/tokenizedasset";
