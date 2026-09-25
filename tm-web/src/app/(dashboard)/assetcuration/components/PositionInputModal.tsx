import React, { useState } from "react";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import { FaX } from "react-icons/fa6";

interface AssetRow {
  id: number;
  position: number;
  assetname: string;
  assettitle: string;
  createdOn: string;
  updated: string;
  curatorName: string;
  curatorRole: string;
}

interface PositionInputModalProps {
  isOpen: boolean;
  onClose: () => void;
  maxPosition: number;
  selectedRow: AssetRow;
  setDataSource: React.Dispatch<React.SetStateAction<AssetRow[]>>;
}

const PositionInputModal: React.FC<PositionInputModalProps> = ({
  isOpen,
  onClose,
  maxPosition,
  selectedRow,
  setDataSource,
}) => {
  const [position, setPosition] = useState<string>("");

  const handleMove = () => {
    const newPosition = parseInt(position, 10);
    if (!isNaN(newPosition) && newPosition >= 1 && newPosition <= maxPosition) {
      setDataSource((items) => {
        const oldIndex = items.findIndex((item) => item.id === selectedRow.id);
        const newIndex = newPosition - 1;

        const newItems = [...items];
        const [movedItem] = newItems.splice(oldIndex, 1);
        newItems.splice(newIndex, 0, movedItem);

        return newItems.map((item, index) => ({
          ...item,
          position: index + 1,
        }));
      });

      onClose();
    } else {
      alert(`Please enter a valid position between 1 and ${maxPosition}`);
    }
  };

  return isOpen ? (
    <ModalContent>
      <Header>
        <h3>Move Asset Token</h3>
        <FaX color="#00225A" size={12} onClick={onClose} cursor="pointer" />
      </Header>
      <div>
        <Label>Enter Position</Label>

        <StyledInput
          type="number"
          placeholder="1"
          value={position}
          onChange={(e) => setPosition(e.target.value)}
          min="1"
          max={maxPosition}
        />
      </div>
      <PrimaryButton onClick={handleMove}>Apply Position</PrimaryButton>
    </ModalContent>
  ) : null;
};

export default PositionInputModal;

const ModalContent = styled.div`
  position: absolute;
  bottom: -60px;
  right: 150px;
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
  min-width: 300px;
  box-shadow: 0px 4px 100px 0px #00000026;
`;

const StyledInput = styled.input`
  width: 100%;
  padding: 8px;
  margin: 12px 0;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-family: inherit;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;

  h3 {
    color: #00225a;
  }
`;
const Label = styled.h2`
  font-weight: 500;
  font-size: 16px;
  line-height: 24px;
  letter-spacing: 0%;
  color: #828282;
  text-align: left;
`;
