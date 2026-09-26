"use client";

import { useParams, useRouter } from "next/navigation";
import React from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaArrowLeft } from "react-icons/fa6";
import { showErrorToast, showSuccessToast } from "@/components";
import { useGetCuratedAssetByIdQuery, useSaveCuratedAssetMutation } from "@/redux/api/curatedAssets";
import CuratedAssetForm, { CuratedAssetFormValues } from "../components/CuratedAssetForm";

// The route segment is still named [userId] (unchanged, to avoid an
// unrelated file-move diff) but identifies a curated asset here, not a
// user.
const EditCuratedAssetPage = () => {
  const params = useParams();
  const router = useRouter();
  const id = Number(params.userId);

  const { data, isLoading } = useGetCuratedAssetByIdQuery(id, { skip: !id });
  const [saveCuratedAsset, { isLoading: isSaving }] = useSaveCuratedAssetMutation();

  const handleSubmit = async (values: CuratedAssetFormValues) => {
    try {
      await saveCuratedAsset({
        action: "update",
        id,
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
      showSuccessToast("Curated asset updated");
      router.push("/assetcuration");
    } catch (e: any) {
      showErrorToast(e?.data?.error || "Could not update curated asset");
    }
  };

  return (
    <Container>
      <BackLink href="/assetcuration">
        <FaArrowLeft />
      </BackLink>

      {isLoading || !data?.data ? (
        <Loading>Loading asset...</Loading>
      ) : (
        <>
          <Header>
            <Title>Edit {data.data.assetCode}</Title>
          </Header>
          <CuratedAssetForm
            initialAsset={data.data}
            submitLabel="Save changes"
            submitting={isSaving}
            onSubmit={handleSubmit}
          />
        </>
      )}
    </Container>
  );
};

export default EditCuratedAssetPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
  gap: 20px;
`;

const Header = styled.div`
  padding-bottom: 20px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const Loading = styled.p`
  color: #828282;
  padding: 40px 0;
  text-align: center;
`;

const BackLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #000000;
  width: 20px;
  padding-bottom: 20px;
`;
