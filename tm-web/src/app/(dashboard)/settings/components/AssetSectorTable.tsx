"use client";

import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import { BiEdit } from "react-icons/bi";
import { FaEllipsisVertical } from "react-icons/fa6";
import { RiDeleteBin6Line } from "react-icons/ri";
import styled from "styled-components";
import DeleteSectorModal from "./DeleteSectorModal";
import TruncatedText from "@/hooks/useTruncate";
import ViewSubSectorModal from "./ViewSubSectorModal";
import viewIcon from "@/assets/images/note-text.svg";
import Image from "next/image";

const AssetSectorTable = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isViewModal, setIsViewModal] = useState(false);
  const columns = [
    {
      title: "Sector",
      dataIndex: "sector",
      key: "sector",
    },
    {
      title: "Sub Sectors",
      dataIndex: "subSector",
      render: (_: any, record: any) => {
        return (
          <SubSection>
            <TruncatedText text={record.subSector} maxLength={20} />
            <BoldText>+3</BoldText>
          </SubSection>
        );
      },
    },

    {
      title: "Date",
      dataIndex: "date",
      key: "date",
    },

    {
      title: "Added By",
      dataIndex: "addedBy",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <Avatar />

            <div>
              <AssetName>{record.admin}</AssetName>
              <AssetSubName>{record.adminEmail}</AssetSubName>
            </div>
          </AssetIdSection>
        );
      },
    },

    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: any, record: any) => {
        const isOpen = activeTulipId === record.id;

        return (
          <ActionsWrapper
            onClick={(e) => {
              e.stopPropagation();
              setActiveTulipId((prev) =>
                prev === record.id ? null : record.id
              );
            }}
          >
            <FaEllipsisVertical />
            {isOpen && (
              <ActionContent>
                <EditLink
                  onClick={() => {
                    setIsViewModal(true);
                    setActiveTulipId(null);
                  }}
                >
                  <Image
                    src={viewIcon}
                    alt="view-icon"
                    width={14}
                    height={14}
                  />
                  View Subsectors
                </EditLink>
                <EditLink
                  onClick={() => {
                    router.push(`/settings/assetsector/edit/${record.id}`);
                    setActiveTulipId(null);
                  }}
                >
                  <BiEdit color="#00225A" size={14} />
                  Edit Sector
                </EditLink>
                <DeleteBtn
                  onClick={() => {
                    setIsDeleteModalOpen(true);
                    setActiveTulipId(null);
                  }}
                >
                  {" "}
                  <RiDeleteBin6Line color="#BE3800" size={14} />
                  Delete Sector
                </DeleteBtn>
              </ActionContent>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      id: index + 1,
      sector: "Real Estate",
      subSector: "industrial properties, Residential",
      date: "23 Sep 2023",
      admin: "John Doe",
      adminEmail: "johndoe@gmail.com",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <Container>
      <CustomTable columns={columns} dataSource={paginatedData} />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {isViewModal && (
        <ViewSubSectorModal
          isOpen={isViewModal}
          onClose={() => setIsViewModal(false)}
        />
      )}

      {isDeleteModalOpen && (
        <DeleteSectorModal
          openModal={isDeleteModalOpen}
          setOpenModal={setIsDeleteModalOpen}
        />
      )}
    </Container>
  );
};

export default AssetSectorTable;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 14%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 14%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 12%;
  }
  th:nth-child(5),
  td:nth-child(5) {
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

const SubSection = styled.div`
  display: flex;
  align-items: center;
`;

const Avatar = styled.div`
  width: 36px;
  height: 36px;
  border-radius: 20px;
  background-color: rebeccapurple;
`;

const BoldText = styled.span`
  background-color: #acd1ef;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #191919;
  font-weight: 600;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
`;

const AssetName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const AssetSubName = styled.p`
  font-size: 10px;
  font-weight: 400;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const EditLink = styled.div`
  font-weight: 400;
  font-size: 14px;
  line-height: 28px;
  letter-spacing: 0%;
  display: flex;
  align-items: center;
  gap: 2px;
  color: #00225a;
  text-decoration: none;
`;

const ActionsWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px;
  margin: -8px;
`;

const ActionContent = styled.div`
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 8px;
  width: 200px;
  background-color: #ffffff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
  z-index: 1000;
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
