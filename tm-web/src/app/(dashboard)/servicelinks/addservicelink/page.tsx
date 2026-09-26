"use client";
import React from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaArrowLeft } from "react-icons/fa6";
import { useRouter } from "next/navigation";
import { showErrorToast, showSuccessToast } from "@/components";
import { useSaveServiceLinkMutation } from "@/redux/api/serviceLinks";
import ServiceLinkForm, { ServiceLinkFormValues } from "../components/ServiceLinkForm";

const AddServiceLinkPage = () => {
  const router = useRouter();
  const [saveServiceLink, { isLoading }] = useSaveServiceLinkMutation();

  const handleSubmit = async (values: ServiceLinkFormValues) => {
    try {
      await saveServiceLink({
        action: "create",
        ownerUsername: values.ownerUsername.trim(),
        shortName: values.shortName.trim(),
        longName: values.longName,
        loginPermission: values.loginPermission,
        paymentPermission: values.paymentPermission,
        tokenInfoPermission: values.tokenInfoPermission,
        authorizationPermission: values.authorizationPermission,
        eventPermission: values.eventPermission,
        allowUserInfo: values.allowUserInfo,
        pushNotificationPermission: values.pushNotificationPermission,
        includePhoneNumbers: values.includePhoneNumbers,
        includeUserBalances: values.includeUserBalances,
        tokenizedAssetAuthorizationPermission: values.tokenizedAssetAuthorizationPermission,
        createUsersPermission: values.createUsersPermission,
        allowReferralForRegisteredUsers: values.allowReferralForRegisteredUsers,
        verified: values.verified,
        inactive: values.inactive,
      }).unwrap();
      showSuccessToast("Service link created");
      router.push("/servicelinks");
    } catch (e: any) {
      showErrorToast(e?.data?.error || "Could not create service link");
    }
  };

  return (
    <PageContainer>
      <BackLink href="/servicelinks">
        <FaArrowLeft />
      </BackLink>

      <Header>
        <TitleSection>
          <Title>Add service link</Title>
          <Text>
            Provision a new white-label partner integration. The owner username
            must already have a Trovo Wallet account - its wallet address is
            copied onto the new service link automatically.
          </Text>
        </TitleSection>
      </Header>

      <ServiceLinkForm submitLabel="Create service link" submitting={isLoading} onSubmit={handleSubmit} />
    </PageContainer>
  );
};

export default AddServiceLinkPage;

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
