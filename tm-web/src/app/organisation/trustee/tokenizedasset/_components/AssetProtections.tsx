"use client";

import React, { useState } from "react";
import { FaAngleDown } from "react-icons/fa6";
import styled from "styled-components";

import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";
import { IStakeholderAssetProtection } from "@/redux/api/sharedstakeholders";

interface AssetItem {
  id: number;
  title: string;
  subTitle: string;
}

interface AssetSection {
  sectionTitle: string;
  item: AssetItem[];
}

interface AssetProtectionProps {
  protection?: IStakeholderAssetProtection;
}

const yesNo = (value?: boolean) =>
  typeof value === "boolean" ? (value ? "Yes" : "No") : "N/A";

const AssetProtection = ({ protection }: AssetProtectionProps) => {
  const [openAccordian, setOpenAccordian] = useState<string[]>([]);

  const handleOpen = (section: string): void => {
    setOpenAccordian((prev) =>
      prev.includes(section)
        ? prev.filter((sec) => sec !== section)
        : [...prev, section],
    );
  };

  const assetData: AssetSection[] = [
    {
      sectionTitle: "Insurance",
      item: [
        {
          id: 1,
          title: "Real Asset Protection",
          subTitle: protection?.methods?.length
            ? protection.methods.join(", ")
            : "N/A",
        },
        {
          id: 2,
          title: "Insurance Company Name",
          subTitle: protection?.insurance_company_name || "N/A",
        },
        {
          id: 3,
          title: "Insurance Policy Number",
          subTitle: protection?.insurance_policy_number || "N/A",
        },
        {
          id: 4,
          title: "Insurance Policy Holder",
          subTitle: protection?.insurance_policy_holder || "N/A",
        },
        {
          id: 5,
          title: "Insurance Percentage Value",
          subTitle:
            protection?.insurance_coverage_percentage != null
              ? `${protection.insurance_coverage_percentage}%`
              : "N/A",
        },
      ],
    },
    {
      sectionTitle: "Contractual Protection",
      item: [
        {
          id: 1,
          title: "Revenue Guarantees",
          subTitle: yesNo(protection?.revenue_guarantees),
        },
        {
          id: 2,
          title: "Performance Bond",
          subTitle: yesNo(protection?.performance_bond),
        },
        {
          id: 3,
          title: "Service Level Agreement",
          subTitle: yesNo(protection?.service_level_agreement),
        },
      ],
    },
    {
      sectionTitle: "Ownership Safeguards",
      item: [
        {
          id: 1,
          title: "Free From Liens and Encumbrances",
          subTitle: yesNo(protection?.free_from_liens_and_encumbrances),
        },
        {
          id: 2,
          title: "Title Transferred to Custodian",
          subTitle: yesNo(protection?.transfer_title_to_custodian),
        },
        {
          id: 3,
          title: "Other Protection",
          subTitle: protection?.other || "N/A",
        },
      ],
    },
    {
      sectionTitle: "Risk Sharing Mechanisms",
      item: [
        { id: 1, title: "Completion Guarantees", subTitle: "N/A" },
        { id: 2, title: "Public Private Partnerships", subTitle: "N/A" },
        { id: 3, title: "Hedge Instruments", subTitle: "N/A" },
      ],
    },
    {
      sectionTitle: "Governance and Oversight",
      item: [
        { id: 1, title: "Independent Monitoring List", subTitle: "N/A" },
      ],
    },
    {
      sectionTitle: "Environmental, Social and Governance (ESG) Safeguards",
      item: [
        { id: 1, title: "ESG Sustainability Certifications", subTitle: "N/A" },
        { id: 2, title: "Community Engineering Plan", subTitle: "N/A" },
      ],
    },
    {
      sectionTitle: "Security Measures",
      item: [
        { id: 1, title: "Surveillance Systems", subTitle: "N/A" },
        { id: 2, title: "On-site Security Personnel", subTitle: "N/A" },
        { id: 3, title: "Perimeter Security", subTitle: "N/A" },
        { id: 4, title: "Critical Infrastructure Protection", subTitle: "N/A" },
      ],
    },
    {
      sectionTitle: "Legal/Financial Counsel",
      item: [
        { id: 1, title: "Legal Counsel", subTitle: "N/A" },
        { id: 2, title: "Financial Counsel", subTitle: "N/A" },
      ],
    },
  ];

  return (
    <>
      <HeadingContent>
        <Heading>Asset Protection</Heading>

        {/* <SecondaryButton
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
        </SecondaryButton> */}
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
                {section.item.map((item) => (
                  <AssetCard
                    key={item.id}
                    label={item.title}
                    value={item.subTitle}
                  />
                ))}
              </AccordianBody>
            )}
          </AccordianContent>
        ))}
      </Wrapper>
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
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24px;
  color: #00225a;
`;

const AccordianContent = styled.div`
  background-color: #f2f6f9;
  border-radius: 8px;
  padding: 16px;
`;

const AccordianHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;

const Title = styled.h2`
  font-weight: 500;
  font-size: 16px;
  color: #00225a;
  margin: 0;
`;

const AccordianBody = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  padding-top: 16px;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
