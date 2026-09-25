"use client";
import { DropdownSelect, Modal, showErrorToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
} from "@/redux/api/assettokenization";
import React, { useEffect, useState } from "react";
import styled from "styled-components";
import { DatePicker } from "antd";
import dayjs, { Dayjs } from "dayjs";

interface EditAssetTokenProps {
  isOpen: boolean;
  asset?: TokenizationRecord;
  onClose: () => void;
  onSave?: (updatedData: Partial<UpdateTokenizationPayload>) => void;
}

const EditAssetTokenSalesModal: React.FC<EditAssetTokenProps> = ({
  isOpen,
  asset,
  onClose,
  onSave,
}) => {
  const [editedAsset, setEditedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >({});
  const [quoteCurrency, setQuoteCurrency] = useState("");
  const [proceedCycle, setProceedCycle] = useState("");
  const [payoutCurrency, setPayoutCurrency] = useState("");
  const [salesStartDate, setSalesStartDate] = useState<Date | null>(null);
  const [salesEndDate, setSalesEndDate] = useState<Date | null>(null);

  useEffect(() => {
    if (asset) {
      setEditedAsset(asset);
      setQuoteCurrency(asset.assetQuoteCurrency || "");
      setProceedCycle(asset.proceedCycle || "");
      setPayoutCurrency(asset.proceedPayoutCurrency || "");
      setSalesStartDate(asset.salesStart ? new Date(asset.salesStart) : null);
      setSalesEndDate(asset.salesEnd ? new Date(asset.salesEnd) : null);
    }
  }, [asset]);

  const formatLocalDateToISOString = (date: Date): string => {
    return date.toISOString(); // This gives "2025-05-29T00:00:00.000Z"
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    let newValue: string | number = value;
    if (
      name === "numberOfTokenToBeIssued" ||
      name === "capDurationInDays" ||
      name === "capAmountInFiat" ||
      name === "numberOfTokenToBeSold" ||
      name === "feeInFiat"
    ) {
      newValue = Number(value);
    }
    setEditedAsset((prev) => ({ ...prev, [name]: newValue }));
  };

  return (
    <Modal
      title="Edit Asset Tokens and Sales"
      isOpen={isOpen}
      onClose={onClose}
      closeIconPosition="right"
    >
      <ModalContent>
        <FormGroup>
          <Label>Tokens to be Issued</Label>
          <Input
            name="numberOfTokenToBeIssued"
            value={editedAsset.numberOfTokenToBeIssued || ""}
            onChange={handleChange}
          />
        </FormGroup>

        <FormGroup>
          <Label>Asset Quote Currency</Label>
          <DropdownSelect
            options={["USD", "NGN", "CNGN"]}
            placeholder="Select Quote Currency"
            value={quoteCurrency}
            onSelect={(item) => {
              setQuoteCurrency(item);
              setEditedAsset((prev) => ({ ...prev, assetQuoteCurrency: item }));
            }}
            borderless
            iconColor="#00225A"
          />
        </FormGroup>

        <FormGroup>
          <Label>Preferred Tokenization Fee (Fiat)</Label>
          <Input
            name="feeInFiat"
            value={editedAsset.feeInFiat || ""}
            onChange={handleChange}
          />
        </FormGroup>

        <FormGroup>
          <Label>Sales Start Date</Label>
          <StyledDatePicker
            value={salesStartDate ? dayjs(salesStartDate) : null}
            onChange={(date) =>
              setSalesStartDate((date as Dayjs | null)?.toDate() || null)
            }
            format="YYYY-MM-DD"
            placeholder="Select Start Date"
            size="large"
          />
        </FormGroup>

        <FormGroup>
          <Label>Sales End Date</Label>
          <StyledDatePicker
            value={salesEndDate ? dayjs(salesEndDate) : null}
            onChange={(date) =>
              setSalesEndDate((date as Dayjs | null)?.toDate() || null)
            }
            format="YYYY-MM-DD"
            placeholder="Select End Date"
            size="large"
          />
        </FormGroup>

        <FormGroup>
          <Label>Cap Duration</Label>
          <Input
            name="capDurationInDays"
            value={editedAsset.capDurationInDays || ""}
            onChange={handleChange}
          />
        </FormGroup>

        <FormGroup>
          <Label>Cap Amount</Label>
          <Input
            name="capAmountInFiat"
            value={editedAsset.capAmountInFiat || ""}
            onChange={handleChange}
          />
        </FormGroup>

        <FormGroup>
          <Label>Proceed Payout Cycle</Label>
          <DropdownSelect
            options={["Weekly", "Monthly", "Yearly"]}
            placeholder="Select Proceed Cycle"
            value={proceedCycle}
            onSelect={(item) => {
              setProceedCycle(item);
              setEditedAsset((prev) => ({ ...prev, proceedCycle: item }));
            }}
            borderless
            iconColor="#00225A"
          />
        </FormGroup>

        <FormGroup>
          <Label>Payout Currency</Label>
          <DropdownSelect
            options={["NGN", "USD"]}
            placeholder="Select Payout Currency"
            value={payoutCurrency}
            onSelect={(item) => {
              setPayoutCurrency(item);
              setEditedAsset((prev) => ({
                ...prev,
                proceedPayoutCurrency: item,
              }));
            }}
            borderless
            iconColor="#00225A"
          />
        </FormGroup>

        <PrimaryButton
          onClick={() => {
            if (onSave) {
              onSave({
                ...editedAsset,
                assetQuoteCurrency: quoteCurrency,
                proceedCycle,
                proceedPayoutCurrency: payoutCurrency,
                ...(salesStartDate && {
                  salesStart: formatLocalDateToISOString(salesStartDate),
                }),
                ...(salesEndDate && {
                  salesEnd: formatLocalDateToISOString(salesEndDate),
                }),
              });
            }
          }}
          buttonStyle={{ width: "100%" }}
        >
          Save Changes
        </PrimaryButton>
      </ModalContent>
    </Modal>
  );
};

export default EditAssetTokenSalesModal;

const ModalContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
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
`;

const Input = styled.input`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  width: 100%;
  outline: none;
  color: #00225a;
`;

const StyledDatePicker = styled(DatePicker)`
  width: 100%;
  .ant-picker-input > input {
    padding: 10px;
    font-size: 14px;
    color: #00225a;
  }
  &.ant-picker {
    border: 1px solid #bdbdbd;
    border-radius: 12px;
    &:hover,
    &.ant-picker-focused {
      border: 1.5px solid #00225a;
    }
  }
`;
