import React from "react";
import styled from "styled-components";

type ButtonProps = {
  children: React.ReactNode;
  onClick?: () => void;
  buttonStyle?: React.CSSProperties;
  disabled?: boolean;
};

const SecondaryButton: React.FC<ButtonProps> = ({
  children,
  disabled,
  onClick,
  buttonStyle,
}) => {
  return (
    <Button onClick={onClick} disabled={disabled} style={buttonStyle}>
      {children}
    </Button>
  );
};

export default SecondaryButton;

const Button = styled.button`
  background: transparent;
  height: 48px;
  width: 240px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  border: 1px solid #007cdf;
  color: #007cdf;
  cursor: pointer;
  margin: 20px auto 0 auto;
  display: flex;
  align-items: center;
  font-family: inherit;
  justify-content: center;
`;
