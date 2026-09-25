"use client";

import React, { useState } from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaAngleDown, FaAngleUp, FaArrowLeft } from "react-icons/fa6";
import deleteIcon from "@/assets/images/deletIcon.svg";
import Image from "next/image";
import { useParams, useRouter } from "next/navigation";
import DeleteCountry from "../components/DeleteCountry";

interface CountryItem {
  title: string;
  subTitle: any;
  field?: any;
}

interface CountrySection {
  sectionTitle: string;
  item: CountryItem[];
}

import { useGetCountryByIdQuery } from "@/redux/api/country";

import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";
import AssetCard from "@/app/(dashboard)/assettokenization/components/AssetCard";
const CountiresDetailsPage = () => {
  const params = useParams();
  const id = Number(params.id);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { data, isLoading } = useGetCountryByIdQuery(id);

  const [openAccordian, setOpenAccordian] = useState<string[]>([]);

  const router = useRouter();
  const handleOpen = (section: string): void => {
    setOpenAccordian((prev) =>
      prev.includes(section)
        ? prev.filter((sec) => sec !== section)
        : [...prev, section]
    );
  };
  const countryData: CountrySection[] = [
    {
      sectionTitle: "Country Information ",
      item: [
        { title: "Country Code", subTitle: data?.country_code || "N/A" },
        { title: "Country Name", subTitle: data?.country_name || "N/A" },
        { title: "Fiat Glyph", subTitle: data?.fiat_glyph || "N/A" },
        { title: "Fiat Label", subTitle: data?.fiat_label || "N/A" },
        {
          title: "Quote Currency Code",
          subTitle: data?.quote_currency_code || "N/A",
        },
        { title: "Region Name", subTitle: data?.region_name || "N/A" },
      ],
    },
    // {
    //   sectionTitle: "Regulator",
    //   item: [
    //     { title: "Regulator Name", subTitle: data?.regulator_name || "N/A" },
    //   ],
    // },
    {
      sectionTitle: "Legal and Professional Fees",
      item: [
        {
          title: "Fee Percentage",
          subTitle: data?.legal_and_professional_fee_percent || "N/A",
        },
        {
          title: "Fixed Fee",
          subTitle: data?.legal_and_professional_fee_fixed
            ? `${data.legal_and_professional_fee_fixed} ${data.quote_currency_code}`
            : "N/A",
        },
      ],
    },
    // {
    //   sectionTitle: "Rating Agency Fees",
    //   item: [
    //     {
    //       title: "Fee Percentage",
    //       subTitle: data?.rating_agency_fee_percent
    //         ? `${data.rating_agency_fee_percent}%`
    //         : "N/A",
    //     },
    //     {
    //       title: "Fixed Fee",
    //       subTitle: data?.rating_agency_fee_fixed
    //         ? `${data.rating_agency_fee_fixed} ${data.quote_currency_code}`
    //         : "N/A",
    //     },
    //   ],
    // },
    {
      sectionTitle: "Tokenization Fees",
      item: [
        {
          title: "Tokenization Fee",
          subTitle: data?.min_tokenization_fee
            ? `${data.min_tokenization_fee} ${data.quote_currency_code}`
            : "N/A",
        },

        {
          title: "Application Fee",
          subTitle: data?.tokenization_application_fee,
        },
        {
          title: " Application Fee Asset",
          subTitle: data?.tokenization_application_fee_asset || "N/A",
        },
      ],
    },

    {
      sectionTitle: "Regulatory Fee",
      item: [
        {
          title: "SEC Trade Fee Fixed",
          subTitle: data?.sec_trade_fee_fixed
            ? `${data.sec_trade_fee_fixed} ${data.quote_currency_code}`
            : "N/A",
        },

        {
          title: "SEC Trade Fee Percent",
          subTitle: data?.sec_trade_fee_percent
            ? `${data.sec_trade_fee_percent}%`
            : "N/A",
        },
        {
          title: "SEC Tokenization Fee",
          subTitle: data?.sec_tokenization_fee
            ? `${data.sec_tokenization_fee} ${data.quote_currency_code}`
            : "N/A",
        },
      ],
    },

    {
      sectionTitle: "VAT Percent",
      item: [
        {
          title: "VAT Percent",
          subTitle: data?.vat_percent ? `${data.vat_percent}%` : "N/A",
        },
      ],
    },
  ];

  return (
    <DetailsWrapper>
      <BackButtonLink href="/countries">
        <FaArrowLeft />
      </BackButtonLink>
      <DetailsHeader>
        <Heading>Country Details</Heading>

        <DetailsActions>
          <SecondaryButton
            onClick={() => router.push(`/countries/edit/${id}`)}
            buttonStyle={{
              width: "100px",
              height: "40px",
              display: "flex",
              alignItems: "center",
              gap: "2px",
            }}
          >
            <FaRegEdit />
            Edit
          </SecondaryButton>

          <SecondaryButton
            onClick={() => setIsModalOpen(!isModalOpen)}
            buttonStyle={{
              width: "100px",

              height: "40px",
              display: "flex",
              alignItems: "center",
              gap: "2px",

              color: "#BE3800",
              border: "1px solid #BE3800",
            }}
          >
            <Image src={deleteIcon} alt="Delete icon" width={16} height={16} />
            Delete
          </SecondaryButton>
        </DetailsActions>
      </DetailsHeader>

      {isLoading ? (
        <p> Loading...</p>
      ) : data ? (
        <>
          <DetailsContent>
            <Avatar />
            <DetailsInfo>
              <StakeholderName>{data?.country_code}</StakeholderName>
              <InfoText>{data?.country_name}</InfoText>
            </DetailsInfo>
          </DetailsContent>

          <CountryContent>
            {countryData.map((section) => (
              <AccordianContent key={section.sectionTitle}>
                <AccordianHeader
                  onClick={() => handleOpen(section.sectionTitle)}
                >
                  <Title>{section.sectionTitle}</Title>
                  <FaAngleUp />
                </AccordianHeader>
                {openAccordian.includes(section.sectionTitle) && (
                  <AccordianBody>
                    {section.item.map((items) => (
                      <AssetCard label={items.title} value={items.subTitle} />
                    ))}
                  </AccordianBody>
                )}
              </AccordianContent>
            ))}
          </CountryContent>
        </>
      ) : (
        <p>Not found</p>
      )}

      {isModalOpen && (
        <DeleteCountry
          openModal={isModalOpen}
          setOpenModal={setIsModalOpen}
          countryName={data?.country_name ?? ""}
          countryId={data?.id ?? 0}
        />
      )}
    </DetailsWrapper>
  );
};

export default CountiresDetailsPage;

const DetailsWrapper = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
`;

const CountryContent = styled.div`
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  padding: 24px;
`;

const DetailsHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
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

const DetailsContent = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 20px 0;
`;

const DetailsInfo = styled.div``;

const StakeholderName = styled.h1`
  font-weight: 600;
  font-size: 24px;
  color: #00225a;
  margin-top: 8px;
`;

const InfoText = styled.p`
  color: #00225a;
  font-weight: 400;
  font-size: 14px;
  cursor: pointer;
  // text-align: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const AccordianContent = styled.div`
  background-color: #f2f6f9;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  margin-bottom: 10px;
`;

const AccordianHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const AccordianBody = styled.div`
  background-color: #f2f6f9;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  border-radius: 8px;
  padding: 20px 0px 16px 0px;
`;

const Title = styled.h2`
  font-weight: 500;
  font-size: 16px;
  color: #00225a;
`;
