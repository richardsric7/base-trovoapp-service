import { SearchBar } from "@/components";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import AdminUserTable from "./AdminUsersTable";
import InviteAdminUser from "./InviteAdminUser";
import AdminFilter from "./AdminFilter";
import { useListAdminsQuery } from "@/redux/api/admin/admin";
const AdminUsersTab = () => {
  const { data } = useListAdminsQuery({});
  const [searchTerm, setSearchTerm] = useState<string>("");

  return (
    <Wrapper>
      <Header>
        <TitleSection>
          <Title>Admin Users</Title>
          <TotalCount>
            Total: <span>{data?.pagination?.total || 0}</span>
          </TotalCount>
        </TitleSection>
        <InviteAdminUser />
      </Header>
      <SubHeader>
        <SearchBar
          value={searchTerm}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setSearchTerm(e.target.value)
          }
        />
        <AdminFilter />
      </SubHeader>
      <AdminUserTable searchTerm={searchTerm} />
    </Wrapper>
  );
};

export default AdminUsersTab;

const Wrapper = styled.section`
  background-color: #ffffff;
  padding: 8px;
  border-radius: 24px;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;

const TitleSection = styled.div``;

const Title = styled.h1`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  line-height: 30px;
`;

const TotalCount = styled.p`
  color: #828282;
  font-size: 16px;
  font-weight: 500;
  margin: 0;
  line-height: 28px;
`;
