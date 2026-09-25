"use client";

import TrusteeAssetDetails from "../../../trustee/tokenizedasset/_components/TokenizedAssetDetails";
import type { IStakeholderAssetDetailResponse } from "@/redux/api/sharedstakeholders";

export default function CustodianTokenizedAssetDetails(props: {
  assetDetail?: IStakeholderAssetDetailResponse;
  isLoading?: boolean;
}) {
  return <TrusteeAssetDetails {...props} showActions={false} />;
}
