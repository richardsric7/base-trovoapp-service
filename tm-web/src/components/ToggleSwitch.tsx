import React from "react";
import styled from "styled-components";

interface ToggleSwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label?: string;
  id?: string;
  disabled?: boolean;
}

const ToggleSwitch: React.FC<ToggleSwitchProps> = ({
  checked,
  onChange,
  label,
  id,
  disabled = false,
}) => {
  return (
    <Wrapper>
      {label && <Label htmlFor={id}>{label}</Label>}
      <SwitchLabel>
        <HiddenCheckbox
          id={id}
          checked={checked}
          disabled={disabled}
          onChange={() => !disabled && onChange(!checked)}
        />
        <Slider checked={checked} disabled={disabled} />
      </SwitchLabel>
    </Wrapper>
  );
};

export default ToggleSwitch;

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const Label = styled.label`
  font-weight: 500;
`;

const SwitchLabel = styled.label`
  position: relative;
  display: inline-block;
  width: 24px;
  height: 13.5px;
`;

const HiddenCheckbox = styled.input.attrs({ type: "checkbox" })`
  opacity: 0;
  width: 0;
  height: 0;
`;

const Slider = styled.span<{ checked: boolean; disabled?: boolean }>`
  position: absolute;
  cursor: ${({ disabled }) => (disabled ? "not-allowed" : "pointer")};
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: ${({ checked }) => (checked ? "#2196F3" : "#ccc")};
  opacity: ${({ disabled }) => (disabled ? 0.5 : 1)};
  transition: 0.4s;
  border-radius: 13.5px;

  &::before {
    position: absolute;
    content: "";
    height: 10px;
    width: 9px;
    left: 2px;
    bottom: 2px;
    background-color: white;
    transition: 0.4s;
    border-radius: 50%;
    transform: ${({ checked }) =>
      checked ? "translateX(10.5px)" : "translateX(0)"};
  }
`;
