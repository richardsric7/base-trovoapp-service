"use client";

import styled from "styled-components";
import DocumentCard from "./DocumentCard";
import PrimaryButton from "@/components/PrimaryButton";
import SecondaryButton from "@/components/SecondaryButton";
import { FaRegEdit } from "react-icons/fa";

interface Props {
  title: string;
  documents: any[];
  showActions?: boolean;
}

const DocumentSection: React.FC<Props> = ({
  title,
  documents,
  showActions = true,
}) => {
  return (
    <Container>
      <Header>
        <SectionTitle>{title}</SectionTitle>

        {showActions && (
          <Buttons>
            <SecondaryButton buttonStyle={{ width: "200px" }}>
              + Add Document
            </SecondaryButton>

            <SecondaryButton
              buttonStyle={{
                width: "100px",
                display: "flex",
                alignItems: "center",
                gap: "4px",
              }}
            >
              <FaRegEdit /> Edit
            </SecondaryButton>
          </Buttons>
        )}
      </Header>
      <Wrapper>
        <Grid>
          {documents.map((doc, i) => (
            <DocumentCard key={i} {...doc} />
          ))}
        </Grid>
      </Wrapper>
    </Container>
  );
};

export default DocumentSection;

const Container = styled.div`
  margin-bottom: 30px;
`;

const Wrapper = styled.div`
  border: 1px solid #e6e6e6;
  border-radius: 20px;
  padding: 24px;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
`;

const SectionTitle = styled.h3`
  font-size: 20px;
  font-weight: 600;
  color: #00225a;
  margin: 0;
  padding: 0;
  line-height: 1.2;
`;

const Buttons = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;

  button {
    margin: 0;
  }
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
`;
