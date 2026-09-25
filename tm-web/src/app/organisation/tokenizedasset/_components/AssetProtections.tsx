"use client";

import React, { useState } from "react";
import { FaAngleDown } from "react-icons/fa6";
import styled from "styled-components";

import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";

interface AssetItem {
  id: number;
  title: string;
  subTitle: string;
}

interface AssetSection {
  sectionTitle: string;
  item: AssetItem[];
}

const AssetProtection = () => {
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
          subTitle: "Full Insurance Coverage",
        },
        {
          id: 2,
          title: "Insurance Company Name",
          subTitle: "AXA Insurance",
        },
        {
          id: 3,
          title: "Insurance Policy Number",
          subTitle: "AXA-098233-INS",
        },
        {
          id: 4,
          title: "Insurance Policy Holder",
          subTitle: "Atlantis Asset Ltd",
        },
        {
          id: 5,
          title: "Insurance Percentage Value",
          subTitle: "80%",
        },
      ],
    },

    {
      sectionTitle: "Contractual Protection",
      item: [
        {
          id: 1,
          title: "Revenue Guarantees",
          subTitle: "Yes",
        },
        {
          id: 2,
          title: "Performance Bond",
          subTitle: "Yes",
        },
        {
          id: 3,
          title: "Service Level Agreement",
          subTitle: "No",
        },
      ],
    },

    {
      sectionTitle: "Risk Sharing Mechanisms",
      item: [
        {
          id: 1,
          title: "Completion Guarantees",
          subTitle: "Yes",
        },
        {
          id: 2,
          title: "Public Private Partnerships",
          subTitle: "Yes",
        },
        {
          id: 3,
          title: "Hedge Instruments",
          subTitle: "No",
        },
      ],
    },

    {
      sectionTitle: "Governance and Oversight",
      item: [
        {
          id: 1,
          title: "Independent Monitoring List",
          subTitle: "Deloitte Risk Monitoring",
        },
      ],
    },

    {
      sectionTitle: "Environmental, Social and Governance (ESG) Safeguards",
      item: [
        {
          id: 1,
          title: "ESG Sustainability Certifications",
          subTitle: "Yes",
        },
        {
          id: 2,
          title: "Community Engineering Plan",
          subTitle: "Yes",
        },
      ],
    },

    {
      sectionTitle: "Security Measures",
      item: [
        {
          id: 1,
          title: "Surveillance Systems",
          subTitle: "Yes",
        },
        {
          id: 2,
          title: "On-site Security Personnel",
          subTitle: "Yes",
        },
        {
          id: 3,
          title: "Perimeter Security",
          subTitle: "Yes",
        },
        {
          id: 4,
          title: "Critical Infrastructure Protection",
          subTitle: "No",
        },
      ],
    },

    {
      sectionTitle: "Legal/Financial Counsel",
      item: [
        {
          id: 1,
          title: "Legal Counsel",
          subTitle: "LexTrust Legal Advisory",
        },
        {
          id: 2,
          title: "Financial Counsel",
          subTitle: "Goldman Financial Advisory",
        },
      ],
    },
  ];

  return (
    <>
      <HeadingContent>
        <Heading>Asset Protection</Heading>

        <SecondaryButton
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
