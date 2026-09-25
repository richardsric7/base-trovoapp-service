"use client";
import React, { useState } from "react";
import styled from "styled-components";
import Image, { StaticImageData } from "next/image";

interface Asset {
  name: string;
  value: string;
  amount: string;
  val: string;
  image: string | StaticImageData;
}

interface SharedAssetsProps {
  assetTokens: Asset[];
  otherAssets: Asset[];
}

const UserAssets: React.FC<SharedAssetsProps> = ({
  assetTokens,
  otherAssets,
}) => {
  const [currentTab, setCurrentTab] = useState("assettoken");

  return (
    <UserAssetsContainer>
      <SectionHeader>
        <SectionTitle>Asset</SectionTitle>
        <TabContainer>
          <TabButton
            active={currentTab === "assettoken"}
            onClick={() => setCurrentTab("assettoken")}
          >
            Asset Token
          </TabButton>
          <TabButton
            active={currentTab === "othertoken"}
            onClick={() => setCurrentTab("othertoken")}
          >
            Other Token
          </TabButton>
        </TabContainer>
      </SectionHeader>

      {currentTab === "assettoken" && (
        <AssetToken>
          {assetTokens.map((asset, index) => (
            <React.Fragment key={index}>
              <AssetContainer>
                <AssetCard>
                  <Avatar>
                    {asset.image && (
                      <StyledImage
                        src={asset.image}
                        alt={asset.name}
                        width={40}
                        height={40}
                      />
                    )}
                  </Avatar>

                  <AssetDetails>
                    <AssetName>{asset.name}</AssetName>
                    <AssetValue>{asset.value}</AssetValue>
                  </AssetDetails>
                </AssetCard>

                <AssetDetails>
                  <AssetName>{asset.amount}</AssetName>
                  <AssetValue>{asset.val}</AssetValue>
                </AssetDetails>
              </AssetContainer>
              {index < assetTokens.length - 1 && <Divider />}
            </React.Fragment>
          ))}
        </AssetToken>
      )}

      {currentTab === "othertoken" && (
        <OtherToken>
          {otherAssets.map((asset, index) => (
            <React.Fragment key={index}>
              <AssetContainer>
                <AssetCard>
                  <Avatar></Avatar>
                  <AssetDetails>
                    <AssetName>{asset.name}</AssetName>
                    <AssetValue>{asset.value}</AssetValue>
                  </AssetDetails>
                </AssetCard>

                <AssetDetails>
                  <AssetName>{asset.amount}</AssetName>
                  <AssetValue>{asset.val}</AssetValue>
                </AssetDetails>
              </AssetContainer>
              {index < otherAssets.length - 1 && <Divider />}
            </React.Fragment>
          ))}
        </OtherToken>
      )}
    </UserAssetsContainer>
  );
};

export default UserAssets;

const UserAssetsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 30px;
  padding: 16px;
  background-color: #fff;
  border-radius: 24px;
`;

const SectionHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const SectionTitle = styled.h2`
  font-size: 20px;
  font-weight: 700;
  line-height: 28px;
  color: #00225a;
  margin: 0;
  white-space: nowrap;
`;

const TabContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: #f5f5f5;
  padding: 4px;
  border-radius: 12px;
`;

const TabButton = styled.button<{ active: boolean }>`
  padding: 6px 12px;
  font-size: 14px;
  font-weight: 400;
  line-height: 25.27px;
  color: ${(props) => (props.active ? "#00225A" : "#828282")};

  background-color: ${(props) => (props.active ? "#fff" : "transparent")};
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;

  &:hover {
    background-color: ${(props) => (props.active ? "#F0F4FF" : "#F9F9F9")};
  }
`;

const AssetToken = styled.div``;

const OtherToken = styled.div``;

const AssetContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
`;

const AssetCard = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const StyledImage = styled(Image)`
  width: 40px;
  height: 40px;
  border-radius: 50%;
`;

const AssetDetails = styled.div``;

const AssetName = styled.h3`
  font-size: 16px;
  font-weight: 500;
  line-height: 24px;
  color: #00225a;

  margin: 0;
`;

const AssetValue = styled.p`
  font-size: 14px;
  font-weight: 400;
  line-height: 25.27px;
  color: #00225a;
  margin: 0;
`;

const Divider = styled.div`
  height: 1px;
  background-color: #e5e5ef;
  margin: 15px 0;
`;
