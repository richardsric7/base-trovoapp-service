import { useMemo } from "react";

interface UsePaginationProps {
  totalCount: number;
  pageSize: number;
  siblingCount?: number; // Number of sibling pages to show around the current page
  currentPage: number;
}

export const DOTS = "...";

export const usePagination = ({
  totalCount,
  pageSize,
  siblingCount = 1,
  currentPage,
}: UsePaginationProps): (number | string)[] => {
  const paginationRange = useMemo(() => {
    const totalPageCount = Math.ceil(totalCount / pageSize);
    const totalPageNumbers = siblingCount * 2 + 5;

    // If the total number of pages is less than the max page numbers we want to show
    if (totalPageNumbers >= totalPageCount) {
      return Array.from({ length: totalPageCount }, (_, i) => i + 1);
    }

    const leftSiblingIndex = Math.max(currentPage - siblingCount, 1);
    const rightSiblingIndex = Math.min(
      currentPage + siblingCount,
      totalPageCount
    );

    const showLeftDots = leftSiblingIndex > 2;
    // const showRightDots = rightSiblingIndex < totalPageCount - 2;
    const showRightDots = rightSiblingIndex < totalPageCount - 1;
    const pagination: (number | string)[] = [];

    if (!showLeftDots && showRightDots) {
      const leftRange = Array.from(
        { length: 3 + 2 * siblingCount },
        (_, i) => i + 1
      );
      return [...leftRange, DOTS, totalPageCount];
    }

    if (showLeftDots && !showRightDots) {
      const rightRange = Array.from(
        { length: 3 + 2 * siblingCount },
        (_, i) => totalPageCount - (3 + 2 * siblingCount) + i + 1
      );
      return [1, DOTS, ...rightRange];
    }

    if (showLeftDots && showRightDots) {
      const middleRange = Array.from(
        { length: 2 * siblingCount + 1 },
        (_, i) => leftSiblingIndex + i
      );
      return [1, DOTS, ...middleRange, DOTS, totalPageCount];
    }

    return [1, DOTS, totalPageCount];
  }, [totalCount, pageSize, siblingCount, currentPage]);

  return paginationRange;
};
