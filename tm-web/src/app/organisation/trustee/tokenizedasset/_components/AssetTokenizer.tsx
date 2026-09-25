"use client";

import React from "react";
import styled from "styled-components";

interface AssetTokenizerProps {
  username?: string;
}

const AssetTokenizer: React.FC<AssetTokenizerProps> = ({ username }) => {
  return (
    <>
      <HeadingContent>
        <Heading>Asset Tokenizer</Heading>
      </HeadingContent>

      <Wrapper>
        <AssetIdSection>
          <Avatar />

          <div>
            <UserName>{username || "N/A"}</UserName>
          </div>
        </AssetIdSection>
      </Wrapper>
    </>
  );
};

export default AssetTokenizer;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  margin: 16px 0px;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const UserName = styled.div`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

