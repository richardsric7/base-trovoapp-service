"use client";

import React, { useState } from "react";
import styled from "styled-components";

import {
  TokenizationDetailResponse,
  UpdateTokenizationPayload,
  useUpdateTokenizationInfoMutation,
} from "@/redux/api/assettokenization";
import { boolean } from "yup";
import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import { showErrorToast, showSuccessToast } from "@/components";
import EditCountires from "./EditCountires";
import { getCleanedUpdatePayload } from "./tokenizationHelpers";

interface ExemptedCountriesProps {
  asset?: TokenizationDetailResponse;
}
const ExemptedCountries: React.FC<ExemptedCountriesProps> = ({ asset }) => {
  const [showEditModal, setShowEditModal] = useState(false);
  const [updatedAsset, setUpdatedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >({});
  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();
  const countries = asset?.exemptedCountries
    ? asset.exemptedCountries
        .trim()
        .split(",")
        .map((item) => item.trim())
        .filter(boolean)
    : [];

  const handleSave = async (
    updatedAsset: Partial<UpdateTokenizationPayload>,
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
        <Heading>Exempted Countries</Heading>

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
        {countries.length > 0 ? (
          countries.map((country, index) => (
            <Title key={index}>{country}</Title>
          ))
        ) : (
          <p>N/A</p>
        )}
      </Wrapper>

      {showEditModal && (
        <EditCountires
          isOpen={showEditModal}
          asset={asset}
          onClose={() => setShowEditModal(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
};

export default ExemptedCountries;

const Wrapper = styled.section`
  padding: 20px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  padding: 16px 0;
`;

const Title = styled.p`
  display: flex;
  align-items: center;

  gap: 4px;
  font-size: 14px;
  font-weight: 400;
  line-height: 16px;
  color: #00225a;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
