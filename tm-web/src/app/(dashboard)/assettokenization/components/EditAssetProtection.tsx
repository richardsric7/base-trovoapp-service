"use client"

import { DropdownSelect, Modal } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
} from "@/redux/api/assettokenization";
import React, { useEffect, useState } from "react";
import { FaAngleDown } from "react-icons/fa";
import styled from "styled-components";

interface AssetProtectionProps {
  isOpen: boolean;
  asset?: TokenizationRecord;
  onClose: () => void;
  onSave?: (updatedData: Partial<UpdateTokenizationPayload>) => void;
}

const EditAssetProtection: React.FC<AssetProtectionProps> = ({
  isOpen,
  asset,
  onClose,
  onSave,
}) => {
  const [editedAsset, setEditedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >(asset || {});

  const [openSections, setOpenSections] = useState<Record<string, boolean>>({
    insurance: true,
    protection: true,
    govt: true,
    risk: true,
    esg: true,
    sec: true,
    legal: true,
  });

  useEffect(() => {
    if (asset) {
      setEditedAsset(asset as Partial<UpdateTokenizationPayload>);
    }
  }, [asset]);

  // Define numeric fields excluding legal/financial counsel and independent monitoring list
  const numericFields: Array<keyof UpdateTokenizationPayload> = [
    "insurancePolicyNumber",
    "percentageValueOfInsurance",
    "contractualProtectionRevGuarantees",
    "contractualProtectionPerfBond",
    "contractualProtectionSLA",
    "riskSharingMechanismPPPs",
    "riskSharingMechanismHedgeInstruments",
    "eSGSafeguardsSusCerts",
    "eSGSafeguardsCommEngPlans",
    "securityMeasuresAccessControl",
    "securityMeasuresSurveilanceSystems",
    "securityMeasuresOnSiteSecurityPersonnel",
    "securityMeasuresPerimeterSecurity",
    "securityMeasuresCriticalInfraProtections",
  ];

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    let newValue: string | number = value;
    // Convert to number if the field is supposed to be numeric

    if (numericFields.includes(name as keyof UpdateTokenizationPayload)) {
      newValue = Number(value);
    }

    setEditedAsset((prev) => ({
      ...prev,
      [name]: newValue,
    }));
  };

  const toogleSection = (section: string) => {
    setOpenSections((prev) => ({
      ...prev,

      [section]: !prev[section],
    }));
  };
  return (
    <>
      <Modal
        title="Edit Asset Protection"
        isOpen={isOpen}
        onClose={onClose}
        closeIconPosition="right"
      >
        <ModalContent>
          <ModalSections>
            <SubHeading onClick={() => toogleSection("insurance")}>
              <Title>INSURANCE</Title>
              <FaAngleDown />
            </SubHeading>
            {openSections.insurance && (
              <>
                <FormGroup>
                  <Label>Insurance Company Name</Label>
                  <Input
                    name="insuranceCompanyName"
                    value={editedAsset.insuranceCompanyName || ""}
                    onChange={handleChange}
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Insurance Policy Number</Label>
                  <Input
                    name="insurancePolicyNumber"
                    value={editedAsset.insurancePolicyNumber || ""}
                    onChange={handleChange}
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Insurance Policy Holder</Label>
                  <Input
                    name="insurancePolicyHolder"
                    value={editedAsset.insurancePolicyHolder || ""}
                    onChange={handleChange}
                  />
                </FormGroup>

                <FormGroup>
                  <Label>Insurance Percentage Value</Label>
                  <Input
                    name="percentageValueOfInsurance"
                    value={editedAsset.percentageValueOfInsurance || ""}
                    onChange={handleChange}
                  />
                </FormGroup>
              </>
            )}
          </ModalSections>

          <ModalSections>
            <SubHeading onClick={() => toogleSection("protection")}>
              <Title>CONTARCTUAL PROTECTION</Title>
              <FaAngleDown />
            </SubHeading>

            {openSections.protection && (
              <>
                <CheckWrapper>
                  {" "}
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="contractualProtectionRevGuarantees"
                      checked={
                        editedAsset.contractualProtectionRevGuarantees === 1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          contractualProtectionPerfBond: e.target.checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Revenue Guarantees
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="contractualProtectionPerfBond"
                      checked={editedAsset.contractualProtectionPerfBond === 1}
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          contractualProtectionPerfBond: e.target.checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Performance Bond
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="contractualProtectionSLA"
                      checked={editedAsset.contractualProtectionSLA === 1}
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          contractualProtectionSLA: e.target.checked ? 1 : 0,
                        }))
                      }
                    />
                    Service Level Agreements
                  </CheckLabel>
                </CheckWrapper>{" "}
              </>
            )}
          </ModalSections>

          <ModalSections>
            <SubHeading onClick={() => toogleSection("risk")}>
              <Title>RISK SHARING MECHANISMS</Title>
              <FaAngleDown />
            </SubHeading>
            {openSections.risk && (
              <>
                <CheckWrapper>
                  {" "}
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="riskSharingMechanismCompletionGuarantees"
                      checked={
                        editedAsset.riskSharingMechanismCompletionGuarantees ===
                        1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          riskSharingMechanismCompletionGuarantees: e.target
                            .checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Completion Guarantees
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="riskSharingMechanismHedgeInstruments"
                      checked={
                        editedAsset.riskSharingMechanismHedgeInstruments === 1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          riskSharingMechanismHedgeInstruments: e.target.checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Hedging Instruments
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="riskSharingMechanismPPPs"
                      checked={editedAsset.riskSharingMechanismPPPs === 1}
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          riskSharingMechanismPPPs: e.target.checked ? 1 : 0,
                        }))
                      }
                    />
                    Public Private Patnerships
                  </CheckLabel>
                </CheckWrapper>{" "}
              </>
            )}
          </ModalSections>

          <ModalSections>
            <SubHeading onClick={() => toogleSection("govt")}>
              <Title>GOVERNANCE AND OVERSIGHT</Title>
              <FaAngleDown />
            </SubHeading>
            {openSections.govt && (
              <>
                {" "}
                <FormGroup>
                  <Label>Independent Monitoring List</Label>
                  <Input
                    name="independentMonitoringList"
                    value={editedAsset.independentMonitoringList || ""}
                    onChange={handleChange}
                  />
                </FormGroup>{" "}
              </>
            )}
          </ModalSections>

          <ModalSections>
            <SubHeading onClick={() => toogleSection("esg")}>
              <Title>
                ENVIROMENTAL, SOCIAL, AND GOVERNANCE (ESG) SAFEGUARDS
              </Title>
              <FaAngleDown />
            </SubHeading>
            {openSections.esg && (
              <>
                <CheckWrapper>
                  {" "}
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="eSGSafeguardsSusCerts"
                      checked={editedAsset.eSGSafeguardsSusCerts === 1}
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          eSGSafeguardsSusCerts: e.target.checked ? 1 : 0,
                        }))
                      }
                    />
                    Sustainability Certifications
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="eSGSafeguardsCommEngPlans"
                      checked={editedAsset.eSGSafeguardsCommEngPlans === 1}
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          eSGSafeguardsCommEngPlans: e.target.checked ? 1 : 0,
                        }))
                      }
                    />
                    Community Engagement Plans
                  </CheckLabel>
                </CheckWrapper>
              </>
            )}
          </ModalSections>
          <ModalSections>
            <SubHeading onClick={() => toogleSection("sec")}>
              <Title>SECURITY MEASURES</Title>
              <FaAngleDown />
            </SubHeading>

            {openSections.sec && (
              <>
                <CheckWrapper>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="securityMeasuresSurveilanceSystems"
                      checked={
                        editedAsset.securityMeasuresSurveilanceSystems === 1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          securityMeasuresSurveilanceSystems: e.target.checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Surveillance Systems
                  </CheckLabel>{" "}
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="securityMeasuresOnSiteSecurityPersonnel"
                      checked={
                        editedAsset.securityMeasuresOnSiteSecurityPersonnel ===
                        1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          securityMeasuresOnSiteSecurityPersonnel: e.target
                            .checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    On-Site Security Personnel
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="securityMeasuresPerimeterSecurity"
                      checked={
                        editedAsset.securityMeasuresPerimeterSecurity === 1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          securityMeasuresPerimeterSecurity: e.target.checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Perimeter Security
                  </CheckLabel>
                  <CheckLabel>
                    <Check
                      type="checkbox"
                      name="securityMeasuresCriticalInfraProtections"
                      checked={
                        editedAsset.securityMeasuresCriticalInfraProtections ===
                        1
                      }
                      onChange={(e) =>
                        setEditedAsset((prev) => ({
                          ...prev,
                          securityMeasuresCriticalInfraProtections: e.target
                            .checked
                            ? 1
                            : 0,
                        }))
                      }
                    />
                    Critical Infrastructure Protections
                  </CheckLabel>
                </CheckWrapper>
              </>
            )}
          </ModalSections>

          <ModalSections>
            <SubHeading onClick={() => toogleSection("insurance")}>
              <Title>LEGAL/FINANCIAL COUNSEL</Title>
              <FaAngleDown color="#00225A" />
            </SubHeading>
            {openSections.legal && (
              <>
                <FormGroup>
                  <Label>Legal Counsel</Label>
                  <Input
                    name="legalAdvisor"
                    value={editedAsset.legalAdvisor || ""}
                    onChange={handleChange}
                  />
                </FormGroup>
                <FormGroup>
                  <Label>Financial Counsel</Label>
                  <Input
                    name="financialAdvisor"
                    value={editedAsset.financialAdvisor || ""}
                    onChange={handleChange}
                  />
                </FormGroup>
              </>
            )}
          </ModalSections>

          <PrimaryButton
            onClick={() =>
              onSave &&
              onSave({
                ...editedAsset,
              })
            }
            buttonStyle={{ width: "100%" }}
          >
            Save Changes
          </PrimaryButton>
        </ModalContent>
      </Modal>
    </>
  );
};

export default EditAssetProtection;

const ModalContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const ModalSections = styled.div`
  background-color: #f2f6f9;
  gap: 16px;
  border-radius: 8px;
  padding: 16px;
`;

const SubHeading = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
`;

const Title = styled.h2`
  font-weight: 500;
  font-size: 16px;
  line-height: 100%;
  letter-spacing: 0%;
  color: #00225a;
  padding: 6px 0 15px 0;
`;
const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  color: #00225a;
  font-weight: 500;
  font-size: 16px;
  line-height: 24px;
`;

const Input = styled.input`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 400;
  width: 100%;
  outline: none;
  color: #00225a;
  font-family: inherit;
`;

const CheckWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Check = styled.input`
  width: 16px;
  height: 16px;
`;

const CheckLabel = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 400;
  font-size: 16px;
  line-height: 25px;
  letter-spacing: 0px;
  vertical-align: middle;
  color: #00225a;
`;
