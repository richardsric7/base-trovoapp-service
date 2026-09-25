"use client";

import React, { useState } from "react";
import styled from "styled-components";

interface InputProps {
  placeholder: string;
  value: string;
  onChange: (val: string) => void;
  suggestions?: string[];
  onSelectSuggestion?: (val: string) => void;
  selectedItems?: string[];
  onRemoveItem?: (val: string) => void;
}

const StakeHoldersInput: React.FC<InputProps> = ({
  placeholder,
  value,
  onChange,
  suggestions = [],
  onSelectSuggestion,
  selectedItems = [],
  onRemoveItem,
}) => {
  const [showDropdown, setShowDropdown] = useState(false);

  return (
    <>
      <Input
        placeholder={placeholder}
        value={value}
        onChange={(e) => {
          onChange(e.target.value);
          setShowDropdown(true);
        }}
        onFocus={() => setShowDropdown(true)}
      />
      {showDropdown && suggestions.length > 0 && onSelectSuggestion && (
        <Dropdown>
          {suggestions.map((item, index) => (
            <DropdownItem
              key={index}
              onClick={() => {
                onSelectSuggestion(item);
                setShowDropdown(false);
              }}
            >
              <BoldText>{item.charAt(0).toUpperCase()}</BoldText>
              {item}
            </DropdownItem>
          ))}
        </Dropdown>
      )}

      {selectedItems.map((name) => (
        <SelectedItem key={name}>
          <NameContent>
            <BoldText>{name.charAt(0).toUpperCase()}</BoldText>
            <span>{name}</span>
          </NameContent>
          <RemoveBtn onClick={() => onRemoveItem?.(name)}>Remove</RemoveBtn>
        </SelectedItem>
      ))}
    </>
  );
};

export default StakeHoldersInput;

const Input = styled.input`
  background-color: transparent;
  border: 1px solid #828282;
  width: 100%;
  height: 48px;
  border-radius: 10px;
  font-weight: 400;
  font-size: 14px;
  line-height: 24px;
  padding: 0 4px;
  font-family: inherit;

  &::placeholder {
    color: #828282;

    font-family: inherit;
  }
`;

const Line = styled.div`
  background-color: #e0e0e0;
  width: 100%;
  height: 1px;
  margin-bottom: 8px;
`;

const Dropdown = styled.ul`
  background-color: #fff;
  border-radius: 6px;
  list-style: none;
  padding: 8px 0;
`;

const DropdownItem = styled.li`
  padding: 8px 12px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #00225a;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const SelectedItem = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;

  color: #191919;
  font-size: 13px;
  padding: 6px 10px;
  border-radius: 8px;

  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;

  margin-top: 8px;
`;

const RemoveBtn = styled.button`
  cursor: pointer;
  font-size: 14px;
  border: 1px solid #acd1ef;
  color: #007cdf;
  font-weight: 400;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  border-radius: 6px;
  padding: 4px;
  font-family: inherit;
  background-color: none;
`;

const BoldText = styled.div`
  background-color: #007cdf;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
`;

const NameContent = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;
