"use client";

import {
  TokenizationRecord,
  UpdateTokenizationPayload,
  useUpdateTokenizationInfoMutation,
} from "@/redux/api/assettokenization";
import React, { useState } from "react";
import { FaAngleDown } from "react-icons/fa6";
import styled from "styled-components";
import AssetCard from "./AssetCard";
import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import { showErrorToast, showSuccessToast } from "@/components";
import EditAssetProtection from "./EditAssetProtection";
import { getCleanedUpdatePayload } from "./tokenizationHelpers";

interface AssetProtectionProps {
  asset?: TokenizationRecord;
}

interface AssetItem {
  id: number;
  title: string;
  subTitle: any;
  field?: any;
}

interface AssetSection {
  sectionTitle: string;
  item: AssetItem[];
}

const AssetProtection: React.FC<AssetProtectionProps> = ({ asset }) => {
  const [openAccordian, setOpenAccordian] = useState<string[]>([]);

  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();
  const [showEditModal, setShowEditModal] = useState(false);
  const [updatedAsset, setUpdatedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >({});

  const handleOpen = (section: string): void => {
    setOpenAccordian((prev) =>
      prev.includes(section)
        ? prev.filter((sec) => sec !== section)
        : [...prev, section]
    );
  };

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

  const assetData: AssetSection[] = [
    {
      sectionTitle: "Insurance",

      item: [
        {
          id: 1,
          title: "Real Asset Protection",

          subTitle: asset?.protectionMethods || "N/A",
        },
        {
          id: 2,
          title: "Insurance Company Name",

          subTitle: asset?.insuranceCompanyName || "N/A",
        },
        {
          id: 3,
          title: "Insurance Policy Number",

          subTitle: asset?.insurancePolicyNumber || "N/A",
        },
        {
          id: 4,
          title: "Insurance Policy Holder",

          subTitle: asset?.insurancePolicyHolder || "N/A",
        },
        {
          id: 5,
          title: "Insurance Percentage Value",

          subTitle: asset?.percentageValueOfInsurance
            ? `${asset.percentageValueOfInsurance}%`
            : "N/A",
        },
      ],
    },

    {
      sectionTitle: "Contractual Protection",

      item: [
        {
          id: 1,
          title: "Contractual Protection RevGuarantees",

          subTitle:
            asset?.contractualProtectionRevGuarantees === 1 ? "Yes" : "No",
        },

        {
          id: 2,
          title: "Contractual Protection PerfBond",
          subTitle: asset?.contractualProtectionPerfBond === 1 ? "Yes" : "No",
        },
        {
          id: 3,
          title: "Contractual Protection SLA",

          subTitle: asset?.contractualProtectionSLA === 1 ? "Yes" : "No",
        },
      ],
    },

    {
      sectionTitle: "Risk Sharing Mechanisms",

      item: [
        {
          id: 1,
          title: "Risk Sharing Mechanism Completion Guarantees",
          subTitle:
            asset?.riskSharingMechanismCompletionGuarantees === 1
              ? "Yes"
              : "No",
        },
        {
          id: 2,
          title: "Risk Sharing Mechanism PPPs",
          subTitle: asset?.riskSharingMechanismPPPs === 1 ? "Yes" : "No",
        },

        {
          id: 3,
          title: "Risk Sharing Mechanism Hedge Instruments",
          subTitle:
            asset?.riskSharingMechanismHedgeInstruments === 1 ? "Yes" : "No",
        },
      ],
    },

    {
      sectionTitle: "Governance and Oversight",

      item: [
        {
          id: 1,
          title: "Independent Monitoring List",

          subTitle: asset?.independentMonitoringList || "N/A",
        },
      ],
    },

    {
      sectionTitle: "Environmental, Social and Governance (ESG) Safeguards",

      item: [
        {
          id: 1,
          title: "ESG Sustainability Certifications",

          subTitle: asset?.eSGSafeguardsSusCerts === 1 ? "Yes" : "No",
        },
        {
          id: 2,
          title: "ESG Community Engineering Plan",

          subTitle: asset?.eSGSafeguardsCommEngPlans === 1 ? "Yes" : "No",
        },
      ],
    },

    {
      sectionTitle: "Security Measures",

      item: [
        {
          id: 1,
          title: "Security Measure Surveilance Systems",

          subTitle:
            asset?.securityMeasuresSurveilanceSystems === 1 ? "Yes" : "No",
        },
        {
          id: 2,
          title: "Security Measure On-site Security Personnel",

          subTitle:
            asset?.securityMeasuresOnSiteSecurityPersonnel === 1 ? "Yes" : "No",
        },
        {
          id: 3,
          title: "Perimeter Security",

          subTitle:
            asset?.securityMeasuresPerimeterSecurity === 1 ? "Yes" : "No",
        },
        {
          id: 4,
          title: "Critical Infrastructure Protections",

          subTitle:
            asset?.securityMeasuresCriticalInfraProtections === 1
              ? "Yes"
              : "No",
        },
      ],
    },

    {
      sectionTitle: "Legal/Financial Counsel",

      item: [
        {
          id: 1,
          title: "Legal Counsel",

          subTitle: asset?.legalAdvisor || "N/A",
        },
        {
          id: 2,
          title: "Financial Counsel",

          subTitle: asset?.financialAdvisor || "N/A",
        },
      ],
    },
  ];

  return (
    <>
      <HeadingContent>
        <Heading>Asset Protection</Heading>

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
        {assetData.map((section) => (
          <AccordianContent key={section.sectionTitle}>
            <AccordianHeader onClick={() => handleOpen(section.sectionTitle)}>
              <Title>{section.sectionTitle}</Title>
              <FaAngleDown />
            </AccordianHeader>
            {openAccordian.includes(section.sectionTitle) && (
              <AccordianBody>
                {section.item.map((items) => (
                  <AssetCard
                    key={items.id}
                    label={items.title}
                    value={items.subTitle}
                  />
                ))}
              </AccordianBody>
            )}
          </AccordianContent>
        ))}
      </Wrapper>

      {showEditModal && (
        <EditAssetProtection
          isOpen={showEditModal}
          asset={asset}
          onClose={() => setShowEditModal(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
};

export default AssetProtection;
const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 20px;
  // margin-bottom: 10px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
`;
const AccordianContent = styled.div`
  background-color: #f2f6f9;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
`;

const AccordianHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.h2`
  font-weight: 500;
  font-size: 16px;
  color: #00225a;
`;

const AccordianBody = styled.div`
  background-color: #f2f6f9;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  border-radius: 8px;
  padding: 20px 0px 16px 0px;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
