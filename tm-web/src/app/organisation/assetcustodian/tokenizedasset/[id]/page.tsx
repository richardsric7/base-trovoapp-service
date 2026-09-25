"use client";

import Tab from "@/components/Tab";
import { useState } from "react";
import styled from "styled-components";
import CustodianTokenizedAssetDetail from "../_components/TokenizedAssetDetails";
import AssetActivity from "../../../components/AssetActivity";
import { FaArrowLeft } from "react-icons/fa6";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useGetStakeholderAssetDetailQuery } from "@/redux/api/sharedstakeholders";
import Loader from "@/components/Loader";
import { Alert, Button } from "antd";

export default function CustodianTokenizedAssetsDetailsPage() {
  const [currentTab, setCurrentTab] = useState("asset-details");
  const params = useParams();
  const id = params?.id as string;

  const { data: assetDetail, isLoading, isError, refetch } = useGetStakeholderAssetDetailQuery(
    id,
    { skip: !id },
  );

  const tabs = [
    { key: "asset-details", label: "Asset Details" },
    { key: "activity", label: "Activity" },
  ];

  if (isLoading) {
    return <Loader />;
  }

  if (isError || !assetDetail?.data) {
    return <Alert type="error" message="Unable to load asset details." action={<Button onClick={() => refetch()}>Retry</Button>} />;
  }

  return (
    <Container>
      <StyledLink href="/organisation/assetcustodian/tokenizedasset">
        <FaArrowLeft size={20} />
      </StyledLink>

      <Tab
        tabs={tabs}
        currentTab={currentTab}
        setCurrentTab={setCurrentTab}
        tabContainerStyle={{
          width: "100%",
          maxWidth: "400px",
        }}
        tabContentBackground="transparent"
        tabContentStyle={{ padding: "0" }}
      >
        {currentTab === "asset-details" && (
          <CustodianTokenizedAssetDetail
            assetDetail={assetDetail}
            isLoading={isLoading}
          />
        )}
        {currentTab === "activity" && (
          <AssetActivity assetId={id} assetCode={assetDetail?.data.assetCode} assetName={assetDetail?.data.assetName} />
        )}
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
