"use client";
import Image from "next/image";
import React from "react";
import styled from "styled-components";

interface SharedAccessProps {
  access: {
    name: string;
    username: string;
    image: string;
  }[];
}

const SharedAccess: React.FC<SharedAccessProps> = ({ access }) => {
  return (
    <Container>
      <Header>Shared Access</Header>
      {access.map((user, index) => (
        <React.Fragment key={index}>
          <Card>
            <CardDetails>
              <Avatar></Avatar>
              <Details>
                <Name>{user.name}</Name>
                <Username>{user.username}</Username>
              </Details>
            </CardDetails>
            <Roles>
              <Role>Initiator</Role>
              <Role>Approver</Role>
            </Roles>
          </Card>
          {index < access.length - 1 && <Divider />}
        </React.Fragment>
      ))}
    </Container>
  );
};

export default SharedAccess;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  background-color: #fff;
  margin-top: 20px;
  border-radius: 24px;
  padding: 16px;
`;

const Header = styled.h2`
  font-size: 20px;
  font-weight: 600;
  line-height: 28px;
  color: #00225a;
  padding: 15px;
`;

const Card = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const CardDetails = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
`;

const Details = styled.div`
  //   flex-grow: 1;
`;

const Name = styled.h3`
  font-size: 14px;
  font-weight: 400;
  line-height: 25.27px;
  color: #00225a;
  margin: 0;
`;

const Username = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 25.27px;
  color: #828282;
  margin: 0;
`;

const Roles = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

const Role = styled.p`
  font-size: 12px;
  font-weight: 400;
  color: #00225a;
  margin: 0;
  background-color: #f2f6f9;
  padding: 6px 8px;
  border-radius: 8px;
`;

const Divider = styled.div`
  height: 1px;
  background-color: #e5e5ef;
  margin: 2px 0;
`;
