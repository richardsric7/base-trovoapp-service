import { Modal } from "@/components";
import {
  AiOutlineCheckCircle,
  AiOutlineClockCircle,
  AiOutlineCloseCircle,
} from "react-icons/ai";
import { MdOutlineContentCopy, MdVerified } from "react-icons/md";
import styled from "styled-components";

type PaymentStatusType = "completed" | "pending" | "failed";

interface PaymentModalProps {
  isOpen: boolean;
  onClose: () => void;
  record: {
    id: string;
    username: string;

    amount: string;
    paymentType: string;
    date: string;
    status: string;
    serviceProvider: string;
  };
}

const PaymentModal: React.FC<PaymentModalProps> = ({
  isOpen,
  onClose,
  record,
}) => {
  const statusConfig: Record<
    PaymentStatusType,
    {
      color: string;
      bg: string;
      bigIcon: JSX.Element;
      smallIcon: JSX.Element;
    }
  > = {
    completed: {
      color: "#00A859",
      bg: "#00A8591A",
      bigIcon: <AiOutlineCheckCircle size={34} color="#00A859" />,
      smallIcon: <AiOutlineCheckCircle size={18} color="#00A859" />,
    },
    pending: {
      color: "#007CDF",
      bg: "#F2F6F9",
      bigIcon: <AiOutlineClockCircle size={34} color="#007CDF" />,
      smallIcon: <AiOutlineClockCircle size={18} color="#007CDF" />,
    },
    failed: {
      color: "#BE3800",
      bg: "#BE38001A",
      bigIcon: <AiOutlineCloseCircle size={34} color="#BE3800" />,
      smallIcon: <AiOutlineCloseCircle size={18} color="#BE3800" />,
    },
  };

  const normalizedStatus: PaymentStatusType =
    record.status === "COMPLETED"
      ? "completed"
      : record.status === "PENDING"
        ? "pending"
        : "failed";

  const config = statusConfig[normalizedStatus];

  const formattedDate = new Date(record.date).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });

  return (
    <>
      <Modal title="" isOpen={isOpen} onClose={onClose}>
        <ModalHeader>
          <HeaderDetails>
            {config.bigIcon}

            <Text>
              {Number(record.amount).toLocaleString("en-US", {
                minimumFractionDigits: 2,
              })}{" "}
              NGN
            </Text>
            <DateValue> {formattedDate}</DateValue>
          </HeaderDetails>
        </ModalHeader>

        <ModalBody>
          <Section>
            <SectionTitle>Payment Details</SectionTitle>
            <DetailsSection>
              <Label>Status</Label>
              <PaymentStatus
                style={{ color: config.color, backgroundColor: config.bg }}
              >
                {config.smallIcon}

                {record.status.charAt(0).toUpperCase() +
                  record.status.slice(1).toLowerCase()}
              </PaymentStatus>
            </DetailsSection>
            <DetailsSection>
              <Label>Amount</Label>
              <CopyableTextBlue>
                {Number(record.amount).toLocaleString("en-US", {
                  minimumFractionDigits: 2,
                })}{" "}
                NGN
              </CopyableTextBlue>
            </DetailsSection>

            <DetailsSection>
              <Label>Date</Label>
              <CopyableTextBlue>{formattedDate}</CopyableTextBlue>
            </DetailsSection>
            <DetailsSection>
              <Label>Payment Gateway</Label>
              <CopyableTextBlue>{record.serviceProvider}</CopyableTextBlue>
            </DetailsSection>
          </Section>
          <SectionDivider />
          <Section>
            <SectionTitle>Received On</SectionTitle>
            <DetailsSection>
              <Label>Account Name</Label>
              <CopyableTextBlue>
                Trovotech Ltd
                <MdOutlineContentCopy size={20} />
              </CopyableTextBlue>
            </DetailsSection>
            <DetailsSection>
              <Label>Account Number</Label>
              <CopyableTextBlue>
                1234567890
                <MdOutlineContentCopy size={20} />
              </CopyableTextBlue>
            </DetailsSection>
            <DetailsSection>
              <Label>Bank</Label>
              <CopyableTextBlue>
                Bank of Africa
                <MdOutlineContentCopy size={20} />
              </CopyableTextBlue>
            </DetailsSection>
          </Section>

          <SectionDivider />
          <Section>
            <SectionTitle>Sender</SectionTitle>
            <UserInfoSection>
              <Avatar />
              <div>
                <UserNameWrapper>
                  <UserName>{record.username}</UserName>
                  {/* <MdVerified color="#007cdf" size={12} /> */}
                </UserNameWrapper>
              </div>
            </UserInfoSection>
          </Section>
        </ModalBody>
      </Modal>
    </>
  );
};

export default PaymentModal;

const ModalTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
  margin-bottom: 10px;
`;

const ModalHeader = styled.div`
  background-color: #f2f6f9;
  padding: 16px;
  gap: 10px;
  border-radius: 12px;
  display: flex;
  margin-bottom: 30px;
  align-items: center;
  justify-content: center;
`;

const HeaderDetails = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  text-align: center;
`;

const ModalBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const DetailsSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  column-gap: 50px;
`;

const Label = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #828282;
`;

const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const Avatar = styled.div`
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: rebeccapurple;
  max-width: 40px;
`;

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const UserInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
  background-color: #f2f6f9;
  border-radius: 12px;
  padding: 16px;
`;

const FullName = styled.p`
  font-size: 11px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
  text-transform: capitalize;
`;

const CopyableTextBlue = styled.p`
  font-size: 14px;
  font-weight: 500;
  color: #00225a;
  display: flex;
  aligin-items: center;
  gap: 4px;
  cursor: pointer;
`;

const SectionDivider = styled.div`
  width: 100%;
  height: 1px;
  background-color: #e5e5ef;
  margin: 12px 0;
`;

const SectionTitle = styled(ModalTitle)`
  margin-bottom: 0;
`;

const Section = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

const Text = styled.p`
  font-weight: 600;
  font-size: 24px;
  letter-spacing: 0px;
  text-align: center;
  color: #00225a;
`;

const DateValue = styled.p`
  font-weight: 400;
  font-size: 14px;
  ext-align: center;
  color: #828282;
`;

const PaymentStatus = styled.p`
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border-radius: 8px;
  text-transform: capitalize;
`;
