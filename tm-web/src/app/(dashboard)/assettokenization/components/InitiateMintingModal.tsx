"use client";

import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  useUpdateTokenizationInfoMutation,
  useUpdateTokenizationSalesDatesMutation,
} from "@/redux/api/assettokenization";
import React, { useEffect, useRef, useState } from "react";
import { DatePicker } from "antd";
import { FaPlus, FaX } from "react-icons/fa6";
import styled from "styled-components";

import dayjs, { Dayjs } from "dayjs";
import "react-datepicker/dist/react-datepicker.css";
import { Modal, showErrorToast, showSuccessToast } from "@/components";

interface InitiateMintingModalProps {
  isInitiateMinting: boolean;
  setIsInitiateMinting: React.Dispatch<React.SetStateAction<boolean>>;
  onSubmit: () => void;
  asset?: TokenizationRecord;
  refetchAsset: () => void;
}

const InitiateMintingModal: React.FC<InitiateMintingModalProps> = ({
  isInitiateMinting,
  setIsInitiateMinting,
  onSubmit,
  asset,
  refetchAsset,
}) => {
  const [isAddingInitiator, setIsAddingInitiator] = useState<boolean>(false);
  const [isAddingApprover, setIsAddingApprover] = useState<boolean>(false);
  const [newInitiator, setNewInitiator] = useState<string>("");
  const [newApprover, setNewApprover] = useState<string>("");
  const [initiators, setInitiators] = useState<string[]>([]);
  const [approvers, setApprovers] = useState<string[]>([]);

  const [startDate, setStartDate] = useState<Date | null>(null);
  const [endDate, setEndDate] = useState<Date | null>(null);

  const [updateTokenizationInfo, { isLoading, isError }] =
    useUpdateTokenizationInfoMutation();

  useEffect(() => {
    if (isInitiateMinting) {
      setStartDate(asset?.salesStart ? new Date(asset.salesStart) : null);
      setEndDate(asset?.salesEnd ? new Date(asset.salesEnd) : null);
    }
  }, [asset?.salesStart, asset?.salesEnd, isInitiateMinting]);

  const [updateSalesDates, { isLoading: isUpdatingDates }] =
    useUpdateTokenizationSalesDatesMutation();

  const handleSubmit = () => {
    onSubmit();
    setIsInitiateMinting(false);
  };

  // Format Date as "YYYY-MM-DD"
  const formatLocalDateYYYYMMDD = (date: Date): string => {
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, "0");
    const day = date.getDate().toString().padStart(2, "0");
    return `${year}-${month}-${day}`;
  };

  // function to update sales start and end date
  const handleUpdateSalesDate = async () => {
    try {
      if (asset?.id && startDate && endDate) {
        if (endDate <= startDate) {
          showErrorToast("Sales end date must be after start date.");
          return;
        }
        const payload: any = {
          tokenizedAssetID: asset.id,
          ...(startDate && {
            salesStart: formatLocalDateYYYYMMDD(startDate),
          }),
          ...(endDate && {
            salesEnd: formatLocalDateYYYYMMDD(endDate),
          }),
        };
        await updateSalesDates(payload);
      }
    } catch (error) {
      console.error("Failed to update sales dates", error);
    }
  };

  useEffect(() => {
    if (asset?.mintingInitators) {
      const parsed = asset.mintingInitators
        .replace(/\n/g, "")
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean);
      setInitiators(parsed);
    }
  }, [asset?.mintingInitators]);

  useEffect(() => {
    if (asset?.mintingApprovers) {
      const parsed = asset.mintingApprovers
        .replace(/\n/g, "")
        .split(",")
        .map((item) => item.trim())
        .filter(Boolean);
      setApprovers(parsed);
    }
  }, [asset?.mintingApprovers]);

  // update Handle MintingInitiators
  const handleAddInitiators = async () => {
    if (!newInitiator.trim()) return;
    const updatedInitiators = [...initiators, newInitiator.trim()];

    try {
      await updateTokenizationInfo({
        tokenizedAssetID: asset?.id as string,
        data: {
          ...asset,
          mintingInitators: updatedInitiators.join(","),
        },
      }).unwrap();

      setNewInitiator("");

      showSuccessToast(`${newInitiator} successfully added as an initiator`);
      refetchAsset();
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
  const handleAddApprovers = async () => {
    if (!newApprover.trim()) return;
    const updatedApprovers = [...approvers, newApprover.trim()];
    try {
      await updateTokenizationInfo({
        tokenizedAssetID: asset?.id as string,
        data: {
          ...asset,
          mintingApprovers: updatedApprovers.join(","),
        },
      }).unwrap();

      setNewApprover("");
      showSuccessToast(`${newApprover} successfully added as an approver`);
      refetchAsset();
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

  const handleRemoveInitiator = async (index: number) => {
    const name = initiators[index];
    const updatedInitiators = initiators.filter((_, i) => i !== index);

    try {
      await updateTokenizationInfo({
        tokenizedAssetID: asset?.id as string,
        data: {
          ...asset,
          mintingInitators: updatedInitiators.join(","),
        },
      }).unwrap();

      setInitiators(updatedInitiators);
      showSuccessToast(
        `${name} successfully removed as an initiator for this asset`
      );
    } catch (error) {
      showErrorToast("Failed to remove initiator");
    }
  };

  const handleRemoveApprover = async (index: number) => {
    const name = approvers[index];
    const updatedApprovers = approvers.filter((_, i) => i !== index);

    try {
      await updateTokenizationInfo({
        tokenizedAssetID: asset?.id as string,
        data: {
          ...asset,
          mintingApprovers: updatedApprovers.join(","),
        },
      }).unwrap();

      setApprovers(updatedApprovers);
      showSuccessToast(
        `${name} successfully removed as an approver for this asset`
      );
    } catch (error) {
      showErrorToast("Failed to remove approver");
    }
  };

  return (
    <>
      <Modal
        title="Confirm Tokenization Information"
        isOpen={isInitiateMinting}
        onClose={() => setIsInitiateMinting(!isInitiateMinting)}
        closeIconPosition="right"
      >
        <ModalContent>
          <Content>
            <div>
              <Heading>Minting Initiators</Heading>
              <Description>
                Remove or add Minting Initiators for this asset on the existing
                list
              </Description>
            </div>
            <AddButton onClick={() => setIsAddingInitiator(!isAddingInitiator)}>
              <FaPlus color="#007CDF" />
              <span>Add Initiator</span>
            </AddButton>
          </Content>
          <Divider />
          {isAddingInitiator && (
            <InputContainer>
              <TextInput
                placeholder="Enter user name or email..."
                value={newInitiator}
                onChange={(e) => setNewInitiator(e.target.value)}
              />
              <ButtonGroup>
                <ActionButton onClick={handleAddInitiators}>Add</ActionButton>
                <CancelButton onClick={() => setIsAddingInitiator(false)}>
                  <FaX color="#00225A" />
                </CancelButton>
              </ButtonGroup>
            </InputContainer>
          )}
          <ListContainer $scrollable={initiators.length > 3}>
            {initiators.map((initiator, index) => (
              <ListItem key={index}>
                <ItemDetails>
                  <Avatar />
                  <div>
                    <ItemName>{initiator}</ItemName>
                    {/* <ItemSubText>obi@gmail.com</ItemSubText> */}
                  </div>
                </ItemDetails>
                <RemoveButton onClick={() => handleRemoveInitiator(index)}>
                  Remove
                </RemoveButton>
              </ListItem>
            ))}
          </ListContainer>
        </ModalContent>

        <ModalContent>
          <Content>
            <div>
              <Heading>Minting Approvers</Heading>
              <Description>
                Remove or add Minting Approvers for this asset on the existing
                list
              </Description>
            </div>
            <AddButton onClick={() => setIsAddingApprover(!isAddingApprover)}>
              <FaPlus color="#007CDF" />
              <span>Add Approvers</span>
            </AddButton>
          </Content>
          <Divider />
          {isAddingApprover && (
            <InputContainer>
              <TextInput
                placeholder="Enter user name or email..."
                value={newApprover}
                onChange={(e) => setNewApprover(e.target.value)}
              />
              <ButtonGroup>
                <ActionButton onClick={handleAddApprovers}>Add</ActionButton>
                <CancelButton onClick={() => setIsAddingApprover(false)}>
                  <FaX color="#00225A" />
                </CancelButton>
              </ButtonGroup>
            </InputContainer>
          )}
          <ListContainer $scrollable={approvers.length > 3}>
            {approvers.map((approver, index) => (
              <ListItem key={index}>
                <ItemDetails>
                  <Avatar />
                  <div>
                    <ItemName>{approver}</ItemName>
                  </div>
                </ItemDetails>
                <RemoveButton onClick={() => handleRemoveApprover(index)}>
                  Remove
                </RemoveButton>
              </ListItem>
            ))}
          </ListContainer>
        </ModalContent>

        <ModalContent>
          <Heading>Primary Sales Dates</Heading>
          <Divider />

          <SalesDatesContainer>
            <SalesDateItem>
              <SalesDateLabel>Sales Start Date</SalesDateLabel>
              <StyledDatePicker
                value={startDate ? dayjs(startDate) : null}
                onChange={(date: unknown) => {
                  const typedDate = date as Dayjs | null;
                  setStartDate(typedDate ? typedDate.toDate() : null);
                }}
                format="YYYY-MM-DD"
                style={{ width: "100%" }}
                placeholder={
                  asset?.salesStart
                    ? dayjs(asset.salesStart).format("YYYY-MM-DD")
                    : "Select Start Date"
                }
                size="middle"
              />
            </SalesDateItem>

            <SalesDateItem>
              <SalesDateLabel>Sales End Date</SalesDateLabel>

              <StyledDatePicker
                value={endDate ? dayjs(endDate) : null}
                onChange={(date: unknown) => {
                  const typedDate = date as Dayjs | null;
                  setEndDate(typedDate ? typedDate.toDate() : null);
                }}
                format="YYYY-MM-DD"
                style={{ width: "100%" }}
                placeholder="Select Start Date"
                size="middle"
              />
            </SalesDateItem>
          </SalesDatesContainer>
          <SalesDatesButton
            onClick={handleUpdateSalesDate}
            disabled={isUpdatingDates}
          >
            {isUpdatingDates ? "Updating..." : "Update Sales Dates"}
          </SalesDatesButton>
        </ModalContent>

        <PrimaryButton onClick={handleSubmit}>Proceed to mint</PrimaryButton>
      </Modal>
    </>
  );
};

export default InitiateMintingModal;

const ModalContent = styled.div`
  gap: 16px;
  border-radius: 12px;
  padding: 16px;
  background-color: #f2f6f9;
  margin-top: 8px;
`;

const Content = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 6px 0;
`;

const Heading = styled.h3`
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #191919;
  margin-bottom: 2px;
`;

const Description = styled.p`
  font-weight: 400;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #828282;
`;

const AddButton = styled.p`
  font-weight: 600;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #007cdf;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
`;

const Divider = styled.div`
  height: 1px;
  width: 100%;
  background-color: #e0e0e0;
  margin: 10px auto;
`;

const InputContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 10px 0;
`;

const TextInput = styled.input`
  font-family: inherit;
  border-radius: 10px;
  padding: 10px 4px;
  width: 80%;
  border: 1px solid #007cdf;
`;

const ButtonGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const ActionButton = styled.button`
  width: 50px;
  height: 36px;
  border-radius: 8px;
  padding: 12px;
  background-color: #007cdf;
  border: none;
  font-family: inherit;
  cursor: pointer;
  color: #fff;
`;

const CancelButton = styled.button`
  width: 40px;
  height: 36px;
  border-radius: 8px;
  padding: 12px 8px;
  background-color: #e0e0e0;
  border: none;
  font-family: inherit;
  cursor: pointer;
`;

const SalesDatesContainer = styled.div``;

const SalesDateItem = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
`;

const SalesDateLabel = styled.p`
  font-weight: 400;
  font-size: 14px;
  line-height: 17.07px;
  color: #828282;
  padding-bottom: 4px;
`;

const ListItem = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 14px;
  margin-top: 4px;
`;

const ItemDetails = styled.div`
  display: flex;
  gap: 4px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const ItemName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const RemoveButton = styled.button`
  border-radius: 8px;
  padding: 4px 12px;
  border: 1px solid #acd1ef;
  background-color: transparent;
  font-family: inherit;
  cursor: pointer;
  color: #007cdf;
`;
const ListContainer = styled.div<{ $scrollable: boolean }>`
  ${({ $scrollable }) =>
    $scrollable &&
    `
      max-height: 150px;
      overflow-y: auto;
    `}
`;

const SalesDatesButton = styled.button`
  margin-top: 8px;
  padding: 8px 16px;
  border: none;
  background-color: #007cdf;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
  font-family: inherit;
  font-size: 14px;
`;

const StyledDatePicker = styled(DatePicker)`
  width: 120px;

  padding: 0 2px !important; /* AntD controls input's padding internally */
  .ant-picker-input > input {
    padding: 10px; /* Matches your Input component */
    font-size: 14px;
    color: #00225a;
    font-family: inherit;
    border-radius: 12px;
  }
  &.ant-picker,
  &.ant-picker-focused {
    border: 1px solid #bdbdbd;
    border-radius: 12px;
    box-shadow: none;
    outline: none;
    transition: border 0.2s;
  }
  &:hover,
  &.ant-picker-focused {
    border: 1.5px solid #00225a;
  }
`;
