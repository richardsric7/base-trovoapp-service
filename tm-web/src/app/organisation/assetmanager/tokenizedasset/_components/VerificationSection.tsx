import styled from "styled-components";
import DocumentCard from "./DocumentCard";
import { AiOutlineCheckCircle } from "react-icons/ai";
import { FiCalendar } from "react-icons/fi";

interface Props {
  title: string;
  button: string;
  verifiedDate?: string;
  verifiedBy?: string;
}

const VerificationSection = ({
  title,
  button,
  verifiedBy,
  verifiedDate,
}: Props) => {
  const documents = [
    {
      id: 1,
      title: "Certificate of Incorporation",
      size: "1.2 MB",
      fileUrl: "/documents/certificate.pdf",
    },
    {
      id: 2,
      title: "Memorandum of Association",
      size: "860 KB",
      fileUrl: "/documents/memorandum.pdf",
    },
    {
      id: 3,
      title: "Tax Clearance Certificate",
      size: "540 KB",
      fileUrl: "",
    },
  ];
  const openPdfViewer = (doc: any) => {
    console.log("Opening PDF:", doc);

    // Example if you want to open the file
    if (doc.fileUrl) {
      window.open(doc.fileUrl, "_blank");
    }
  };

  return (
    <Section>
      <SectionHeader>
        <RightSection>
          <SectionTitle>{title}</SectionTitle>

          {verifiedDate && verifiedBy && (
            <VerifiedMeta>
              <DateWrap>
                <FiCalendar size={12} />
                <span>Verified on {verifiedDate}</span>
              </DateWrap>

              <VerifierWrap>
                <Avatar>A</Avatar>
                <VerifierText>By {verifiedBy}</VerifierText>
              </VerifierWrap>
            </VerifiedMeta>
          )}
        </RightSection>

        <SeeMore>
          <AiOutlineCheckCircle /> {button}
        </SeeMore>
      </SectionHeader>

      <CardGrid>
        {documents.map((doc) => (
          <DocumentCard
            key={doc.id}
            title={doc.title}
            size={doc.size}
            fileUrl={doc.fileUrl}
            onClick={() => openPdfViewer(doc)}
          />
        ))}
      </CardGrid>
    </Section>
  );
};

export default VerificationSection;

const Section = styled.div`
  background: white;
  border-radius: 10px;
  padding: 18px;
  margin-bottom: 16px;
`;

const SectionHeader = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  align-items: center;
  gap: 4px;
`;

const SectionTitle = styled.h4`
  font-size: 14px;
  font-weight: 600;
  color: #00225a;
`;

const SeeMore = styled.button`
  border: 1px solid #007cdf;
  background: #fff;
  padding: 6px 10px;
  color: #007cdf;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  font-family: inherit;
  display: flex;
  gap: 4px;
  align-items: center;
`;

const CardGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;

  border: 1px solid #e0e0e0;
  border-radius: 12px;
  padding: 24px;
`;

const RightSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: 6px;
`;

const VerifiedMeta = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
`;

const DateWrap = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  color: #667085;
  font-size: 11px;
`;

const VerifierWrap = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const Avatar = styled.div`
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #f2f4f7;
  color: #344054;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 600;
`;

const VerifierText = styled.span`
  color: #667085;
  font-size: 11px;
`;
