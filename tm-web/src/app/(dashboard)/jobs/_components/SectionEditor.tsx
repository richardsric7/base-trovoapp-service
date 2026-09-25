"use client";

import { SectionKey } from "@/redux/api/jobs";
import { Button, Input } from "antd";
import { useState } from "react";
import { FaPlus } from "react-icons/fa6";
import { TbTrash } from "react-icons/tb";
import styled from "styled-components";

const DEFAULT_SECTIONS: Record<SectionKey, string[]> = {
  "About The Role": [],
  "What We're Looking For": [],
  "What You'll Own": [],
  "What Success Looks Like": [],
};

const SectionEditor = ({
  title,
  items,
  onChange,
  placeholder,
}: {
  title: SectionKey;
  items: string[];
  onChange: (next: string[]) => void;
  placeholder?: string;
}) => {
  const [draft, setDraft] = useState("");

  function sanitizeBullet(text: string) {
    return text.replace(/\s+/g, " ").trim();
  }

  const addItem = () => {
    const cleaned = sanitizeBullet(draft);
    if (!cleaned) return;
    onChange([...items, cleaned]);
    setDraft("");
  };

  const removeItem = (index: number) => {
    onChange(items.filter((_, i) => i !== index));
  };

  return (
    <Block>
      <Label>{title}</Label>

      <AddRow>
        <StyledInput
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          placeholder={placeholder || "Enter a point"}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addItem();
            }
          }}
        />
        <AddButton onClick={addItem} icon={<FaPlus />}>
          Add
        </AddButton>
      </AddRow>

      {items.map((item, idx) => (
        <ItemContainer>
          <ItemRow key={idx}>
            <span>{item}</span>
          </ItemRow>
          <RemoveBtn danger type="text" onClick={() => removeItem(idx)}>
            <TbTrash size={30} color="#BE3800" />
          </RemoveBtn>
        </ItemContainer>
      ))}
    </Block>
  );
};

export default SectionEditor;
const Block = styled.div`
  margin-bottom: 20px;
`;

const Label = styled.label`
  display: block;
  font-weight: 500;
  margin-bottom: 8px;
  color: #00225a;
  font-size: 14px;
`;

const StyledInput = styled(Input)`
  height: 48px;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
`;

const AddRow = styled.div`
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
`;

const AddButton = styled(Button)`
  min-width: 80px;
  background-color: #007cdf;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  border-radius: 8px;
  height: 48px;
`;

const ItemRow = styled.div`
  color: #00225a;
  background: #f2f6f9;
  padding: 12px;
  border-radius: 10px;
  flex-grow: 1;
  // height: 48px;
`;

const ItemContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
`;
const RemoveBtn = styled(Button)`
  border: 1px solid #be3800;
  border-radius: 8px;
  height: 48px;
  width: 48px;
`;
