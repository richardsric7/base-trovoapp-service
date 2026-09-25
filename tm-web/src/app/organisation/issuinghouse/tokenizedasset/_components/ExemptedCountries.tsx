"use client";

import React from "react";
import styled from "styled-components";
import ReactCountryFlag from "react-country-flag";

const dummyCountries = [
  { code: "US", name: "United States" },
  { code: "CN", name: "China" },
  { code: "RU", name: "Russia" },
  { code: "IR", name: "Iran" },
  { code: "KP", name: "North Korea" },
];

interface ExemptedCountriesProps {
  countries?: string[];
}

const ExemptedCountries = ({ countries }: ExemptedCountriesProps) => {
  const displayNames = new Intl.DisplayNames(["en"], { type: "region" });
  const countryData = Array.isArray(countries)
    ? countries.map((code) => ({
        code,
        name: displayNames.of(code.toUpperCase()) || code,
      }))
    : dummyCountries;

  return (
    <>
      <HeadingContent>
        <Heading>Exempted Countries</Heading>
      </HeadingContent>

      <Wrapper>
        {countryData.length === 0 && (
          <EmptyText>No exempted countries</EmptyText>
        )}
        {countryData.map((country) => (
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

const EmptyText = styled.p`
  grid-column: 1 / -1;
  margin: 0;
  color: #828282;
  font-size: 14px;
`;
