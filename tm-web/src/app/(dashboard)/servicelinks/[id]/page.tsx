"use client";

import { useParams, useRouter } from "next/navigation";
import React from "react";
import styled from "styled-components";
import Link from "next/link";
import { FaArrowLeft } from "react-icons/fa6";
import { showErrorToast, showSuccessToast } from "@/components";
import { useGetServiceLinkByIdQuery, useSaveServiceLinkMutation } from "@/redux/api/serviceLinks";
import ServiceLinkForm, { ServiceLinkFormValues } from "../components/ServiceLinkForm";

const EditServiceLinkPage = () => {
  const params = useParams();
  const router = useRouter();
  const id = Array.isArray(params.id) ? params.id[0] : params.id;

  const { data, isLoading } = useGetServiceLinkByIdQuery(id as string, { skip: !id });
  const [saveServiceLink, { isLoading: isSaving }] = useSaveServiceLinkMutation();

  const handleSubmit = async (values: ServiceLinkFormValues) => {
    try {
      await saveServiceLink({
        action: "update",
        id,
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
      showSuccessToast("Service link updated");
      router.push("/servicelinks");
    } catch (e: any) {
      showErrorToast(e?.data?.error || "Could not update service link");
    }
  };

  return (
    <Container>
      <BackLink href="/servicelinks">
        <FaArrowLeft />
      </BackLink>

      {isLoading || !data?.data ? (
        <Loading>Loading service link...</Loading>
      ) : (
        <>
          <Header>
            <Title>Edit {data.data.shortName}</Title>
            <Meta>
              Owner: <MetaValue>{data.data.ownerUsername}</MetaValue>
              {data.data.owner?.email ? <MetaValue> ({data.data.owner.email})</MetaValue> : null}
            </Meta>
            <Meta>
              API key: <MetaValue $mono>{data.data.apiKey}</MetaValue>
            </Meta>
            {data.data.suspended === 1 && (
              <SuspendedNotice>
                This service link is suspended{data.data.suspensionReason ? `: ${data.data.suspensionReason}` : "."}
              </SuspendedNotice>
            )}
          </Header>
          <ServiceLinkForm
            initialLink={data.data}
            submitLabel="Save changes"
            submitting={isSaving}
            onSubmit={handleSubmit}
          />
        </>
      )}
    </Container>
  );
};

export default EditServiceLinkPage;

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
  margin: 0 0 8px 0;
  color: #00225a;
`;

const Meta = styled.p`
  font-size: 13px;
  color: #828282;
  margin: 0 0 4px 0;
`;

const MetaValue = styled.span<{ $mono?: boolean }>`
  color: #00225a;
  font-weight: 500;
  font-family: ${(props) => (props.$mono ? "monospace" : "inherit")};
`;

const SuspendedNotice = styled.p`
  font-size: 13px;
  color: #be3800;
  background-color: #be38001a;
  border-radius: 8px;
  padding: 8px 12px;
  margin: 8px 0 0 0;
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
