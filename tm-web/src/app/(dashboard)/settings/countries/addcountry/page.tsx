"use client";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import React, { useState } from "react";
import styled from "styled-components";
import { FaAngleUp, FaArrowLeft } from "react-icons/fa6";
import { DropdownSelect, showErrorToast } from "@/components";
import Link from "next/link";
import {
  Country,
  ICountryPayload,
  useAddOrUpdateCountryMutation,
} from "@/redux/api/country";
import { useRouter } from "next/navigation";
import RenderField from "../components/RenderField";
type FormFieldsTypes = {
  [K in keyof Omit<Country, "id">]: string;
};

const AddCountryPage = () => {
  const router = useRouter();
  const [formFields, setFormFields] = useState<FormFieldsTypes>({
    country_name: "",
    country_code: "",
    fiat_glyph: "",
    fiat_label: "",
    quote_currency_code: "",
    region_name: "",
    regulator_name: "",
    legal_and_professional_fee_fixed: "",
    legal_and_professional_fee_percent: "",
    rating_agency_fee_fixed: "",
    rating_agency_fee_percent: "",
    min_tokenization_fee: "",
    sec_tokenization_fee: "",
    sec_tokenization_fee_fixed: "",
    sec_tokenization_fee_percent: "",

    sec_trade_fee: "", //  NEW
    sec_trade_fee_fixed: "", // NEW
    sec_trade_fee_percent: "", //  NEW
    tokenization_application_fee: "", //  NEW
    tokenization_application_fee_asset: "", //  NEW
    vat_percent: "",
  });

  const [addCountry, { isLoading }] = useAddOrUpdateCountryMutation();

  const [selectCountry, setSelectCountry] = useState("");
  const [addedCountryName, setAddedCountryName] = useState("");
  const [openAccordian, setOpenAccordian] = useState<string[]>([]);

  const [addSuccess, setAddSuccess] = useState(false);
  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormFields((prevState) => ({ ...prevState, [name]: value }));
  };

  const resetForm = () => {
    setFormFields({
      country_name: "",
      country_code: "",
      fiat_glyph: "",
      fiat_label: "",
      quote_currency_code: "",
      region_name: "",
      regulator_name: "",
      legal_and_professional_fee_fixed: "",
      legal_and_professional_fee_percent: "",
      rating_agency_fee_fixed: "",
      rating_agency_fee_percent: "",
      min_tokenization_fee: "",
      sec_tokenization_fee: "",
      sec_tokenization_fee_fixed: "",
      sec_tokenization_fee_percent: "",

      sec_trade_fee: "", //  NEW
      sec_trade_fee_fixed: "", // NEW
      sec_trade_fee_percent: "", //  NEW
      tokenization_application_fee: "", //  NEW
      tokenization_application_fee_asset: "", //  NEW
      vat_percent: "",
    });
    setSelectCountry("");
  };

  const countires = ["Nigeria", " Canada"];

  const toNumber = (value: string): number => {
    const parsed = Number(value);
    return isNaN(parsed) ? 0 : parsed; // fallback to 0 or handle as needed
  };

  const preprocessPayload = (fields: typeof formFields): ICountryPayload => ({
    action: "create",

    country_code: fields.country_code.trim(),
    country_name: fields.country_name.trim(),
    fiat_glyph: fields.fiat_glyph.trim(),
    fiat_label: fields.fiat_label.trim(),
    quote_currency_code: fields.quote_currency_code.trim(),
    region_name: fields.region_name.trim(),
    regulator_name: fields.regulator_name.trim(),

    legal_and_professional_fee_fixed: toNumber(
      fields.legal_and_professional_fee_fixed
    ),
    legal_and_professional_fee_percent: toNumber(
      fields.legal_and_professional_fee_percent
    ),

    rating_agency_fee_fixed: toNumber(fields.rating_agency_fee_fixed),
    rating_agency_fee_percent: toNumber(fields.rating_agency_fee_percent),

    min_tokenization_fee: toNumber(fields.min_tokenization_fee),

    sec_tokenization_fee: toNumber(fields.sec_tokenization_fee),
    sec_tokenization_fee_fixed: toNumber(fields.sec_tokenization_fee_fixed),
    sec_tokenization_fee_percent: toNumber(fields.sec_tokenization_fee_percent),

    sec_trade_fee: toNumber(fields.sec_trade_fee),
    sec_trade_fee_fixed: toNumber(fields.sec_trade_fee_fixed),
    sec_trade_fee_percent: toNumber(fields.sec_trade_fee_percent),

    tokenization_application_fee: toNumber(fields.tokenization_application_fee),
    tokenization_application_fee_asset:
      fields.tokenization_application_fee_asset.trim(),

    vat_percent: toNumber(fields.vat_percent),
  });

  const handleAddCountry = async () => {
    try {
      const payload = preprocessPayload(formFields);

      await addCountry(payload).unwrap();
      setAddedCountryName(formFields.country_name);
      setAddSuccess(true);
      resetForm();
    } catch (error: any) {
      showErrorToast(error?.data?.message || "Failed to add country");
    }
  };

  const handleOpen = (section: string): void => {
    setOpenAccordian((prev) =>
      prev.includes(section)
        ? prev.filter((sec) => sec !== section)
        : [...prev, section]
    );
  };

  const handleDropdownSelect = (name: keyof FormFieldsTypes, value: string) => {
    setFormFields((prev) => ({ ...prev, [name]: value }));
  };

  // --- CONFIG ---
  const countryFormSections = [
    {
      title: "Country Information",
      fields: [
        {
          label: "Country",
          name: "country_name",
          placeholder: "Eg Nigeria",
          type: "dropdown",
          options: ["Nigeria", "Canada"],
        },
        {
          label: "Country Code",
          name: "country_code",
          placeholder: "Eg USA",
          type: "text",
        },
        {
          label: "Fiat Glyph",
          name: "fiat_glyph",
          placeholder: "Eg $",
          type: "text",
        },
        {
          label: "Fiat Label",
          name: "fiat_label",
          placeholder: "Eg Dollar",
          type: "text",
        },
        {
          label: "Quote Currency Code",
          name: "quote_currency_code",
          placeholder: "Eg USD",
          type: "text",
        },
        {
          label: "Region Name",
          name: "region_name",
          placeholder: "",
          type: "text",
        },
      ],
    },
    {
      title: "Regulator",
      fields: [
        {
          label: "Regulator Name",
          name: "regulator_name",
          placeholder: "",
          type: "text",
        },
      ],
    },
    {
      title: "Legal Fees",
      fields: [
        {
          label: "Legal and Professional Fixed Fee",
          name: "legal_and_professional_fee_fixed",
          placeholder: "0",
          type: "fee",
          suffix: "CNGN",
        },
        {
          label: "Legal and Professional Fee Percentage",
          name: "legal_and_professional_fee_percent",
          placeholder: "0",
          type: "fee",
          suffix: "%",
        },
      ],
    },
    {
      title: "Rating Agency Fees",
      fields: [
        {
          label: "Rating Agency Fixed Fee",
          name: "rating_agency_fee_fixed",
          placeholder: "0",
          type: "fee",
          suffix: "CNGN",
        },
        {
          label: "Rating Agency Fee Percentage",
          name: "rating_agency_fee_percent",
          placeholder: "0",
          type: "fee",
          suffix: "%",
        },
      ],
    },
    {
      title: "Tokenization Fees",
      fields: [
        {
          label: "Minimum Tokenization Fee",
          name: "min_tokenization_fee",
          placeholder: "0",
          type: "fee",
          suffix: "CNGN",
        },
      ],
    },

    {
      title: "Tokenization Application Fees",
      fields: [
        {
          label: " Application Fee",
          name: "tokenization_application_fee",
          placeholder: "0",
          type: "fee",
          suffix: "CNGN",
        },
        {
          label: " Application Fee  Asset",
          name: "tokenization_application_fee_asset",
          placeholder: "",
          type: "text",
        },
      ],
    },
    {
      title: "Regulatory Fee",
      fields: [
        {
          label: "SEC Tokenization Fee",
          name: "sec_tokenization_fee",
          placeholder: "0",
          type: "fee",
          suffix: "CNGN",
        },
        {
          label: "SEC Tokenization Fee Percentage",
          name: "sec_tokenization_fee_percent",
          placeholder: "0",
          type: "fee",
          suffix: "%",
        },
        {
          label: "SEC Trade Fee",
          name: "sec_trade_fee",
          placeholder: "0",
          type: "fee",
          suffix: "CNGN",
        },
        {
          label: "SEC Trade Fee Fixed",
          name: "sec_trade_fee_fixed",
          placeholder: "0",
          type: "text",
        },
        {
          label: "SEC Trade Fee Percent",
          name: "sec_trade_fee_percent",
          placeholder: "0",
          type: "fee",
          suffix: "%",
        },
      ],
    },

    {
      title: "VAT Percent",
      fields: [
        {
          label: "VAT Percent",
          name: "vat_percent",
          placeholder: "1.3",
          type: "fee",
          suffix: "%",
        },
      ],
    },
  ];

  return (
    <AddCountryWrapper>
      <BackButtonLink href="/countries">
        <FaArrowLeft />
      </BackButtonLink>

      <CountryFormCard>
        <FormTitle>Add Country</FormTitle>
        {countryFormSections.map((section) => (
          <div key={section.title}>
            <AccordianHeader onClick={() => handleOpen(section.title)}>
              <SectionTitle>{section.title}</SectionTitle>

              <FaAngleUp />
            </AccordianHeader>
            {openAccordian.includes(section.title) && (
              <>
                {section.fields.map((field) => (
                  <FormFieldGroup key={field.name}>
                    <FormLabel>{field.label}</FormLabel>
                    <RenderField
                      field={field}
                      value={formFields[field.name as keyof FormFieldsTypes]}
                      handleInputChange={handleInputChange}
                      handleDropdownSelect={
                        field.type === "dropdown"
                          ? (val: string) =>
                              handleDropdownSelect(
                                field.name as keyof FormFieldsTypes,
                                val
                              )
                          : undefined
                      }
                    />
                  </FormFieldGroup>
                ))}
              </>
            )}
          </div>
        ))}

        <PrimaryButton
          onClick={handleAddCountry}
          buttonStyle={{ width: "100%" }}
        >
          {isLoading ? "Adding..." : "Add Country"}
        </PrimaryButton>
      </CountryFormCard>

      {addSuccess && (
        <SuccessMessage
          isOpen={addSuccess}
          setIsOpen={(val) => {
            setAddSuccess(val);
            if (!val) {
              router.push("/countries");
            }
          }}
          heading="Success!"
          message="You have successfully to country list."
          email={addedCountryName || "Unnamed country"}
        />
      )}
    </AddCountryWrapper>
  );
};

export default AddCountryPage;
const AddCountryWrapper = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const CountryFormCard = styled.div`
  width: 480px;
  display: flex;
  margin: auto;
  flex-direction: column;
  gap: 16px;
`;

const AccordianHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;
const FormTitle = styled.h2`
  color: #00225a;
  text-align: center;
  font-weight: 700;
  font-size: 24px;
  line-height: 28px;
`;

const SectionTitle = styled.p`
  color: #00225a;
  font-weight: 600;
  font-size: 16px;
  line-height: 24px;
  padding: 10px 0;
`;

const FormFieldGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const FormLabel = styled.label`
  font-weight: 500;
  font-size: 14px;
  color: #828282;
  padding: 4px 0;
`;

const BackButtonLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;
