import React, { useEffect, useState } from "react";
import styled from "styled-components";
import {
  FaAngleRight,
  FaAngleLeft,
  FaAnglesLeft,
  FaAnglesRight,
  FaAngleDown,
} from "react-icons/fa6";
import { usePagination, DOTS } from "@/hooks/usePagination";
import Loader from "./Loader";

type DropdownOption = number;

interface CustomDropdownProps {
  options: DropdownOption[];
  value: DropdownOption;
  onSelect: (value: DropdownOption) => void;
}

// Custom Dropdown Component
const CustomDropdown: React.FC<CustomDropdownProps> = ({
  options,
  value,
  onSelect,
}) => {
  const [isOpen, setIsOpen] = useState<boolean>(false);

  return (
    <DropdownContainer>
      <SelectedValue onClick={() => setIsOpen(!isOpen)}>
        {value}
        <FaAngleDown />
      </SelectedValue>
      {isOpen && (
        <OptionsList>
          {options.map((option) => (
            <OptionItem
              key={option}
              onClick={() => {
                onSelect(option);
                setIsOpen(false);
              }}
            >
              {option}
            </OptionItem>
          ))}
        </OptionsList>
      )}
    </DropdownContainer>
  );
};

interface PaginationProps {
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  totalCount: number;
  siblingCount?: number;
  currentPage: number;
  pageSize: number;
  className?: string;
  isFetching?: boolean;
}

const Pagination: React.FC<PaginationProps> = ({
  onPageChange,
  totalCount,
  siblingCount = 1,
  currentPage,
  pageSize,
  className,
  onPageSizeChange,
  isFetching,
}) => {
  const paginationRange = usePagination({
    currentPage,
    totalCount,
    siblingCount,
    pageSize,
  });

  const [isPageChanging, setIsPageChanging] = useState<boolean>(false);
  useEffect(() => {
    if (!isFetching) setIsPageChanging(false); // Reset loader when fetch completes
  }, [isFetching]);

  if (!paginationRange) {
    return null;
  }

  const lastPage = Math.ceil(totalCount / pageSize);

  const pageDrop = [10, 20, 50, 100];

  const handlePageSizeChange = (newPageSize: number): void => {
    onPageSizeChange(newPageSize);
    onPageChange(currentPage);
  };
  return (
    <>
      {(isPageChanging || isFetching) && <Loader />}
      <Container>
        <Text>
          Total <StyledSpan> {totalCount}</StyledSpan>{" "}
          {totalCount === 1 ? "item" : "items"}
        </Text>
        <PaginationWrapper className={className}>
          <PaginationItem
            disabled={currentPage === 1}
            onClick={() => onPageChange(1)}
          >
            <FaAnglesLeft size={12} />
          </PaginationItem>
          <PaginationItem
            disabled={currentPage === 1}
            onClick={() => onPageChange(currentPage - 1)}
          >
            <FaAngleLeft size={12} />
          </PaginationItem>
          {paginationRange.map((pageNumber, idx) =>
            pageNumber === DOTS ? (
              <Dots key={idx}>&#8230;</Dots>
            ) : (
              <PaginationItem
                key={`page-${pageNumber}`}
                selected={pageNumber === currentPage}
                onClick={() => onPageChange(pageNumber as number)}
              >
                {pageNumber}
              </PaginationItem>
            )
          )}
          <PaginationItem
            disabled={currentPage === lastPage}
            onClick={() => onPageChange(currentPage + 1)}
          >
            <FaAngleRight size={12} />
          </PaginationItem>
          <PaginationItem
            disabled={currentPage === lastPage}
            onClick={() => onPageChange(lastPage)}
          >
            <FaAnglesRight size={12} />
          </PaginationItem>
        </PaginationWrapper>
        <Text>
          {" "}
          <StyledBtn>
            <CustomDropdown
              options={pageDrop}
              value={pageSize}
              onSelect={handlePageSizeChange}
            />
          </StyledBtn>
          Items per page
        </Text>
      </Container>
    </>
  );
};

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
`;
const PaginationWrapper = styled.ul`
  display: flex;
  list-style: none;
  padding: 0;
  margin: 20px 0;
  justify-content: center;
  align-items: center;
`;

const PaginationItem = styled.li<{ selected?: boolean; disabled?: boolean }>`
  margin: 0 5px;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background-color: ${({ selected, disabled }) =>
    disabled ? "#F2F6F9" : selected ? "#007CDF" : "#F2F6F9"};
  color: ${({ selected, disabled }) =>
    disabled ? "#a0a0a0" : selected ? "#fff" : "#00225A"};
  cursor: ${({ disabled }) => (disabled ? "not-allowed" : "pointer")};
  pointer-events: ${({ disabled }) => (disabled ? "none" : "auto")};
  transition: background-color 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  font-size: 14px;
  &:hover {
    background-color: ${({ selected, disabled }) =>
      disabled ? "#e0e0e0" : selected ? "#007CDF" : "#D1E3F3"};
  }

  svg {
    font-size: 14px;
  }
`;

const Dots = styled.li`
  padding: 8px 12px;
  color: #00225a;
  background-color: #f2f6f9;
  border-radius: 8px;
`;

const Text = styled.div`
  font-size: 14px;
  font-weight: 400;
  line-height: 17.07px;
  color: #00225a;
  display: flex;
  align-items: center;
  gap: 4px;
`;

const StyledSpan = styled.span`
  // padding: 0 2px;
`;

const StyledBtn = styled.span`
  background-color: #f2f6f9;
  padding: 0 12px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 14px;
  color: #00225a;
  gap: 6px;
  cursor: pointer;
`;

const DropdownContainer = styled.div`
  position: relative;
  // width: 40px;
`;

const SelectedValue = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 10px;
  gap: 6px;
  border-radius: 4px;
  cursor: pointer;
`;

const OptionsList = styled.div`
  position: absolute;
  top: 100%;
  left: 0;
  width: 100%;
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #ccc;
  border-top: none;
  background-color: white;
  z-index: 10;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background-color: #e0e0e0;
    border-radius: 10px;
  }
`;

const OptionItem = styled.div`
  padding: 5px 10px;
  cursor: pointer;
  color: #313131;
  &:hover {
    background-color: #f0f0f0;
  }
`;

export default Pagination;
