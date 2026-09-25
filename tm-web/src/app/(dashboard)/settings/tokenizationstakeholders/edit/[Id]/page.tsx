"use client";

import { DropdownSelect, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import {
  useCreatePartnerMutation,
  useGetSinglePartnerQuery,
} from "@/redux/api/tokenizationstakeholders";
import { skipToken } from "@reduxjs/toolkit/query";
import Link from "next/link";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import React, { useEffect, useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";

const EditStakeholderPage = () => {
  const [addSuccess, setAddSuccess] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [percentageFee, setPercentageFee] = useState("");
  const [fixedFee, setFixedFee] = useState("");

  const params = useParams();

  const searchParams = useSearchParams();
  const stakeholderType = searchParams.get("type") as
    | "asset_manager"
    | "asset_issuing_house"
    | "approved_asset_custodian"
    | "legal_and_professionals"
    | "rating_agency"
    | "trustees"
    | "legal_adviser"
    | "financial_adviser";

  // Fix: Extract the ID using the correct case (capital 'I' in 'Id')
  const id = params?.Id || undefined;

  // Ensure we have a numeric value
  const idNumber = id ? Number(id) : undefined;

  const { data: stakeholderData } = useGetSinglePartnerQuery(
    idNumber !== undefined && !isNaN(idNumber) && stakeholderType
      ? { id: idNumber, type: stakeholderType }
      : skipToken
  );

  const [updatePartner, { isLoading }] = useCreatePartnerMutation();
  const router = useRouter();
  // console.log("idNumber:", idNumber);
  // console.log("stakeholderType:", stakeholderType);

  // console.log("stakeholderData: ", stakeholderData);
  useEffect(() => {
    if (stakeholderData && stakeholderType) {
      setSelectValues({ Category: stakeholderType });

      if ("asset_manager_name" in stakeholderData) {
        setName(stakeholderData.asset_manager_name || "");
        setAddress(stakeholderData.asset_manager_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      } else if ("asset_issuing_house_name" in stakeholderData) {
        setName(stakeholderData.asset_issuing_house_name || "");
        setAddress(stakeholderData.asset_issuing_house_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      } else if ("asset_custodian_name" in stakeholderData) {
        setName(stakeholderData.asset_custodian_name || "");
        setAddress(stakeholderData.asset_custodian_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      } else if ("partner_name" in stakeholderData) {
        setName(stakeholderData.partner_name || "");
        setAddress(stakeholderData.partner_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      } else if ("agency_name" in stakeholderData) {
        setName(stakeholderData.agency_name || "");
        setAddress(stakeholderData.agency_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      } else if ("trustee_name" in stakeholderData) {
        setName(stakeholderData.trustee_name || "");
        setAddress(stakeholderData.trustee_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      } else if ("adviser_name" in stakeholderData) {
        setName(stakeholderData.adviser_name || "");
        setAddress(stakeholderData.adviser_address || "");
        setPercentageFee(stakeholderData.fee_percent?.toString() || "");
        setFixedFee(stakeholderData.fee_fixed?.toString() || "");
      }
    }
  }, [stakeholderData, stakeholderType]);

  const handleUpdateStakeholder = async () => {
    try {
      if (!idNumber || !selectValues.Category) {
        throw new Error("Missing stakeholder ID or type.");
      }

      const payload: any = {
        id: idNumber,
      };

      if (selectValues.Category === "asset_manager") {
        payload.asset_manager_name = name;
        payload.asset_manager_address = address;
        payload.asset_manager_country = "NG";
        payload.fee_percent = parseFloat(percentageFee);
        payload.fee_fixed = parseFloat(fixedFee);
        payload.requirement_document = "<h1>NA</h1>";
      } else if (selectValues.Category === "asset_issuing_house") {
        payload.asset_issuing_house_name = name;
        payload.asset_issuing_house_address = address;
        payload.asset_issuing_house_country = "NG";
        payload.fee_percent = parseFloat(percentageFee);
        payload.fee_fixed = parseFloat(fixedFee);
        payload.requirement_document = "<h1>NA</h1>";
      } else if (selectValues.Category === "approved_asset_custodian") {
        payload.asset_custodian_name = name;
        payload.asset_custodian_address = address;
        payload.asset_custodian_country = "NG";
        payload.fee_percent = percentageFee;
        payload.fee_fixed = fixedFee;
        payload.requirement_document = "<h1>NA</h1>";
      } else if (selectValues.Category === "legal_and_professionals") {
        payload.partner_name = name;
        payload.partner_address = address;
        payload.partner_country = "NG";
        payload.fee_percent = percentageFee;
        payload.fee_fixed = fixedFee;
      } else if (selectValues.Category === "rating_agency") {
        payload.agency_name = name;
        payload.agency_address = address;
        payload.agency_country = "NG";
        payload.fee_percent = percentageFee;
        payload.fee_fixed = fixedFee;
      } else if (selectValues.Category === "trustees") {
        payload.trustee_name = name;
        payload.trustee_address = address;
        payload.trustee_country = "NG";
        payload.fee_percent = percentageFee;
        payload.fee_fixed = fixedFee;
      } else if (selectValues.Category === "legal_adviser") {
        payload.adviser_name = name;
        payload.adviser_address = address;
        payload.adviser_country = "NG";
        payload.fee_percent = percentageFee;
        payload.fee_fixed = fixedFee;
      } else if (selectValues.Category === "financial_adviser") {
        payload.adviser_name = name;
        payload.adviser_address = address;
        payload.adviser_country = "NG";
        payload.fee_percent = percentageFee;
        payload.fee_fixed = fixedFee;
      }
      await updatePartner({
        type: selectValues.Category,
        action: "update",

        data: payload,
      }).unwrap();

      setAddSuccess(true);
    } catch (error: any) {
      // Retrieve the error string from the response:
      const errorString =
        error?.response?.data?.message ||
        error?.data?.message ||
        error?.message ||
        error?.response?.data?.error ||
        error?.data?.error ||
        error?.error ||
        "An error occurred while updating.";
      const messageMatch = errorString.match(/message:(.*?)(\]$|$)/);
      let finalErrorMessage = errorString;
      if (messageMatch && messageMatch[1]) {
        finalErrorMessage = messageMatch[1].trim();
      }
      showErrorToast(finalErrorMessage);
    }
  };

  const stakeholderOptions = [
    { label: "Asset Manager", value: "asset_manager" },
    { label: "Issuing House", value: "asset_issuing_house" },
    { label: "Asset Custodian", value: "approved_asset_custodian" },
    { label: "Legal & Professional", value: "legal_and_professionals" },
    { label: "Rating Agency", value: "rating_agency" },
    { label: "Trustee", value: "trustees" },
    { label: "Legal Adviser", value: "legal_adviser" },
    { label: "Financial Adviser", value: "financial_adviser" },
  ];

  const labelFromValue = (val: string | undefined) =>
    stakeholderOptions.find((opt) => opt.value === val)?.label || "";

  const valueFromLabel = (label: string) =>
    stakeholderOptions.find((opt) => opt.label === label)?.value || "";

  return (
    <EditStakeholderWrapper>
      <BackButtonLink href="/settings/tokenizationstakeholders">
        <FaArrowLeft />
      </BackButtonLink>

      <StakeholderFormCard>
        <FormTitle>Edit Tokenization Stakeholders</FormTitle>

        <div>
          <SectionTitle>Stakeholder Information</SectionTitle>

          <FormFieldGroup>
            <FormLabel>Stakeholder Type</FormLabel>
            <DropdownSelect
              options={stakeholderOptions.map((opt) => opt.label)}
              value={labelFromValue(selectValues.Category || stakeholderType)}
              onSelect={(label) => {
                const value = valueFromLabel(label);
                if (value) {
                  setSelectValues({ ...selectValues, Category: value });
                }
              }}
            />
          </FormFieldGroup>

          <FormFieldGroup>
            <FormLabel>Stakeholder Name</FormLabel>
            <TextInput value={name} onChange={(e) => setName(e.target.value)} />
          </FormFieldGroup>

          <FormFieldGroup>
            <FormLabel>Stakeholder&apos;s Address</FormLabel>
            <TextInput
              value={address}
              onChange={(e) => setAddress(e.target.value)}
            />
          </FormFieldGroup>
        </div>

        <div>
          <SectionTitle>Fees</SectionTitle>

          <FormFieldGroup>
            <FormLabel>Fee Percentage</FormLabel>
            <FeeInputWrapper>
              <FeeInputField
                value={percentageFee}
                onChange={(e) => setPercentageFee(e.target.value)}
              />
              <DividerGroup>
                <VerticalDivider />
                <FeeLabel>%</FeeLabel>
              </DividerGroup>
            </FeeInputWrapper>
          </FormFieldGroup>
          <FormFieldGroup>
            <FormLabel>Fixed Fee</FormLabel>
            <FeeInputWrapper>
              <FeeInputField
                value={fixedFee}
                onChange={(e) => setFixedFee(e.target.value)}
              />
              <DividerGroup>
                <VerticalDivider />
                <FeeLabel>CNGN</FeeLabel>
              </DividerGroup>
            </FeeInputWrapper>
          </FormFieldGroup>
        </div>

        <PrimaryButton
          onClick={handleUpdateStakeholder}
          buttonStyle={{ width: "100%" }}
        >
          {isLoading ? "Updating..." : "Update Stakeholder"}
        </PrimaryButton>
      </StakeholderFormCard>

      {addSuccess && (
        <SuccessMessage
          isOpen={addSuccess}
          setIsOpen={(val) => {
            setAddSuccess(val);
            if (!val) {
              router.push("/settings/tokenizationstakeholders");
            }
          }}
          heading="Success!"
          message={
            <>
              You have successfully updated <strong>{name}</strong> as a
              Tokenization Stakeholder
            </>
          }
          email=""
        />
      )}
    </EditStakeholderWrapper>
  );
};

export default EditStakeholderPage;

const EditStakeholderWrapper = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const StakeholderFormCard = styled.section`
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

const TextInput = styled.input`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
`;

const FeeInputField = styled.input`
  font-size: 14px;
  width: 100%;
  outline: none;
  font-family: inherit;
  border: none;
`;

const FeeInputWrapper = styled.div`
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
`;

const DividerGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const VerticalDivider = styled.div`
  width: 1px;
  height: 24px;
  background-color: #e0e0e0;
`;

const FeeLabel = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #191919;
`;

const BackButtonLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
`;
