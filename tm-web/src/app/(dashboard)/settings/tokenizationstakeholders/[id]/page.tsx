"use client";

import Image from "next/image";
import { useParams, useSearchParams } from "next/navigation";
import React, { useState } from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaArrowLeft } from "react-icons/fa6";
import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
import deleteIcon from "@/assets/images/deletIcon.svg";
import DeleteStakeholder from "../components/DeletStakeholdersModal";
import {
  StakeholderData,
  useGetSinglePartnerQuery,
} from "@/redux/api/tokenizationstakeholders";

const StakeHoldersDetailsPage = () => {
  const params = useParams();
  const searchParams = useSearchParams();

  const id = Number(params.id);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const partnerType = searchParams.get("type") as
    | "asset_manager"
    | "asset_issuing_house"
    | "approved_asset_custodian"
    | "legal_and_professionals"
    | "rating_agency"
    | "trustees"
    | "legal_adviser"
    | "financial_adviser";
  const { data, isLoading } = useGetSinglePartnerQuery({
    id,
    type: partnerType,
  });

  const getStakeholderName = (data: StakeholderData) => {
    if ("asset_manager_name" in data) return data.asset_manager_name;
    if ("asset_issuing_house_name" in data)
      return data.asset_issuing_house_name;
    if ("asset_custodian_name" in data) return data.asset_custodian_name;
    if ("partner_name" in data) return data.partner_name;
    if ("agency_name" in data) return data.agency_name;

    if ("trustee_name" in data) return data.trustee_name;
    if ("adviser_name" in data) return data.adviser_name;
    return "N/A";
  };

  const getStakeholderInitial = (data: StakeholderData) => {
    const name = getStakeholderName(data);
    return name?.charAt(0).toUpperCase() || "P";
  };

  const formatStakeholderType = (type: string): string => {
    const typeMap: Record<string, string> = {
      asset_manager: "Asset Manager",
      asset_issuing_house: "Issuing House",
      approved_asset_custodian: "Asset Custodian",
      legal_and_professionals: "Legal & Professional",
      rating_agency: "Rating Agency",
      trustees: "Trustee",
      legal_adviser: "Legal Adviser",
      financial_adviser: "Financial Adviser",
    };

    return typeMap[type] || "Unknown Type";
  };

  const getStakeholderAddress = (data: StakeholderData) => {
    if ("asset_manager_address" in data) return data.asset_manager_address;
    if ("asset_issuing_house_address" in data)
      return data.asset_issuing_house_address;
    if ("asset_custodian_address" in data) return data.asset_custodian_address;
    if ("partner_address" in data) return data.partner_address;
    if ("agency_address" in data) return data.agency_address;
    if ("trustee_address" in data) return data.trustee_address;
    if ("adviser_address" in data) return data.adviser_address;

    return "N/A";
  };

  return (
    <DetailsWrapper>
      <DetailsHeader>
        <BackButtonLink href="/settings/tokenizationstakeholders">
          <FaArrowLeft />
        </BackButtonLink>

        <DetailsActions>
          <EditStakeholderLink
            href={`/settings/tokenizationstakeholders/edit/${id}?type=${partnerType}`}
          >
            <Image src={editIcon} alt="Edit Token" width={16} height={16} />
          </EditStakeholderLink>

          <DeleteStakeholderIcon
            src={deleteIcon}
            alt="Delete Token"
            width={16}
            height={16}
            onClick={() => setIsModalOpen(!isModalOpen)}
          />
        </DetailsActions>
      </DetailsHeader>

      {isLoading ? (
        <p> Loading...</p>
      ) : data ? (
        <DetailsContent>
          <DetailsInfo>
            <AvatarInitial>
              {data ? getStakeholderInitial(data) : "P"}
            </AvatarInitial>

            <StakeholderName>
              <StakeholderName>
                {data ? getStakeholderName(data) : "N/A"}
              </StakeholderName>
            </StakeholderName>
          </DetailsInfo>

          <DetailsInfo>
            <InfoLabel>Stakeholder Type</InfoLabel>
            <TypeBadge>{formatStakeholderType(partnerType)}</TypeBadge>
          </DetailsInfo>

          <DetailsInfo>
            <InfoLabel>Stakeholder&apos;s Address</InfoLabel>
            <FeeText>{data ? getStakeholderAddress(data) : "N/A"}</FeeText>
          </DetailsInfo>

          <DetailsInfo>
            <InfoLabel>Fees</InfoLabel>
            <FeeText>
              Percentage:{" "}
              <HighlightText>
                {" "}
                {data?.fee_percent != null ? `${data.fee_percent}%` : "N/A"}
              </HighlightText>
            </FeeText>
            <FeeText>
              Fixed Fee:{" "}
              <HighlightText>
                {data?.fee_fixed != null
                  ? `${Number(data.fee_fixed).toLocaleString()} CNGN`
                  : "N/A"}
              </HighlightText>
            </FeeText>
          </DetailsInfo>
        </DetailsContent>
      ) : (
        <p>Not found</p>
      )}

      {isModalOpen && (
        <DeleteStakeholder
          openModal={isModalOpen}
          setOpenModal={setIsModalOpen}
          stakeholderName={data ? getStakeholderName(data) : "N/A"}
          stakeholderType={partnerType}
          stakeholderId={id}
        />
      )}
    </DetailsWrapper>
  );
};

export default StakeHoldersDetailsPage;

const DetailsWrapper = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const DetailsHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0;
`;

const BackButtonLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;

const DetailsActions = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const EditStakeholderLink = styled(Link)`
  display: flex;
  align-items: center;
`;

const DeleteStakeholderIcon = styled(Image)`
  cursor: pointer;
`;

const DetailsContent = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding-bottom: 30px;
`;

const DetailsInfo = styled.div`
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 8px;
`;

const AvatarInitial = styled.div`
  background-color: #007cdf;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
`;

const StakeholderName = styled.h1`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
`;

const InfoLabel = styled.h2`
  font-weight: 600;
  font-size: 16px;
  color: #00225a;
  text-align: center;
`;

const TypeBadge = styled.div`
  font-weight: 400;
  font-size: 12px;
  text-align: center;
  color: #007cdf;
  background-color: #007cdf1a;
  border-radius: 8px;
  padding: 8px;
  // width: 50%;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 4px;
`;

const FeeText = styled.p`
  color: #00225a;
  font-weight: 400;
  font-size: 14px;
  text-align: center;
  padding-bottom: 4px;
`;

const HighlightText = styled.span`
  color: #00225a;
  font-weight: 600;
  font-size: 14px;
`;
