"use client";

import { Modal } from "@/components";
import { useGetAuditActivityQuery } from "@/redux/api/auditTrail/api";
import { HiOutlineCheckCircle, HiOutlineXCircle } from "react-icons/hi2";
import { MdVerified } from "react-icons/md";
import styled from "styled-components";

interface ActivityModalProps {
  id: string;
  isOpen: boolean;
  onClose: () => void;
}

const ActivityModal: React.FC<ActivityModalProps> = ({
  id,
  isOpen,
  onClose,
}) => {
  const { data, isLoading } = useGetAuditActivityQuery(
    { id },
    { skip: !id },
  );
  const record = data?.data;

  const isSuccessful = record?.status === "successful";

  const statusLabel = record ? (
    <StatusLabel $ok={isSuccessful}>
      {isSuccessful ? (
        <HiOutlineCheckCircle size={16} color="#0DB57A" />
      ) : (
        <HiOutlineXCircle size={16} color="#BE3800" />
      )}
      {record.status}
    </StatusLabel>
  ) : (
    "—"
  );

  const activityData = [
    { text: "Action", label: record?.action || "—" },
    { text: "Date/Time", label: record?.date || "—" },
    { text: "IP Address", label: record?.ip_address || "—" },
    { text: "Status", label: statusLabel },
  ];

  // "More Details" is variable-length: only show rows we actually have data for.
  const moreDetails = [
    record?.login_method
      ? { text: "Login Method", label: prettyMethod(record.login_method) }
      : null,
    record?.location ? { text: "Location", label: record.location } : null,
    record?.target ? { text: "Target", label: record.target } : null,
    record?.role ? { text: "Role", label: record.role } : null,
    record?.detail ? { text: "Detail", label: record.detail } : null,
  ].filter(Boolean) as { text: string; label: string }[];

  return (
    <Modal
      title="View Activity"
      isOpen={isOpen}
      onClose={onClose}
      closeIconPosition="left"
    >
      <ModalBody>
        <Section>
          <UserInfoSection>
            <Avatar />
            <div>
              <UserNameWrapper>
                <UserName>
                  {isLoading ? "Loading…" : record?.fullname || record?.username || "—"}
                </UserName>
                <MdVerified color="#007cdf" size={14} />
              </UserNameWrapper>
              <FullName>{record?.email || record?.username || ""}</FullName>
            </div>
          </UserInfoSection>
        </Section>

        <ActivityDataContainer>
          {activityData.map((item, index) => (
            <GridItem key={index}>
              <ActivityText>{item.text}</ActivityText>
              <ActivityLabel>{item.label}</ActivityLabel>
            </GridItem>
          ))}
        </ActivityDataContainer>

        {moreDetails.length > 0 && (
          <>
            <SectionDivider />
            <Section>
              <SectionTitle>More Details</SectionTitle>
              <ActivityDataContainer>
                {moreDetails.map((item, index) => (
                  <GridItem key={index}>
                    <ActivityText>{item.text}</ActivityText>
                    <ActivityLabel>{item.label}</ActivityLabel>
                  </GridItem>
                ))}
              </ActivityDataContainer>
            </Section>
          </>
        )}
      </ModalBody>
    </Modal>
  );
};

export default ActivityModal;

// prettyMethod turns the stored login method into a human label. Admin login is
// delegated (push / QR approval) — there is no password method.
function prettyMethod(method: string): string {
  switch (method) {
    case "push_approval":
      return "Push approval";
    case "qr":
      return "QR approval";
    default:
      return method;
  }
}

const ModalBody = styled.div``;

const Section = styled.div`
  margin-bottom: 24px;
`;

const SectionTitle = styled.h4`
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 10px;
  color: #00225a;
`;

const SectionDivider = styled.hr`
  border: none;
  border-top: 1px solid #e0e0e0;
  margin: 24px 0;
`;

const UserInfoSection = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const Avatar = styled.div`
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: rebeccapurple;
  max-width: 40px;
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const UserName = styled.span`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const FullName = styled.div`
  font-size: 14px;
  color: #4b5563;
`;

const ActivityDataContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  row-gap: 16px;
  column-gap: 24px;
  padding: 10px 0;
`;
const GridItem = styled.div``;

const ActivityText = styled.div`
  color: #828282;
  font-weight: 400;
  padding-bottom: 4px;
`;

const ActivityLabel = styled.div`
  font-weight: 600;
  color: #007cdf;
  font-size: 14px;
`;

const StatusLabel = styled.div<{ $ok?: boolean }>`
  display: flex;
  align-items: center;
  gap: 6px;
  color: ${({ $ok }) => ($ok ? "#0db57a" : "#BE3800")};
  font-weight: 600;
  text-transform: capitalize;
`;
