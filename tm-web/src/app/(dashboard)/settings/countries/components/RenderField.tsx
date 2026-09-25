import { DropdownSelect } from "@/components";
import React from "react";
import styled from "styled-components";

interface RenderFieldProps {
  field: any;
  value: string;
  handleInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleDropdownSelect?: (val: any) => void;
}

const RenderField = ({
  field,
  value,
  handleDropdownSelect,
  handleInputChange,
}: RenderFieldProps) => {
  if (field.type === "dropdown" && handleDropdownSelect)
    return (
      <>
        <DropdownSelect
          options={field.options}
          placeholder={field.placeholder}
          labelText=""
          value={value}
          onSelect={handleDropdownSelect}
        />
      </>
    );

  if (field.type === "fee")
    return (
      <FeeInputWrapper>
        <FeeInputField
          name={field.name}
          placeholder="0"
          value={value}
          onChange={handleInputChange}
        />
        <DividerGroup>
          <VerticalDivider />
          <FeeLabel>{field.suffix}</FeeLabel>
        </DividerGroup>
      </FeeInputWrapper>
    );

  if (field.type === "text")
    return (
      <TextInput
        name={field.name}
        placeholder={field.placeholder}
        value={value}
        onChange={handleInputChange}
      />
    );
};

export default RenderField;

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

const DividerGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
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
