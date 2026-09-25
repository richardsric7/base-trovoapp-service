import React, { ReactNode } from "react";
import styled, { keyframes } from "styled-components";
import { FaX } from "react-icons/fa6";
interface ModalProps {
  title: string;
  children: ReactNode;
  isOpen: boolean;
  onClose: () => void;
  actions?: ReactNode;
  closeIconPosition?: "left" | "right";
  style?: React.CSSProperties;
}

export const Modal: React.FC<ModalProps> = ({
  title,
  children,
  isOpen,
  onClose,
  actions,
  closeIconPosition = "right",
  style,
}) => {
  if (!isOpen) return null;

  return (
    <Overlay>
      <Content style={style}>
        <ModalHeader iconPosition={closeIconPosition}>
          {/* <FaX onClick={onClose} cursor="pointer" /> */}
          {closeIconPosition === "left" && <CloseIcon onClick={onClose} />}
          <Title>{title}</Title>
          {closeIconPosition === "right" && <CloseIcon onClick={onClose} />}

          {/* <Title>{title}</Title> */}
        </ModalHeader>
        <Body>{children}</Body>
        {actions && <Footer>{actions}</Footer>}
      </Content>
      {/* <Content>{children}</Content> */}
    </Overlay>
  );
};

const fadeIn = keyframes`
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
`;

const Overlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
`;

const Content = styled.div`
  background-color: #ffffff;
  border-radius: 24px;
  padding: 20px;
  width: 500px;
  animation: ${fadeIn} 300ms ease-out;
`;

const ModalHeader = styled.div<{ iconPosition: "left" | "right" }>`
  display: flex;
  align-items: center;
  justify-content: ${({ iconPosition }) =>
    iconPosition === "left" ? "flex-start" : "space-between"};
  gap: ${({ iconPosition }) => (iconPosition === "left" ? "10rem" : "10px")};

  margin: 6px 0 20px 0;
`;

const CloseIcon = styled(FaX)`
  cursor: pointer;
  font-size: 18px;
  color: #333;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  text-align: center;
`;

const Body = styled.div`
  margin: 20px 0;
  flex: 1;
  overflow-y: auto;
  max-height: 80vh;
  padding-right: 8px;

  scrollbar-width: thin; /* Firefox */
  scrollbar-color: #828282 #f2f6f9;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-thumb {
    background: #828282;
    border-radius: 8px;
  }

  &::-webkit-scrollbar-track {
    background: #f2f6f9;
  }
`;

const Footer = styled.div`
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
`;
