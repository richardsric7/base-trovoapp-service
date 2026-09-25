"use client";
import Pagination from "@/components/CustomPagination";
import CustomTable from "@/components/CustomTable";
import TruncatedText from "@/hooks/useTruncate";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import { FaEllipsisVertical } from "react-icons/fa6";
import styled from "styled-components";
import Link from "next/link";
import { RiDeleteBin6Line } from "react-icons/ri";
import { BiEdit } from "react-icons/bi";
import DeleteStakeholder from "./DeletStakeholdersModal";
import { useGetPartnersListQuery } from "@/redux/api/tokenizationstakeholders";

const StackholdersTable = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const router = useRouter();
  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedStakeholder, setSelectedStakeholder] = useState<any | null>(
    null
  );
  const { data = [], isLoading } = useGetPartnersListQuery({});

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      render: (_: any, record: any) => {
        return (
          <AssetIdSection>
            <BoldText>
              {record.name ? record.name.charAt(0).toUpperCase() : "?"}
            </BoldText>

            <StackholderName>
              {/* {record.name} */}
              <TruncatedText text={record.name} maxLength={27} />
            </StackholderName>
          </AssetIdSection>
        );
      },
    },

    {
      title: "Type",
      dataIndex: "type",
      render: (_: any, record: any) => {
        return <TypeWrapper>{record.type}</TypeWrapper>;
      },
      key: "type",
    },

    {
      title: "Fee",
      dataIndex: "fee",
      key: "fee",
    },

    {
      title: "Actions",
      dataIndex: "actions",
      render: (_: any, record: any) => {
        const isOpen = activeTulipId === record.compositeId;

        return (
          <ActionsWrapper
            onClick={(e) => {
              e.stopPropagation();
              setActiveTulipId((prev) =>
                prev === record.compositeId ? null : record.compositeId
              );
            }}
          >
            <FaEllipsisVertical />
            {isOpen && (
              <ActionContent onClick={(e) => e.stopPropagation()}>
                <EditLink
                  href={`/settings/tokenizationstakeholders/edit/${record.Id}?type=${record.rawType}`}
                >
                  <BiEdit />
                  Edit Stakeholder
                </EditLink>
                <DeleteBtn
                  onClick={() => {
                    setSelectedStakeholder(record);
                    setIsDeleteModalOpen(true);
                    setActiveTulipId(null);
                  }}
                >
                  {" "}
                  <RiDeleteBin6Line />
                  Delete Stakeholder
                </DeleteBtn>
              </ActionContent>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    if (!data || Array.isArray(data)) return [];
    const custodian = (data?.approved_asset_custodian || []).map(
      (item: any) => ({
        Id: item.id,
        name: item.asset_custodian_name,
        // type: "approved_asset_custodian",
        type: "Asset Custodian",
        rawType: "approved_asset_custodian",
        fee: `${item.fee_percent}%`,
        compositeId: `custodian-${item.id}`,
      })
    );

    const issuingHouse = (data?.asset_issuing_house || []).map((item: any) => ({
      Id: item.id,
      name: item.asset_issuing_house_name,
      // type: "asset_issuing_house",
      type: "Issuing House",
      rawType: "asset_issuing_house",
      fee: `${item.fee_percent}%`,
      compositeId: `asset_issuing_house-${item.id}`,
    }));

    const manager = (data?.asset_manager || []).map((item: any) => ({
      Id: item.id,
      name: item.asset_manager_name,

      type: "Asset Manager",
      rawType: "asset_manager",
      fee: `${item.fee_percent}%`,
      compositeId: `manager-${item.id}`,
    }));

    const rating = (data?.rating_agency || []).map((item: any) => ({
      Id: item.id,
      name: item.agency_name,
      // type: "asset_manager",
      type: "Rating Agency",
      rawType: "rating_agency",
      fee: `${item.fee_percent}%`,
      compositeId: `rating-${item.id}`,
    }));

    const legal = (data?.legal_and_professionals || []).map((item: any) => ({
      Id: item.id,
      name: item.partner_name,
      type: "Legal Agency",
      rawType: "legal_and_professionals",
      fee: `${item.fee_percent}%`,
      compositeId: `
  legal_and_professionals-${item.id}`,
    }));

    const trustee = (data?.trustees || []).map((item: any) => ({
      Id: item.id,
      name: item.trustee_name,
      type: "Trustee",
      rawType: "trustees",
      fee: `${item.fee_percent}%`,
      compositeId: ` trustees-${item.id}`,
    }));

    const legalAdviser = (data?.legal_adviser || []).map((item: any) => ({
      Id: item.id,
      name: item.adviser_name,
      type: "Legal Adviser",
      rawType: "legal_adviser",
      fee: `${item.fee_percent}%`,
      compositeId: `legal_adviser-${item.id}`,
    }));

    const financialAdviser = (data?.financial_adviser || []).map(
      (item: any) => ({
        Id: item.id,
        name: item.adviser_name,
        type: "Financial Adviser",
        rawType: "financial_adviser",
        fee: `${item.fee_percent}%`,
        compositeId: `financial_adviser-${item.id}`,
      })
    );

    return [
      ...custodian,
      ...issuingHouse,
      ...manager,
      ...rating,
      ...legal,
      ...trustee,
      ...legalAdviser,
      ...financialAdviser,
    ].sort((a, b) => b.Id - a.Id);
  }, [data]);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  if (isLoading) {
    return <p>Loading Stakeholders...</p>;
  }
  return (
    <Container>
      <CustomTable
        columns={columns}
        dataSource={paginatedData}
        onRowClick={(record) => {
          router.push(
            `/settings/tokenizationstakeholders/${record.Id}?type=${record.rawType}`
          );
        }}
      />
      <Pagination
        currentPage={currentPage}
        totalCount={dataSource.length}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {isDeleteModalOpen && (
        <DeleteStakeholder
          openModal={isDeleteModalOpen}
          setOpenModal={setIsDeleteModalOpen}
          stakeholderName={selectedStakeholder?.name}
          stakeholderType={selectedStakeholder.rawType}
          stakeholderId={selectedStakeholder.Id}
        />
      )}
    </Container>
  );
};

export default StackholdersTable;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 16%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 8%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 4%;
  }
  // th:nth-child(5),
  // td:nth-child(5) {
  //   width: 8%;
  //   text-align: center;
  // }

  th:nth-child(6),
  td:nth-child(6) {
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

const BoldText = styled.div`
  background-color: #007cdf;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
`;
const AssetIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const TypeWrapper = styled.div`
  font-weight: 400;
  font-size: 12px;
  line-height: 16px;
  letter-spacing: 0.1px;
  text-align: center;
  color: #007cdf;
  background-color: #007cdf1a;
  // width: 150px;
  // height: 32px;
  display: inline-flex;
  gap: 10px;
  border-radius: 8px;
  padding: 8px;
`;

const StackholderName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const EditLink = styled(Link)`
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
const ActionsWrapper = styled.div`
  position: relative;
  display: inline-block;
  cursor: pointer;
  padding: 12px; // Increases clickable area
  margin: -8px; // Cancels visual space shift
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
  z-index: 1000;
`;
