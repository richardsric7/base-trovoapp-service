"use client";

import { useEffect } from "react";
import { Form, Input, Modal, Select, message } from "antd";
import StepUpAuthorizationModal, {
  Content,
  Title,
  Description,
  PrimaryButton,
} from "@/app/organisation/trustee/_components/StepUpAuthorizationModal";
import FundReleaseAuthorizationDetails from "@/app/organisation/components/FundReleaseAuthorizationDetails";
import { IFundReleaseRecord } from "@/redux/api/trustees";
import {
  FundReleaseExecutionPayload,
  useExecuteCustodianFundReleaseMutation,
  useUpdateCustodianFundReleaseStatusMutation,
} from "@/redux/api/org";

const errorMessage = (error: unknown) => {
  const value = error as {
    message?: string;
    error?: string;
    data?: { error?: string; message?: string };
  };
  return (
    value?.data?.message ||
    value?.data?.error ||
    value?.message ||
    value?.error ||
    "Unable to update fund release. Please try again."
  );
};

export default function ExecutionModal({
  open,
  onClose,
  record,
}: {
  open: boolean;
  onClose: () => void;
  record?: IFundReleaseRecord;
}) {
  const [form] = Form.useForm<FundReleaseExecutionPayload>();
  const status = Form.useWatch("status", form);
  const [execute, execution] = useExecuteCustodianFundReleaseMutation();
  const [update, updating] = useUpdateCustodianFundReleaseStatusMutation();
  const isExecuting = record?.status === "execution_pending";

  useEffect(() => {
    if (!open) {
      form.resetFields();
    }
  }, [open, form]);

  const submit = async () => {
    let values: FundReleaseExecutionPayload;
    try {
      values = await form.validateFields();
    } catch {
      return;
    }
    if (!record || !["execution_pending", "processing"].includes(record.status))
      return;
    const next = {
      ...values,
      execution_reference: values.execution_reference?.trim(),
      failure_reason:
        values.status === "failed" ? values.failure_reason?.trim() : undefined,
    };
    if (isExecuting) return;
    try {
      await update({ requestId: record.id, payload: next }).unwrap();
      message.success("Fund release status updated.");
      onClose();
    } catch (error) {
      message.error(errorMessage(error));
    }
  };

  if (isExecuting)
    return (
      <StepUpAuthorizationModal
        open={open}
        onClose={() => {
          if (!execution.isLoading) onClose();
        }}
        action="fund_release.execute"
        entityType="fund_release_request"
        entityId={record?.id}
        title="Authorize Fund Release"
        description="Confirm the details, then release these funds with your linked Trovo Wallet."
        confirmLabel="Release Funds"
        isAuthorizing={execution.isLoading}
        onAuthorize={async (challengeId) => {
          if (!record) return;
          try {
            await execute({
              requestId: record.id,
              payload: { status: "processing", challenge_id: challengeId },
            }).unwrap();
            message.success("Fund release execution recorded.");
            onClose();
          } catch (error) {
            message.error(errorMessage(error));
          }
        }}
      >
        <FundReleaseAuthorizationDetails record={record} />
      </StepUpAuthorizationModal>
    );

  return (
    <>
      <Modal
        open={open}
        centered
        width={540}
        footer={null}
        destroyOnClose
        onCancel={() => {
          if (!updating.isLoading) onClose();
        }}
      >
        <Content>
          <Title>Update Status</Title>
          <Description>
            Confirm the details, then update the fund release execution status.
          </Description>
          <FundReleaseAuthorizationDetails record={record} />
          <Form
            form={form}
            layout="vertical"
            disabled={updating.isLoading}
            initialValues={{
              status: "completed",
              execution_reference: record?.execution_reference,
            }}
          >
            <Form.Item
              name="status"
              label="Execution Status"
              rules={[{ required: true }]}
            >
              <Select
                options={["completed", "failed"].map((value) => ({
                  value,
                  label: value.charAt(0).toUpperCase() + value.slice(1),
                }))}
              />
            </Form.Item>
            <Form.Item name="execution_reference" label="Execution Reference">
              <Input />
            </Form.Item>
            {status === "failed" && (
              <Form.Item
                name="failure_reason"
                label="Failure Reason"
                rules={[{ required: true, whitespace: true }]}
              >
                <Input.TextArea />
              </Form.Item>
            )}
          </Form>
          <PrimaryButton
            type="button"
            onClick={() => void submit()}
            disabled={updating.isLoading}
          >
            {updating.isLoading ? "Updating..." : "Update Status"}
          </PrimaryButton>
        </Content>
      </Modal>
    </>
  );
}
