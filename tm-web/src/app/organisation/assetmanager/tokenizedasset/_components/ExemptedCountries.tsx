"use client";

import React from "react";
import styled from "styled-components";
import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import ReactCountryFlag from "react-country-flag";

interface ExemptedCountriesProps {
  countries?: string[];
}

const ExemptedCountries = ({ countries }: ExemptedCountriesProps) => {
  const displayNames = new Intl.DisplayNames(["en"], { type: "region" });
  const countryData = (countries ?? [])
    .map((country) => country.trim())
    .filter(Boolean)
    .map((country) => {
      const code = country.toUpperCase();
      let name = country;
      try {
        name = displayNames.of(code) || country;
      } catch {
        // Preserve a backend-provided country name if it is not an ISO code.
      }
      return { code, name };
    });

  return (
    <>
      <HeadingContent>
        <Heading>Exempted Countries</Heading>
      </HeadingContent>

      <Wrapper>
        {countryData.length === 0 && (
          <EmptyState>No exempted countries provided.</EmptyState>
        )}
        {countryData.map((country) => (
          <CountryItem key={country.code}>
            {country.code.length === 2 && (
              <ReactCountryFlag
                countryCode={country.code}
                svg
                style={{
                  width: "20px",
                  height: "20px",
                }}
              />
            )}

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

const EmptyState = styled.p`
  grid-column: 1 / -1;
  color: #828282;
  font-size: 14px;
`;
