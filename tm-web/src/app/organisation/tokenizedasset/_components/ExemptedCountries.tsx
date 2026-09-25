"use client";

import React from "react";
import styled from "styled-components";
import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import ReactCountryFlag from "react-country-flag";

const dummyCountries = [
  { code: "US", name: "United States" },
  { code: "CN", name: "China" },
  { code: "RU", name: "Russia" },
  { code: "IR", name: "Iran" },
  { code: "KP", name: "North Korea" },
];

const ExemptedCountries = () => {
  return (
    <>
      <HeadingContent>
        <Heading>Exempted Countries</Heading>

        <SecondaryButton
          buttonStyle={{
            width: "100px",
            display: "flex",
            alignItems: "center",
            gap: "4px",
            marginBottom: "16px",
            marginRight: "20px",
          }}
        >
          <FaRegEdit />
          Edit
        </SecondaryButton>
      </HeadingContent>

      <Wrapper>
        {dummyCountries.map((country) => (
          <CountryItem key={country.code}>
            <ReactCountryFlag
              countryCode={country.code}
              svg
              style={{
                width: "20px",
                height: "20px",
              }}
            />

            <Title>{country.name}</Title>
          </CountryItem>
        ))}
      </Wrapper>
    </>
  );
};

export default ExemptedCountries;
const Wrapper = styled.section`
  padding: 20px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  //   padding: 16px 0;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const CountryItem = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const Title = styled.p`
  font-size: 14px;
  font-weight: 400;
  color: #00225a;
`;
