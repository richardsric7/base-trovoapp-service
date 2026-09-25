"use client";
import { Modal, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import { TokenizationRecord } from "@/redux/api/assettokenization";
import { useState } from "react";
import styled from "styled-components";

import {
  Organization,
  useGetOrganizationsListQuery,
} from "@/redux/api/organizations";
import RenderStakeHoldersSelection from "./RenderStakeHoldersSelection";

interface AsignStakeHoldersProps {
  isOpen: boolean;
  onClose: () => void;
  onAssignSuccess: (ids: {
    assetManagerId?: number;
    approvedAssetCustodianId?: number;
    assetIssuingHouseId?: number;
    legalAndProfesionalPartnerId?: number;
    ratingAgencyId?: number;
    trusteeId?: number;
    legalAdviserId?: number;
    financialAdviserId?: number;
  }) => void;
  asset?: TokenizationRecord;
}

const AsignStakeHolders: React.FC<AsignStakeHoldersProps> = ({
  isOpen,
  onClose,
  onAssignSuccess,
  asset,
}) => {
  const { data, isLoading } = useGetOrganizationsListQuery({
    page: 1,
    pageSize: 100,
  });

  const organisations = data?.organizations ?? [];

  const assetManagersData = organisations.filter(
    (org) => org.type === "ASSET_MANAGER" || org.type === "asset_manager"
  );

  const assetCustodiansData = organisations.filter(
    (org) =>
      org.type === "ASSET_CUSTODIAN" ||
      org.type === "APPROVED_ASSET_CUSTODIAN" ||
      org.type === "approved_asset_custodian"
  );

  const issuingHousesData = organisations.filter(
    (org) =>
      org.type === "ASSET_ISSUING_HOUSE" || org.type === "asset_issuing_house"
  );

  const legalAgencyData = organisations.filter(
    (org) =>
      org.type === "LEGAL_AGENCY" ||
      org.type === "PROFESSIONAL_AGENCY" ||
      org.type === "legal_and_professionals"
  );

  const ratingAgencyData = organisations.filter(
    (org) => org.type === "RATING_AGENCY" || org.type === "rating_agency"
  );

  const trusteeData = organisations.filter((org) => org.type === "trustees");

  // Legal Adviser and Financial Adviser are distinct A5 roles (optional to assign).
  const legalAdviserData = organisations.filter(
    (org) => org.type === "LEGAL_ADVISER" || org.type === "legal_adviser"
  );
  const financialAdviserData = organisations.filter(
    (org) =>
      org.type === "FINANCIAL_ADVISER" || org.type === "financial_adviser"
  );

  const [managerInput, setManagerInput] = useState("");
  const [custodianInput, setCustodainInput] = useState("");
  const [issuingHouse, setIssuingHouse] = useState("");
  const [legalAdviserInput, setLegalAdviserInput] = useState("");
  const [ratingAgencyInput, setRatingAgencyInput] = useState("");
  const [trusteeInput, setTrusteeInput] = useState("");
  const [legalAdviserSearchInput, setLegalAdviserSearchInput] = useState("");
  const [financialAdviserInput, setFinancialAdviserInput] = useState("");

  const [selectedManager, setSelectedManager] = useState<Organization | null>(
    null
  );
  const [selectedCustodian, setSelectedCustodian] =
    useState<Organization | null>(null);
  const [selectedIssuingHouse, setSelectedIssuingHouse] =
    useState<Organization | null>(null);
  const [selectedLegalAgency, setSelectedLegalAgency] =
    useState<Organization | null>(null);
  const [selectedRatingAgency, setSelectedRatingAgency] =
    useState<Organization | null>(null);
  const [selectedTrustee, setSelectedTrustee] = useState<Organization | null>(
    null
  );
  const [selectedLegalAdviser, setSelectedLegalAdviser] =
    useState<Organization | null>(null);
  const [selectedFinancialAdviser, setSelectedFinancialAdviser] =
    useState<Organization | null>(null);

  const matchesSearch = (input: string, target: string) =>
    target?.toLowerCase().includes(input.toLowerCase());

  const managerSuggestions = assetManagersData.filter((m) =>
    matchesSearch(managerInput, m.name)
  );

  const custodianSuggestions = assetCustodiansData.filter((c) =>
    matchesSearch(custodianInput, c.name)
  );

  const issuingHouseSuggestions = issuingHousesData.filter((i) =>
    matchesSearch(issuingHouse, i.name)
  );

  const legalAgencySuggestions = legalAgencyData.filter((l) =>
    matchesSearch(legalAdviserInput, l.name)
  );

  const ratingAgencySuggestions = ratingAgencyData.filter((r) =>
    matchesSearch(ratingAgencyInput, r.name)
  );

  const trusteeSuggestions = trusteeData.filter((t) =>
    matchesSearch(trusteeInput, t.name)
  );

  const legalAdviserSuggestions = legalAdviserData.filter((l) =>
    matchesSearch(legalAdviserSearchInput, l.name)
  );

  const financialAdviserSuggestions = financialAdviserData.filter((f) =>
    matchesSearch(financialAdviserInput, f.name)
  );

  const handleSelectManager = (name: string) => {
    const found = assetManagersData.find((m) => m.name === name);
    if (found) setSelectedManager(found);
  };

  const handleSelectCustodian = (name: string) => {
    const found = assetCustodiansData.find((c) => c.name === name);
    if (found) setSelectedCustodian(found);
  };

  const handleSelectIssuingHouse = (name: string) => {
    const found = issuingHousesData.find((i) => i.name === name);
    if (found) setSelectedIssuingHouse(found);
  };

  const handleLegalAgency = (name: string) => {
    const found = legalAgencyData.find((l) => l.name === name);
    if (found) setSelectedLegalAgency(found);
  };

  const handleRatingAgency = (name: string) => {
    const found = ratingAgencyData.find((r) => r.name === name);
    if (found) setSelectedRatingAgency(found);
  };

  const handleTrustee = (name: string) => {
    const found = trusteeData.find((t) => t.name === name);
    if (found) setSelectedTrustee(found);
  };

  const handleLegalAdviser = (name: string) => {
    const found = legalAdviserData.find((l) => l.name === name);
    if (found) setSelectedLegalAdviser(found);
  };

  const handleFinancialAdviser = (name: string) => {
    const found = financialAdviserData.find((f) => f.name === name);
    if (found) setSelectedFinancialAdviser(found);
  };

  const handleRemoveManager = () => setSelectedManager(null);
  const handleRemoveCustodian = () => setSelectedCustodian(null);
  const handleRemoveIssuingHouse = () => setSelectedIssuingHouse(null);
  const handleRemoveLegalAgency = () => setSelectedLegalAgency(null);
  const handleRemoveRatingAgency = () => setSelectedRatingAgency(null);
  const handleRemoveTrustee = () => setSelectedTrustee(null);
  const handleRemoveLegalAdviser = () => setSelectedLegalAdviser(null);
  const handleRemoveFinancialAdviser = () => setSelectedFinancialAdviser(null);

  const handleAssign = () => {
    // Validate that required stakeholders are selected
    if (
      !selectedManager ||
      !selectedCustodian ||
      !selectedIssuingHouse ||
      !selectedRatingAgency ||
      !selectedTrustee ||
      !selectedLegalAgency
    ) {
      showErrorToast("Please select stakeholders before assigning.");
      return;
    }

    // Build the payload with proper type conversion.
    // Legal/Financial Adviser are optional A5 roles: only included when selected.
    const payload: {
      assetManagerId: number;
      approvedAssetCustodianId: number;
      assetIssuingHouseId: number;
      ratingAgencyId: number;
      legalAndProfesionalPartnerId: number;
      trusteeId: number;
      legalAdviserId?: number;
      financialAdviserId?: number;
    } = {
      assetManagerId: Number(selectedManager.stakeholder_id),
      approvedAssetCustodianId: Number(selectedCustodian.stakeholder_id),
      assetIssuingHouseId: Number(selectedIssuingHouse.stakeholder_id),
      ratingAgencyId: Number(selectedRatingAgency?.stakeholder_id),
      legalAndProfesionalPartnerId: Number(selectedLegalAgency?.stakeholder_id),
      trusteeId: Number(selectedTrustee?.stakeholder_id),
    };

    if (selectedLegalAdviser) {
      payload.legalAdviserId = Number(selectedLegalAdviser.stakeholder_id);
    }
    if (selectedFinancialAdviser) {
      payload.financialAdviserId = Number(
        selectedFinancialAdviser.stakeholder_id
      );
    }

    // Debug logging
    console.log(" Assignment Payload ");
    console.log("Payload:", payload);

    // Verify required IDs are valid numbers (advisers are optional).
    if (
      isNaN(payload.assetManagerId) ||
      isNaN(payload.approvedAssetCustodianId) ||
      isNaN(payload.assetIssuingHouseId) ||
      isNaN(payload.legalAndProfesionalPartnerId) ||
      isNaN(payload.ratingAgencyId) ||
      isNaN(payload.trusteeId) ||
      (payload.legalAdviserId !== undefined &&
        isNaN(payload.legalAdviserId)) ||
      (payload.financialAdviserId !== undefined &&
        isNaN(payload.financialAdviserId))
    ) {
      showErrorToast(
        "Invalid organization IDs. Please try selecting the stakeholders again."
      );
      return;
    }

    onAssignSuccess(payload);
    onClose();
  };

  return (
    <Modal
      title="Assign Tokenization Stakeholders"
      isOpen={isOpen}
      onClose={onClose}
      closeIconPosition="right"
    >
      <ModalContent>
        {isLoading ? (
          <LoadingText>Loading organizations...</LoadingText>
        ) : (
          <>
            <RenderStakeHoldersSelection
              title="Asset Custodian *"
              description="Search for custodian by name to select and assign them."
              placeholder="Search custodian"
              value={custodianInput}
              onChange={setCustodainInput}
              suggestions={custodianSuggestions.map((s) => s.name)}
              onSelect={handleSelectCustodian}
              selectedItems={selectedCustodian ? [selectedCustodian.name] : []}
              onRemove={handleRemoveCustodian}
              hasOrganizations={assetCustodiansData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Asset Manager *"
              description="Search for Manager by name to select and assign them."
              placeholder="Search manager"
              value={managerInput}
              onChange={setManagerInput}
              suggestions={managerSuggestions.map((s) => s.name)}
              onSelect={handleSelectManager}
              selectedItems={selectedManager ? [selectedManager.name] : []}
              onRemove={handleRemoveManager}
              hasOrganizations={assetManagersData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Issuing House *"
              description="Search for issuing house by name to select and assign them."
              placeholder="Search issuing house"
              value={issuingHouse}
              onChange={setIssuingHouse}
              suggestions={issuingHouseSuggestions.map((s) => s.name)}
              onSelect={handleSelectIssuingHouse}
              selectedItems={
                selectedIssuingHouse ? [selectedIssuingHouse.name] : []
              }
              onRemove={handleRemoveIssuingHouse}
              hasOrganizations={issuingHousesData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Legal and Professional Partner *"
              description="Search for legal partner by name to select and assign them."
              placeholder="Search legal partner"
              value={legalAdviserInput}
              onChange={setLegalAdviserInput}
              suggestions={legalAgencySuggestions.map((s) => s.name)}
              onSelect={handleLegalAgency}
              selectedItems={
                selectedLegalAgency ? [selectedLegalAgency.name] : []
              }
              onRemove={handleRemoveLegalAgency}
              hasOrganizations={legalAgencyData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Rating Agency *"
              description="Search for rating agency by name to select and assign them."
              placeholder="Search rating agency"
              value={ratingAgencyInput}
              onChange={setRatingAgencyInput}
              suggestions={ratingAgencySuggestions.map((s) => s.name)}
              onSelect={handleRatingAgency}
              selectedItems={
                selectedRatingAgency ? [selectedRatingAgency.name] : []
              }
              onRemove={handleRemoveRatingAgency}
              hasOrganizations={ratingAgencyData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Trustees *"
              description="Search for trustee by name to select and assign them."
              placeholder="Search trustee"
              value={trusteeInput}
              onChange={setTrusteeInput}
              suggestions={trusteeSuggestions.map((s) => s.name)}
              onSelect={handleTrustee}
              selectedItems={selectedTrustee ? [selectedTrustee.name] : []}
              onRemove={handleRemoveTrustee}
              hasOrganizations={trusteeData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Legal Adviser"
              description="Optional: search for a legal adviser by name to assign them."
              placeholder="Search legal adviser"
              value={legalAdviserSearchInput}
              onChange={setLegalAdviserSearchInput}
              suggestions={legalAdviserSuggestions.map((s) => s.name)}
              onSelect={handleLegalAdviser}
              selectedItems={
                selectedLegalAdviser ? [selectedLegalAdviser.name] : []
              }
              onRemove={handleRemoveLegalAdviser}
              hasOrganizations={legalAdviserData.length > 0}
            />
            <RenderStakeHoldersSelection
              title="Financial Adviser"
              description="Optional: search for a financial adviser by name to assign them."
              placeholder="Search financial adviser"
              value={financialAdviserInput}
              onChange={setFinancialAdviserInput}
              suggestions={financialAdviserSuggestions.map((s) => s.name)}
              onSelect={handleFinancialAdviser}
              selectedItems={
                selectedFinancialAdviser ? [selectedFinancialAdviser.name] : []
              }
              onRemove={handleRemoveFinancialAdviser}
              hasOrganizations={financialAdviserData.length > 0}
            />
            <PrimaryButton onClick={handleAssign}>
              Assign Stakeholders
            </PrimaryButton>{" "}
          </>
        )}
      </ModalContent>
    </Modal>
  );
};
export default AsignStakeHolders;

const LoadingText = styled.div`
  color: #666;
  font-size: 14px;
  padding: 8px;
`;

const ModalContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const Content = styled.div`
  display: flex;
  flex-direction: column;
  background-color: #f2f6f9;
  border-radius: 12px;
  padding: 16px;
  gap: 4px;
`;
const Title = styled.h2`
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #191919;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 11px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #828282;
`;

const Line = styled.div`
  background-color: #e0e0e0;
  width: 100%;
  height: 1px;
  margin-bottom: 8px;
`;
