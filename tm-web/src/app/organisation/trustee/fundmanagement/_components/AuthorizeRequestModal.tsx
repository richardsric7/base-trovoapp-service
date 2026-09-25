"use client";

import StepUpAuthorizationModal from "@/app/organisation/trustee/_components/StepUpAuthorizationModal";
import { useApproveTrusteeFundReleaseMutation } from "@/redux/api/trustees";
import { IFundReleaseRecord } from "@/redux/api/trustees/interface";
import { message } from "antd";
import FundReleaseAuthorizationDetails from "@/app/organisation/components/FundReleaseAuthorizationDetails";

interface Props {
  open: boolean;
  onClose: () => void;
  record?: IFundReleaseRecord;
}

export const AuthorizeRequestModal = ({ open, onClose, record }: Props) => {
  const [approve, { isLoading }] = useApproveTrusteeFundReleaseMutation();

  const handleAuthorize = async (challengeId: string) => {
    if (!record?.id) return;
    try {
      await approve({
        requestId: record.id,
        payload: { challenge_id: challengeId },
      }).unwrap();
      message.success("Fund release authorized successfully.");
      onClose();
    } catch (error: any) {
      message.error(
        error?.data?.message ??
          error?.error ??
          "Failed to authorize fund release. Please try again.",
      );
      throw error;
    }
  };

  return (
    <StepUpAuthorizationModal
      open={open}
      onClose={onClose}
      action="fund_release.approve"
      entityType="fund_release_request"
      entityId={record?.id}
      title="Authorize Fund Release"
      description="Confirm the details, then approve this fund release with your linked Trovo Wallet."
      confirmLabel="Authorize Fund Release"
      onAuthorize={handleAuthorize}
      isAuthorizing={isLoading}
    >
      <FundReleaseAuthorizationDetails record={record} />
    </StepUpAuthorizationModal>
  );
};
