"use client";

import {
  DropdownSelect,
  Modal,
  showErrorToast,
  showSuccessToast,
} from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import {
  TokenizationRecord,
  UpdateTokenizationPayload,
  useGetTokenizationParamsQuery,
  useUpdateLogoMutation,
} from "@/redux/api/assettokenization";
import Image from "next/image";
import React, { useEffect, useRef, useState } from "react";
import { FaX } from "react-icons/fa6";
import { GoTrash } from "react-icons/go";
import { MdOutlineModeEditOutline } from "react-icons/md";
import styled, { keyframes } from "styled-components";

interface EditInfoModalProps {
  isOpen: boolean;
  asset?: TokenizationRecord;
  onClose: () => void;
  onSave?: (updatedData: Partial<UpdateTokenizationPayload>) => void;
  isLoading?: boolean;
}

const EditInfoModal: React.FC<EditInfoModalProps> = ({
  isOpen,
  asset,
  onClose,
  onSave,
  isLoading,
}) => {
  // Single state object for all asset fields
  const [editedAsset, setEditedAsset] = useState<
    Partial<UpdateTokenizationPayload>
  >(asset || {});

  const { data: tokenizationParams } = useGetTokenizationParamsQuery();

  const getAssetTypeName = (typeId: string | number | undefined) => {
    if (!typeId || !tokenizationParams?.assetTypes) return "";
    const match = tokenizationParams.assetTypes.find(
      (item) => item.id === Number(typeId),
    );
    return match?.assetType || "";
  };

  const handleTypeSelect = (typeName: string) => {
    const match = tokenizationParams?.assetTypes?.find(
      (item) => item.assetType === typeName,
    );
    if (match) {
      setEditedAsset((prev) => ({ ...prev, assetType: String(match.id) }));
    }
  };

  useEffect(() => {
    if (asset) {
      // Partial makes every field optional so you only update what you need.

      //  our payload includes many properties but only a few are being updated, using Partial allows you to construct an object with only those properties without TypeScript complaining about missing fields..
      setEditedAsset(asset as Partial<UpdateTokenizationPayload>);
    }
  }, [asset]);

  const [updateLogo] = useUpdateLogoMutation();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleLogoReplace = () => {
    fileInputRef.current?.click();
  };

  const handleUpdateLogo = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Check file type
    const validTypes = ["image/jpeg", "image/jpg", "image/png", "image/gif"];
    if (!validTypes.includes(file.type)) {
      showErrorToast("Only jpg, jpeg, png, and gif formats are supported");
      return;
    }

    // Check file size (900KB = 900000 bytes)
    if (file.size > 900000) {
      showErrorToast("Image size must be less than 900KB");
      return;
    }

    if (asset?.id) {
      try {
        await updateLogo({
          tokenizedAssetID: asset.id,
          documentFile: file,
        }).unwrap();

        showSuccessToast("Logo Updated Successfully");
      } catch (error: any) {
        console.error("Upload error:", error);

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
    }
  };

  // Generic change handler using the input name attribute
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;
    setEditedAsset((prev) => ({ ...prev, [name]: value }));
  };

  return (
    <>
      <Input
        type="file"
        ref={fileInputRef}
        onChange={handleUpdateLogo}
        style={{ display: "none" }}
        accept="image/*"
      />
      <Modal
        title="Edit Asset Information"
        isOpen={isOpen}
        onClose={onClose}
        closeIconPosition="right"
      >
        <Card>
          <FormGroup>
            <Label>Asset Logo</Label>
            <LogoUpload>
              <LogoPlaceholder>
                {asset?.assetLogo ? (
                  <Image
                    src={asset.assetLogo}
                    alt="Asset Logo"
                    width={40}
                    height={40}
                    style={{ borderRadius: "50%" }}
                    unoptimized={true}
                    onError={(e) => {
                      e.currentTarget.onerror = null;
                      e.currentTarget.src =
                        "https://media.istockphoto.com/id/1300845620/vector/user-icon-flat-isolated-on-white-background-user-symbol-vector-illustration.jpg?s=612x612&w=0&k=20&c=yBeyba0hUkh14_jgv1OKqIH0CCSWU_4ckRkAoy2p73o=";
                    }}
                  />
                ) : (
                  <Avatar />
                )}
              </LogoPlaceholder>
              <div>
                <Text>{`${asset?.assetCode}.png`}</Text>

                <ButtonGroup>
                  <Button onClick={handleLogoReplace}>
                    <MdOutlineModeEditOutline />
                    Replace
                  </Button>
                  <RemoveButton>
                    <GoTrash />
                    Remove
                  </RemoveButton>
                </ButtonGroup>
              </div>
            </LogoUpload>
          </FormGroup>

          <FormGroup>
            <Label>Asset Name</Label>
            <Input
              name="assetName"
              value={editedAsset.assetName || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Description</Label>
            <TextArea
              name="assetDescription"
              value={editedAsset.assetDescription || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Asset Code</Label>
            <Input
              name="assetCode"
              value={editedAsset.assetCode || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <DropdownSelect
              options={
                tokenizationParams?.assetSectors?.map((s) => s.sector) || []
              }
              placeholder="Select Sector"
              labelText="Sector"
              value={editedAsset.assetSector || ""}
              onSelect={(item) =>
                setEditedAsset((prev) => ({ ...prev, assetSector: item }))
              }
            />
          </FormGroup>

          <FormGroup>
            <DropdownSelect
              options={
                tokenizationParams?.assetSubSectors?.map((s) => s.subSector) ||
                []
              }
              placeholder="Select SubSector"
              labelText="SubSector"
              value={editedAsset.assetSubSector || ""}
              onSelect={(item) =>
                setEditedAsset((prev) => ({ ...prev, assetSubSector: item }))
              }
            />
          </FormGroup>

          <FormGroup>
            <DropdownSelect
              options={
                tokenizationParams?.assetTypes?.map((t) => t.assetType) || []
              }
              placeholder="Select Asset Type"
              labelText="Asset Type"
              value={getAssetTypeName(editedAsset.assetType)}
              onSelect={handleTypeSelect}
            />
          </FormGroup>

          <FormGroup>
            <Label>Asset Status</Label>
            <RadioWrapper>
              <RadioLabel>
                <Radio
                  type="radio"
                  name="assetStatus"
                  checked={editedAsset.assetAlreadyExists === 1}
                  onChange={() =>
                    setEditedAsset((prev) => ({
                      ...prev,
                      assetAlreadyExists: 1,
                    }))
                  }
                />
                Existing
              </RadioLabel>
              <RadioLabel>
                <Radio
                  type="radio"
                  name="assetStatus"
                  checked={editedAsset.assetAlreadyExists === 0}
                  onChange={() =>
                    setEditedAsset((prev) => ({
                      ...prev,
                      assetAlreadyExists: 0,
                    }))
                  }
                />
                Non-Existing
              </RadioLabel>
            </RadioWrapper>
          </FormGroup>
          <FormGroup>
            <Label>Offering Type</Label>
            <RadioWrapper>
              <RadioLabel>
                <Radio
                  type="radio"
                  name="offeringType"
                  value="public"
                  checked={editedAsset.offeringType === "public"}
                  onChange={(e) =>
                    setEditedAsset((prev) => ({
                      ...prev,
                      offeringType: e.target.value,
                    }))
                  }
                />
                Public
              </RadioLabel>
              <RadioLabel>
                <Radio
                  type="radio"
                  name="offeringType"
                  value="private"
                  checked={editedAsset.offeringType === "private"}
                  onChange={(e) =>
                    setEditedAsset((prev) => ({
                      ...prev,
                      offeringType: e.target.value,
                    }))
                  }
                />
                Private
              </RadioLabel>
            </RadioWrapper>
          </FormGroup>

          <FormGroup>
            <Label>Country</Label>
            <Input
              name="assetCountryLocation"
              value={editedAsset.assetCountryLocation || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Physical Asset Address</Label>
            <Input
              name="assetPhysicalAddress"
              value={editedAsset.assetPhysicalAddress || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Longitude</Label>
            <Input
              name="assetLongitude"
              value={editedAsset.assetLongitude || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <FormGroup>
            <Label>Latitude</Label>
            <Input
              name="assetLatitude"
              value={editedAsset.assetLatitude || ""}
              onChange={handleChange}
            />
          </FormGroup>

          <PrimaryButton
            onClick={() => {
              if (onSave) {
                onSave(editedAsset as UpdateTokenizationPayload); // ✅ FULL clean object
              }
            }}
            buttonStyle={{ width: "100%" }}

            // onClick={() => onSave && onSave(editedAsset)}
            // buttonStyle={{ width: "100%" }}
          >
            {isLoading ? "Saving..." : "      Save Changes"}
          </PrimaryButton>
        </Card>
      </Modal>
    </>
  );
};

export default EditInfoModal;

const Card = styled.section`
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: auto;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 14px;
  color: #00225a;
  padding-bottom: 6px;
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

const TextArea = styled.textarea`
  padding: 10px;
  border: 1px solid #bdbdbd;
  border-radius: 12px;
  font-size: 14px;
  width: 100%;
  outline: none;
  font-weight: 400;
  color: #00225a;
  height: 80px;
  resize: none;
  font-family: inherit;
`;

const RadioWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const RadioLabel = styled.label`
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
`;

const Radio = styled.input`
  width: 16px;
  height: 16px;
`;

const LogoUpload = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: 12px;
  background-color: #f2f6f9;
  border: 1px solid #007cdf;
  padding: 14px 16px;
`;

const LogoPlaceholder = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  overflow: hidden;
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 8px;
`;

const Button = styled.button`
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #0066ff;
  background: none;
  cursor: pointer;
  font-family: inherit;
  color: #007cdf;
  display: flex;
  align-items: center;
  gap: 2px;
  &:hover {
    background: #f0f6ff;
  }
`;

const RemoveButton = styled(Button)`
  border-color: #ff4d4f;
  color: #ff4d4f;
  &:hover {
    background: #fff0f0;
  }
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;
