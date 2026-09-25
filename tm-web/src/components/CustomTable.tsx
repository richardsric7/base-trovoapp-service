"use client";

import { useState, useMemo, useCallback } from "react";
import styled from "styled-components";
import EmptyTableData from "./EmptyTableData";

interface TableColumn {
  title: string;
  dataIndex: string;
  render?: (value: any, record: any) => React.ReactNode;
}

interface TableProps<T> {
  columns: TableColumn[];
  dataSource?: T[];
  totalItems?: number; // Total number of items for pagination
  pageSize?: number; // Number of items per page
  isLoading?: boolean; // Add isLoading prop
  onPageChange?: (page: number) => void;
  onclick?: () => void;
  onRowClick?: (record: T) => void;
}

const CustomTable = <T extends Record<string, any>>({
  columns,
  dataSource,
  totalItems = 0,
  pageSize = 10,
  isLoading = false,
  onPageChange,
  onclick,
  onRowClick,
}: TableProps<T>) => {
  const [currentPage, setCurrentPage] = useState(1);

  const handlePageChange = useCallback(
    (pageNumber: number) => {
      setCurrentPage(pageNumber);
      if (onPageChange) {
        onPageChange(pageNumber);
      }
    },
    [onPageChange]
  );

  const totalPages = useMemo(
    () => Math.ceil(totalItems / pageSize),
    [totalItems, pageSize]
  );

  return (
    <>
      {/* Conditional rendering of Spinner */}
      {isLoading && <h2>Loading ...</h2>}
      {!isLoading && (
        <Table>
          <THead>
            <Tr>
              {columns.map((column) => (
                <Th key={column.dataIndex}>{column.title}</Th>
              ))}
            </Tr>
          </THead>
          {dataSource && dataSource.length > 0 ? (
            <TBody>
              {dataSource.map((record, index) => (
                <Tr
                  // onClick={onclick}
                  key={index}
                  onClick={() => {
                    // console.log("Clicked row:", record);
                    onRowClick?.(record);
                  }}
                  style={{ cursor: "pointer" }}
                >
                  {columns.map((column) => (
                    <Td key={column.dataIndex}>
                      {column.render
                        ? column.render(record[column.dataIndex], record)
                        : record[column.dataIndex]}
                    </Td>
                  ))}
                </Tr>
              ))}
            </TBody>
          ) : (
            <TBody>
              <Tr>
                <EmptyCell colSpan={columns.length}>
                  <EmptyTableData />
                </EmptyCell>
              </Tr>
            </TBody>
          )}
        </Table>
      )}
      {/* Pagination component */}
      {/* {!isLoading && totalItems > pageSize && (
        <Pagination label="truncated pagination navigation">
          <PaginationItems>
            <PaginationArrow
              label="Go to previous page"
              variant="back"
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage === 1}
            />
            <PaginationNumbers>
              {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                (pageNumber) => (
                  <PaginationNumber
                    label=""
                    key={pageNumber}
                    isCurrent={pageNumber === currentPage}
                    onClick={() => handlePageChange(pageNumber)}
                  >
                    {pageNumber}
                  </PaginationNumber>
                )
              )}
            </PaginationNumbers>
            <PaginationArrow
              label="Go to next page"
              variant="forward"
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
            />
          </PaginationItems>
        </Pagination>
      )} */}
    </>
  );
};

export default CustomTable;

const Table = styled.table`
  width: 100%;
`;
const Tr = styled.tr`
  border-bottom: 1px solid #e5e5ef;
  /* display: block; */
  /* justify-content: space-between; */
`;
const TBody = styled.tbody``;
const THead = styled.thead`
  border-bottom: 1px solid #e5e5ef;
`;
const Td = styled.td`
  color: #00225a;
  font-size: 12px;
  font-weight: 400;
  line-height: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #e5e5ef;
`;
const EmptyCell = styled(Td)`
  padding: 0;
`;
const Th = styled.th`
  color: #828282;
  font-weight: 500;
  font-size: 14px;
  text-align: left;
  line-height: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #e5e5ef;
`;
