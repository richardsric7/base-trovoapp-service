"use client";
import React, { useState } from "react";
import styled from "styled-components";
import { showErrorToast, showSuccessToast } from "@/components";
import PrimaryButton from "@/components/PrimaryButton";
import CustomTable from "@/components/CustomTable";
import {
  IFeeExemptUser,
  useAddFeeExemptUserMutation,
  useGetFeeExemptUsersQuery,
  useRemoveFeeExemptUserMutation,
} from "@/redux/api/feeExemptions";

const errorMessage = (err: any, fallback: string) =>
  err?.data?.error || err?.data?.message || err?.error || fallback;

const FeeExemptionsPage = () => {
  const { data, isLoading } = useGetFeeExemptUsersQuery();
  const [addExemption, { isLoading: adding }] = useAddFeeExemptUserMutation();
  const [removeExemption] = useRemoveFeeExemptUserMutation();
  const [username, setUsername] = useState("");
  const [reason, setReason] = useState("");

  const rows = data?.data ?? [];

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !reason.trim()) {
      showErrorToast("Enter the account's username and the reason for the exemption.");
      return;
    }
    try {
      await addExemption({ username: username.trim(), reason: reason.trim() }).unwrap();
      showSuccessToast(`${username.trim()} no longer pays service fees`);
      setUsername("");
      setReason("");
    } catch (err) {
      showErrorToast(errorMessage(err, "Could not add the exemption"));
    }
  };

  const handleRemove = async (record: IFeeExemptUser) => {
    if (!window.confirm(`Make ${record.username} pay service fees again?`)) return;
    try {
      await removeExemption(record.username).unwrap();
      showSuccessToast(`${record.username} pays service fees again`);
    } catch (err) {
      showErrorToast(errorMessage(err, "Could not remove the exemption"));
    }
  };

  const columns = [
    { title: "Username", dataIndex: "username", render: (v: string) => <Strong>{v}</Strong> },
    { title: "Reason", dataIndex: "reason" },
    { title: "Added by", dataIndex: "addedBy", render: (v: string) => v || "—" },
    {
      title: "Since",
      dataIndex: "createdAt",
      render: (v: string) => (v ? new Date(v).toLocaleDateString() : "—"),
    },
    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: any, record: IFeeExemptUser) => (
        <RemoveButton type="button" onClick={() => handleRemove(record)}>
          Remove
        </RemoveButton>
      ),
    },
  ];

  return (
    <PageContainer>
      <Title>Fee Exemptions</Title>
      <Text>
        Accounts listed here pay no platform service fees: swap, payment, patron, account recovery,
        sub-wallet creation, tokenization application and closed-group fees. Use it for the
        platform&apos;s own trading or operations accounts. The tokenization issuing profile is always
        exempt and is not listed. Changes apply to the account&apos;s next request and are recorded in
        the audit trail.
      </Text>

      <AddForm onSubmit={handleAdd}>
        <Input placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
        <Input
          placeholder="Reason (e.g. platform market maker)"
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          $wide
        />
        <PrimaryButton disabled={adding} buttonStyle={{ width: "auto", padding: "10px 20px" }}>
          {adding ? "Adding..." : "Exempt account"}
        </PrimaryButton>
      </AddForm>

      <CustomTable
        columns={columns}
        dataSource={rows}
        totalItems={rows.length}
        pageSize={Math.max(rows.length, 1)}
        isLoading={isLoading}
        onPageChange={() => {}}
      />
    </PageContainer>
  );
};

export default FeeExemptionsPage;

const PageContainer = styled.section`
  background: #ffffff;
  padding: 32px;
  border-radius: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const Text = styled.p`
  font-size: 14px;
  color: #828282;
  margin: 8px 0 0;
  line-height: 22px;
  max-width: 820px;
`;

const AddForm = styled.form`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  margin: 28px 0;
`;

const Input = styled.input<{ $wide?: boolean }>`
  flex: ${(p) => (p.$wide ? "1 1 320px" : "0 1 220px")};
  min-width: 0;
  padding: 10px 14px;
  border: 1px solid #d9e1ec;
  border-radius: 8px;
  font-size: 14px;
  color: #00225a;
  background: #f2f6f9;
`;

const Strong = styled.span`
  font-weight: 600;
  color: #00225a;
`;

const RemoveButton = styled.button`
  background: none;
  border: none;
  color: #d92d20;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  padding: 0;
`;
