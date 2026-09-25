"use client";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import React, { useEffect, useMemo, useState } from "react";
import styled from "styled-components";

import { FaEllipsisVertical } from "react-icons/fa6";
// import { MdVerified } from "react-icons/md";
import {
  useGetMintingUsersQuery,
  useSearchMintingUsersQuery,
} from "@/redux/api/minting";

import { RiDeleteBin6Line } from "react-icons/ri";
import DeleteModal from "./DeleteModal";
interface InitiatorsTableProps {
  search: string;
}

const InitiatorsTable: React.FC<InitiatorsTableProps> = ({ search }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedId, setSelectedId] = useState<any | null>(null);
  const [activeRowId, setActiveRowId] = useState<number | null>(null);
  useEffect(() => {
    setCurrentPage(1);
  }, [search]);

  const isSearching = !!search && search.trim().length > 0;

  const searchResult = useSearchMintingUsersQuery(
    { query: search, role: "initiator", page: currentPage, page_size: pageSize },
    { skip: !isSearching },
  );
  const listResult = useGetMintingUsersQuery(
    { role: "initiator", page: currentPage, page_size: pageSize },
    { skip: isSearching },
  );
  const { data, isLoading } = isSearching ? searchResult : listResult;

  // const { data, isLoading } = useGetMintingUsersQuery({ role: "initiator" });

  const columns = [
    {
      title: "User",
      dataIndex: "user",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <Avatar />

            <div>
              <AssetContent>
                <UserName>{record.username}</UserName>
                {/* <MdVerified color="#007CDF" /> */}
              </AssetContent>
            </div>
          </AssetIdSection>
        );
      },
    },
    {
      title: "Email Address",
      dataIndex: "username",
      key: "email",
      render: (_: any, record: any) => <FullName> --</FullName>,
    },

    {
      title: "Added On",
      dataIndex: "addedOn",
      render: (created_at: number) =>
        new Date(created_at).toLocaleDateString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
        }),
      key: "addedOn",
    },

    {
      title: "Actions",
      dataIndex: "actions",

      render: (_: any, record: any) => {
        const isOpen = activeTulipId === record.userId;
        return (
          <ActionsWrapper
            onClick={(e) => {
              e.stopPropagation();
              setActiveTulipId((prev) =>
                prev === record.userId ? null : record.userId
              );
            }}
          >
            <FaEllipsisVertical style={{ cursor: "pointer" }} />
            {isOpen && (
              <ActionContent>
                <DeleteBtn
                  onClick={() => {
                    setSelectedId(record);
                    setIsDeleteModalOpen(true);
                    setActiveTulipId(null);
                  }}
                >
                  {" "}
                  <RiDeleteBin6Line />
                  Delete Initiator
                </DeleteBtn>
              </ActionContent>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  const mapUser = (item: any) => ({
    userId: item.id,
    username: item.username || "Obi",
    addedOn: item.created_at || "23 Sep 2023, ",
    lastactivity: item.updated_at,
    role: item.role,
  });

  const dataSource = useMemo(() => {
    if (!data) return [];
    if ("data" in data && Array.isArray(data.data?.data)) {
      return data.data.data.map(mapUser);
    }
    if ("users" in data && Array.isArray(data.users)) {
      return data.users.map(mapUser);
    }
    return [];
  }, [data]);

  if (isLoading) {
    return <UserName>Loading Initiators...</UserName>;
  }

  if (!isLoading && isSearching && dataSource.length === 0) {
    return <UserName>No results found.</UserName>;
  }

  return (
    <Container>
      <CustomTable columns={columns} dataSource={dataSource} />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />
      {isDeleteModalOpen && (
        <DeleteModal
          openModal={isDeleteModalOpen}
          setOpenModal={setIsDeleteModalOpen}
          id={selectedId?.userId ?? 0}
          username={selectedId?.username ?? ""}
          role={selectedId?.role ?? ""}
        />
      )}
    </Container>
  );
};

export default InitiatorsTable;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 20%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 18%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 16%;
    text-align: center;
  }

  th:nth-child(4),
  td:nth-child(4) {
    width: 4%;
    text-align: center;
  }

  th {
    color: #828282;
    font-size: 14px;
    font-weight: 500;
  }
  td {
    color: #00225a;
    font-weight: 400;
  }
`;

const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const UserName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const FullName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const ActionsWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const AssetContent = styled.div`
  display: flex;
  align-items: center;
  gap: 1px;
`;
const ActionContent = styled.div`
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 8px;
  width: 180px;
  background-color: #ffffff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
  z-index: 1;
`;

const DeleteBtn = styled.div`
  font-weight: 400;
  font-size: 14px;
  line-height: 28px;
  letter-spacing: 0%;
  display: flex;
  align-items: center;
  gap: 2px;
  color: #be3800;
`;
