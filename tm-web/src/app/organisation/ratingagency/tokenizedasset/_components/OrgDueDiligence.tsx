"use client";

import React, { useState } from "react";
import styled from "styled-components";
import VerificationSection from "./VerificationSection";
import DueDiligenceScore from "./DueDiligenceScoreCard";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { BsThreeDots } from "react-icons/bs";
import { AiOutlineClockCircle } from "react-icons/ai";
import { Badge, Dropdown, MenuProps, Skeleton, message } from "antd";
import {
  useGetTrusteeDueDiligenceQuery,
  useApproveTrusteeDueDiligenceMutation,
} from "@/redux/api/trustees";
import { IDueDiligenceItem } from "@/redux/api/trustees/interface";
import { RejectDueDiligenceModal } from "./RejectDueDiligenceModal";
import { VerifyCategoryModal } from "./VerifyCategoryModal";

interface Props {
  assetId: string;
}

const CATEGORY_LABELS: Record<string, string> = {
  legal: "Legal Verification",
  financial: "Financial Verification",
  operational: "Operational Verification",
  technical: "Technical Due Diligence",
  risk: "Risk Assessment",
};

const STATUS_COLOR: Record<string, string> = {
  pending: "#007cdf",
  approved: "#12b76a",
  rejected: "#be3800",
  in_review: "#f79009",
};

const DueDiligenceTab = ({ assetId }: Props) => {
  const [rejectModalOpen, setRejectModalOpen] = useState(false);
  const [verifyCategoryState, setVerifyCategoryState] = useState<{
    open: boolean;
    categoryKey: string;
    categoryTitle: string;
  }>({
    open: false,
    categoryKey: "",
    categoryTitle: "",
  });

  const { data, isLoading, refetch } = useGetTrusteeDueDiligenceQuery(assetId, {
    skip: !assetId,
  });

  const [approveDueDiligence, { isLoading: isApproving }] =
    useApproveTrusteeDueDiligenceMutation();

  const checklist = data?.data;
  const items: IDueDiligenceItem[] = checklist?.items ?? [];

  // Group items by category
  const grouped = items.reduce<Record<string, IDueDiligenceItem[]>>(
    (acc, item) => {
      if (!acc[item.category]) acc[item.category] = [];
      acc[item.category].push(item);
      return acc;
    },
    {},
  );

  const overallStatus = checklist?.status ?? "pending";
  const statusLabel = overallStatus.replace(/_/g, " ");
  const isApproved = overallStatus.toLowerCase() === "approved";

  const handleApprove = async () => {
    try {
      await approveDueDiligence(assetId).unwrap();
      message.success("Due diligence approved successfully.");
      refetch();
    } catch {
      message.error("Failed to approve due diligence. Please try again.");
    }
  };

  const handleOpenVerifyCategory = (catKey: string, catTitle: string) => {
    setVerifyCategoryState({
      open: true,
      categoryKey: catKey,
      categoryTitle: catTitle,
    });
  };

  const menuItems: MenuProps["items"] = [
    // { key: "1", label: "Add Notes" },
    // { key: "2", label: "Upload Document" },
    // { key: "3", label: "Download Documents" },

    // { key: "5", label: "Generate a Report" },

    {
      key: "7",
      label: <span style={{ color: "#BE3800" }}>Reject Asset</span>,
      onClick: () => setRejectModalOpen(true),
    },
  ];

  return (
    <>
      <Container>
        <HeaderRow>
          <div>
            <Title>Due Diligence</Title>
            <Status $color={STATUS_COLOR[overallStatus] ?? "#007cdf"}>
              <AiOutlineClockCircle />
              {statusLabel.charAt(0).toUpperCase() + statusLabel.slice(1)}
            </Status>
          </div>

          <ButtonContainer>
            {!isApproved && (
              <PrimaryButton
                buttonStyle={{ width: "154px" }}
                onClick={handleApprove}
                disabled={isApproving}
              >
                {isApproving ? "Approving…" : "Approve Asset"}
              </PrimaryButton>
            )}
            <DropdownStyles>
              <Dropdown
                menu={{ items: menuItems }}
                trigger={["click"]}
                placement="bottomRight"
                overlayClassName="due-diligence-dropdown"
              >
                <div>
                  <DotBadge dot>
                    <SecondaryButton buttonStyle={{ width: "48px" }}>
                      <BsThreeDots />
                    </SecondaryButton>
                  </DotBadge>
                </div>
              </Dropdown>
            </DropdownStyles>
          </ButtonContainer>
        </HeaderRow>

        <DueDiligenceScore score={items.length} />

        {isLoading ? (
          <>
            <Skeleton
              active
              paragraph={{ rows: 3 }}
              style={{ marginBottom: 16 }}
            />
            <Skeleton
              active
              paragraph={{ rows: 3 }}
              style={{ marginBottom: 16 }}
            />
            <Skeleton
              active
              paragraph={{ rows: 3 }}
              style={{ marginBottom: 16 }}
            />
          </>
        ) : Object.keys(grouped).length === 0 ? (
          <>
            {/* Fallback to static sections when no API data yet */}
            <VerificationSection
              title="Legal Verification"
              status="pending"
              onVerify={() =>
                handleOpenVerifyCategory("legal", "Legal Verification")
              }
            />
            <VerificationSection
              title="Financial Verification"
              status="pending"
              onVerify={() =>
                handleOpenVerifyCategory("financial", "Financial Verification")
              }
            />
            <VerificationSection
              title="Operational Verification"
              status="pending"
              onVerify={() =>
                handleOpenVerifyCategory(
                  "operational",
                  "Operational Verification",
                )
              }
            />
            <VerificationSection
              title="Technical Due Diligence"
              status="pending"
              onVerify={() =>
                handleOpenVerifyCategory("technical", "Technical Due Diligence")
              }
            />
            <VerificationSection
              title="Risk Assessment"
              status="pending"
              onVerify={() =>
                handleOpenVerifyCategory("risk", "Risk Assessment")
              }
            />
          </>
        ) : (
          Object.entries(grouped).map(([category, categoryItems]) => {
            const allComplete = categoryItems.every(
              (ci) => ci.status === "approved" || ci.status === "complete",
            );
            const anyFailed = categoryItems.some(
              (ci) => ci.status === "failed" || ci.status === "rejected",
            );
            const anyNotApplicable = categoryItems.some(
              (ci) => ci.status === "not_applicable",
            );

            let categoryStatus = "pending";
            if (allComplete) {
              categoryStatus = "complete";
            } else if (anyFailed) {
              categoryStatus = "failed";
            } else if (
              anyNotApplicable &&
              categoryItems.every(
                (ci) =>
                  ci.status === "complete" ||
                  ci.status === "approved" ||
                  ci.status === "not_applicable",
              )
            ) {
              categoryStatus = "not_applicable";
            } else {
              categoryStatus = categoryItems[0]?.status ?? "pending";
            }

            const sortedItems = [...categoryItems].sort(
              (a, b) =>
                new Date(b.updated_at || b.created_at || 0).getTime() -
                new Date(a.updated_at || a.created_at || 0).getTime(),
            );
            const latestDate =
              sortedItems[0]?.updated_at || sortedItems[0]?.created_at;

            const title =
              CATEGORY_LABELS[category] ??
              `${
                category.charAt(0).toUpperCase() + category.slice(1)
              } Verification`;

            return (
              <VerificationSection
                key={category}
                title={title}
                status={categoryStatus}
                updatedAt={latestDate}
                checklistItems={categoryItems}
                onVerify={() => handleOpenVerifyCategory(category, title)}
              />
            );
          })
        )}
      </Container>

      {/* Reject Modal */}
      <RejectDueDiligenceModal
        open={rejectModalOpen}
        onClose={() => setRejectModalOpen(false)}
        assetId={assetId}
        assetCode={checklist?.asset_code}
        onSuccess={refetch}
      />

      {/* Verify Category Modal */}
      <VerifyCategoryModal
        open={verifyCategoryState.open}
        onClose={() =>
          setVerifyCategoryState((prev) => ({ ...prev, open: false }))
        }
        assetId={assetId}
        categoryKey={verifyCategoryState.categoryKey}
        categoryTitle={verifyCategoryState.categoryTitle}
        onSuccess={refetch}
      />
    </>
  );
};

export default DueDiligenceTab;

const Container = styled.div`
  padding: 24px;
  background-color: #fff;
  border-radius: 24px;
  background-color: #fff;
`;

const HeaderRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
`;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
`;

const Status = styled.div<{ $color: string }>`
  background-color: #f2f6f9;
  color: ${({ $color }) => $color};
  font-size: 12px;
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  text-transform: capitalize;
`;

const DotBadge = styled(Badge)`
  .ant-badge-dot {
    top: 20px;
    right: 2px;
    width: 8px;
    height: 8px;
  }
`;

const DropdownStyles = styled.div`
  .due-diligence-dropdown .ant-dropdown-menu-item {
    font-family: "Montserrat", sans-serif !important;
    font-weight: 400 !important;
    font-size: 14px !important;
    line-height: 28px !important;
    letter-spacing: 0% !important;
    color: #00225a !important;
  }

  .due-diligence-dropdown .ant-dropdown-menu-item:hover {
    background: #f2f6f9;
    color: #00225a;
  }
`;
