"use client";
import React, { useState } from "react";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import ToggleSwitch from "@/components/ToggleSwitch";
import { DropdownSelect } from "@/components/Dropdown";
import { showErrorToast } from "@/components";
import { ICuratedAsset, useGetAssetClassesQuery } from "@/redux/api/curatedAssets";

export type CuratedAssetFormValues = {
  assetCode: string;
  assetName: string;
  contractAddress: string;
  assetClassId: number;
  decimalPlaces: number;
  priority: number;
  assetLimit: number;
  description: string;
  website: string;
  organization: string;
  contactEmail: string;
  withdrawable: boolean;
  generateDepositAddress: boolean;
  inactive: boolean;
  p2pEnabled: boolean;
};

const fromAsset = (asset?: ICuratedAsset): CuratedAssetFormValues => ({
  assetCode: asset?.assetCode ?? "",
  assetName: asset?.assetName ?? "",
  contractAddress: asset?.contractAddress ?? "",
  assetClassId: asset?.assetClassId ?? 0,
  decimalPlaces: asset?.decimalPlaces ?? 7,
  priority: asset?.priority ?? 0,
  assetLimit: asset?.assetLimit ?? 0,
  description: asset?.description ?? "",
  website: asset?.website ?? "",
  organization: asset?.organization ?? "",
  contactEmail: asset?.contactEmail ?? "",
  withdrawable: !!asset?.withdrawable,
  generateDepositAddress: !!asset?.generateDepositAddress,
  inactive: !!asset?.inactive,
  p2pEnabled: !!asset?.p2pEnabled,
});

const CuratedAssetForm = ({
  initialAsset,
  submitLabel,
  submitting,
  onSubmit,
}: {
  initialAsset?: ICuratedAsset;
  submitLabel: string;
  submitting: boolean;
  onSubmit: (values: CuratedAssetFormValues) => void;
}) => {
  const [values, setValues] = useState<CuratedAssetFormValues>(fromAsset(initialAsset));
  const { data: assetClassesData } = useGetAssetClassesQuery();
  const assetClasses = assetClassesData?.data ?? [];
  const selectedClass = assetClasses.find((c) => c.id === values.assetClassId);

  const set = <K extends keyof CuratedAssetFormValues>(key: K, value: CuratedAssetFormValues[K]) =>
    setValues((prev) => ({ ...prev, [key]: value }));

  const handleTextChange = (key: keyof CuratedAssetFormValues) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    set(key, e.target.value as never);

  const handleNumberChange = (key: keyof CuratedAssetFormValues) => (e: React.ChangeEvent<HTMLInputElement>) =>
    set(key, (Number(e.target.value) || 0) as never);

  const handleSubmit = () => {
    if (!values.assetCode.trim()) {
      showErrorToast("Asset code is required");
      return;
    }
    onSubmit(values);
  };

  return (
    <FormGrid>
      <Field>
        <Label>Asset code *</Label>
        <Input
          value={values.assetCode}
          onChange={handleTextChange("assetCode")}
          disabled={!!initialAsset}
          placeholder="e.g. USDC"
        />
      </Field>
      <Field>
        <Label>Asset name</Label>
        <Input value={values.assetName} onChange={handleTextChange("assetName")} placeholder="e.g. USD Coin" />
      </Field>
      <Field>
        <Label>Contract address</Label>
        <Input value={values.contractAddress} onChange={handleTextChange("contractAddress")} />
      </Field>
      <Field>
        <Label>Asset class</Label>
        <DropdownSelect
          value={selectedClass?.assetClass ?? ""}
          options={assetClasses.map((c) => c.assetClass)}
          placeholder="Select a class"
          labelText=""
          onSelect={(_item, index) => set("assetClassId", assetClasses[index].id as never)}
          backgroundColor="#ffffff"
          borderColor="#e0e0e0"
        />
      </Field>
      <Field>
        <Label>Decimal places</Label>
        <Input type="number" value={values.decimalPlaces} onChange={handleNumberChange("decimalPlaces")} />
      </Field>
      <Field>
        <Label>Priority (display order)</Label>
        <Input type="number" value={values.priority} onChange={handleNumberChange("priority")} />
      </Field>
      <Field>
        <Label>Asset limit (0 = unlimited)</Label>
        <Input type="number" value={values.assetLimit} onChange={handleNumberChange("assetLimit")} />
      </Field>
      <Field>
        <Label>Organization</Label>
        <Input value={values.organization} onChange={handleTextChange("organization")} />
      </Field>
      <Field>
        <Label>Contact email</Label>
        <Input value={values.contactEmail} onChange={handleTextChange("contactEmail")} />
      </Field>
      <Field>
        <Label>Website</Label>
        <Input value={values.website} onChange={handleTextChange("website")} />
      </Field>
      <FullWidthField>
        <Label>Description</Label>
        <Textarea value={values.description} onChange={handleTextChange("description")} rows={3} />
      </FullWidthField>

      <TogglesRow>
        <ToggleSwitch label="Withdrawable" checked={values.withdrawable} onChange={(v) => set("withdrawable", v as never)} />
        <ToggleSwitch
          label="Generate deposit address"
          checked={values.generateDepositAddress}
          onChange={(v) => set("generateDepositAddress", v as never)}
        />
        <ToggleSwitch label="Inactive" checked={values.inactive} onChange={(v) => set("inactive", v as never)} />
      </TogglesRow>

      <P2PField>
        <ToggleSwitch
          id="p2pEnabled"
          label="Available on P2P marketplace"
          checked={values.p2pEnabled}
          onChange={(v) => set("p2pEnabled", v as never)}
        />
        <P2PHint>
          When off, merchants cannot create a P2P offer for this asset, and any existing offer for
          it is hidden from marketplace search.
        </P2PHint>
      </P2PField>

      <SubmitRow>
        <PrimaryButton onClick={handleSubmit} disabled={submitting} buttonStyle={{ width: "auto", padding: "10px 32px" }}>
          {submitting ? "Saving..." : submitLabel}
        </PrimaryButton>
      </SubmitRow>
    </FormGrid>
  );
};

export default CuratedAssetForm;

const FormGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  max-width: 800px;
`;

const Field = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const FullWidthField = styled(Field)`
  grid-column: 1 / -1;
`;

const Label = styled.label`
  font-size: 14px;
  color: #828282;
`;

const Input = styled.input`
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  outline: none;
  font-size: 14px;

  &:disabled {
    background-color: #f5f5f5;
    color: #828282;
  }
`;

const Textarea = styled.textarea`
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #e0e0e0;
  outline: none;
  font-size: 14px;
  resize: vertical;
`;

const TogglesRow = styled.div`
  grid-column: 1 / -1;
  display: flex;
  gap: 32px;
  padding: 12px 0;
  border-top: 1px solid #e5e5ef;
  border-bottom: 1px solid #e5e5ef;
`;

const P2PField = styled.div`
  grid-column: 1 / -1;
  background-color: #f2f6f9;
  border-radius: 12px;
  padding: 16px;
`;

const P2PHint = styled.p`
  font-size: 12px;
  color: #828282;
  margin: 8px 0 0 0;
`;

const SubmitRow = styled.div`
  grid-column: 1 / -1;
  padding-top: 8px;
`;
