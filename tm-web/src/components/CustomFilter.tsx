"use client";
import React, { useState } from "react";
import FilterIcon from "@/assets/images/filter.svg";
import styled from "styled-components";
import Image from "next/image";
import { FaX } from "react-icons/fa6";

interface PositionProps {
  top?: string;
  bottom?: string;
  right?: string;
  left?: string;
}

interface filterProps {
  children: React.ReactNode;
  position: PositionProps;
  open: boolean;
  onClose: (open: boolean) => void;
}

const CustomFilter: React.FC<filterProps> = ({
  children,
  position = { top: "280px", right: "40px" },
  open,
  onClose,
}) => {
  return (
    <>
      <FilterButton onClick={() => onClose(!open)} active={open}>
        <StyledImage src={FilterIcon} alt="filter-icon" active={open} />
        <Text active={open}>Filter</Text>
      </FilterButton>

      {open && (
        <FilterContainer position={position}>
          <FilterRow>
            <FilterText>Filter</FilterText>
            <FaX onClick={() => onClose(false)} cursor="pointer" size={12} />
          </FilterRow>
          <Divider />

          <FilterContent>{children}</FilterContent>
        </FilterContainer>
      )}
    </>
  );
};

export default CustomFilter;

const FilterButton = styled.button<{ active: boolean }>`
  border: 1px solid ${(props) => (props.active ? "#007CDF" : "#E0E0E0")};
  background-color: #ffffff;
  width: 110px;
  height: 40px;
  padding: 8px 16px;
  gap: 15px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-family: inherit;
`;

const StyledImage = styled(Image)<{ active: boolean }>`
  filter: ${(props) =>
    props.active
      ? "brightness(0) saturate(100%) invert(35%) sepia(97%) saturate(1064%) hue-rotate(183deg) brightness(99%) contrast(101%)"
      : "none"};
`;

const Text = styled.h1<{ active: boolean }>`
  color: ${(props) => (props.active ? "#007CDF" : "#00225A")};
  font-size: 14px;
  font-weight: 400;
  line-height: 20px;
  letter-spacing: 0.1px;
  margin: 0;
`;

const FilterContainer = styled.div<{ position: PositionProps }>`
  position: absolute;
  ${(props) => props.position.top && `top: ${props.position.top};`}
  ${(props) => props.position.right && `right: ${props.position.right};`}
  ${(props) => props.position.bottom && `bottom: ${props.position.bottom};`}
  ${(props) => props.position.left && `left: ${props.position.left};`}
 background-color: #ffffff;
  box-shadow: 0px 4px 100px 0px #00000026;
  width: 280px;
  padding: 30px 24px;
  border-radius: 24px;
  z-index: 3000;
  margin-bottom: 20px;
  z-index: 10;
`;

const FilterRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const FilterText = styled.h3`
  font-size: 20px;
  font-weight: 700;
  color: #00225a;
`;

const Divider = styled.div`
  background-color: #e5e5ef;
  height: 1px;
  width: 100%;
  margin: 6px 0;
`;

const FilterContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
  marign-top: 15px;
`;
