import React from "react";
import { Navbar } from "../navbar";
import { Footer } from "../footer";
import styled from "styled-components";

const AccountDeletion = () => {
  const handleContactClick = () => {
    window.location.href = "/#contact";
  };

  return (
    <main>
      <Navbar />

      <Title>Help Center</Title>
      <MainContent>
        <Section>
          <SectionTitle>Deleting your Trovo Account</SectionTitle>
          <Paragraph>
            If you wish to delete your account on Trovo App, follow the
            instructions below.
          </Paragraph>
          <List>
            <ListItem>
              Login to the <StyledSpan>Trovo App</StyledSpan> on your mobile
              device
            </ListItem>
            <ListItem>
              Tap the <StyledSpan>hamburger menu</StyledSpan> icon tab at the
              top-left of the screen
            </ListItem>
            <ListItem>
              Tap on your <StyledSpan>Profile picture</StyledSpan> at the
              top-left corner to launch your profile
            </ListItem>
            <ListItem>
              Under your profile, tap on{" "}
              <StyledSpan>“Delete Account”</StyledSpan>
            </ListItem>
            <ListItem>
              This will bring up a confirmation screen for you to confirm if
              you’d like to permanently delete your account.
            </ListItem>
          </List>
          <Note>
            <NoteTitle>
              Note:{" "}
              <StyledNoteSpan>
                {" "}
                Before permanently deleting your account, ensure you have done
                the following:
              </StyledNoteSpan>
            </NoteTitle>
            <NoteContent>
              <NoteItem>Disable account recovery on your account</NoteItem>
              <NoteItem>
                Disable the Initiator access you have to all wallets on your
                account
              </NoteItem>
            </NoteContent>
          </Note>
        </Section>

        <Section>
          <SectionTitle>How Account Deletion Works</SectionTitle>
          <Paragraph>When you permanently delete your account:</Paragraph>
          <NoteContent>
            <NoteItem>
              Your account deletion request will be subject to approval, which
              may take up to 30 days. During this period, you will not be able
              to access your account.
            </NoteItem>
            <NoteItem>
              All data associated with your account will be permanently wiped
              from our servers upon approval of your request. This action cannot
              be undone.
            </NoteItem>
          </NoteContent>
        </Section>

        <HelpFooter>
          <FooterText>Need more help?</FooterText>
          <HelpText>Contact our support team for quick assistance</HelpText>
          <ContactButton onClick={handleContactClick}>Contact Us</ContactButton>
        </HelpFooter>
      </MainContent>
      <Footer />
    </main>
  );
};

export default AccountDeletion;

const MainContent = styled.div`
  margin-top: 40px;
  padding: 0px 150px 0px 150px;
  font-family: Montserrat;
  color: #000000;

  @media (max-width: 1024px) {
    padding: 0px 80px;
  }

  @media (max-width: 768px) {
    padding: 0px 40px;
  }

  @media (max-width: 600px) {
    padding: 0px 30px;
  }
`;

const Title = styled.h1`
  background-color: #004988;
  color: #ffffff;
  margin-top: 106px;
  padding: 0px 150px 0px 150px;
  font-family: Montserrat;
  font-size: 32px;
  font-weight: 700;
  line-height: 72px;

  @media (max-width: 768px) {
    font-size: 28px;
    line-height: 48px;
    padding: 0px 40px;
  }

  @media (max-width: 600px) {
    font-size: 24px;
    line-height: 36px;
    padding: 20px 30px;
    margin-top: 80px;
  }
`;

const Section = styled.div`
  margin-bottom: 24px;
  @media (max-width: 768px) {
    margin-bottom: 16px;
  }
`;

const SectionTitle = styled.h2`
  font-size: 24px;
  font-weight: 600;
  line-height: 33.6px;
  letter-spacing: 0.16875000298023224px;
  text-align: left;
  text-decoration-line: underline;
  text-decoration-style: solid;
  text-underline-position: from-font;
  text-decoration-skip-ink: auto;

  color: #000000;
  margin-bottom: 8px;
  padding-bottom: 4px;

  @media (max-width: 768px) {
    font-size: 20px;
    line-height: 28px;
  }

  @media (max-width: 600px) {
    font-size: 18px;
    line-height: 24px;
  }
`;

const Paragraph = styled.p`
  font-family: Montserrat;
  font-size: 18px;
  font-weight: 400;
  line-height: 25.2px;
  letter-spacing: 0.16875000298023224px;
  text-align: left;

  color: #000000;
  line-height: 1.5;

  @media (max-width: 768px) {
    font-size: 16px;
    line-height: 22.4px;
  }

  @media (max-width: 600px) {
    font-size: 14px;
    line-height: 20px;
  }
`;

const List = styled.ol`
  font-size: 18px;
  font-weight: 400;
  line-height: 36px;

  color: #000000;
  line-height: 1.5;
  padding-left: 20px;
  margin-bottom: 16px;

  @media (max-width: 768px) {
    font-size: 16px;
    line-height: 28px;
    padding-left: 16px;
  }

  @media (max-width: 600px) {
    font-size: 14px;
    line-height: 24px;
  }
`;

const ListItem = styled.li`
  margin-bottom: 8px;
  @media (max-width: 600px) {
    margin-bottom: 6px;
  }
`;

const StyledSpan = styled.span`
  font-size: 18px;
  font-weight: 600;
  line-height: 36px;

  @media (max-width: 768px) {
    font-size: 16px;
  }

  @media (max-width: 600px) {
    font-size: 14px;
  }
`;
const Note = styled.div`
  background-color: #f2f8fd;

  font-size: 20px;
  font-weight: 600;
  line-height: 28px;

  padding: 12px;
  border-radius: 4px;
  margin-bottom: 16px;
  @media (max-width: 768px) {
    font-size: 18px;
    line-height: 24px;
  }

  @media (max-width: 600px) {
    font-size: 16px;
    padding: 8px;
  }
`;
const StyledNoteSpan = styled.span`
  font-size: 18px;
  font-weight: 400;
  line-height: 28px;
  color: #1b1d21;

  @media (max-width: 768px) {
    font-size: 16px;
  }

  @media (max-width: 600px) {
    font-size: 14px;
  }
`;

const NoteTitle = styled.div`
  font-weight: bold;
  color: #1b1d21;
  margin-bottom: 8px;
  @media (max-width: 600px) {
    font-size: 16px;
  }
`;

const NoteContent = styled.ul`
  list-style-type: disc;
  padding-left: 20px;
  margin: 0;

  @media (max-width: 768px) {
    padding-left: 16px;
  }

  @media (max-width: 600px) {
    padding-left: 12px;
  }
`;

const NoteItem = styled.li`
  color: #000000;
  font-size: 18px;
  font-weight: 400;
  padding-bottom: 8px;

  line-height: 25.2px;
  @media (max-width: 768px) {
    font-size: 16px;
    padding-bottom: 6px;
  }

  @media (max-width: 600px) {
    font-size: 14px;
    padding-bottom: 4px;
  }
`;

const HelpFooter = styled.div`
  margin: 32px auto;
  text-align: center;
  background-color: #f2f8fd;
  border-radius: 8px;
  width: 50%;
  padding: 24px 40px 24px 40px;

  @media (max-width: 768px) {
    width: 80%;
    margin: 32px auto;
    padding: 20px 30px;
  }

  @media (max-width: 600px) {
    width: 80%;
    margin: 32px auto;
    padding: 16px 30px;
  }
`;

const FooterText = styled.h2`
  color: #00225a;
  font-family: Poppins;
  font-size: 32px;
  font-weight: 600;
  line-height: 16px;

  @media (max-width: 768px) {
    font-size: 28px;
  }

  @media (max-width: 600px) {
    font-size: 24px;
  }
`;

const HelpText = styled.p`
  font-family: Inter;
  font-size: 20px;
  font-weight: 400;
  text-align: center;
  color: #00225a;
  @media (max-width: 768px) {
    font-size: 18px;
  }

  @media (max-width: 600px) {
    font-size: 16px;
  }
`;

const ContactButton = styled.button`
  width: 188px;
  height: 56px;
  border-radius: 10px;
  border: none;
  font-family: Inter;
  font-size: 18px;
  font-weight: 500;
  line-height: 21.78px;
  text-align: center;
  background-color: #004988;
  color: #ffffff;
  cursor: pointer;

  @media (max-width: 768px) {
    width: 160px;
    height: 48px;
    font-size: 16px;
  }

  @media (max-width: 600px) {
    width: 140px;
    height: 40px;
    font-size: 14px;
  }
`;
