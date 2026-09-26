"use client";
import React from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaArrowLeft } from "react-icons/fa6";
import { useRouter } from "next/navigation";
import { showErrorToast, showSuccessToast } from "@/components";
import { useSaveCuratedAssetMutation } from "@/redux/api/curatedAssets";
import CuratedAssetForm, { CuratedAssetFormValues } from "../components/CuratedAssetForm";

const AddCurationPage = () => {
  const router = useRouter();
  const [saveCuratedAsset, { isLoading }] = useSaveCuratedAssetMutation();

  const handleSubmit = async (values: CuratedAssetFormValues) => {
    try {
      await saveCuratedAsset({
        action: "create",
        assetCode: values.assetCode.trim().toUpperCase(),
        assetName: values.assetName,
        contractAddress: values.contractAddress,
        assetClassId: values.assetClassId,
        decimalPlaces: values.decimalPlaces,
        priority: values.priority,
        assetLimit: values.assetLimit,
        description: values.description,
        website: values.website,
        organization: values.organization,
        contactEmail: values.contactEmail,
        withdrawable: values.withdrawable,
        generateDepositAddress: values.generateDepositAddress,
        inactive: values.inactive,
        p2pEnabled: values.p2pEnabled,
      }).unwrap();
      showSuccessToast("Curated asset created");
      router.push("/assetcuration");
    } catch (e: any) {
      showErrorToast(e?.data?.error || "Could not create curated asset");
    }
  };

  return (
    <PageContainer>
      <BackLink href="/assetcuration">
        <FaArrowLeft />
      </BackLink>

      <Header>
        <TitleSection>
          <Title>Add curated asset</Title>
          <Text>
            Add a new asset to the platform catalog. Turn on &quot;Available on P2P
            marketplace&quot; if merchants should be able to trade it right away.
          </Text>
        </TitleSection>
      </Header>

      <CuratedAssetForm submitLabel="Create asset" submitting={isLoading} onSubmit={handleSubmit} />
    </PageContainer>
  );
};

export default AddCurationPage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  gap: 20px;
  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 25px;
`;

const TitleSection = styled.div`
  flex-grow: 1;
`;
const Text = styled.p`
  font-size: 16px;
  font-weight: 400;
  margin: 0;
  color: #828282;
`;
const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const BackLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  padding-bottom: 20px;
`;
