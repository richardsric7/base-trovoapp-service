"use client";

import { DropdownSelect, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import SuccessMessage from "@/components/SuccessMessage";
import { useCreatePartnerMutation } from "@/redux/api/tokenizationstakeholders";
import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { FaArrowLeft } from "react-icons/fa6";
import styled from "styled-components";

const AddStakeholder = () => {
  const [addSuccess, setAddSuccess] = useState(false);
  const [selectValues, setSelectValues] = useState<any>({});
  const [formValues, setFormValues] = useState({
    name: "",
    address: "",
    percentageFee: "",
    fixedFee: "",
  });
  const [submittedName, setSubmittedName] = useState("");
  const router = useRouter();
  const [createPartner, { isLoading }] = useCreatePartnerMutation();

  const filterCriteria = {
    Category: [
      "asset_manager",
      "asset_issuing_house",
      "approved_asset_custodian",
      "legal_and_professionals",
      "rating_agency",
      "trustees",
      "legal_adviser",
      "financial_adviser",
    ],
  };

  const displayLabelMap: Record<string, string> = {
    asset_manager: "Asset Manager",
    asset_issuing_house: "Issuing House",
    approved_asset_custodian: "Asset Custodian",
    legal_and_professionals: "Legal & Professionals",
    rating_agency: "Rating Agency",
    trustees: "Trustee",
    legal_adviser: "Legal Adviser",
    financial_adviser: "Financial Adviser",
  };

  const handleAddStakeholder = async () => {
    const { name, address, percentageFee, fixedFee } = formValues;
    const stakeholderType = selectValues.Category;

    if (!stakeholderType) {
      showErrorToast("Please select stakeholder type.");
      return;
    }

    if (!percentageFee) {
      showErrorToast("Percentage fee is required!");
      return;
    }

    if (!fixedFee) {
      showErrorToast("Fixed fee is required!");
      return;
    }

    if (!address) {
      showErrorToast("Address is required!");
      return;
    }
    const payload: any = {};
    console.log("pay", payload);

    if (stakeholderType === "asset_manager") {
      payload.asset_manager_name = name;
      payload.asset_manager_address = address;
      payload.asset_manager_country = "NG";
      payload.fee_percent = parseFloat(percentageFee);
      payload.fee_fixed = parseFloat(fixedFee);

      payload.requirement_document = "<h1>NA</h1>";
    } else if (stakeholderType === "asset_issuing_house") {
      payload.asset_issuing_house_name = name;
      payload.asset_issuing_house_address = address;
      payload.asset_issuing_house_country = "NG";
      payload.fee_percent = parseFloat(percentageFee);
      payload.fee_fixed = parseFloat(fixedFee);
      payload.requirement_document = "<h1>NA</h1>";
    } else if (stakeholderType === "approved_asset_custodian") {
      payload.asset_custodian_name = name;
      payload.asset_custodian_address = address;
      payload.asset_custodian_country = "NG";
      payload.fee_percent = percentageFee;
      payload.fee_fixed = fixedFee;
      payload.requirement_document = "<h1>NA</h1>";
    } else if (stakeholderType === "legal_and_professionals") {
      payload.partner_name = name;
      payload.partner_address = address;
      payload.partner_country = "NG";
      payload.fee_percent = percentageFee;
      payload.fee_fixed = fixedFee;
    } else if (stakeholderType === "rating_agency") {
      payload.agency_name = name;
      payload.agency_address = address;
      payload.agency_country = "NG";
      payload.fee_percent = percentageFee;
      payload.fee_fixed = fixedFee;
    } else if (stakeholderType === "trustees") {
      payload.trustee_name = name;
      payload.trustee_address = address;
      payload.trustee_country = "NG";
      payload.fee_percent = percentageFee;
      payload.fee_fixed = fixedFee;
    } else if (
      stakeholderType === "legal_adviser" ||
      stakeholderType === "financial_adviser"
    ) {
      payload.adviser_name = name;
      payload.adviser_address = address;
      payload.adviser_country = "NG";
      payload.fee_percent = percentageFee;
      payload.fee_fixed = fixedFee;
    }

    try {
      await createPartner({
        type: stakeholderType,
        data: payload,
        action: "create",
      }).unwrap();
      setSubmittedName(name);
      setAddSuccess(true);
      resetForm();
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
      // Show the error toast.
      showErrorToast(finalErrorMessage);
    }
  };

  const resetForm = () => {
    setSelectValues({});
    setFormValues({
      name: "",
      address: "",
      percentageFee: "",
      fixedFee: "",
    });
  };
  return (
    <AddStakeholderWrapper>
      <BackButtonLink href="/settings/tokenizationstakeholders">
        <FaArrowLeft />
      </BackButtonLink>

      <StakeholderFormCard>
        <FormTitle>Add Tokenization Stakeholders</FormTitle>

        <div>
          <SectionTitle>Stakeholder Information</SectionTitle>

          <FormFieldGroup>
            <FormLabel>Stakeholder Type</FormLabel>
            <DropdownSelect
              options={filterCriteria.Category.map(
                (value) => displayLabelMap[value]
              )}
              placeholder=""
              labelText=""
              value={displayLabelMap[selectValues["Category"]] || ""}
              onSelect={(label) => {
                // Find the corresponding value
                const value = Object.keys(displayLabelMap).find(
                  (key) => displayLabelMap[key] === label
                );
                if (value) {
                  setSelectValues({ ...selectValues, Category: value });
                }
              }}
            />
          </FormFieldGroup>

          <FormFieldGroup>
            <FormLabel>Stakeholder Name</FormLabel>
            <TextInput
              placeholder="Enter entity name"
              value={formValues.name}
              onChange={(e) =>
                setFormValues({ ...formValues, name: e.target.value })
              }
            />
          </FormFieldGroup>

          <FormFieldGroup>
            <FormLabel>Stakeholder&apos;s Address</FormLabel>
            <TextInput
              placeholder="Enter entity address"
              value={formValues.address}
              onChange={(e) =>
                setFormValues({ ...formValues, address: e.target.value })
              }
            />
          </FormFieldGroup>
        </div>

        <div>
          <SectionTitle>Fees</SectionTitle>

          <FormFieldGroup>
            <FormLabel>Fee Percentage</FormLabel>
            <FeeInputWrapper>
              <FeeInputField
                placeholder="0"
                value={formValues.percentageFee}
                onChange={(e) =>
                  setFormValues({
                    ...formValues,
                    percentageFee: e.target.value,
                  })
                }
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
                placeholder="0"
                value={formValues.fixedFee}
                onChange={(e) =>
                  setFormValues({ ...formValues, fixedFee: e.target.value })
                }
              />
              <DividerGroup>
                <VerticalDivider />
                <FeeLabel>CNGN</FeeLabel>
              </DividerGroup>
            </FeeInputWrapper>
          </FormFieldGroup>
        </div>

        <PrimaryButton
          onClick={handleAddStakeholder}
          buttonStyle={{ width: "100%" }}
        >
          {isLoading ? "Adding..." : "Add Stakeholder"}
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
          // setIsOpen={setAddSuccess}
          heading="Success!"
          message="You have successfully added the Tokenization Stakeholder."
          email={submittedName}
        />
      )}
    </AddStakeholderWrapper>
  );
};

export default AddStakeholder;

const AddStakeholderWrapper = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const StakeholderFormCard = styled.div`
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
  border-radius: 8px;
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
