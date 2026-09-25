import React from "react";
import styled from "styled-components";

type ButtonProps = {
  children: React.ReactNode;
  onClick?: () => void;
  buttonStyle?: React.CSSProperties;
  disabled?: boolean;
};

const PrimaryButton: React.FC<ButtonProps> = ({
  children,
  onClick,
  buttonStyle,
  disabled,
}) => {
  return (
    <Button onClick={onClick} style={buttonStyle} disabled={disabled}>
      {children}
    </Button>
  );
};

export default PrimaryButton;

const Button = styled.button`
  background-color: #007cdf;
  height: 48px;
  width: 240px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  color: #ffffff;
  cursor: pointer;
  margin: 20px auto 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: inherit;
  transition: opacity 0.15s;

  &:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }
`;
