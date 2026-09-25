"use client";
import CustomTable from "@/components/CustomTable";
import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { useMemo, useState } from "react";
import styled from "styled-components";
import nigflag from "@/assets/images/emojione_flag-for-nigeria.svg";
import Image from "next/image";
import { BiEdit } from "react-icons/bi";
import { RiDeleteBin6Line } from "react-icons/ri";
import Pagination from "@/components/CustomPagination";
import { FaEllipsisVertical, FaEye } from "react-icons/fa6";
import DeleteCountry from "./DeleteCountry";

import { useGetCountryListQuery } from "@/redux/api/country";
import { IoEyeOutline } from "react-icons/io5";
const CountriesTable = () => {
  const { data = [], isLoading } = useGetCountryListQuery();
  const router = useRouter();
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [activeTulipId, setActiveTulipId] = useState<number | null>(null);

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedCountry, setSelectedCountry] = useState<any | null>(null);
  const columns = [
    {
      title: "Country",
      dataIndex: "id",
      render: (_: any, record: any) => {
        return (
          <CountryIdSection>
            <Avatar />
            <div>
              <CountryCode>{record.countryCode}</CountryCode>
              <CountryName>{record.countryName}</CountryName>
            </div>
          </CountryIdSection>
        );
      },
    },
    // {
    //   title: "ID",
    //   dataIndex: "countryId",
    //   key: "countryId",
    // },
    {
      title: "Fiat Label",
      dataIndex: "FiatLabel",
      key: "FiatLabel",
    },

    {
      title: "Fiat Glyph",
      dataIndex: "FiatGlyph",
      key: "FiatGlyph",
    },

    {
      title: "Quote Currency Code",
      dataIndex: "currencyCode",
      key: "currencyCode",
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
                <ViewDetails href={`/countries/${record.id}`}>
                  <IoEyeOutline color="#292D32" />
                  View Details
                </ViewDetails>
                <EditLink href={`/countries/edit/${record.id}`}>
                  <BiEdit />
                  Edit Country
                </EditLink>
                <DeleteBtn
                  onClick={() => {
                    setSelectedCountry(record);
                    setIsDeleteModalOpen(true);
                    setActiveTulipId(null);
                  }}
                >
                  {" "}
                  <RiDeleteBin6Line />
                  Delete Country
                </DeleteBtn>
              </ActionContent>
            )}
          </ActionsWrapper>
        );
      },
    },
  ];

  const dataSource = useMemo(() => {
    if (!data) return [];

    return data.map((item) => ({
      id: item.id,
      countryName: item.country_name,
      countryCode: item.country_code,
      countryId: item.id,
      FiatLabel: item.fiat_label,
      FiatGlyph: item.fiat_glyph,
      currencyCode: item.quote_currency_code,
    }));
  }, [data]);

  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return dataSource.slice(startIndex, endIndex);
  }, [currentPage, pageSize, dataSource]);

  if (isLoading) return <div>Loading countries...</div>;

  return (
    <Container>
      <CustomTable
        columns={columns}
        dataSource={paginatedData}
        onRowClick={(record) => {
          router.push(`/countries/${record.id}`);
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
        <DeleteCountry
          openModal={isDeleteModalOpen}
          setOpenModal={setIsDeleteModalOpen}
          countryName={selectedCountry?.countryName}
          countryId={selectedCountry?.id}
        />
      )}
    </Container>
  );
};

export default CountriesTable;

const Container = styled.section`
  table {
    width: 100%;
    table-layout: fixed;
  }

  th:nth-child(1),
  td:nth-child(1) {
    width: 15%;
  }
  th:nth-child(2),
  td:nth-child(2) {
    width: 12%;
  }
  th:nth-child(3),
  td:nth-child(3) {
    width: 14%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
  }
  th:nth-child(5),
  td:nth-child(5) {
    width: 6%;
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

const CountryIdSection = styled.div`
  display: flex;
  column-gap: 4px;
  align-items: center;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rebeccapurple;
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

const ViewDetails = styled(Link)`
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
const CountryName = styled.p`
  font-size: 14px;
  font-weight: 500;
  line-height: 20px;
  text-align: left;
  margin: 0;
  color: #00225a;
`;

const CountryCode = styled.p`
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
