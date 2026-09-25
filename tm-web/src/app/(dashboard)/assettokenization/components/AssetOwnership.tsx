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
import EditAssetOwner from "./EditAssetOwner";
import { getCleanedUpdatePayload } from "./tokenizationHelpers";

interface AssetOwnershipProps {
  asset?: TokenizationRecord;
}
const AssetOwnership: React.FC<AssetOwnershipProps> = ({ asset }) => {
  const [showEditModal, setShowEditModal] = useState(false);
  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();

  const isIndividualAndDirect =
    asset?.ownershipType === "DIRECT" && asset?.ownershipKind === "INDIVIDUAL";

  const assestData = isIndividualAndDirect
    ? [
        {
          id: 1,
          title: "Ownership",
          subTitle: asset?.ownershipType,
        },
        {
          id: 2,
          title: "Owner’s Name",
          subTitle: asset?.assetOwnerName || "N/A",
        },
      ]
    : [
        {
          id: 1,
          title: "Ownership",
          subTitle: asset?.ownershipKind,
        },
        {
          id: 2,
          title: "Type",

          subTitle: asset?.ownershipType,
        },
        {
          id: 3,
          title: "Organization Name",
          subTitle: asset?.assetOwnerName || "N/A",
        },
        {
          id: 4,
          title: "Organization Address",
          subTitle: asset?.assetOwnerAddress || "N/A",
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

      // console.log("✅ Success:", response);
      showSuccessToast("Update sent successfully!");
      setShowEditModal(false);
    } catch (error: any) {
      // Retrieve the error string from the response:
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";
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
        <Heading>Asset Owner</Heading>

        <SecondaryButton
          onClick={() => setShowEditModal(true)}
          buttonStyle={{
            width: "100px",
            marginRight: "20px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginBottom: "16px",
          }}
        >
          {" "}
          <FaRegEdit />
          Edit
        </SecondaryButton>
      </HeadingContent>
      <Wrapper>
        {assestData.map((data) => (
          <AssetCard key={data.id} label={data.title} value={data.subTitle} />
        ))}
      </Wrapper>

      {showEditModal && (
        <EditAssetOwner
          isOpen={showEditModal}
          asset={asset}
          onClose={() => setShowEditModal(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
};

export default AssetOwnership;

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
