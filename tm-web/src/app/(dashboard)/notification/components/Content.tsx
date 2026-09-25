"use client";

import { DropdownSelect } from "@/components";
import React, { useState } from "react";
import styled from "styled-components";
import TextEditor from "./TextEditor";

interface SelectValues {
  NoticeType?: string;
}

const Content = () => {
  const [selectValues, setSelectValues] = useState<SelectValues>({});
  const [content, setContent] = useState<string>("");

  const noticeTypes = ["Update", "Upgrade", "Feature"];
  const variableOptions = [
    "User's First Name",
    "User's Last Name",
    "Current Time",
    "User's Trovo Plan",
    "Notification Title",
  ];

  return (
    <Container>
      <Title>Content</Title>
      <DropdownSelect
        options={noticeTypes}
        placeholder="Update"
        labelText="Select Notification Type"
        labelColor="#828282"
        placeholderColor="#00225A"
        value={selectValues.NoticeType || ""}
        onSelect={(item: string) =>
          setSelectValues((prev) => ({ ...prev, NoticeType: item }))
        }
      />
      <InputLabel>Notification Title</InputLabel>
      <TextInput placeholder="Enter notification title" />
      <VariableContainer>
        <VariableText>
          You can pick and drop the variables at the desired spot on your text
          editor
        </VariableText>
        <VariableContent>
          {variableOptions.map((item) => (
            <VariableItem key={item}>{item}</VariableItem>
          ))}
        </VariableContent>
      </VariableContainer>
      <TextEditor content={content} setContent={setContent} />
    </Container>
  );
};

export default Content;

const Container = styled.section``;
const Title = styled.h1`
  font-size: 20px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  color: #00225a;
`;

const TextInput = styled.input`
  border: 1px solid #e0e0e0;
  padding: 8px 12px;
  border-radius: 10px;
  width: 95%;
  outline: none;
  margin-bottom: 20px;

  &::placeholder {
    font-size: 14px;
    font-weight: 500;
  }
`;

const InputLabel = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 24px;
  color: #828282;
`;
const VariableContainer = styled.div``;
const VariableText = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #828282;
  margin-bottom: 8px;
  margin-top: 0;
`;

const VariableContent = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 20px;
  margin-top: 10px;
`;

const VariableItem = styled.div`
  padding: 6px 10px;
  background-color: #007cdf1a;
  border-radius: 16px;
  border: 1px solid #007cdf;
  font-size: 12px;
  cursor: pointer;
  color: #007cdf;
  user-select: none;
`;
