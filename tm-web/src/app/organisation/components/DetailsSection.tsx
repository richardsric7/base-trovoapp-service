"use client";
import React from "react";
import styled from "styled-components";
import { FaRegEdit } from "react-icons/fa";
import SecondaryButton from "@/components/SecondaryButton";
import AssestCard from "@/app/(dashboard)/assettokenization/components/AssetCard";

interface DetailItem {
  label: string;
  value: string;
}

interface DetailsSectionProps {
  title: string;
  items: DetailItem[];
  showEdit?: boolean;
}

const DetailsSection: React.FC<DetailsSectionProps> = ({
  title,
  items,
  showEdit = true,
}) => {
  return (
    <>
      <HeadingContent>
        <Heading>{title}</Heading>

        {showEdit && (
          <SecondaryButton
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
        )}
      </HeadingContent>
      <Wrapper>
        {items.map((item, index) => (
          <AssestCard
            key={index}
            label={item.label}
            value={item.value || "N/A"}
          />
        ))}
      </Wrapper>
    </>
  );
};

export default DetailsSection;

const Wrapper = styled.section`
  padding: 20px;
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
  padding: 16px 0;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;
