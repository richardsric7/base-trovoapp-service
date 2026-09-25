"use client";

import { useEffect, useState } from "react";
import styled from "styled-components";
import { DropdownSelect } from "@/components";
import CustomFilter from "@/components/CustomFilter";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { ITrusteeFundReleaseListParams } from "@/redux/api/trustees";
import { useGetStakeholderAssetsQuery } from "@/redux/api/sharedstakeholders";

interface FundFilterProps {
  filters: ITrusteeFundReleaseListParams;
  onApply: (filters: ITrusteeFundReleaseListParams) => void;
}

const statuses = [
  "submitted",
  "trustee_approved",
  "execution_pending",
  "processing",
  "completed",
  "trustee_rejected",
  "failed",
];

const formatStatus = (status: string) =>
  status
    .replaceAll("_", " ")
    .replace(/\b\w/g, (letter) => letter.toUpperCase());

const FundFilter = ({ filters, onApply }: FundFilterProps) => {
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [status, setStatus] = useState(filters.status ?? "");
  const [assetId, setAssetId] = useState(filters.asset_id ?? "");
  const { data: assetsResponse, isLoading: isLoadingAssets } =
    useGetStakeholderAssetsQuery({ page: 1, limit: 100 });
  const assets = assetsResponse?.data?.records ?? [];
  const assetOptions = assets.map((asset) =>
    asset.assetName
      ? `${asset.assetName} (${asset.assetCode})`
      : asset.assetCode || asset.id,
  );
  const selectedAsset = assets.find((asset) => asset.id === assetId);
  const selectedAssetLabel = selectedAsset
    ? selectedAsset.assetName
      ? `${selectedAsset.assetName} (${selectedAsset.assetCode})`
      : selectedAsset.assetCode || selectedAsset.id
    : "";

  useEffect(() => {
    if (!isFilterOpen) {
      setStatus(filters.status ?? "");
      setAssetId(filters.asset_id ?? "");
    }
  }, [filters, isFilterOpen]);

  const handleApply = () => {
    const nextFilters: ITrusteeFundReleaseListParams = {};
    if (status) nextFilters.status = status;
    if (assetId.trim()) nextFilters.asset_id = assetId.trim();
    onApply(nextFilters);
    setIsFilterOpen(false);
  };

  const handleReset = () => {
    setStatus("");
    setAssetId("");
    onApply({});
    setIsFilterOpen(false);
  };

  return (
    <CustomFilter
      position={{ top: "400px", right: "30px" }}
      open={isFilterOpen}
      onClose={setIsFilterOpen}
    >
      <FilterContent>
        <DropdownSelect
          options={["All", ...statuses.map(formatStatus)]}
          labelText="Status"
          placeholder="All"
          value={status ? formatStatus(status) : "All"}
          onSelect={(item) =>
            setStatus(
              item === "All"
                ? ""
                : (statuses.find((value) => formatStatus(value) === item) ??
                    ""),
            )
          }
        />

        {/* <DropdownSelect
          options={
            isLoadingAssets
              ? ["Loading assets..."]
              : assets.length
                ? ["All", ...assetOptions]
                : ["No asset"]
          }
          labelText="Asset"
          placeholder={isLoadingAssets ? "Loading assets..." : "All"}
          value={
            isLoadingAssets
              ? "Loading assets..."
              : !assets.length
                ? "No asset"
                : selectedAssetLabel || (assetId ? assetId : "All")
          }
          onSelect={(item, index) => {
            if (isLoadingAssets || !assets.length) return;
            setAssetId(item === "All" ? "" : assets[index - 1]?.id ?? "");
          }}
        /> */}

        <ButtonContainer>
          <SecondaryButton onClick={handleReset}>Reset</SecondaryButton>
          <PrimaryButton onClick={handleApply}>Apply</PrimaryButton>
        </ButtonContainer>
      </FilterContent>
    </CustomFilter>
  );
};

export default FundFilter;

const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
`;

const ButtonContainer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 20px;
`;
