"use client";

import Link from "next/link";
import React, { useState } from "react";
import { FaArrowLeft, FaEllipsisVertical } from "react-icons/fa6";
import styled from "styled-components";
import DeleteOrganizationModal from "../components/DeleteOrgModal";
import { RiDeleteBin6Line } from "react-icons/ri";
import { BiEdit } from "react-icons/bi";
import EditOrganizationModal from "../components/EditOrgModal";
import MembersTable from "../components/MembersTable";

import Image from "next/image";

import calIcon from "@/assets/images/solar_calendar-linear.svg";
import { useGetOrganizationDetailsQuery } from "@/redux/api/organizations";
import { useParams } from "next/navigation";
import { usePrettyType } from "../components/TypeFormatter";

const OrganizationDetailsPage = () => {
  const params = useParams();
  const organizationId = Array.isArray(params.id) ? params.id[0] : params.id;

  const {
    data: resp,
    isLoading,
    error,
  } = useGetOrganizationDetailsQuery(organizationId ?? "", {
    skip: !organizationId, // prevents firing before ID is available
  });

  // const organizationId = Array.isArray(params.id) ? params.id[0] : params.id;
  // const {
  //   data: resp,
  //   isLoading,
  //   error,
  // } = useGetOrganizationDetailsQuery(organizationId);

  const organization = resp?.data?.organization;
  const stakeholder = resp?.data?.stakeholder;
  const organizationAdmin = resp?.data?.organization_admin;
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isActionsOpen, setIsActionsOpen] = useState(false);
  const [selectedOrg, setSelectedOrg] = useState<any | null>(null);
  const PrettyType = usePrettyType();
  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>Failed to load organization details.</p>;

  return (
    <DetailsWrapper>
      <BackButton href="/organizations">
        <FaArrowLeft />
      </BackButton>
      <Header>
        <OrganizationDetailSection>
          <OrganizationInfo>
            <LogoCircle>
              {organization?.name
                ? organization.name.charAt(0).toUpperCase()
                : "?"}
            </LogoCircle>
            <Details>
              <OrgName>{organization?.name}</OrgName>
              <OrgType>{PrettyType(organization?.type || "")}</OrgType>
            </Details>
          </OrganizationInfo>

          <OrganizationFees>
            <Added>
              {" "}
              <Image src={calIcon} alt="calender-icon" width={17} height={17} />
              <DateText>
                Added on{" "}
                {new Date(organization?.created_at || "").toLocaleDateString(
                  "en-GB",
                  {
                    day: "numeric",
                    month: "short",
                    year: "numeric",
                  },
                )}
              </DateText>
            </Added>

            <FeesContainer>
              <FeesContent>
                Percentage Fees:{" "}
                <StyledSpan>{stakeholder?.fee_percent} %</StyledSpan>
              </FeesContent>

              <FeesContent>
                Fixed Fees:{" "}
                <StyledSpan>{stakeholder?.fee_fixed} CNGN</StyledSpan>
              </FeesContent>
            </FeesContainer>
          </OrganizationFees>
        </OrganizationDetailSection>

        <ActionsContainer
          onClick={(e) => {
            e.stopPropagation();
            setIsActionsOpen(!isActionsOpen);
          }}
        >
          <FaEllipsisVertical />
          {isActionsOpen && (
            <ActionsDropdown>
              <DropdownItem
                onClick={(e) => {
                  e.stopPropagation();
                  setIsActionsOpen(false);
                  const mappedOrgData = {
                    id: organization?.id,
                    name: organization?.name,
                    email: organizationAdmin?.email,
                    address: stakeholder?.agency_address || "",
                    country: stakeholder?.agency_country || "",
                    stakeholder_type: organization?.stakeholder_type,
                    stakeholder_id: stakeholder?.id,
                    admin: organizationAdmin
                      ? `${organizationAdmin.first_name} ${organizationAdmin.last_name}`
                      : "",
                    adminEmail: organizationAdmin?.email,
                    fee_fixed: String(stakeholder?.fee_fixed || ""),
                    fee_percent: String(stakeholder?.fee_percent || ""),
                  };
                  setSelectedOrg(mappedOrgData);
                  setIsEditModalOpen(true);
                }}
              >
                <BiEdit /> Edit Organization
              </DropdownItem>
              <DropdownItem
                className="danger"
                onClick={(e) => {
                  e.stopPropagation();
                  setIsActionsOpen(false);
                  setIsDeleteModalOpen(true);
                }}
              >
                <RiDeleteBin6Line /> Deactivate Organization
              </DropdownItem>
            </ActionsDropdown>
          )}
        </ActionsContainer>
      </Header>

      <MembersTable organizationId={organizationId} />

      {isEditModalOpen && selectedOrg && (
        <EditOrganizationModal
          isEditModal={isEditModalOpen}
          setIsEditModal={setIsEditModalOpen}
          selectedOrg={selectedOrg}
        />
      )}
      {isDeleteModalOpen && organization && (
        <DeleteOrganizationModal
          setIsDeleteModalOpen={setIsDeleteModalOpen}
          isDeleteModalOpen={isDeleteModalOpen}
          orgId={organization.id}
          orgName={organization.name}
        />
      )}
    </DetailsWrapper>
  );
};

export default OrganizationDetailsPage;

const DetailsWrapper = styled.section`
  padding: 24px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  background: #fff;
  border-radius: 24px;
  padding: 10px;
`;

const BackButton = styled(Link)`
  color: #000;
  width: 20px;
  display: block;
  margin-bottom: 8px;
`;

const OrganizationDetailSection = styled.div`
  display: grid;
  align-items: center;

  gap: 10px;
`;

const OrganizationInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const OrganizationFees = styled.div`
  margin-left: 3rem;
`;

const LogoCircle = styled.div`
  background: #007cdf;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-weight: 600;
  font-size: 14px;
`;

const OrgName = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
`;

const Details = styled.div`
  display: flex;
  gap: 4px;
  flex-direction: column;
`;

const OrgType = styled.div`
  font-size: 12px;
  color: #007cdf;
  background: #007cdf1a;
  border-radius: 8px;
  padding: 6px 8px;
  display: inline-flex;
  text-align: center;
`;

const ActionsContainer = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const ActionsDropdown = styled.div`
  position: absolute;
  right: 50%;
  top: 25%;
  margin-top: 8px;
  width: 220px;
  background: #fff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
  z-index: 1000;
`;

const DropdownItem = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  padding: 6px 0;
  color: #00225a;
  cursor: pointer;
  &.danger {
    color: #be3800;
  }
`;

const Added = styled.div`
  display: flex;
  align-items: center;
  gap: 3px;
  margin-botton: 2px;
`;

const DateText = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 100%;
  letter-spacing: 0%;
  color: #00225a;
`;

const FeesContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const FeesContent = styled.p`
  background-color: #f2f6f9;
  font-weight: 400;
  font-size: 14px;
  line-height: 25.27px;
  letter-spacing: 1.9%;

  border-radius: 12px;
  padding: 8px 10px;
  margin-top: 8px;
  color: #00225a;
`;

const StyledSpan = styled.span`
  font-weight: 600;
  font-size: 14px;
  color: #00225a;
`;
