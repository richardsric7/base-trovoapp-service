"use client";

import Tab from "@/components/Tab";
import { useState } from "react";
import styled from "styled-components";
import OrgTokenizedAssetDetail from "../_components/TokenizedAssetDetails";
import AssetActivity from "../../../components/AssetActivity";
import OrgDueDiligence from "../_components/OrgDueDiligence";
import OrgFundManagement from "../_components/OrgFundManagement";
import { FaArrowLeft } from "react-icons/fa6";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useGetStakeholderAssetDetailQuery } from "@/redux/api/sharedstakeholders";

export default function TokenizedAssetsDetailsPage() {
  const [currentTab, setCurrentTab] = useState("asset-details");
  const params = useParams();
  const id = params?.id as string;

  const { data: assetDetail, isLoading } = useGetStakeholderAssetDetailQuery(
    id,
    { skip: !id },
  );

  const tabs = [
    { key: "asset-details", label: "Asset Details" },
    { key: "activity", label: "Activity" },
    // { key: "due-diligence", label: "Due Diligence" },
    // { key: "fund-management", label: "Fund Management" },
  ];
  return (
    <Container>
      <StyledLink href="/organisation/issuinghouse/tokenizedasset">
        <FaArrowLeft size={20} />
      </StyledLink>

      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "250px",
        }}
        tabContentBackground="transparent"
        tabContentStyle={{ padding: "0" }}
      >
        {currentTab === "asset-details" && (
          <OrgTokenizedAssetDetail
            assetDetail={assetDetail}
            isLoading={isLoading}
          />
        )}
        {currentTab === "activity" && (
          <AssetActivity assetId={id} assetCode={assetDetail?.data.assetCode} assetName={assetDetail?.data.assetName} />
        )}

        {currentTab === "due-diligence" && <OrgDueDiligence assetId={id} />}
        {/* {currentTab === "fund-management" && <OrgFundManagement />} */}
      </Tab>
    </Container>
  );
}

const Container = styled.section`
  padding: 20px;
`;

const StyledLink = styled(Link)`
  text-decoration: none;
  color: #00225a;
  font-size: 20px;
  cursor: pointer;
  border: none;
  margin-bottom: 100px;
`;
