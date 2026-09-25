"use client";

import Image from "next/image";
import React, { useState } from "react";
import styled from "styled-components";
import { FaRegEdit } from "react-icons/fa";
import SecondaryButton from "@/components/SecondaryButton";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";

const dummyAsset = {
  id: 1,
  assetCode: "IKY001",
  createdAt: "2026-02-20",
  assetSector: "Real Estate",
  assetSubSector: "Residential",
  assetType: "Single-family homes",
  assetAlreadyExists: 1,
  offeringType: "Public",
  assetCountryLocation: "Nigeria",
  assetPhysicalAddress: "12 Clover Road, Ikoyi",
  assetLongitude: "6.8454548",
  assetLatitude: "5.128468",
  assetName: "12 Clover Road Ikoyi - 5 Bedroom Duplex",
  assetDescription:
    "Luxury residential duplex located in Ikoyi, Lagos with premium amenities and high investment potential.",
  assetLogo: null,
};

const AssetInformation: React.FC = () => {
  const [showEditModal, setShowEditModal] = useState(false);

  const assetDetails = [
    { id: 1, title: "Asset Code", text: dummyAsset.assetCode },
    {
      id: 2,
      title: "Date Created",
      text: new Date(dummyAsset.createdAt).toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
      }),
    },
    { id: 3, title: "Sector", text: dummyAsset.assetSector },
    { id: 4, title: "SubSector", text: dummyAsset.assetSubSector },
    { id: 5, title: "Type", text: dummyAsset.assetType },
    {
      id: 6,
      title: "Asset Status",
      text: dummyAsset.assetAlreadyExists ? "Existing" : "Non-Existing",
    },
    { id: 7, title: "Offering Type", text: dummyAsset.offeringType },
    { id: 8, title: "Country", text: dummyAsset.assetCountryLocation },
    {
      id: 9,
      title: "Physical Asset Address",
      text: dummyAsset.assetPhysicalAddress,
    },
    { id: 10, title: "Longitude", text: dummyAsset.assetLongitude },
    { id: 11, title: "Latitude", text: dummyAsset.assetLatitude },
  ];

  return (
    <ContentDetails>
      <HeadingContent>
        <Heading>Asset Information</Heading>

        <SecondaryButton
          onClick={() => setShowEditModal(true)}
          buttonStyle={{
            width: "100px",
            display: "flex",
            alignItems: "center",
            gap: "4px",
          }}
        >
          <FaRegEdit />
          Edit
        </SecondaryButton>
      </HeadingContent>

      <PropertyInfoSection>
        <Cont>
          <AssestInfoSection>
            {dummyAsset.assetLogo ? (
              <Image
                src={dummyAsset.assetLogo}
                alt="Asset Logo"
                width={40}
                height={40}
                style={{ borderRadius: "50%" }}
              />
            ) : (
              <Avatar />
            )}

            <div>
              <AssestSubTitle>{dummyAsset.assetName}</AssestSubTitle>
            </div>
          </AssestInfoSection>
        </Cont>

        <Description>{dummyAsset.assetDescription}</Description>

        <AssetDetails>
          {assetDetails.map((data) => (
            <AssetCard key={data.id} label={data.title} value={data.text} />
          ))}
        </AssetDetails>
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
  color: #00225a;
`;

const Description = styled.p`
  font-size: 14px;
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
`;
