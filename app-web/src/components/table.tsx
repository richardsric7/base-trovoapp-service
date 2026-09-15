import { ReactNode, useEffect, useState } from 'react';
import styles from './table.module.css';

export interface ColumnDef<T> {
  key: keyof T | string;
  header: string;
  width?: string;
  render?: (row: T) => ReactNode;
  className?: string;
}
// interface PaginationProps {
//   hasNextPage: boolean;
//   totalPages: number;
//   totalCount: number;
//   nextPage: number;
//   hasPreviousPage: boolean;
// }

export interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  onRowClick?: (row: T) => void;
  onRowAction?: (row: T, action: string) => void;
  itemsPerPageOptions?: number[];
  defaultItemsPerPage?: number;
  paginate?: boolean;
  pagination?: {
    currentPage: number;
    totalItems: number;
    itemsPerPage: number;
    onPageChange: (page: number) => void;
  };
  showSearch?: boolean;
  searchPlaceholder?: string;
  className?: string;
  emptyMessage?: string;
  onSearch?: (query: string) => void;
  // pagination?: PaginationProps;
  // currentPage?: number;
  // limit?: number;
  // onPageChange?: (page: number) => void;
  // onLimitChange?: (limit: number) => void;
}

export const CustomTable = <T = ReactNode,>({
  data,
  columns,
  onRowClick,
  className = '',
  emptyMessage = 'No data available',
  paginate = false,
  defaultItemsPerPage = 5,
  pagination,
  // pagination,
  // currentPage,
  // limit,
  // onPageChange,
  // onLimitChange,
}: DataTableProps<T>) => {
  const [currentPage, setCurrentPage] = useState(1);
  const isPaginated = paginate || Boolean(pagination);
  const itemsPerPage = pagination?.itemsPerPage ?? defaultItemsPerPage;
  const page = pagination?.currentPage ?? currentPage;
  const totalItems = pagination?.totalItems ?? data.length;
  const totalPages = Math.max(1, Math.ceil(totalItems / itemsPerPage));
  const displayedData = paginate && !pagination
    ? data.slice(
        (page - 1) * itemsPerPage,
        page * itemsPerPage,
      )
    : data;

  useEffect(() => {
    if (!pagination) {
      setCurrentPage((current) => Math.min(current, totalPages));
    }
  }, [pagination, totalPages]);

  const changePage = (nextPage: number) => {
    if (pagination) {
      pagination.onPageChange(nextPage);
      return;
    }
    setCurrentPage(nextPage);
  };

  return (
    <div className={`p-4 bg-[#F7FAFC] rounded-lg ${className}`}>
      <div className="overflow-x-auto">
        <table className="min-w-full text-left text-sm">
          <thead>
            <tr className="text-gray-500 capitalize text-xs border-b">
              {columns.map((column) => (
                <th
                  key={column.key as string}
                  style={column.width ? { width: column.width } : {}}
                  className={`px-4 py-3 font-medium ${column.className || ''}`}
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className={`divide-y divide-gray-100 ${styles.tTody}`}>
            {displayedData.length > 0 ? (
              displayedData.map((row, index) => (
                <tr
                  key={index}
                  onClick={() => onRowClick && onRowClick(row)}
                  className={`group bg-white hover:bg-primary-100 transition rounded-lg ${
                    onRowClick ? 'cursor-pointer' : ''
                  }`}
                >
                  {columns.map((column) => (
                    <td
                      key={`${index}-${column.key as string}`}
                      className={`px-4 py-4 whitespace-nowrap group-hover:bg-primary-100 ${
                        column.className || ''
                      }`}
                    >
                      {column.render
                        ? column.render(row)
                        : (row[column.key as keyof T] as ReactNode)}
                    </td>
                  ))}
                </tr>
              ))
            ) : (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-4 py-6 text-center text-gray-400"
                >
                  {emptyMessage}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {isPaginated && totalItems > itemsPerPage && (
        <div className="flex items-center justify-between px-4 py-3 text-sm text-[#3a6787]">
          <span>
            Showing {(page - 1) * itemsPerPage + 1}–
            {Math.min(page * itemsPerPage, totalItems)} of {totalItems}
          </span>
          <div className="flex gap-2">
            <button
              type="button"
              className="rounded border border-[#d3dfe9] px-3 py-1 disabled:cursor-not-allowed disabled:opacity-50"
              disabled={page === 1}
              onClick={() => changePage(page - 1)}
            >
              Previous
            </button>
            <button
              type="button"
              className="rounded border border-[#d3dfe9] px-3 py-1 disabled:cursor-not-allowed disabled:opacity-50"
              disabled={page === totalPages}
              onClick={() => changePage(page + 1)}
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
