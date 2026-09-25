"use client";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
  useGetTokenizationParamsQuery,
  useUpdateTokenizationInfoMutation,
} from "@/redux/api/assettokenization";
import Image from "next/image";
import React, { useState, useEffect } from "react";
import styled from "styled-components";
import AssetCard from "./AssetCard";
import { showErrorToast, showSuccessToast } from "@/components";
import { FaRegEdit } from "react-icons/fa";

import SecondaryButton from "@/components/SecondaryButton";
import EditInfoModal from "./EditInfoModal";
import { getCleanedUpdatePayload } from "./tokenizationHelpers";

interface AssetInformationProps {
  asset?: TokenizationRecord;
}

const AssetInformation: React.FC<AssetInformationProps> = ({ asset }) => {
  const [showEditModal, setShowEditModal] = useState(false);
  //update Query mutation hook
  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();

  const { data: tokenizationParams } = useGetTokenizationParamsQuery();

  const getAssetTypeName = (typeId: string | number | undefined) => {
    if (!typeId || !tokenizationParams?.assetTypes) return "N/A";
    const match = tokenizationParams.assetTypes.find(
      (item) => item.id === Number(typeId),
    );
    return match?.assetType || `Unknown Type (${typeId})`;
  };

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

      // console.log(" Success:", response);
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

  const assetDetails = [
    {
      id: 1,
      title: "Asset Code",

      text: asset?.assetCode || "N/A",
    },
    {
      id: 2,
      title: "Date Created",

      text: asset?.createdAt
        ? new Date(asset.createdAt).toLocaleDateString("en-US", {
            month: "short",
            day: "numeric",
            year: "numeric",
          })
        : "N/A",
    },

    {
      id: 3,
      title: "Sector",

      text: asset?.assetSector || "No sector provided",
    },
    {
      id: 4,
      title: "SubSector",

      text: asset?.assetSubSector || "N/A",
    },
    {
      id: 5,
      title: "Type",
      text: getAssetTypeName(asset?.assetType),
    },
    {
      id: 6,
      title: "Asset Status",

      text: asset?.assetAlreadyExists === 1 ? " Existing" : "Non-Existing",
    },
    {
      id: 7,
      title: "Offering Type",

      text: asset?.offeringType || "N/A",
    },
    {
      id: 8,
      title: "Country",

      text: asset?.assetCountryLocation || "N/A",
    },
    {
      id: 9,
      title: "Physical Asset Address",

      text: asset?.assetPhysicalAddress || "No address provided",
    },
    {
      id: 10,
      title: "Longitude",

      text: asset?.assetLongitude || "N/A",
    },
    {
      id: 11,
      title: "Latitude",

      text: asset?.assetLatitude || "N/A",
    },
  ];

  return (
    <ContentDetails>
      <HeadingContent>
        <Heading>Asset Information</Heading>

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
      <PropertyInfoSection>
        <Cont>
          <AssestInfoSection>
            {asset?.assetLogo ? (
              <Image
                src={asset.assetLogo}
                alt="Asset Logo"
                width={40}
                height={40}
                style={{ borderRadius: "50%" }}
                unoptimized={true}
                onError={(e) => {
                  e.currentTarget.onerror = null;
                  e.currentTarget.src =
                    "https://media.istockphoto.com/id/1300845620/vector/user-icon-flat-isolated-on-white-background-user-symbol-vector-illustration.jpg?s=612x612&w=0&k=20&c=yBeyba0hUkh14_jgv1OKqIH0CCSWU_4ckRkAoy2p73o=";
                }}
              />
            ) : (
              <Avatar />
            )}

            <div>
              <AssestSubTitle>{asset?.assetName}</AssestSubTitle>
            </div>
          </AssestInfoSection>
        </Cont>

        <Description>{asset?.assetDescription}</Description>

        <AssetDetails>
          {assetDetails.map((data) => (
            <AssetCard key={data.id} label={data.title} value={data.text} />
          ))}
        </AssetDetails>
        {showEditModal && (
          <EditInfoModal
            isOpen={showEditModal}
            asset={asset}
            onClose={() => setShowEditModal(false)}
            isLoading={isLoading}
            onSave={handleSave}
          />
        )}
      </PropertyInfoSection>
    </ContentDetails>
  );
};

export default AssetInformation;

const ContentDetails = styled.div`
  margin-top: 40px;
`;

const Cont = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  margin-bottom: 0;
  color: #00225a;
`;

const PropertyInfoSection = styled.div`
  flex: 1;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  padding: 24px;
`;

const AssestInfoSection = styled.div`
  display: flex;
  column-gap: 8px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;
const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const AssestSubTitle = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  margin: 4px 0;
  color: #00225a;
  text-transform: capitalize;
`;

const Description = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #828282;
  margin-top: 8px;
  width: 90%;
  line-height: 25px;
`;

const AssetDetails = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
  padding-top: 10px;
  width: 100%;
`;
