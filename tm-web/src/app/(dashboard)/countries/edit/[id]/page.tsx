"use client";

import { DropdownSelect, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import {
  Country,
  ICountryPayload,
  useAddOrUpdateCountryMutation,
  useGetCountryByIdQuery,
} from "@/redux/api/country";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { FaAngleUp, FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";
import RenderField from "../../components/RenderField";
type FormFieldsTypes = {
  [K in keyof Omit<Country, "id">]: string;
};
const EditCountryPage = () => {
  const router = useRouter();
  const { id } = useParams();
  const countryId = Number(id);
  const [openAccordian, setOpenAccordian] = useState<string[]>([]);
  const { data: countryData, isLoading: isFetching } =
    useGetCountryByIdQuery(countryId);
  const [updateCountry, { isLoading: isUpdating }] =
    useAddOrUpdateCountryMutation();

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
    sec_trade_fee: "",
    sec_trade_fee_fixed: "",
    sec_trade_fee_percent: "",
    tokenization_application_fee: "",
    tokenization_application_fee_asset: "",
    vat_percent: "",
  });
  const [addSuccess, setAddSuccess] = useState(false);
  const [selectCountry, setSelectCountry] = useState("");

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormFields((prevValue) => ({ ...prevValue, [name]: value }));
  };
  useEffect(() => {
    if (countryData) {
      setFormFields({
        country_name: countryData.country_name || "",
        country_code: countryData.country_code || "",
        fiat_glyph: countryData.fiat_glyph || "",
        fiat_label: countryData.fiat_label || "",
        quote_currency_code: countryData.quote_currency_code || "",
        region_name: countryData.region_name || "",
        regulator_name: countryData.regulator_name || "",
        legal_and_professional_fee_fixed:
          countryData.legal_and_professional_fee_fixed?.toString() || "",
        legal_and_professional_fee_percent:
          countryData.legal_and_professional_fee_percent?.toString() || "",
        rating_agency_fee_fixed:
          countryData.rating_agency_fee_fixed?.toString() || "",
        rating_agency_fee_percent:
          countryData.rating_agency_fee_percent?.toString() || "",
        min_tokenization_fee:
          countryData.min_tokenization_fee?.toString() || "",
        sec_tokenization_fee:
          countryData.sec_tokenization_fee?.toString() || "",
        sec_tokenization_fee_fixed:
          countryData.sec_tokenization_fee_fixed?.toString() || "",
        sec_tokenization_fee_percent:
          countryData.sec_tokenization_fee_percent?.toString() || "",
        sec_trade_fee: countryData.sec_trade_fee?.toString() || "",
        sec_trade_fee_fixed: countryData.sec_trade_fee_fixed?.toString() || "",
        sec_trade_fee_percent:
          countryData.sec_trade_fee_percent?.toString() || "",
        tokenization_application_fee:
          countryData.tokenization_application_fee?.toString() || "",
        tokenization_application_fee_asset:
          countryData.tokenization_application_fee_asset || "",
        vat_percent: countryData.vat_percent?.toString() || "",
      });
    }
  }, [countryData]);

  const toNumber = (value: string) =>
    isNaN(Number(value)) ? 0 : Number(value);

  const handleUpdateCountry = async () => {
    try {
      const payload: ICountryPayload = {
        action: "update",
        id: countryId,
        country_code: formFields.country_code.trim(),
        country_name: formFields.country_name.trim(),
        fiat_glyph: formFields.fiat_glyph.trim(),
        fiat_label: formFields.fiat_label.trim(),
        quote_currency_code: formFields.quote_currency_code.trim(),
        region_name: formFields.region_name.trim(),
        regulator_name: formFields.regulator_name.trim(),
        legal_and_professional_fee_fixed: toNumber(
          formFields.legal_and_professional_fee_fixed
        ),
        legal_and_professional_fee_percent: toNumber(
          formFields.legal_and_professional_fee_percent
        ),
        rating_agency_fee_fixed: toNumber(formFields.rating_agency_fee_fixed),
        rating_agency_fee_percent: toNumber(
          formFields.rating_agency_fee_percent
        ),
        min_tokenization_fee: toNumber(formFields.min_tokenization_fee),
        sec_tokenization_fee: toNumber(formFields.sec_tokenization_fee),
        sec_tokenization_fee_fixed: toNumber(
          formFields.sec_tokenization_fee_fixed
        ),
        sec_tokenization_fee_percent: toNumber(
          formFields.sec_tokenization_fee_percent
        ),
        sec_trade_fee: toNumber(formFields.sec_trade_fee),
        sec_trade_fee_fixed: toNumber(formFields.sec_trade_fee_fixed),
        sec_trade_fee_percent: toNumber(formFields.sec_trade_fee_percent),
        tokenization_application_fee: toNumber(
          formFields.tokenization_application_fee
        ),
        tokenization_application_fee_asset:
          formFields.tokenization_application_fee_asset.trim(),
        vat_percent: toNumber(formFields.vat_percent),
      };

      const response = await updateCountry(payload).unwrap();

      setAddSuccess(true);
    } catch (error: any) {
      showErrorToast(error?.data?.message || "Failed to add country");
    }
  };

  const countires = ["Nigeria"];

  const handleDropdownSelect = (name: keyof FormFieldsTypes, value: string) => {
    setFormFields((prev) => ({ ...prev, [name]: value }));
  };

  const handleOpen = (section: string): void => {
    setOpenAccordian((prev) =>
      prev.includes(section)
        ? prev.filter((sec) => sec !== section)
        : [...prev, section]
    );
  };

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
        <FormTitle>Edit Country</FormTitle>
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
                    {field.name === "country_name" ? (
                      <NonEditableValue>
                        {formFields.country_name}
                      </NonEditableValue>
                    ) : (
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
                    )}
                  </FormFieldGroup>
                ))}
              </>
            )}
          </div>
        ))}

        <PrimaryButton
          onClick={handleUpdateCountry}
          buttonStyle={{ width: "100%" }}
        >
          {isUpdating ? "Updating..." : "Update Country"}
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
          message="You have successfully updated"
          email={formFields.country_name || "Unnamed country"}
        />
      )}
    </AddCountryWrapper>
  );
};

export default EditCountryPage;

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

const AccordianHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;
const BackButtonLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;
const NonEditableValue = styled.div`
  padding: 10px 12px;
  background: #f6f6f6;
  color: #828282;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  font-size: 15px;
  font-weight: 500;
`;
