"use client";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import CustomTable from "@/components/CustomTable";
import { MdVerified } from "react-icons/md";
import { SearchBar } from "@/components";
import { useRouter } from "next/navigation";
import UsersFilter from "../components/UsersFilter";
import Pagination from "@/components/CustomPagination";
import nigFlag from "@/assets/images/twemoji_flag-nigeria.svg";
import Image from "next/image";
import { useGetUsersQuery } from "@/redux/api/users";
import { useAppSelector } from "@/lib/hooks";
import TruncatedText from "@/hooks/useTruncate";
import goldTag from "@/assets/images/goldtag.svg";
import platinumTag from "@/assets/images/platinumtag.svg";
import diamondTag from "@/assets/images/diamondtag.svg";
import Loader from "@/components/Loader";
import SelectedFiltersComponent from "@/components/SelectedFilter";

interface UserStatusProps {
  isActive: boolean;
}

const UsersTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [selectedFilters, setSelectedFilters] = useState<
    Record<string, string>
  >({});
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [searchInput, setSearchInput] = useState<string>("");
  const {
    data: usersList,
    isLoading,
    isFetching,
    error,
    refetch,
  } = useGetUsersQuery(
    {
      page: currentPage,
      pageSize,
      search: searchTerm,
    },
    {
      skip: false, // Only fetch when needed
      refetchOnMountOrArgChange: false, // Prevent automatic refetching
    }
  );

  const auth = useAppSelector((state) => state.auth);

  const handleRemoveFilter = (key: string) => {
    setSelectedFilters((prev) => {
      const newFilters = { ...prev };
      delete newFilters[key];
      return newFilters;
    });
  };

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchInput(value);

    if (value.trim() === "") {
      setSearchTerm("");
      setCurrentPage(1);
    }
  };

  const handleSearch = () => {
    setSearchTerm(searchInput);
    setCurrentPage(1);
    refetch();
  };
  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber); // This will trigger data fetching for the specific page
  };

  const router = useRouter();

  const columns = [
    {
      title: "Users",
      dataIndex: "user_name",
      render: (_: any, record: any) => {
        return (
          <UserInfoSection>
            <Avatar />
            <div>
              <UserNameWrapper>
                <UserName>
                  {" "}
                  <TruncatedText text={record.user_name} maxLength={15} />
                  {record.kyc_level === 0 ? (
                    ""
                  ) : (
                    <VerifiedIcon>
                      <MdVerified />
                    </VerifiedIcon>
                  )}
                </UserName>
              </UserNameWrapper>
              <UserEmail>
                <TruncatedText
                  text={`${record.first_name} ${record.last_name}`}
                  maxLength={25}
                />
              </UserEmail>
            </div>
          </UserInfoSection>
        );
      },
    },
    {
      title: "Email Address",
      dataIndex: "email",
      key: "email",
      render: (_: any, record: any) => {
        return <TruncatedText text={record.email} maxLength={25} />;
      },
    },
    {
      title: " Location",
      dataIndex: "location",
      render: (_: any, record: any) => {
        return (
          <LocationWrapper>
            {/* <Image src={nigFlag} alt="Nigeria-flag" /> */}
            {record.city}
          </LocationWrapper>
        );
      },
      key: "location",
    },
    {
      title: "Phone number",
      dataIndex: "contact_phone",
      key: "number",
    },
    {
      title: "Registration Date",
      dataIndex: "created_at",
      render: (created_at: number) =>
        new Date(created_at).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
      key: "created_at",
    },

    {
      title: "Acccount Status",
      dataIndex: "suspended",
      key: "suspended",
      render: (_: any, record: any) => (
        <UserStatus isActive={record.suspended === 0}>
          {record.suspended === 0 ? "Active" : "Suspended"}
        </UserStatus>
      ),
    },

    {
      title: "Trovo Plan",
      dataIndex: "trovoPlan",
      // render: (_: any, record: any) => {
      //   const planTagImages = {
      //     "Gold Patron": goldTag,
      //     "Platinum Patron": platinumTag,
      //     "Diamond Patron": diamondTag,
      //   } as const;

      //   type PlanType = keyof typeof planTagImages;

      //   return (
      //     <PlanTag plan={record.trovoPlan}>
      //       <Image
      //         src={planTagImages[record.trovoPlan as PlanType] || goldTag}
      //         alt={`${record.trovoPlan} tag`}
      //       />
      //       {record.trovoPlan || "Gold"}
      //     </PlanTag>
      //   );
      // },
      key: "trovoPlan",
    },
  ];

  const totalUsers = usersList?.data?.total || 0;

  const dataSource = useMemo(() => {
    return (usersList?.data?.data || []).map((user) => ({
      ...user,
      key: user.id,
      user_name: user.user_name,
      email: user.email,
      location: user.city,
      contact_phone: user.mobile,
      created_at: user.created_at,
      suspended: user.suspended,
      trovoPlan: "N/A",
      image_thumbnail: user.image_thumbnail,
    }));
  }, [usersList?.data?.data]);

  if (isLoading) {
    return <UserName>Loading users...</UserName>;
  }

  return (
    <Container>
      <Content>
        <div>
          <PageTitle>Users</PageTitle>
          <SearchBar
            customWidth="320px"
            value={searchInput}
            onChange={handleSearchInputChange}
            handleSearch={handleSearch}
          />
        </div>
        <UsersFilter setSelectedFilters={setSelectedFilters} />
      </Content>

      <SelectedFiltersComponent
        selectedFilters={selectedFilters}
        handleRemoveFilter={handleRemoveFilter}
      />

      {/* <TableContainer> */}
      {searchTerm && isFetching ? (
        <>
          <UserName>Searching users...</UserName>
        </>
      ) : searchTerm && totalUsers === 0 ? (
        <>
          <UserName>No users match your search.</UserName>
        </>
      ) : (
        <>
          <CustomTable
            columns={columns}
            dataSource={dataSource}
            onRowClick={(record) => {
              const destination =
                record.suspended === 1
                  ? `/users/suspendUser/${record.user_name}`
                  : `/users/${record.user_name}`;
              router.push(destination);
            }}
          />

          <Pagination
            totalCount={totalUsers}
            onPageChange={handlePageChange}
            currentPage={currentPage}
            pageSize={pageSize}
            onPageSizeChange={setPageSize}
            isFetching={isFetching}
          />
        </>
      )}
      {/* </TableContainer> */}
    </Container>
  );
};

export default UsersTable;

const Container = styled.section`


  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 15%;
    //  text-align: center;
   }

  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
    //  text-align: center;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
     text-align: center;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 8%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 6%;
    text-align: center;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 8%;
    text-align: center;
  }

  th:nth-child(7),
  td:nth-child(7) {
    width: 4%;
    text-align: center;
  }

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
    text-align center;
  }:
`;

const SelectedFilters = styled.div`
  margin: 20px 0;
  padding: 10px;
  // background-color: #f8f9fa;
  // border-radius: 8px;

  h4 {
    margin: 0 0 8px 0;
    color: #00225a;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    color: #333;
  }
`;

const PageTitle = styled.h2`
  color: #00225a;
  font-size: 24px;
  font-weight: 700;
  line-height: 28px;
  text-align: left;
  margin-bottom: 10px;
`;
const UserNameWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const VerifiedIcon = styled.div`
  color: #007cdf;
  display: flex;
  align-items: center;
  svg {
    width: 12px;
    height: 12px;
  }
`;

const TableAddon = styled.div``;
const Content = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 10px;
`;

const UserStatus = styled.p<UserStatusProps>`
  color: ${(props) => (props.isActive ? "#00A859" : "#BE3800")};
  background-color: ${(props) => (props.isActive ? "#00A8591A" : "#BE38001A")};
  font-weight: 500;
  display: inline-flex;
  padding: 8px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
`;

const StyledButton = styled.div`
  background: none;
  border: none;
  padding: 0;
  margin: 0;
  text-align: left;
  cursor: pointer;
  text-decoration: none;
  display: block;
  color: #00225a;
`;
const UserInfoSection = styled.div`
  display: flex;
  column-gap: 10px;
  align-items: center;
  cursor: pointer;
`;
const Avatar = styled.div`
  width: 40px;
  height: 40px;
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
  text-transform: capitalize;
`;
const UserEmail = styled.p`
  font-size: 12px;
  font-weight: 400;
  line-height: 20px;
  text-align: left;
  margin: 0;
  text-transform: uppercase;
`;

const LocationWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  justify-content: center;
`;

const PlanTag = styled.p<{ plan: string }>`
  color: ${({ plan }) =>
    plan === "Gold Patron"
      ? "#CA8503" // Gold color
      : plan === "Platinum Patron"
      ? "#00A859" // Green for Platinum
      : plan === "Diamond Patron"
      ? "#007CDF" // Blue for Diamond
      : "#CA8503"};

  font-weight: 500;

  background-color: ${({ plan }) =>
    plan === "Gold Patron"
      ? "#FFC54D33"
      : plan === "Platinum Patron"
      ? "#00A8591A"
      : plan === "Diamond Patron"
      ? "#5DADEC26"
      : "#FFC54D33"};

  padding: 4px 8px;
  width: 90px;
  border-radius: 8px;
  text-align: center;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
`;

const StyledSearch = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;
const SingleSearch = styled.div`
  border: 1px solid #007cdf;
  background-color: #acd1ef;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: center;
  border-radius: 16px;
  padding: 8px 16px;
  display: inline-flex;
  align-items: center;
  gap: 6px;

  transition: all 0.2s ease-in-out;
`;

const KeyText = styled.span`
  color: #4f4f4f;
  font-weight: 500;
`;

const ValueText = styled.span`
  color: #007cdf;
`;
const RemoveButton = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  padding: 0;
  margin-left: 4px;
  cursor: pointer;
  color: #007cdf;
  width: 10px;
  height: 10px;
`;
