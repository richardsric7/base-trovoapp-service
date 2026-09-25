"use client";
import React, { useEffect, useState } from "react";
import styled from "styled-components";
import AssestCard from "./AssetCard";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
  useUpdateTokenizationInfoMutation,
} from "@/redux/api/assettokenization";
import useFormatDate from "@/hooks/useFormatDate";
import { FaRegEdit } from "react-icons/fa";

import SecondaryButton from "@/components/SecondaryButton";
import { showErrorToast, showSuccessToast } from "@/components";
import EditAssetTokenSalesModal from "./EditAssetTokenSalesModal";
import { getCleanedUpdatePayload } from "./tokenizationHelpers";

interface AssetTokenProps {
  asset?: TokenizationRecord;
}

interface AssetData {
  id: number;
  title: string;
  field?: keyof UpdateTokenizationPayload;
  subTitle: any;
}

const AssetToken: React.FC<AssetTokenProps> = ({ asset }) => {
  const formattedDate = useFormatDate();
  const [showEditModal, setShowEditModal] = useState(false);
  const [updatedAsset, setUpdatedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >({});
  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();

  useEffect(() => {
    if (asset) {
      setUpdatedAsset(asset as Partial<UpdateTokenizationPayload>);
    }
  }, [asset]);

  const assestData: AssetData[] = [
    {
      id: 1,
      title: "Tokens to be Issued",
      field: "numberOfTokenToBeIssued",
      subTitle: asset?.numberOfTokenToBeIssued
        ? `${Number(
            asset.numberOfTokenToBeIssued.toFixed(7),
          ).toLocaleString()} ${asset.assetCode}`
        : "N/A",
    },

    {
      id: 2,
      title: "Token for Sale",
      field: "numberOfTokenToBeSold",

      subTitle: asset?.numberOfTokenToBeSold
        ? `${Number(asset.numberOfTokenToBeSold.toFixed(7)).toLocaleString()} ${
            asset.assetCode
          }`
        : "N/A",
    },

    {
      id: 3,
      title: "Token not for Sale",
      field: "totalTokenHeldByManager",
      subTitle: asset?.totalTokenHeldByManager
        ? `${asset.totalTokenHeldByManager.toLocaleString()} ${asset.assetCode}`
        : "N/A",
    },

    {
      id: 4,
      title: "Asset Quote Currency",
      field: "assetQuoteCurrency",
      subTitle: asset?.assetQuoteCurrency ?? "N/A",
    },

    {
      id: 5,
      title: "Price Per Token",
      field: "pricePerToken",
      subTitle: asset?.pricePerToken
        ? `${Number(asset.pricePerToken.toFixed(7)).toLocaleString()} NGN`
        : "N/A",
    },
    {
      id: 6,
      title: "Value of Total Tokens",
      field: "valueOfTokenizedAsset",
      subTitle: asset?.valueOfTokenizedAsset
        ? `${Number(
            asset.valueOfTokenizedAsset.toFixed(7),
          ).toLocaleString()} NGN`
        : "N/A",
    },

    {
      id: 7,
      title: "Preferred Tokenization Fee",
      subTitle:
        asset?.feeInFiat || asset?.feeInAsset
          ? `${
              asset?.feeInFiat
                ? `${Number(asset.feeInFiat.toFixed(7)).toLocaleString()} NGN`
                : ""
            }${asset?.feeInFiat && asset?.feeInAsset ? " + " : ""}${
              asset?.feeInAsset
                ? `${Number(asset.feeInAsset.toFixed(7)).toLocaleString()} ${
                    asset?.assetCode
                  }`
                : ""
            }`
          : "N/A",
    },

    {
      id: 8,
      title: "Sales Start Date",
      subTitle: formattedDate(asset?.salesStart),
    },
    {
      id: 9,
      title: "Sales End Date",
      subTitle: formattedDate(asset?.salesEnd),
    },
    {
      id: 10,
      title: "Cap Quantity",
      field: "capQuantity",

      subTitle: asset?.capQuantity
        ? `${Number(asset.capQuantity.toFixed(7)).toLocaleString()} ${
            asset.assetCode
          }`
        : "N/A",
    },
    {
      id: 11,
      title: "Cap Amount",
      subTitle: asset?.capAmountInFiat
        ? `${Number(asset.capAmountInFiat.toFixed(7)).toLocaleString()} NGN`
        : "N/A",
    },
    {
      id: 12,
      title: "Cap Duration",
      field: "capDurationInDays",

      subTitle:
        asset?.capDurationInDays != null
          ? `${asset.capDurationInDays} ${
              asset.capDurationInDays === 1 ? "day" : "days"
            }`
          : "N/A",
    },
    {
      id: 13,
      title: "Proceed Payout Cycle",
      field: "proceedCycle",

      subTitle: asset?.proceedCycle,
    },
    {
      id: 14,
      title: "Payout Currency",
      field: "proceedPayoutCurrency",

      subTitle: asset?.proceedPayoutCurrency ?? "N/A",
    },

    {
      id: 15,
      title: "Amount to be raised",
      subTitle:
        asset?.numberOfTokenToBeIssued && asset?.pricePerToken
          ? `${Number(
              (asset.numberOfTokenToBeSold * asset.pricePerToken).toFixed(7),
            ).toLocaleString()} NGN`
          : "N/A",
    },
  ];

  // Handle update submission
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
      showErrorToast(errorString);
    }
  };

  return (
    <>
      <HeadingContent>
        <Heading>Asset Token & Sales</Heading>

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
        {assestData.map((data) => (
          <AssestCard key={data.id} label={data.title} value={data.subTitle} />
        ))}
      </Wrapper>
      {showEditModal && (
        <EditAssetTokenSalesModal
          isOpen={showEditModal}
          asset={asset}
          onClose={() => setShowEditModal(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
};

export default AssetToken;
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
