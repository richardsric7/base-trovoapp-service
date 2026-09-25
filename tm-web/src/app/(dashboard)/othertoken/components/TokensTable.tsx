import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import editIcon from "@/assets/images/fluent_edit-16-regular.svg";
import deletIcon from "@/assets/images/deletIcon.svg";
import Link from "next/link";
import React, { useMemo, useState, useEffect } from "react";
import styled from "styled-components";
import Image from "next/image";
import DeleteAssetModal from "../../assetcuration/components/DeleteAssetModal";
import trov from "@/assets/images/TROVTokenicon.svg";
import { useRouter } from "next/navigation";

const TokensTable = () => {
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [deleteModal, setDeleteModal] = useState<boolean>(false);
  const columns = [
    {
      title: "Asset",
      dataIndex: "asset",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <Image src={trov} alt="trov-icon" width={40} height={40} />
            <div>
              <AssetName>{record.name}</AssetName>
              <AssetSubName>{record.subname}</AssetSubName>
            </div>
          </AssetIdSection>
        );
      },
    },
    {
      title: "Updated",
      dataIndex: "updatedOn",
      key: "updatedOn",
    },
    {
      title: "Added On",
      dataIndex: "addedOn",
      key: "addedOn",
    },
    {
      title: "Added By",
      dataIndex: "addedBy",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <Avatar />
            <div>
              <AssetName>{record.addedByName}</AssetName>
              <AssetSubName>{record.addedByRole}</AssetSubName>
            </div>
          </AssetIdSection>
        );
      },
    },

    {
      title: "",
      dataIndex: "actions",
      width: "20%",
      render: (_: any, record: any) => {
        return (
          <ActionsWrapper onClick={(e) => e.stopPropagation()}>
            <StyledLink href="/othertoken/edittoken">
              <Image src={editIcon} alt="edit" width={16} height={16} />
            </StyledLink>
            <Image
              src={deletIcon}
              alt="delete-icon"
              width={16}
              height={16}
              onClick={() => setDeleteModal(true)}
            />
          </ActionsWrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    return Array.from({ length: 100 }, (_, index) => ({
      tokenId: index + 1,
      name: `Trov ${index + 1}`,
      subname: `25 NGN`,
      updatedOn: "2h ago",
      addedOn: "23 Oct, 2023",
      addedByName: `Obi Enechi ${index + 1}`,
      addedByRole: "admin",
    }));
  }, []);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  return (
    <Container>
      <CustomTable
        columns={columns}
        dataSource={paginatedData}
        onRowClick={(record) => {
          router.push(`/othertoken/${record.tokenId}`);
        }}
      />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {deleteModal && (
        <DeleteAssetModal
          openModal={deleteModal}
          setOpenModal={setDeleteModal}
        />
      )}
    </Container>
  );
};

export default TokensTable;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 10%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 10%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
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

const StyledLink = styled(Link)`
  text-decoration: none;
  display: block;
  color: #00225a;
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

const ActionButtons = styled.div`
  display: flex;
  gap: 8px;

  button {
    padding: 4px 8px;
    border: none;
    cursor: pointer;
    background-color: #00225a;
    color: white;
    border-radius: 4px;
    font-size: 12px;
  }
`;
const ActionsWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;
