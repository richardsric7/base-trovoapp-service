"use client";

import React, { useState } from "react";
import styled from "styled-components";
import AssetCard from "./AssetCard";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
  useUpdateTokenizationInfoMutation,
} from "@/redux/api/assettokenization";

import { FaRegEdit } from "react-icons/fa";
import SecondaryButton from "@/components/SecondaryButton";
import { showErrorToast, showSuccessToast } from "@/components";
import EditAssetValue from "./EditAssetValue";
import { getCleanedUpdatePayload } from "./tokenizationHelpers";

interface AssetValueProps {
  asset?: TokenizationRecord;
}

const AssetValue: React.FC<AssetValueProps> = ({ asset }) => {
  const [showEditModal, setShowEditModal] = useState(false);
  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();

  const assetData = [
    {
      id: 1,
      title: "Current Asset Vaue",

      subTitle: asset?.assetCurrentValue
        ? `${asset.assetCurrentValue.toLocaleString()}NGN`
        : "N/A",
    },
    {
      id: 2,
      title: "Cost Outside Valuation ",

      subTitle: asset?.assetMscCostOutisdeOfValuation
        ? `${asset.assetMscCostOutisdeOfValuation.toLocaleString()} NGN`
        : "N/A",
    },
    {
      id: 3,
      title: "Value to be retained",

      subTitle: asset?.assetOwnerRetainedOrContributedValue
        ? `${asset.assetOwnerRetainedOrContributedValue.toLocaleString()}NGN`
        : "N/A",
    },
  ];

  const handleSave = async (
    updatedAsset: Partial<UpdateTokenizationPayload>
  ) => {
    if (!asset?.id) {
      showErrorToast("Asset data is missing");
      return;
    }

    try {
      const cleanedData = getCleanedUpdatePayload({
        ...asset,
        ...updatedAsset,
      });

      const response = await updateTokenizationInfo({
        tokenizedAssetID: asset.id,
        data: cleanedData,
      }).unwrap();

      console.log("✅ Success:", response);
      showSuccessToast("Update sent successfully!");
      setShowEditModal(false);
    } catch (error: any) {
      // console.error("Update error caught:", error);

      // Retrieve the error string from the response:
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";

      // For debugging: display an alert to ensure we see the error message.
      // alert("Error occurred: " + errorString);

      // Extract the specific message after "message:" if present.
      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }
      // Show the error toast.
      showErrorToast(finalErrorMessage);
    }
  };

  return (
    <>
      <HeadingContent>
        <Heading>Asset Value</Heading>

        <SecondaryButton
          onClick={() => {
            if (asset) {
              setShowEditModal(true);
            }
          }}
          buttonStyle={{
            width: "100px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginRight: "20px",
            marginBottom: "16px",
          }}
        >
          <FaRegEdit />
          Edit
        </SecondaryButton>
      </HeadingContent>
      <Wrapper>
        {assetData.map((data) => (
          <AssetCard key={data.id} label={data.title} value={data.subTitle} />
        ))}
      </Wrapper>

      {showEditModal && (
        <EditAssetValue
          isOpen={showEditModal}
          asset={asset}
          onClose={() => setShowEditModal(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
};

export default AssetValue;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
